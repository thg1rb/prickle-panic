package session

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/netip"
	"sort"
	"sync"
	"time"

	"github.com/thg1rb/prickle-panic/apps/server/internal/clock"
)

var (
	ErrUnauthenticated = errors.New("session is not authorized")
	ErrRateLimited     = errors.New("session creation rate limited")
	ErrUnavailable     = errors.New("session service unavailable")
	ErrBackpressure    = errors.New("session expiry delivery is backpressured")
	ErrClosed          = errors.New("session manager is closed")
)

type ExpiryReason string

const (
	IdleTimeout     ExpiryReason = "IDLE_TIMEOUT"
	AbsoluteTimeout ExpiryReason = "ABSOLUTE_TIMEOUT"
)

type ExpiryEvent struct {
	ID           string
	GuestID      string
	Reason       ExpiryReason
	ExpiredAtUTC time.Time
	DeadlineTick time.Duration
}

type ExpirySink interface {
	Deliver(context.Context, ExpiryEvent) error
}

type Config struct {
	Clock          clock.Clock
	Sink           ExpirySink
	IdleTTL        time.Duration
	AbsoluteTTL    time.Duration
	CreationLimit  int
	CreationWindow time.Duration
	MaxSourceKeys  int
	OutboxCapacity int
	RetryMin       time.Duration
	RetryMax       time.Duration
	Random         io.Reader
}

type Identity struct {
	GuestID      string
	Username     string
	CreatedAtUTC time.Time
	ExpiresAtUTC time.Time
}
type Credential struct{ value string }

func (c Credential) CookieValue() string         { return c.value }
func (Credential) String() string                { return "[REDACTED]" }
func (Credential) GoString() string              { return "[REDACTED]" }
func (Credential) Format(s fmt.State, verb rune) { _, _ = io.WriteString(s, "[REDACTED]") }

type Created struct {
	Identity   Identity
	Credential Credential
}

type record struct {
	identity         Identity
	hash             [32]byte
	createdTick      time.Duration
	lastActivityTick time.Duration
	deadline         time.Duration
	heapIndex        int
	expired          bool
	pending          bool
	event            ExpiryEvent
}
type deadlineHeap []*record

func (h deadlineHeap) Len() int { return len(h) }
func (h deadlineHeap) Less(i, j int) bool {
	if h[i].deadline == h[j].deadline {
		return h[i].identity.GuestID < h[j].identity.GuestID
	}
	return h[i].deadline < h[j].deadline
}
func (h deadlineHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i]; h[i].heapIndex = i; h[j].heapIndex = j }
func (h *deadlineHeap) Push(x any)   { r := x.(*record); r.heapIndex = len(*h); *h = append(*h, r) }
func (h *deadlineHeap) Pop() any {
	old := *h
	n := len(old)
	r := old[n-1]
	old[n-1] = nil
	r.heapIndex = -1
	*h = old[:n-1]
	return r
}

type allocation struct{ ticks []time.Duration }
type queuedEvent struct {
	rec   *record
	event ExpiryEvent
}

type Manager struct {
	cfg          Config
	mu           sync.Mutex
	byHash       map[[32]byte]*record
	byGuest      map[string]*record
	deadlines    deadlineHeap
	allocations  map[string]allocation
	pending      int
	queue        chan queuedEvent
	wake         chan struct{}
	ctx          context.Context
	cancel       context.CancelFunc
	done         chan struct{}
	dispatchDone chan struct{}
	closed       bool
}

