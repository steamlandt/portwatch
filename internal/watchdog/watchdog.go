package watchdog

import (
	"sync"
	"time"
)

// Status represents the health state of a watchdog.
type Status string

const (
	StatusHealthy  Status = "healthy"
	StatusStale    Status = "stale"
	StatusUnknown  Status = "unknown"
)

// Watchdog tracks whether a monitored loop is alive by requiring
// periodic heartbeats within a configured timeout window.
type Watchdog struct {
	mu        sync.Mutex
	timeout   time.Duration
	lastBeat  time.Time
	started   bool
}

// New creates a Watchdog with the given staleness timeout.
func New(timeout time.Duration) *Watchdog {
	return &Watchdog{timeout: timeout}
}

// Beat records a heartbeat at the current time.
func (w *Watchdog) Beat() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.lastBeat = time.Now()
	w.started = true
}

// Status returns the current health status.
func (w *Watchdog) Status() Status {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.started {
		return StatusUnknown
	}
	if time.Since(w.lastBeat) > w.timeout {
		return StatusStale
	}
	return StatusHealthy
}

// LastBeat returns the time of the most recent heartbeat.
func (w *Watchdog) LastBeat() time.Time {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.lastBeat
}

// Reset clears the watchdog state as if it was never started.
func (w *Watchdog) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.started = false
	w.lastBeat = time.Time{}
}
