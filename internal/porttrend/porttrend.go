// Package porttrend tracks how frequently ports appear across scans
// and surfaces ports that are consistently open versus transient.
package porttrend

import (
	"sync"
	"time"
)

// Entry records observation counts and timing for a single port.
type Entry struct {
	Port      string
	Seen      int
	FirstSeen time.Time
	LastSeen  time.Time
}

// Tracker accumulates port observations over time.
type Tracker struct {
	mu      sync.Mutex
	entries map[string]*Entry
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[string]*Entry),
	}
}

// Observe records that port was seen at the given time.
func (t *Tracker) Observe(port string, at time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	e, ok := t.entries[port]
	if !ok {
		e = &Entry{
			Port:      port,
			FirstSeen: at,
		}
		t.entries[port] = e
	}
	e.Seen++
	e.LastSeen = at
}

// Get returns the Entry for the given port and whether it exists.
func (t *Tracker) Get(port string) (Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	e, ok := t.entries[port]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// All returns a snapshot of every tracked entry.
func (t *Tracker) All() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()

	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, *e)
	}
	return out
}

// Transient returns ports seen fewer than minSeen times.
func (t *Tracker) Transient(minSeen int) []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()

	var out []Entry
	for _, e := range t.entries {
		if e.Seen < minSeen {
			out = append(out, *e)
		}
	}
	return out
}

// Reset removes all tracked data.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[string]*Entry)
}