func NewManager(cfg Config) (*Manager, error) {
	if cfg.Clock == nil || cfg.Sink == nil {
		return nil, fmt.Errorf("clock and expiry sink are required")
	}
	if cfg.IdleTTL <= 0 {
		cfg.IdleTTL = 24 * time.Hour
	}
	if cfg.AbsoluteTTL <= 0 {
		cfg.AbsoluteTTL = 7 * 24 * time.Hour
	}
	if cfg.CreationLimit <= 0 {
		cfg.CreationLimit = 10
	}
	if cfg.CreationWindow <= 0 {
		cfg.CreationWindow = 10 * time.Minute
	}
	if cfg.MaxSourceKeys <= 0 {
		cfg.MaxSourceKeys = 65536
	}
	if cfg.OutboxCapacity <= 0 {
		cfg.OutboxCapacity = 1024
	}
	if cfg.RetryMin <= 0 {
		cfg.RetryMin = 100 * time.Millisecond
	}
	if cfg.RetryMax <= 0 {
		cfg.RetryMax = 30 * time.Second
	}
	if cfg.RetryMax < cfg.RetryMin {
		return nil, fmt.Errorf("expiry retry maximum must be greater than or equal to retry minimum")
	}
	if cfg.Random == nil {
		cfg.Random = rand.Reader
	}
	ctx, cancel := context.WithCancel(context.Background())
	m := &Manager{cfg: cfg, byHash: make(map[[32]byte]*record), byGuest: make(map[string]*record), allocations: make(map[string]allocation), queue: make(chan queuedEvent, cfg.OutboxCapacity), wake: make(chan struct{}, 1), ctx: ctx, cancel: cancel, done: make(chan struct{}), dispatchDone: make(chan struct{})}
	go m.schedule()
	go m.dispatch()
	return m, nil
}

// Create atomically validates quota and allocates a new identity. sourceIP must
// be the canonical address resolved by the trusted-proxy boundary.
func (m *Manager) Create(rawUsername, sourceIP string) (Created, error) {
	username, err := NormalizeUsername(rawUsername)
	if err != nil {
		return Created{}, err
	}
	ip, err := netip.ParseAddr(sourceIP)
	if err != nil {
		return Created{}, ErrUnavailable
	}
	sourceIP = ip.Unmap().String()
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return Created{}, ErrClosed
	}
	if m.pending > 0 {
		return Created{}, ErrBackpressure
	}
	now := m.cfg.Clock.Now()
	m.pruneAllocations(now.Tick)
	a, exists := m.allocations[sourceIP]
	if !exists && len(m.allocations) >= m.cfg.MaxSourceKeys {
		return Created{}, ErrUnavailable
	}
	if len(a.ticks) >= m.cfg.CreationLimit {
		wait := m.cfg.CreationWindow - (now.Tick - a.ticks[0])
		if wait < 0 {
			wait = 0
		}
		return Created{}, RateLimitError{RetryAfter: wait}
	}
	var id string
	for i := 0; i < 3; i++ {
		id, err = m.randomID("gst_", 16)
		if err != nil {
			return Created{}, ErrUnavailable
		}
		if m.byGuest[id] == nil {
			break
		}
		id = ""
	}
	if id == "" {
		return Created{}, ErrUnavailable
	}
	var credential string
	var hash [32]byte
	for i := 0; i < 3; i++ {
		credential, err = m.randomID("ppg_", 32)
		if err != nil {
			return Created{}, ErrUnavailable
		}
		hash = sha256.Sum256([]byte(credential))
		if m.byHash[hash] == nil {
			break
		}
		credential = ""
	}
	if credential == "" {
		return Created{}, ErrUnavailable
	}
	absTick := now.Tick + m.cfg.AbsoluteTTL
	idleTick := now.Tick + m.cfg.IdleTTL
	deadline := minTick(absTick, idleTick)
	identity := Identity{GuestID: id, Username: username, CreatedAtUTC: now.UTC, ExpiresAtUTC: now.UTC.Add(deadline - now.Tick)}
	r := &record{identity: identity, hash: hash, createdTick: now.Tick, lastActivityTick: now.Tick, deadline: deadline, heapIndex: -1}
	m.byHash[hash] = r
	m.byGuest[id] = r
	heapPush(&m.deadlines, r)
	a.ticks = append(a.ticks, now.Tick)
	m.allocations[sourceIP] = a
	m.signal()
	return Created{Identity: identity, Credential: Credential{value: credential}}, nil
}

type RateLimitError struct{ RetryAfter time.Duration }

func (e RateLimitError) Error() string        { return ErrRateLimited.Error() }
func (e RateLimitError) Is(target error) bool { return target == ErrRateLimited }

// Lookup authorizes a credential without extending idle expiry. Call
// RecordActivity only after an authenticated application operation succeeds.
func (m *Manager) Lookup(credential string) (Identity, error) { return m.lookup(credential, false) }

