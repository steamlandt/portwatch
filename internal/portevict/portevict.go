// Package portevict tracks ports that have been evicted (closed after being
// open for a sustained period) and provides a rolling eviction log.
package portevict

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry records a single eviction event.
type Entry struct {
	Port      scanner.Port `json:"port"`
	OpenedAt  time.Time    `json:"opened_at"`
	ClosedAt  time.Time    `json:"closed_at"`
	Duration  string       `json:"duration"`
}

// Tracker maintains an in-memory eviction log backed by an optional file.
type Tracker struct {
	mu      sync.Mutex
	entries []Entry
	path    string
}

// New returns a Tracker. If path is non-empty the log is persisted to disk.
func New(path string) (*Tracker, error) {
	t := &Tracker{path: path}
	if path != "" {
		if err := t.load(); err != nil {
			return nil, err
		}
	}
	return t, nil
}

// Record appends an eviction event for the given port.
func (t *Tracker) Record(p scanner.Port, openedAt, closedAt time.Time) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	e := Entry{
		Port:     p,
		OpenedAt: openedAt,
		ClosedAt: closedAt,
		Duration: closedAt.Sub(openedAt).String(),
	}
	t.entries = append(t.entries, e)

	if t.path != "" {
		return t.save()
	}
	return nil
}

// All returns a copy of all recorded eviction entries.
func (t *Tracker) All() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, len(t.entries))
	copy(out, t.entries)
	return out
}

// Reset clears all entries from memory (and disk if configured).
func (t *Tracker) Reset() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = nil
	if t.path != "" {
		return t.save()
	}
	return nil
}

func (t *Tracker) load() error {
	data, err := os.ReadFile(t.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &t.entries)
}

func (t *Tracker) save() error {
	data, err := json.Marshal(t.entries)
	if err != nil {
		return err
	}
	return os.WriteFile(t.path, data, 0o644)
}
