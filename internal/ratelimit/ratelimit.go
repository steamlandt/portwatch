package ratelimit

import (
	"sync"
	"time"
)

// Limiter suppresses repeated alerts for the same port within a cooldown window.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
}

// New creates a Limiter with the given cooldown duration.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
	}
}

// Allow returns true if the key has not been seen within the cooldown window.
// If allowed, it records the current time for that key.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if t, ok := l.last[key]; ok && now.Sub(t) < l.cooldown {
		return false
	}
	l.last[key] = now
	return true
}

// Reset clears the recorded time for the given key.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, key)
}

// ResetAll clears all recorded times.
func (l *Limiter) ResetAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = make(map[string]time.Time)
}
