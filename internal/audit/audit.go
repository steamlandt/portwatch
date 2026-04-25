// Package audit records a structured log of every port change event
// observed by the monitor, providing a queryable trail for review.
package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Entry represents a single audited port event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"` // "opened" | "closed"
	Port      uint16    `json:"port"`
	Protocol  string    `json:"protocol"`
	Label     string    `json:"label,omitempty"`
}

// Log is a persistent, append-only audit log.
type Log struct {
	mu   sync.Mutex
	path string
}

// New creates a new Log that persists entries to path.
func New(path string) *Log {
	return &Log{path: path}
}

// Record appends an entry to the audit log.
func (l *Log) Record(event, protocol string, port uint16, label string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	e := Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Port:      port,
		Protocol:  protocol,
		Label:     label,
	}

	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}

	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("audit: open: %w", err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s\n", data)
	if err != nil {
		return fmt.Errorf("audit: write: %w", err)
	}
	return nil
}

// All reads and returns all entries from the audit log.
// Returns an empty slice if the file does not exist.
func (l *Log) All() ([]Entry, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	raw, err := os.ReadFile(l.path)
	if os.IsNotExist(err) {
		return []Entry{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("audit: read: %w", err)
	}

	var entries []Entry
	dec := json.NewDecoder(bytes.NewReader(raw))
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