// Recover authorizes an explicit session-recovery request, which qualifies as
// authenticated activity under the approved session policy.
func (m *Manager) Recover(credential string) (Identity, error) { return m.lookup(credential, true) }

func (m *Manager) lookup(credential string, qualifyingActivity bool) (Identity, error) {
	if len(credential) > 128 {
		return Identity{}, ErrUnauthenticated
	}
	hash := sha256.Sum256([]byte(credential))
	m.mu.Lock()
	defer m.mu.Unlock()
	r := m.byHash[hash]
	if r == nil {
		return Identity{}, ErrUnauthenticated
	}
	now := m.cfg.Clock.Now()
	if now.Tick >= r.deadline {
		m.expireLocked(r, now)
		return Identity{}, ErrUnauthenticated
	}
	if m.pending > 0 {
		return Identity{}, ErrBackpressure
	}
	if qualifyingActivity {
		m.touchLocked(r, now)
	}
	return r.identity, nil
}

// RecordActivity extends only the idle deadline and must be called after the
// application action succeeds. It cannot revive expired or revoked sessions.
func (m *Manager) RecordActivity(guestID string) (Identity, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	r := m.byGuest[guestID]
	if r == nil || r.expired {
		return Identity{}, ErrUnauthenticated
	}
	now := m.cfg.Clock.Now()
	if now.Tick >= r.deadline {
		m.expireLocked(r, now)
		return Identity{}, ErrUnauthenticated
	}
	if m.pending > 0 {
		return Identity{}, ErrBackpressure
	}
	m.touchLocked(r, now)
	return r.identity, nil
}

func (m *Manager) touchLocked(r *record, now clock.Sample) {
	r.lastActivityTick = now.Tick
	idleDeadline := now.Tick + m.cfg.IdleTTL
	absDeadline := r.createdTick + m.cfg.AbsoluteTTL
	r.deadline = minTick(idleDeadline, absDeadline)
	r.identity.ExpiresAtUTC = now.UTC.Add(r.deadline - now.Tick)
	heapFix(&m.deadlines, r.heapIndex)
	m.signal()
}

func (m *Manager) Revoke(guestID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	r := m.byGuest[guestID]
	if r == nil || r.expired {
		return false
	}
	delete(m.byHash, r.hash)
	delete(m.byGuest, guestID)
	heapRemove(&m.deadlines, r.heapIndex)
	m.signal()
	return true
}

func (m *Manager) Close() error {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return nil
	}
	m.closed = true
	m.cancel()
	m.mu.Unlock()
	<-m.done
	<-m.dispatchDone
	return nil
}

func (m *Manager) schedule() {
	defer close(m.done)
	for {
		m.mu.Lock()
		var wait time.Duration = time.Hour
		if len(m.deadlines) > 0 {
			wait = m.deadlines[0].deadline - m.cfg.Clock.Now().Tick
			if wait < 0 {
				wait = 0
			}
		}
		m.mu.Unlock()
		t := m.cfg.Clock.NewTimer(wait)
		select {
		case <-m.ctx.Done():
			t.Stop()
			return
		case <-m.wake:
			t.Stop()
			continue
		case <-t.C():
			m.expireDue()
		}
	}
}

func (m *Manager) expireDue() {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := m.cfg.Clock.Now()
	for len(m.deadlines) > 0 && now.Tick >= m.deadlines[0].deadline {
		r := heapPop(&m.deadlines)
		if r.expired {
			continue
		}
		m.expireLocked(r, now)
	}
}

func (m *Manager) expireLocked(r *record, now clock.Sample) {
	if r.expired {
		return
	}
	r.expired = true
	heapRemove(&m.deadlines, r.heapIndex)
	delete(m.byHash, r.hash)
	reason := IdleTimeout
	if r.createdTick+m.cfg.AbsoluteTTL <= r.lastActivityTick+m.cfg.IdleTTL {
		reason = AbsoluteTimeout
	}
	deadlineUTC := now.UTC.Add(r.deadline - now.Tick)
	r.event = ExpiryEvent{ID: "expiry_" + r.identity.GuestID, GuestID: r.identity.GuestID, Reason: reason, ExpiredAtUTC: deadlineUTC, DeadlineTick: r.deadline}
	item := queuedEvent{rec: r, event: r.event}
	select {
	case m.queue <- item:
	default:
		r.pending = true
		m.pending++
	}
	m.signal()
}

