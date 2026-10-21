package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/thg1rb/prickle-panic/apps/server/internal/clientip"
	"github.com/thg1rb/prickle-panic/apps/server/internal/clock"
	"github.com/thg1rb/prickle-panic/apps/server/internal/session"
)

type sink struct{}

func (sink) Deliver(context.Context, session.ExpiryEvent) error { return nil }

type testSetup struct {
	handler *Handler
	manager *session.Manager
}

func setup(t *testing.T, mode Mode, limit int, resolver ClientIPResolver) testSetup {
	t.Helper()
	if resolver == nil {
		resolver, _ = clientip.NewResolver(nil, 0)
	}
	manager, err := session.NewManager(session.Config{
		Clock:          clock.New(),
		Sink:           sink{},
		IdleTTL:        time.Hour,
		AbsoluteTTL:    24 * time.Hour,
		CreationLimit:  limit,
		CreationWindow: 10 * time.Minute,
		OutboxCapacity: 4,
		RetryMin:       time.Millisecond,
		RetryMax:       10 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = manager.Close() })
	origin := "http://localhost:5173"
	if mode == Production {
		origin = "https://game.example.test"
	}
	handler, err := New(Config{Mode: mode, Origin: origin, Sessions: manager, ClientIP: resolver})
	if err != nil {
		t.Fatal(err)
	}
	return testSetup{handler: handler, manager: manager}
}

func postRequest(body string, cookie *http.Cookie) *http.Request {
	r := httptest.NewRequest(http.MethodPost, guestSessionPath, strings.NewReader(body))
	r.Header.Set("Origin", "http://localhost:5173")
	r.Header.Set("Content-Type", "application/json")
	if cookie != nil {
		r.AddCookie(cookie)
	}
	return r
}

func assertError(t *testing.T, rec *httptest.ResponseRecorder, status int, code string) {
	t.Helper()
	if rec.Code != status {
		t.Fatalf("status = %d, want %d: %s", rec.Code, status, rec.Body.String())
	}
	if got := rec.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", got)
	}
	var body errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Error.Code != code {
		t.Fatalf("error code = %q, want %q", body.Error.Code, code)
	}
}

func TestCreateAndProductionCookiePolicy(t *testing.T) {
	s := setup(t, Production, 10, nil)
	r := httptest.NewRequest(http.MethodPost, guestSessionPath, strings.NewReader(`{"username":"โคโค่"}`))
	r.Header.Set("Origin", "https://game.example.test")
	r.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, r)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d: %s", rec.Code, rec.Body.String())
	}
	if rec.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("Cache-Control = %q", rec.Header().Get("Cache-Control"))
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("session API unexpectedly enabled CORS")
	}
	var body sessionResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.GuestID == "" || body.Username != "โคโค่" || body.Expires.IsZero() {
		t.Fatalf("unexpected response: %+v", body)
	}
	if strings.Contains(rec.Body.String(), "ppg_") || strings.Contains(rec.Body.String(), "credential") || strings.Contains(rec.Body.String(), "hash") {
		t.Fatalf("response exposed sensitive material: %s", rec.Body.String())
	}
	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("got %d cookies, want one", len(cookies))
	}
	c := cookies[0]
	if c.Name != prodCookieName || c.Value == "" || !c.HttpOnly || !c.Secure || c.SameSite != http.SameSiteLaxMode || c.Path != "/" || c.Domain != "" {
		t.Fatalf("unsafe production cookie: %+v", c)
	}
	if c.MaxAge <= 0 || c.MaxAge > int((24*time.Hour)/time.Second) || c.Expires.After(time.Now().Add(24*time.Hour)) {
		t.Fatalf("cookie exceeds configured absolute lifetime: %+v", c)
	}
}

