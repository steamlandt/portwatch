package state

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// PortState represents the last known state of a port.
type PortState struct {
	Port     int       `json:"port"`
	Open     bool      `json:"open"`
	LastSeen time.Time `json:"last_seen"`
}

// Store holds the known port states and persists them to disk.
type Store struct {
	mu    sync.RWMutex
	ports map[int]PortState
	path  string
}

// New creates a new Store, loading existing state from path if it exists.
func New(path string) (*Store, error) {
	s := &Store{
		ports: make(map[int]PortState),
		path:  path,
	}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

// Get returns the PortState for a given port and whether it was found.
func (s *Store) Get(port int) (PortState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ps, ok := s.ports[port]
	return ps, ok
}

// Set updates the state for a port and persists to disk.
func (s *Store) Set(ps PortState) error {
	s.mu.Lock()
	s.ports[ps.Port] = ps
	s.mu.Unlock()
	return s.save()
}

// All returns a copy of all stored port states.
func (s *Store) All() []PortState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]PortState, 0, len(s.ports))
	for _, ps := range s.ports {
		out = append(out, ps)
	}
	return out
}

func (s *Store) load() error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	var states []PortState
	if err := json.NewDecoder(f).Decode(&states); err != nil {
		return err
	}
	for _, ps := range states {
		s.ports[ps.Port] = ps
	}
	return nil
}

func (s *Store) save() error {
	s.mu.RLock()
	states := make([]PortState, 0, len(s.ports))
	for _, ps := range s.ports {
		states = append(states, ps)
	}
	s.mu.RUnlock()

	tmp := s.path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if err := json.NewEncoder(f).Encode(states); err != nil {
		f.Close()
		return err
	}
	f.Close()
	return os.Rename(tmp, s.path)
}
