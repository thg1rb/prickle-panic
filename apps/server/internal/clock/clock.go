// Package clock provides paired wall-clock and monotonic time for server domains.
package clock

import "time"

// Sample keeps diagnostic UTC time paired with an enforcement tick. Ticks from
// a Clock are comparable only with other ticks from that same Clock instance.
type Sample struct {
	UTC  time.Time
	Tick time.Duration
}

// Timer is the subset of time.Timer used by the session expiry scheduler.
type Timer interface {
	C() <-chan time.Time
	Stop() bool
	Reset(time.Duration) bool
}

// Clock provides one process-local monotonic time domain.
type Clock interface {
	Now() Sample
	NewTimer(time.Duration) Timer
}

type systemClock struct{ origin time.Time }

func New() Clock { return &systemClock{origin: time.Now()} }
func (c *systemClock) Now() Sample {
	now := time.Now()
	return Sample{UTC: now.UTC(), Tick: now.Sub(c.origin)}
}
func (c *systemClock) NewTimer(d time.Duration) Timer { return realTimer{time.NewTimer(d)} }

type realTimer struct{ timer *time.Timer }

func (t realTimer) C() <-chan time.Time        { return t.timer.C }
func (t realTimer) Stop() bool                 { return t.timer.Stop() }
func (t realTimer) Reset(d time.Duration) bool { return t.timer.Reset(d) }
