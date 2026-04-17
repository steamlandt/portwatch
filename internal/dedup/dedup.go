// Package dedup provides deduplication of port change events to avoid
// firing repeated alerts for the same port state.
package dedup

import "sync"

// Entry tracks the last known state for a port key.
type Entry struct {
	Proto  string
	Port   int
	State  string
}

// Deduplicator suppresses repeated events for the same port+state combination.
type Deduplicator struct {
	mu   sync.Mutex
	seen map[string]string // key -> last state
}

// New returns a new Deduplicator.
func New() *Deduplicator {
	return &Deduplicator{seen: make(map[string]string)}
}

// IsDuplicate returns true if the same port+state was already reported.
// If not a duplicate, it records the new state and returns false.
func (d *Deduplicator) IsDuplicate(proto string, port int, state string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	key := buildKey(proto, port)
	if prev, ok := d.seen[key]; ok && prev == state {
		return true
	}
	d.seen[key] = state
	return false
}

// Reset clears the recorded state for a specific port.
func (d *Deduplicator) Reset(proto string, port int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.seen, buildKey(proto, port))
}

// ResetAll clears all recorded states.
func (d *Deduplicator) ResetAll() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]string)
}

func buildKey(proto string, port int) string {
	return proto + ":" + itoa(port)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
