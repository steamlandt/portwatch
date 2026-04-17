// Package suppress provides a time-based suppression mechanism to avoid
// repeated alerts for ports that remain in the same changed state.
package suppress

import (
	"sync"
	"time"
)

// Suppressor tracks when a key was last alerted and suppresses repeats
// within a configurable window.
type Suppressor struct {
	mu       sync.Mutex
	window   time.Duration
	lastSeen map[string]time.Time
	now      func() time.Time
}

// New returns a Suppressor with the given suppression window.
func New(window time.Duration) *Suppressor {
	return &Suppressor{
		window:   window,
		lastSeen: make(map[string]time.Time),
		now:      time.Now,
	}
}

// IsSuppressed returns true if the key has been seen within the window.
// If not suppressed, it records the current time for the key.
func (s *Suppressor) IsSuppressed(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if last, ok := s.lastSeen[key]; ok {
		if s.now().Sub(last) < s.window {
			return true
		}
	}
	s.lastSeen[key] = s.now()
	return false
}

// Reset removes the suppression record for a key.
func (s *Suppressor) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.lastSeen, key)
}

// ResetAll clears all suppression records.
func (s *Suppressor) ResetAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastSeen = make(map[string]time.Time)
}
