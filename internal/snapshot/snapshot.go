package snapshot

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot holds a point-in-time capture of open ports.
type Snapshot struct {
	Ports     []scanner.Port `json:"ports"`
	CapturedAt time.Time     `json:"captured_at"`
}

// Store persists and retrieves port snapshots.
type Store struct {
	mu   sync.RWMutex
	path string
}

// New returns a Store backed by the given file path.
func New(path string) *Store {
	return &Store{path: path}
}

// Save writes the current port list as a snapshot to disk.
func (s *Store) Save(ports []scanner.Port) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	snap := Snapshot{
		Ports:      ports,
		CapturedAt: time.Now().UTC(),
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

// Load reads the latest snapshot from disk.
// Returns an empty Snapshot if the file does not exist.
func (s *Store) Load() (Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return Snapshot{}, nil
	}
	if err != nil {
		return Snapshot{}, err
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}

// Diff returns ports added and removed between prev and next snapshots.
func Diff(prev, next Snapshot) (added, removed []scanner.Port) {
	prevSet := make(map[string]struct{}, len(prev.Ports))
	for _, p := range prev.Ports {
		prevSet[p.String()] = struct{}{}
	}
	nextSet := make(map[string]struct{}, len(next.Ports))
	for _, p := range next.Ports {
		nextSet[p.String()] = struct{}{}
	}
	for _, p := range next.Ports {
		if _, ok := prevSet[p.String()]; !ok {
			added = append(added, p)
		}
	}
	for _, p := range prev.Ports {
		if _, ok := nextSet[p.String()]; !ok {
			removed = append(removed, p)
		}
	}
	return
}
