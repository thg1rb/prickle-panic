// Package httpapi provides the same-origin Guest Session HTTP adapter.
package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/thg1rb/prickle-panic/apps/server/internal/clientip"
	"github.com/thg1rb/prickle-panic/apps/server/internal/session"
)

const (
	guestSessionPath = "/api/v1/guest-session"
	maxRequestBytes  = 4 << 10
	prodCookieName   = "__Host-prickle_guest"
	devCookieName    = "prickle_guest_dev"
)

type Mode string

const (
	Development Mode = "development"
	Production  Mode = "production"
)

type SessionService interface {
	Create(rawUsername, sourceIP string) (session.Created, error)
	CheckCreationCapacity(sourceIP string) error
	LookupAuthorized(credential string) (session.Authorization, error)
	RecoverAuthorized(credential string) (session.Authorization, error)
}

type ClientIPResolver interface {
	Resolve(remoteAddr, forwardedFor string) (netip.Addr, error)
}

// Config requires explicit environment and browser-origin selection. Cookie
// names and security attributes are derived from Mode and cannot be weakened
// by callers.
type Config struct {
	Mode     Mode
	Origin   string
	Sessions SessionService
	ClientIP ClientIPResolver
}

type Handler struct {
	mode     Mode
	origin   string
	sessions SessionService
	clientIP ClientIPResolver
}

type sessionResponse struct {
	GuestID  string    `json:"guestId"`
	Username string    `json:"username"`
	Expires  time.Time `json:"expiresAt"`
}

type errorResponse struct {
	Error struct {
		Code string `json:"code"`
	} `json:"error"`
}

func New(cfg Config) (*Handler, error) {
	if cfg.Sessions == nil || cfg.ClientIP == nil {
		return nil, fmt.Errorf("session service and client IP resolver are required")
	}
	if cfg.Mode != Development && cfg.Mode != Production {
		return nil, fmt.Errorf("APP_ENV must be explicitly set to development or production")
	}
	u, err := url.Parse(cfg.Origin)
	if err != nil || u == nil || u.Host == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" || u.Opaque != "" {
		return nil, fmt.Errorf("APP_ORIGIN must be an absolute origin without path, query, or fragment")
	}
	switch cfg.Mode {
	case Production:
		if !strings.EqualFold(u.Scheme, "https") {
			return nil, fmt.Errorf("production APP_ORIGIN must use https")
		}
	case Development:
		if !strings.EqualFold(u.Scheme, "http") || !isLoopbackHost(u.Hostname()) {
			return nil, fmt.Errorf("development APP_ORIGIN must use HTTP on localhost or a loopback address")
		}
	}
	if cfg.Origin != u.Scheme+"://"+u.Host {
		return nil, fmt.Errorf("APP_ORIGIN must use canonical scheme and host spelling")
	}
	return &Handler{mode: cfg.Mode, origin: cfg.Origin, sessions: cfg.Sessions, clientIP: cfg.ClientIP}, nil
}

func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	address, err := netip.ParseAddr(host)
	return err == nil && address.IsLoopback()
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if r.URL.Path != guestSessionPath {
		h.writeError(w, http.StatusNotFound, "REQUEST_INVALID")
		return
	}
	switch r.Method {
	case http.MethodPost:
		h.createOrReuse(w, r)
	case http.MethodGet:
		h.recover(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		h.writeError(w, http.StatusMethodNotAllowed, "REQUEST_INVALID")
	}
}

