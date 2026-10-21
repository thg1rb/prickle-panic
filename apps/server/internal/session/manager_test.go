package session

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"testing/iotest"
	"time"

	"github.com/thg1rb/prickle-panic/apps/server/internal/clock"
)

type testClock struct {
	mu     sync.Mutex
	tick   time.Duration
	utc    time.Time
	timers []*testTimer
}
type testTimer struct {
	owner  *testClock
	ch     chan time.Time
	due    time.Duration
	active bool
}

func newTestClock() *testClock { return &testClock{utc: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)} }
func (c *testClock) Now() clock.Sample {
	c.mu.Lock()
	defer c.mu.Unlock()
	return clock.Sample{UTC: c.utc, Tick: c.tick}
}
func (c *testClock) NewTimer(d time.Duration) clock.Timer {
	c.mu.Lock()
	defer c.mu.Unlock()
	t := &testTimer{owner: c, ch: make(chan time.Time, 1), due: c.tick + d, active: true}
	if d <= 0 {
		t.active = false
		t.ch <- c.utc
	}
	c.timers = append(c.timers, t)
	return t
}
func (t *testTimer) C() <-chan time.Time { return t.ch }
func (t *testTimer) Stop() bool {
	t.owner.mu.Lock()
	defer t.owner.mu.Unlock()
	was := t.active
	t.active = false
	return was
}
func (t *testTimer) Reset(d time.Duration) bool {
	t.owner.mu.Lock()
	defer t.owner.mu.Unlock()
	was := t.active
	t.active = true
	t.due = t.owner.tick + d
	return was
}
func (c *testClock) Advance(d time.Duration) {
	c.mu.Lock()
	c.tick += d
	c.utc = c.utc.Add(d)
	for _, t := range c.timers {
		if t.active && c.tick >= t.due {
			t.active = false
			select {
			case t.ch <- c.utc:
			default:
			}
		}
	}
	c.mu.Unlock()
}
func (c *testClock) ShiftWall(d time.Duration) { c.mu.Lock(); c.utc = c.utc.Add(d); c.mu.Unlock() }

type sinkFunc func(context.Context, ExpiryEvent) error

func (f sinkFunc) Deliver(ctx context.Context, e ExpiryEvent) error { return f(ctx, e) }
func newManager(t *testing.T, c clock.Clock, sink ExpirySink, opts func(*Config)) *Manager {
	t.Helper()
	cfg := Config{Clock: c, Sink: sink, IdleTTL: 24 * time.Hour, AbsoluteTTL: 7 * 24 * time.Hour, CreationLimit: 10, CreationWindow: 10 * time.Minute, OutboxCapacity: 2, RetryMin: time.Millisecond, RetryMax: 10 * time.Millisecond}
	if opts != nil {
		opts(&cfg)
	}
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = m.Close() })
	return m
}

func TestManagerRejectsInvalidRetryBounds(t *testing.T) {
	_, err := NewManager(Config{Clock: newTestClock(), Sink: sinkFunc(func(context.Context, ExpiryEvent) error { return nil }), RetryMin: time.Second, RetryMax: time.Millisecond})
	if err == nil {
		t.Fatal("invalid retry bounds accepted")
	}
}

