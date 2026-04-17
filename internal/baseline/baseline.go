package baseline

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Baseline represents a known-good set of open ports captured at a point in time.
type Baseline struct {
	mu      sync.RWMutex
	path    string
	Ports   []scanner.Port `json:"ports"`
	CapturedAt time.Time   `json:"captured_at"`
}

// New returns a Baseline backed by the given file path.
func New(path string) *Baseline {
	return &Baseline{path: path}
}

// Capture records ports as the current baseline and persists to disk.
func (b *Baseline) Capture(ports []scanner.Port) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Ports = ports
	b.CapturedAt = time.Now().UTC()
	return b.save()
}

// Load reads the baseline from disk. Returns an empty baseline if the file does not exist.
func (b *Baseline) Load() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	data, err := os.ReadFile(b.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, b)
}

// Deviations returns ports present in current that are absent from the baseline,
// and ports in the baseline that are absent from current.
func (b *Baseline) Deviations(current []scanner.Port) (added, removed []scanner.Port) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	base := toSet(b.Ports)
	cur := toSet(current)
	for _, p := range current {
		if !base[key(p)] {
			added = append(added, p)
		}
	}
	for _, p := range b.Ports {
		if !cur[key(p)] {
			removed = append(removed, p)
		}
	}
	return
}

func (b *Baseline) save() error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(b.path, data, 0644)
}

func toSet(ports []scanner.Port) map[string]bool {
	m := make(map[string]bool, len(ports))
	for _, p := range ports {
		m[key(p)] = true
	}
	return m
}

func key(p scanner.Port) string {
	return p.Protocol + ":" + p.String()
}