func (h *Handler) createOrReuse(w http.ResponseWriter, r *http.Request) {
	if !h.allowedOrigin(r, true) {
		h.writeError(w, http.StatusForbidden, "ORIGIN_NOT_ALLOWED")
		return
	}
	cookieValue, _, cookieValid := h.readCredential(r)
	var existing session.Authorization
	if cookieValid {
		var err error
		existing, err = h.sessions.LookupAuthorized(cookieValue)
		if err != nil {
			if !errors.Is(err, session.ErrUnauthenticated) {
				h.writeSessionError(w, err)
				return
			}
			cookieValid = false
		}
	}

	var source string
	if !cookieValid {
		var err error
		source, err = h.resolveSource(r)
		if err != nil {
			h.writeError(w, http.StatusServiceUnavailable, "SESSION_UNAVAILABLE")
			return
		}
		if err = h.sessions.CheckCreationCapacity(source); err != nil {
			h.writeSessionError(w, err)
			return
		}
	}

	rawUsername, err := h.decodeUsername(w, r)
	if err != nil {
		return
	}
	username, err := session.NormalizeUsername(rawUsername)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "USERNAME_INVALID")
		return
	}

	if cookieValid {
		// Reauthorize after body processing so an expired or revoked session
		// cannot be returned from the earlier preflight lookup.
		existing, err = h.sessions.LookupAuthorized(cookieValue)
		if err == nil {
			h.writeSession(w, http.StatusOK, existing.Identity, existing.AbsoluteRemaining, existing.AbsoluteExpiresAt, cookieValue)
			return
		}
		if !errors.Is(err, session.ErrUnauthenticated) {
			h.writeSessionError(w, err)
			return
		}
		source, err = h.resolveSource(r)
		if err != nil {
			h.writeError(w, http.StatusServiceUnavailable, "SESSION_UNAVAILABLE")
			return
		}
		if err = h.sessions.CheckCreationCapacity(source); err != nil {
			h.writeSessionError(w, err)
			return
		}
	}

	created, err := h.sessions.Create(username, source)
	if err != nil {
		h.writeSessionError(w, err)
		return
	}
	credential := created.Credential.CookieValue()
	// Re-read authorization to account for time spent validating and allocating;
	// cookie expiry is taken from the manager's immutable absolute deadline.
	authorized, err := h.sessions.LookupAuthorized(credential)
	if err != nil {
		h.writeSessionError(w, err)
		return
	}
	h.writeSession(w, http.StatusCreated, authorized.Identity, authorized.AbsoluteRemaining, authorized.AbsoluteExpiresAt, credential)
}

func (h *Handler) recover(w http.ResponseWriter, r *http.Request) {
	if !h.allowedOrigin(r, false) {
		h.writeError(w, http.StatusForbidden, "ORIGIN_NOT_ALLOWED")
		return
	}
	credential, present, valid := h.readCredential(r)
	if !valid {
		if present {
			h.clearCookie(w)
		}
		h.writeError(w, http.StatusUnauthorized, "SESSION_REQUIRED")
		return
	}
	authorized, err := h.sessions.RecoverAuthorized(credential)
	if err != nil {
		if errors.Is(err, session.ErrUnauthenticated) {
			h.clearCookie(w)
			h.writeError(w, http.StatusUnauthorized, "SESSION_REQUIRED")
			return
		}
		h.writeSessionError(w, err)
		return
	}
	h.writeSession(w, http.StatusOK, authorized.Identity, authorized.AbsoluteRemaining, authorized.AbsoluteExpiresAt, credential)
}

func (h *Handler) decodeUsername(w http.ResponseWriter, r *http.Request) (string, error) {
	body := http.MaxBytesReader(w, r.Body, maxRequestBytes)
	data, err := io.ReadAll(body)
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			h.writeError(w, http.StatusRequestEntityTooLarge, "REQUEST_TOO_LARGE")
		} else {
			h.writeError(w, http.StatusBadRequest, "REQUEST_INVALID")
		}
		return "", errRequestWritten
	}
	if !utf8.Valid(data) {
		h.writeError(w, http.StatusBadRequest, "REQUEST_INVALID")
		return "", errRequestWritten
	}
	username, err := parseUsernameObject(data)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "REQUEST_INVALID")
		return "", errRequestWritten
	}
	return username, nil
}

var errRequestWritten = errors.New("request error response written")

func parseUsernameObject(data []byte) (string, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	tok, err := dec.Token()
	if err != nil || tok != json.Delim('{') {
		return "", errors.New("expected JSON object")
	}
	seen := false
	var username string
	for dec.More() {
		tok, err = dec.Token()
		if err != nil {
			return "", err
		}
		key, ok := tok.(string)
		if !ok || key != "username" || seen {
			return "", errors.New("unexpected or duplicate field")
		}
		seen = true
		if err = dec.Decode(&username); err != nil {
			return "", err
		}
	}
	tok, err = dec.Token()
	if err != nil || tok != json.Delim('}') || !seen {
		return "", errors.New("incomplete JSON object")
	}
	var extra any
	if err = dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return "", errors.New("trailing JSON value")
	}
	return username, nil
}