func TestCreateLookupRevokeAndRedaction(t *testing.T) {
	c := newTestClock()
	m := newManager(t, c, sinkFunc(func(context.Context, ExpiryEvent) error { return nil }), nil)
	created, err := m.Create("Coco", "192.0.2.1")
	if err != nil {
		t.Fatal(err)
	}
	token := created.Credential.CookieValue()
	if len(token) < 40 {
		t.Fatalf("credential too short: %d", len(token))
	}
	raw, err := base64.RawURLEncoding.DecodeString(token[len("ppg_"):])
	if err != nil || len(raw) != 32 {
		t.Fatalf("credential is not 256 bits: %v length=%d", err, len(raw))
	}
	if fmt.Sprint(created.Credential) != "[REDACTED]" || fmt.Sprintf("%+v", created.Credential) != "[REDACTED]" {
		t.Fatal("credential formatting exposed secret")
	}
	encoded, err := json.Marshal(created)
	if err != nil || strings.Contains(string(encoded), token) {
		t.Fatalf("JSON serialization exposed credential: %s err=%v", encoded, err)
	}
	if _, ok := m.byHash[sha256.Sum256([]byte(token))]; !ok {
		t.Fatal("session store is not indexed by credential hash")
	}
	got, err := m.Lookup(created.Credential.CookieValue())
	if err != nil || got.GuestID != created.Identity.GuestID {
		t.Fatalf("lookup=%+v err=%v", got, err)
	}
	if _, err = m.Lookup("bad"); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("bad credential: %v", err)
	}
	if !m.Revoke(got.GuestID) || m.Revoke(got.GuestID) {
		t.Fatal("revoke must be idempotent")
	}
	if _, err = m.Lookup(created.Credential.CookieValue()); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("revoked lookup: %v", err)
	}
}

func TestCredentialUniquenessCollisionRetryAndEntropyFailure(t *testing.T) {
	c := newTestClock()
	sequence := append(bytes.Repeat([]byte{1}, 16), bytes.Repeat([]byte{2}, 32)...)
	sequence = append(sequence, bytes.Repeat([]byte{3}, 16)...)
	sequence = append(sequence, bytes.Repeat([]byte{2}, 32)...)
	sequence = append(sequence, bytes.Repeat([]byte{4}, 32)...)
	m := newManager(t, c, sinkFunc(func(context.Context, ExpiryEvent) error { return nil }), func(cfg *Config) { cfg.Random = bytes.NewReader(sequence) })
	a, err := m.Create("AA", "192.0.2.1")
	if err != nil {
		t.Fatal(err)
	}
	b, err := m.Create("BB", "192.0.2.1")
	if err != nil {
		t.Fatal(err)
	}
	if a.Identity.GuestID == b.Identity.GuestID || a.Credential.CookieValue() == b.Credential.CookieValue() {
		t.Fatal("allocated identities or credentials were reused")
	}
	if _, err := m.Lookup(a.Credential.CookieValue()); err != nil {
		t.Fatalf("first credential lookup: %v", err)
	}
	if _, err := m.Lookup(b.Credential.CookieValue()); err != nil {
		t.Fatalf("retried credential lookup: %v", err)
	}
	failing := newManager(t, c, sinkFunc(func(context.Context, ExpiryEvent) error { return nil }), func(cfg *Config) {
		cfg.Random = io.MultiReader(bytes.NewReader(bytes.Repeat([]byte{5}, 16)), iotest.ErrReader(errors.New("credential entropy unavailable")))
	})
	if _, err := failing.Create("CC", "192.0.2.2"); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("entropy failure must fail closed: %v", err)
	}
	if len(failing.allocations) != 0 {
		t.Fatal("failed credential generation consumed an allocation quota")
	}
}

func TestIdleAndAbsoluteDeadlineUseMonotonicTicks(t *testing.T) {
	c := newTestClock()
	m := newManager(t, c, sinkFunc(func(context.Context, ExpiryEvent) error { return nil }), func(cfg *Config) { cfg.IdleTTL = 10 * time.Second; cfg.AbsoluteTTL = 25 * time.Second })
	created, err := m.Create("Coco", "192.0.2.1")
	if err != nil {
		t.Fatal(err)
	}
	c.Advance(8 * time.Second)
	if _, err = m.Recover(created.Credential.CookieValue()); err != nil {
		t.Fatal(err)
	}
	c.ShiftWall(48 * time.Hour)
	c.Advance(8 * time.Second)
	if _, err = m.RecordActivity(created.Identity.GuestID); err != nil {
		t.Fatalf("successful activity did not refresh idle expiry: %v", err)
	}
	c.Advance(8 * time.Second)
	if _, err = m.RecordActivity(created.Identity.GuestID); err != nil {
		t.Fatalf("activity should still be valid immediately before absolute deadline: %v", err)
	}
	c.Advance(time.Second)
	if _, err = m.RecordActivity(created.Identity.GuestID); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("absolute expiry must not extend: %v", err)
	}
}

