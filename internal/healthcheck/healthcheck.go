package healthcheck

import (
	"sync"
	"time"
)

// Status represents the health of the daemon.
type Status int

const (
	StatusOK Status = iota
	StatusDegraded
	StatusDown
)

func (s Status) String() string {
	switch s {
	case StatusOK:
		return "ok"
	case StatusDegraded:
		return "degraded"
	case StatusDown:
		return "down"
	default:
		return "unknown"
	}
}

// Check holds the result of a single health probe.
type Check struct {
	Name    string
	Status  Status
	Message string
	At      time.Time
}

// HealthCheck aggregates named health probes.
type HealthCheck struct {
	mu     sync.RWMutex
	checks map[string]Check
}

// New returns an initialised HealthCheck.
func New() *HealthCheck {
	return &HealthCheck{checks: make(map[string]Check)}
}

// Set records the result of a named probe.
func (h *HealthCheck) Set(name string, status Status, msg string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checks[name] = Check{Name: name, Status: status, Message: msg, At: time.Now()}
}

// Get returns the last recorded result for a probe.
func (h *HealthCheck) Get(name string) (Check, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	c, ok := h.checks[name]
	return c, ok
}

// Overall returns the worst status across all registered probes.
func (h *HealthCheck) Overall() Status {
	h.mu.RLock()
	defer h.mu.RUnlock()
	worst := StatusOK
	for _, c := range h.checks {
		if c.Status > worst {
			worst = c.Status
		}
	}
	return worst
}

// All returns a copy of all registered checks.
func (h *HealthCheck) All() []Check {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]Check, 0, len(h.checks))
	for _, c := range h.checks {
		out = append(out, c)
	}
	return out
}
