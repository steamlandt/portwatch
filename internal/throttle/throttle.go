package throttle

import (
	"sync"
	"time"
)

// Throttler limits how frequently alerts fire per port key.
type Throttler struct {
	mu       sync.Mutex
	last     map[string]time.Time
	interval time.Duration
}

// New creates a Throttler with the given minimum interval between alerts.
func New(interval time.Duration) *Throttler {
	return &Throttler{
		last:     make(map[string]time.Time),
		interval: interval,
	}
}

// Allow returns true if enough time has passed since the last alert for key.
func (t *Throttler) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	if last, ok := t.last[key]; ok {
		if now.Sub(last) < t.interval {
			return false
		}
	}
	t.last[key] = now
	return true
}

// Reset clears the throttle state for a specific key.
func (t *Throttler) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}

// ResetAll clears all throttle state.
func (t *Throttler) ResetAll() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.last = make(map[string]time.Time)
}

// Interval returns the configured throttle interval.
func (t *Throttler) Interval() time.Duration {
	return t.interval
}