func TestPassiveLookupDoesNotExtendIdleTimeout(t *testing.T) {
	c := newTestClock()
	m := newManager(t, c, sinkFunc(func(context.Context, ExpiryEvent) error { return nil }), func(cfg *Config) { cfg.IdleTTL = 5 * time.Second; cfg.AbsoluteTTL = time.Minute })
	created, err := m.Create("AB", "192.0.2.1")
	if err != nil {
		t.Fatal(err)
	}
	c.Advance(4 * time.Second)
	if _, err := m.Lookup(created.Credential.CookieValue()); err != nil {
		t.Fatal(err)
	}
	c.Advance(time.Second)
	if _, err := m.Lookup(created.Credential.CookieValue()); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("passive lookup extended idle timeout: %v", err)
	}
}

func TestCreationQuotaOnlyCountsSuccessfulAllocationsAndIsAtomic(t *testing.T) {
	c := newTestClock()
	m := newManager(t, c, sinkFunc(func(context.Context, ExpiryEvent) error { return nil }), func(cfg *Config) { cfg.CreationLimit = 10 })
	if _, err := m.Create("x", "192.0.2.1"); !errors.Is(err, ErrUsernameInvalid) {
		t.Fatalf("invalid username: %v", err)
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	success := 0
	limited := 0
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := m.Create("AB", "192.0.2.1")
			mu.Lock()
			defer mu.Unlock()
			if err == nil {
				success++
			} else if errors.Is(err, ErrRateLimited) {
				var limitedErr RateLimitError
				if !errors.As(err, &limitedErr) || limitedErr.RetryAfter != 10*time.Minute {
					t.Errorf("unexpected retry delay: %#v", err)
				}
				limited++
			} else {
				t.Errorf("unexpected create error: %v", err)
			}
		}()
	}
	wg.Wait()
	if success != 10 || limited != 10 {
		t.Fatalf("success=%d limited=%d", success, limited)
	}
	c.Advance(10 * time.Minute)
	if _, err := m.Create("AB", "192.0.2.1"); err != nil {
		t.Fatalf("window did not expire: %v", err)
	}
}

func TestSourceKeyCapacityFailsClosedAndCleansExpiredKeys(t *testing.T) {
	c := newTestClock()
	m := newManager(t, c, sinkFunc(func(context.Context, ExpiryEvent) error { return nil }), func(cfg *Config) { cfg.MaxSourceKeys = 1 })
	if _, err := m.Create("AB", "192.0.2.1"); err != nil {
		t.Fatal(err)
	}
	if _, err := m.Create("CD", "192.0.2.2"); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("source key capacity should fail closed: %v", err)
	}
	c.Advance(10 * time.Minute)
	if _, err := m.Create("CD", "192.0.2.2"); err != nil {
		t.Fatalf("stale source key was not cleaned: %v", err)
	}
}

func TestExpirationEventAndProactiveScheduler(t *testing.T) {
	c := newTestClock()
	events := make(chan ExpiryEvent, 2)
	m := newManager(t, c, sinkFunc(func(_ context.Context, e ExpiryEvent) error { events <- e; return nil }), func(cfg *Config) { cfg.IdleTTL = 5 * time.Second; cfg.AbsoluteTTL = time.Minute })
	created, err := m.Create("AB", "192.0.2.1")
	if err != nil {
		t.Fatal(err)
	}
	c.Advance(5 * time.Second)
	select {
	case e := <-events:
		if e.GuestID != created.Identity.GuestID || e.Reason != IdleTimeout || e.DeadlineTick != 5*time.Second {
			t.Fatalf("bad expiry event: %+v", e)
		}
	case <-time.After(time.Second):
		t.Fatal("proactive expiry was not delivered")
	}
	if _, err = m.Lookup(created.Credential.CookieValue()); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("expired credential authorized: %v", err)
	}
}