func (m *Manager) dispatch() {
	defer close(m.dispatchDone)
	for {
		select {
		case <-m.ctx.Done():
			return
		case item := <-m.queue:
			backoff := m.cfg.RetryMin
			for {
				if err := m.cfg.Sink.Deliver(m.ctx, item.event); err == nil {
					m.delivered(item.rec)
					break
				}
				t := m.cfg.Clock.NewTimer(backoff)
				select {
				case <-m.ctx.Done():
					t.Stop()
					return
				case <-t.C():
				}
				backoff *= 2
				if backoff > m.cfg.RetryMax {
					backoff = m.cfg.RetryMax
				}
			}
		}
	}
}

func (m *Manager) delivered(r *record) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if r.pending {
		r.pending = false
		m.pending--
	}
	delete(m.byGuest, r.identity.GuestID)
	// Promote pending expiry records in deterministic deadline/guest order.
	items := make([]*record, 0, m.pending)
	for _, x := range m.byGuest {
		if x.pending {
			items = append(items, x)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].deadline == items[j].deadline {
			return items[i].identity.GuestID < items[j].identity.GuestID
		}
		return items[i].deadline < items[j].deadline
	})
	for _, x := range items {
		select {
		case m.queue <- queuedEvent{rec: x, event: x.event}:
			x.pending = false
			m.pending--
		default:
			m.signal()
			return
		}
	}
	if m.pending == 0 {
		m.signal()
	}
}

func (m *Manager) pruneAllocations(now time.Duration) {
	for ip, a := range m.allocations {
		n := 0
		for _, tick := range a.ticks {
			if now-tick < m.cfg.CreationWindow {
				a.ticks[n] = tick
				n++
			}
		}
		a.ticks = a.ticks[:n]
		if n == 0 {
			delete(m.allocations, ip)
		} else {
			m.allocations[ip] = a
		}
	}
}
func (m *Manager) randomID(prefix string, n int) (string, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(m.cfg.Random, b); err != nil {
		return "", err
	}
	return prefix + base64.RawURLEncoding.EncodeToString(b), nil
}
func (m *Manager) signal() {
	select {
	case m.wake <- struct{}{}:
	default:
	}
}
func minTick(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

// local heap helpers keep the record index synchronized without exposing heap internals.
func heapPush(h *deadlineHeap, r *record) {
	*h = append(*h, r)
	r.heapIndex = len(*h) - 1
	for i := r.heapIndex; i > 0; {
		p := (i - 1) / 2
		if !(*h).Less(i, p) {
			break
		}
		(*h).Swap(i, p)
		i = p
	}
}
func heapPop(h *deadlineHeap) *record {
	if len(*h) == 0 {
		return nil
	}
	n := len(*h) - 1
	(*h).Swap(0, n)
	r := (*h)[n]
	(*h)[n] = nil
	*h = (*h)[:n]
	r.heapIndex = -1
	if len(*h) > 0 {
		heapDown(*h, 0)
	}
	return r
}
func heapRemove(h *deadlineHeap, i int) {
	if i < 0 || i >= len(*h) {
		return
	}
	n := len(*h) - 1
	if i != n {
		(*h).Swap(i, n)
	}
	r := (*h)[n]
	(*h)[n] = nil
	*h = (*h)[:n]
	r.heapIndex = -1
	if i < len(*h) {
		heapFix(h, i)
	}
}
func heapFix(h *deadlineHeap, i int) {
	if i > 0 && (*h).Less(i, (i-1)/2) {
		for i > 0 {
			p := (i - 1) / 2
			if !(*h).Less(i, p) {
				break
			}
			(*h).Swap(i, p)
			i = p
		}
	} else {
		heapDown(*h, i)
	}
}
func heapDown(h deadlineHeap, i int) {
	for {
		l := 2*i + 1
		if l >= len(h) {
			return
		}
		j := l
		r := l + 1
		if r < len(h) && h.Less(r, l) {
			j = r
		}
		if !h.Less(j, i) {
			return
		}
		h.Swap(i, j)
		i = j
	}
}