func TestPostReusesExistingSessionWithoutRenaming(t *testing.T) {
	s := setup(t, Development, 1, nil)
	first := httptest.NewRecorder()
	s.handler.ServeHTTP(first, postRequest(`{"username":"Coco"}`, nil))
	if first.Code != http.StatusCreated {
		t.Fatalf("first status = %d: %s", first.Code, first.Body.String())
	}
	cookie := first.Result().Cookies()[0]
	second := httptest.NewRecorder()
	s.handler.ServeHTTP(second, postRequest(`{"username":"Other"}`, cookie))
	if second.Code != http.StatusOK {
		t.Fatalf("reuse status = %d: %s", second.Code, second.Body.String())
	}
	var initial, reused sessionResponse
	_ = json.Unmarshal(first.Body.Bytes(), &initial)
	_ = json.Unmarshal(second.Body.Bytes(), &reused)
	if reused.GuestID != initial.GuestID || reused.Username != initial.Username {
		t.Fatalf("POST changed session identity/name: initial=%+v reused=%+v", initial, reused)
	}
	if got := second.Result().Cookies(); len(got) != 1 || got[0].Value != cookie.Value {
		t.Fatal("valid session reuse rotated or omitted the credential cookie")
	}
	if !reused.Expires.Equal(initial.Expires) {
		t.Fatal("POST reuse changed the session expiry")
	}
}

func TestRecoveryQualifiesActivityAndExpiredCredentialIsCleared(t *testing.T) {
	s := setup(t, Development, 10, nil)
	first := httptest.NewRecorder()
	s.handler.ServeHTTP(first, postRequest(`{"username":"Coco"}`, nil))
	cookie := first.Result().Cookies()[0]
	var before sessionResponse
	_ = json.Unmarshal(first.Body.Bytes(), &before)
	r := httptest.NewRequest(http.MethodGet, guestSessionPath, nil)
	r.AddCookie(cookie)
	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, r)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET recovery status = %d: %s", rec.Code, rec.Body.String())
	}
	if got := rec.Result().Cookies(); len(got) != 1 || got[0].Value != cookie.Value {
		t.Fatal("successful recovery did not preserve the credential cookie")
	}
	var after sessionResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &after)
	if !after.Expires.After(before.Expires) {
		t.Fatal("successful GET recovery did not extend idle expiry")
	}
	_ = s.manager.Revoke(before.GuestID)
	stale := httptest.NewRequest(http.MethodGet, guestSessionPath, nil)
	stale.AddCookie(cookie)
	cleared := httptest.NewRecorder()
	s.handler.ServeHTTP(cleared, stale)
	assertError(t, cleared, http.StatusUnauthorized, "SESSION_REQUIRED")
	if got := cleared.Result().Cookies(); len(got) != 1 || got[0].MaxAge != -1 || got[0].Name != devCookieName {
		t.Fatalf("stale cookie was not safely cleared: %+v", got)
	}
}

func TestGetWithoutCredentialRequiresSession(t *testing.T) {
	s := setup(t, Development, 10, nil)
	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, guestSessionPath, nil))
	assertError(t, rec, http.StatusUnauthorized, "SESSION_REQUIRED")
	if len(rec.Result().Cookies()) != 0 {
		t.Fatal("missing credential caused a cookie to be set or cleared")
	}
}

func TestStrictRequestValidation(t *testing.T) {
	cases := []struct {
		name string
		body []byte
		want int
		code string
	}{
		{"malformed", []byte(`{"username":`), 400, "REQUEST_INVALID"},
		{"unknown field", []byte(`{"username":"Coco","extra":1}`), 400, "REQUEST_INVALID"},
		{"duplicate field", []byte(`{"username":"Coco","username":"Other"}`), 400, "REQUEST_INVALID"},
		{"trailing JSON", []byte(`{"username":"Coco"} {}`), 400, "REQUEST_INVALID"},
		{"missing field", []byte(`{}`), 400, "REQUEST_INVALID"},
		{"wrong type", []byte(`{"username":4}`), 400, "REQUEST_INVALID"},
		{"invalid UTF-8", []byte{'{', '"', 'u', 's', 'e', 'r', 'n', 'a', 'm', 'e', '"', ':', '"', 0xff, '"', '}'}, 400, "REQUEST_INVALID"},
		{"invalid username", []byte(`{"username":"<script>"}`), 400, "USERNAME_INVALID"},
		{"oversized", []byte(`{"username":"` + strings.Repeat("a", maxRequestBytes) + `"}`), 413, "REQUEST_TOO_LARGE"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := setup(t, Development, 10, nil)
			r := httptest.NewRequest(http.MethodPost, guestSessionPath, strings.NewReader(string(tc.body)))
			r.Header.Set("Origin", "http://localhost:5173")
			rec := httptest.NewRecorder()
			s.handler.ServeHTTP(rec, r)
			assertError(t, rec, tc.want, tc.code)
			if len(rec.Result().Cookies()) != 0 {
				t.Fatal("invalid request issued a credential cookie")
			}
		})
	}
}