func TestAbsoluteExpirationEventReason(t *testing.T) {
	c := newTestClock()
	events := make(chan ExpiryEvent, 1)
	m := newManager(t, c, sinkFunc(func(_ context.Context, e ExpiryEvent) error { events <- e; return nil }), func(cfg *Config) { cfg.IdleTTL = time.Minute; cfg.AbsoluteTTL = 5 * time.Second })
	if _, err := m.Create("AB", "192.0.2.1"); err != nil {
		t.Fatal(err)
	}
	c.Advance(5 * time.Second)
	select {
	case event := <-events:
		if event.Reason != AbsoluteTimeout {
			t.Fatalf("reason=%s want=%s", event.Reason, AbsoluteTimeout)
		}
	case <-time.After(time.Second):
		t.Fatal("absolute expiry event not delivered")
	}
}

func TestExpiryDeliveryRetriesSameStableEvent(t *testing.T) {
	c := clock.New()
	events := make(chan ExpiryEvent, 2)
	var mu sync.Mutex
	calls := 0
	sink := sinkFunc(func(_ context.Context, e ExpiryEvent) error {
		mu.Lock()
		defer mu.Unlock()
		calls++
		events <- e
		if calls == 1 {
			return errors.New("temporary failure")
		}
		return nil
	})
	m := newManager(t, c, sink, func(cfg *Config) {
		cfg.IdleTTL = 20 * time.Millisecond
		cfg.AbsoluteTTL = time.Second
		cfg.RetryMin = time.Millisecond
		cfg.RetryMax = 2 * time.Millisecond
	})
	if _, err := m.Create("AB", "192.0.2.1"); err != nil {
		t.Fatal(err)
	}
	var first, second ExpiryEvent
	select {
	case first = <-events:
	case <-time.After(time.Second):
		t.Fatal("first delivery missing")
	}
	select {
	case second = <-events:
	case <-time.After(time.Second):
		t.Fatal("retry delivery missing")
	}
	if first.ID == "" || first.ID != second.ID {
		t.Fatalf("retry event ID changed: %q vs %q", first.ID, second.ID)
	}
}

type blockedSink struct {
	entered chan struct{}
	release chan struct{}
}

func (s blockedSink) Deliver(ctx context.Context, _ ExpiryEvent) error {
	select {
	case s.entered <- struct{}{}:
	default:
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.release:
		return nil
	}
}

func TestOutboxBackpressureNeverReauthorizesExpiredSession(t *testing.T) {
	c := newTestClock()
	sink := blockedSink{entered: make(chan struct{}, 1), release: make(chan struct{})}
	m := newManager(t, c, sink, func(cfg *Config) {
		cfg.IdleTTL = 2 * time.Second
		cfg.AbsoluteTTL = time.Minute
		cfg.OutboxCapacity = 1
	})
	a, err := m.Create("AA", "192.0.2.1")
	if err != nil {
		t.Fatal(err)
	}
	b, err := m.Create("BB", "192.0.2.1")
	if err != nil {
		t.Fatal(err)
	}
	third, err := m.Create("CC", "192.0.2.1")
	if err != nil {
		t.Fatal(err)
	}
	c.Advance(2 * time.Second)
	select {
	case <-sink.entered:
	case <-time.After(time.Second):
		t.Fatal("sink not invoked")
	}
	// Let a second expired event occupy the one buffered slot, then force one
	// additional event into per-record pending state.
	if _, err = m.Lookup(third.Credential.CookieValue()); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("expired credential authorized: %v", err)
	}
	if _, err = m.Create("DD", "192.0.2.1"); !errors.Is(err, ErrBackpressure) {
		t.Fatalf("creation should be backpressured: %v", err)
	}
	if _, err = m.Lookup(a.Credential.CookieValue()); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("expired session returned backpressure instead of unauthorized: %v", err)
	}
	close(sink.release)
	_ = b
}
