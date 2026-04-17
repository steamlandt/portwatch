package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry records a port change event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Port      int       `json:"port"`
	Proto     string    `json:"proto"`
	Event     string    `json:"event"` // "opened" or "closed"
}

// History maintains an in-memory log of port change events with optional persistence.
type History struct {
	mu      sync.RWMutex
	entries []Entry
	path    string
}

// New creates a History, loading existing entries from path if it exists.
func New(path string) (*History, error) {
	h := &History{path: path}
	if path == "" {
		return h, nil
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return h, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &h.entries); err != nil {
		return nil, err
	}
	return h, nil
}

// Record appends a new entry and persists if a path is configured.
func (h *History) Record(port int, proto, event string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = append(h.entries, Entry{
		Timestamp: time.Now().UTC(),
		Port:      port,
		Proto:     proto,
		Event:     event,
	})
	if h.path == "" {
		return nil
	}
	return h.flush()
}

// All returns a copy of all recorded entries.
func (h *History) All() []Entry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]Entry, len(h.entries))
	copy(out, h.entries)
	return out
}

func (h *History) flush() error {
	data, err := json.MarshalIndent(h.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0644)
}
