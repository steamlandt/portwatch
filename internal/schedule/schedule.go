package schedule

import (
	"time"
)

// Scheduler triggers a callback at a fixed interval.
type Scheduler struct {
	interval time.Duration
	stop     chan struct{}
}

// New creates a Scheduler with the given interval.
func New(interval time.Duration) *Scheduler {
	return &Scheduler{
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Start begins the scheduler loop, calling fn on each tick.
// It runs in the current goroutine and blocks until Stop is called.
func (s *Scheduler) Start(fn func()) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fn()
		case <-s.stop:
			return
		}
	}
}

// Stop halts the scheduler.
func (s *Scheduler) Stop() {
	close(s.stop)
}

// Interval returns the configured tick interval.
func (s *Scheduler) Interval() time.Duration {
	return s.interval
}