func TestOriginAndFetchMetadataProtection(t *testing.T) {
	tests := []struct {
		name   string
		method string
		origin string
		fetch  string
		want   int
	}{
		{"post missing origin", http.MethodPost, "", "same-origin", 403},
		{"post mismatched origin", http.MethodPost, "https://evil.example", "same-origin", 403},
		{"post cross-site", http.MethodPost, "http://localhost:5173", "cross-site", 403},
		{"post malformed fetch metadata", http.MethodPost, "http://localhost:5173", "cross-site, same-origin", 403},
		{"get mismatched origin", http.MethodGet, "https://evil.example", "same-origin", 403},
		{"get cross-site", http.MethodGet, "", "cross-site", 403},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := setup(t, Development, 10, nil)
			var body io.Reader
			if tc.method == http.MethodPost {
				body = strings.NewReader(`{"username":"Coco"}`)
			}
			r := httptest.NewRequest(tc.method, guestSessionPath, body)
			if tc.origin != "" {
				r.Header.Set("Origin", tc.origin)
			}
			if tc.fetch != "" {
				r.Header.Set("Sec-Fetch-Site", tc.fetch)
			}
			rec := httptest.NewRecorder()
			s.handler.ServeHTTP(rec, r)
			assertError(t, rec, tc.want, "ORIGIN_NOT_ALLOWED")
		})
	}
}

func TestCreationQuotaInvalidInputReuseAndRetryAfter(t *testing.T) {
	s := setup(t, Development, 1, nil)
	invalid := httptest.NewRecorder()
	s.handler.ServeHTTP(invalid, postRequest(`{"username":"!"}`, nil))
	assertError(t, invalid, http.StatusBadRequest, "USERNAME_INVALID")
	first := httptest.NewRecorder()
	s.handler.ServeHTTP(first, postRequest(`{"username":"Coco"}`, nil))
	if first.Code != http.StatusCreated {
		t.Fatalf("valid request after invalid input got %d", first.Code)
	}
	cookie := first.Result().Cookies()[0]
	reuse := httptest.NewRecorder()
	s.handler.ServeHTTP(reuse, postRequest(`{"username":"Other"}`, cookie))
	if reuse.Code != http.StatusOK {
		t.Fatalf("existing session reuse got %d: %s", reuse.Code, reuse.Body.String())
	}
	limited := httptest.NewRecorder()
	s.handler.ServeHTTP(limited, postRequest(`{"username":"Other"}`, nil))
	assertError(t, limited, http.StatusTooManyRequests, "RATE_LIMITED")
	if got := limited.Header().Get("Retry-After"); got != "600" {
		t.Fatalf("Retry-After = %q, want 600 seconds", got)
	}
	malformedAtLimit := httptest.NewRecorder()
	s.handler.ServeHTTP(malformedAtLimit, postRequest(`not-json`, nil))
	assertError(t, malformedAtLimit, http.StatusTooManyRequests, "RATE_LIMITED")
}

func TestConcurrentCreationCannotExceedLimit(t *testing.T) {
	s := setup(t, Development, 10, nil)
	var wg sync.WaitGroup
	statuses := make(chan int, 16)
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rec := httptest.NewRecorder()
			s.handler.ServeHTTP(rec, postRequest(`{"username":"Coco"}`, nil))
			statuses <- rec.Code
		}()
	}
	wg.Wait()
	close(statuses)
	created, limited := 0, 0
	for status := range statuses {
		switch status {
		case http.StatusCreated:
			created++
		case http.StatusTooManyRequests:
			limited++
		default:
			t.Fatalf("unexpected concurrent status %d", status)
		}
	}
	if created != 10 || limited != 6 {
		t.Fatalf("created=%d limited=%d, want 10 and 6", created, limited)
	}
}

