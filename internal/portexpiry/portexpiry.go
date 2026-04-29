// Package portexpiry tracks how long a port has been continuously open
// and emits a warning when it exceeds a configured maximum duration.
package portexpiry

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry records when a port was first seen open.
type Entry struct {
	Port      scanner.Port
	FirstSeen time.Time
	Duration  time.Duration
}

// Expiry holds open-port first-seen timestamps and checks for long-lived ports.
type Expiry struct {
	mu      sync.Mutex
	entries map[string]time.Time
	maxAge  time.Duration
	now     func() time.Time
}

// New creates an Expiry that warns when a port has been open longer than maxAge.
func New(maxAge time.Duration) *Expiry {
	return &Expiry{
		entries: make(map[string]time.Time),
		maxAge:  maxAge,
		now:     time.Now,
	}
}

func key(p scanner.Port) string {
	return fmt.Sprintf("%s:%d", p.Proto, p.Number)
}

// Observe records the first time a port is seen. Subsequent calls for the
// same port are no-ops until the port is removed via Remove.
func (e *Expiry) Observe(p scanner.Port) {
	e.mu.Lock()
	defer e.mu.Unlock()
	k := key(p)
	if _, exists := e.entries[k]; !exists {
		e.entries[k] = e.now()
	}
}

// Remove deletes the tracking record for a port (e.g. when it closes).
func (e *Expiry) Remove(p scanner.Port) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.entries, key(p))
}

// Expired returns all ports whose open duration exceeds the configured maxAge.
func (e *Expiry) Expired(ports []scanner.Port) []Entry {
	e.mu.Lock()
	defer e.mu.Unlock()
	now := e.now()
	var out []Entry
	for _, p := range ports {
		k := key(p)
		first, ok := e.entries[k]
		if !ok {
			continue
		}
		d := now.Sub(first)
		if d >= e.maxAge {
			out = append(out, Entry{Port: p, FirstSeen: first, Duration: d})
		}
	}
	return out
}

// Age returns how long the given port has been tracked, and whether it is known.
func (e *Expiry) Age(p scanner.Port) (time.Duration, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	first, ok := e.entries[key(p)]
	if !ok {
		return 0, false
	}
	return e.now().Sub(first), true
}