func (h *Handler) readCredential(r *http.Request) (string, bool, bool) {
	name := h.cookieName()
	var values []string
	for _, cookie := range r.Cookies() {
		if cookie.Name == name {
			values = append(values, cookie.Value)
		}
	}
	if len(values) == 0 {
		return "", false, false
	}
	if len(values) != 1 || values[0] == "" {
		return "", true, false
	}
	return values[0], true, true
}

func (h *Handler) resolveSource(r *http.Request) (string, error) {
	ip, err := h.clientIP.Resolve(r.RemoteAddr, r.Header.Get("X-Forwarded-For"))
	if err != nil || !ip.IsValid() {
		return "", clientip.ErrInvalidSource
	}
	return ip.String(), nil
}

func (h *Handler) allowedOrigin(r *http.Request, requireOrigin bool) bool {
	fetchSite := r.Header.Values("Sec-Fetch-Site")
	if len(fetchSite) > 1 {
		return false
	}
	if len(fetchSite) == 1 {
		site := strings.ToLower(strings.TrimSpace(fetchSite[0]))
		switch site {
		case "same-origin", "same-site", "none":
		case "cross-site", "":
			return false
		default:
			return false
		}
	}
	values, present := r.Header[http.CanonicalHeaderKey("Origin")]
	if !present {
		return !requireOrigin
	}
	if len(values) != 1 || values[0] != h.origin {
		return false
	}
	return true
}

func (h *Handler) writeSession(w http.ResponseWriter, status int, identity session.Identity, absoluteRemaining time.Duration, absoluteExpiresAt time.Time, credential string) {
	h.setCookie(w, credential, absoluteRemaining, absoluteExpiresAt)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(sessionResponse{GuestID: identity.GuestID, Username: identity.Username, Expires: identity.ExpiresAtUTC.UTC()})
}

func (h *Handler) writeSessionError(w http.ResponseWriter, err error) {
	var limited session.RateLimitError
	if errors.As(err, &limited) || errors.Is(err, session.ErrRateLimited) {
		retryAfter := 1
		if limited.RetryAfter > 0 {
			retryAfter = int(limited.RetryAfter / time.Second)
			if limited.RetryAfter%time.Second != 0 {
				retryAfter++
			}
		}
		w.Header().Set("Retry-After", fmt.Sprintf("%d", retryAfter))
		h.writeError(w, http.StatusTooManyRequests, "RATE_LIMITED")
		return
	}
	if errors.Is(err, session.ErrUnauthenticated) {
		h.writeError(w, http.StatusUnauthorized, "SESSION_REQUIRED")
		return
	}
	if errors.Is(err, session.ErrUnavailable) || errors.Is(err, session.ErrBackpressure) || errors.Is(err, session.ErrClosed) || errors.Is(err, clientip.ErrInvalidSource) {
		h.writeError(w, http.StatusServiceUnavailable, "SESSION_UNAVAILABLE")
		return
	}
	if errors.Is(err, session.ErrUsernameInvalid) {
		h.writeError(w, http.StatusBadRequest, "USERNAME_INVALID")
		return
	}
	h.writeError(w, http.StatusServiceUnavailable, "SESSION_UNAVAILABLE")
}

func (h *Handler) writeError(w http.ResponseWriter, status int, code string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	var response errorResponse
	response.Error.Code = code
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) cookieName() string {
	if h.mode == Production {
		return prodCookieName
	}
	return devCookieName
}

func (h *Handler) cookie(value string, absoluteRemaining time.Duration, absoluteExpiresAt time.Time) *http.Cookie {
	maxAge := int(absoluteRemaining / time.Second)
	return &http.Cookie{Name: h.cookieName(), Value: value, Path: "/", HttpOnly: true, Secure: h.mode == Production, SameSite: http.SameSiteLaxMode, MaxAge: maxAge, Expires: absoluteExpiresAt.UTC()}
}

func (h *Handler) setCookie(w http.ResponseWriter, value string, absoluteRemaining time.Duration, absoluteExpiresAt time.Time) {
	http.SetCookie(w, h.cookie(value, absoluteRemaining, absoluteExpiresAt))
}

func (h *Handler) clearCookie(w http.ResponseWriter) {
	c := h.cookie("", -time.Second, time.Unix(1, 0).UTC())
	c.MaxAge = -1
	c.Expires = time.Unix(1, 0).UTC()
	http.SetCookie(w, c)
}

var _ ClientIPResolver = clientip.Resolver{}