func TestTrustedProxyFailureAndUntrustedForwardedHeader(t *testing.T) {
	trusted, err := clientip.NewResolver([]netip.Prefix{netip.MustParsePrefix("10.0.0.0/24")}, 32)
	if err != nil {
		t.Fatal(err)
	}
	s := setup(t, Development, 10, trusted)
	r := postRequest(`{"username":"Coco"}`, nil)
	r.RemoteAddr = "10.0.0.2:1000"
	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, r)
	assertError(t, rec, http.StatusServiceUnavailable, "SESSION_UNAVAILABLE")

	direct := setup(t, Development, 1, nil)
	first := postRequest(`{"username":"Coco"}`, nil)
	first.RemoteAddr = "192.0.2.1:1000"
	first.Header.Set("X-Forwarded-For", "198.51.100.1")
	one := httptest.NewRecorder()
	direct.handler.ServeHTTP(one, first)
	if one.Code != http.StatusCreated {
		t.Fatalf("first direct request got %d", one.Code)
	}
	second := postRequest(`{"username":"Coco"}`, nil)
	second.RemoteAddr = "192.0.2.2:1000"
	second.Header.Set("X-Forwarded-For", "192.0.2.1")
	two := httptest.NewRecorder()
	direct.handler.ServeHTTP(two, second)
	if two.Code != http.StatusCreated {
		t.Fatalf("untrusted X-Forwarded-For affected source quota: status %d", two.Code)
	}
}

func TestTrustedProxyRejectsMultipleForwardedFieldLines(t *testing.T) {
	trusted, err := clientip.NewResolver([]netip.Prefix{netip.MustParsePrefix("10.0.0.0/24")}, 32)
	if err != nil {
		t.Fatal(err)
	}
	s := setup(t, Development, 10, trusted)
	r := postRequest(`{"username":"Coco"}`, nil)
	r.RemoteAddr = "10.0.0.2:1000"
	r.Header.Add("X-Forwarded-For", "198.51.100.9")
	r.Header.Add("X-Forwarded-For", "203.0.113.8")
	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, r)
	assertError(t, rec, http.StatusServiceUnavailable, "SESSION_UNAVAILABLE")
	if len(rec.Result().Cookies()) != 0 {
		t.Fatal("ambiguous forwarded metadata issued a session cookie")
	}
}

func TestDevelopmentCookieIsExplicitAndProductionRequiresHTTPS(t *testing.T) {
	s := setup(t, Development, 10, nil)
	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, postRequest(`{"username":"Coco"}`, nil))
	c := rec.Result().Cookies()[0]
	if c.Name != devCookieName || c.Secure || !c.HttpOnly || c.SameSite != http.SameSiteLaxMode || c.Path != "/" || c.Domain != "" {
		t.Fatalf("unexpected development cookie policy: %+v", c)
	}
	direct, _ := clientip.NewResolver(nil, 0)
	if _, err := New(Config{Mode: Production, Origin: "http://game.example.test", Sessions: s.manager, ClientIP: direct}); err == nil {
		t.Fatal("production accepted an HTTP origin")
	}
	if _, err := New(Config{Mode: Development, Origin: "http://game.example.test", Sessions: s.manager, ClientIP: direct}); err == nil {
		t.Fatal("development accepted a non-loopback origin")
	}
}

func TestUnknownPathAndMethodAreNoStore(t *testing.T) {
	s := setup(t, Development, 10, nil)
	for _, request := range []*http.Request{
		httptest.NewRequest(http.MethodGet, "/api/other", nil),
		httptest.NewRequest(http.MethodDelete, guestSessionPath, nil),
	} {
		rec := httptest.NewRecorder()
		s.handler.ServeHTTP(rec, request)
		if got := rec.Header().Get("Cache-Control"); got != "no-store" {
			t.Fatalf("status %d response Cache-Control = %q", rec.Code, got)
		}
	}
}

func TestHTTPErrorMappingDoesNotExposeInternalText(t *testing.T) {
	s := setup(t, Development, 10, nil)
	if errors.Is(session.ErrUnavailable, session.ErrClosed) {
		t.Fatal("unrelated manager errors compare equal")
	}
	// Existing source exhaustion is covered by the quota test; this verifies
	// that unexpected manager errors are mapped without diagnostic leakage.
	rec := httptest.NewRecorder()
	rec.Header().Set("Cache-Control", "no-store")
	s.handler.writeSessionError(rec, errors.New("secret verifier diagnostic"))
	assertError(t, rec, http.StatusServiceUnavailable, "SESSION_UNAVAILABLE")
	if strings.Contains(rec.Body.String(), "secret") {
		t.Fatal("internal diagnostic leaked in HTTP error")
	}
}
