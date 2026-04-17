package monitor

import (
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// Monitor periodically scans ports and emits alerts on state changes.
type Monitor struct {
	scanner  *scanner.Scanner
	store    *state.Store
	aleriter  *alert.Alerter
	interval time.Duration
	stop     chan struct{}
}

// Config holds Monitor configuration.
type Config struct {
	StartPort int
	EndPort   int
	Interval  time.Duration
	StatePath string
}

// New creates a Monitor from the given config.
func New(cfg Config, alerter *alert.Alerter) (*Monitor, error) {
	sc, err := scanner.New(cfg.StartPort, cfg.EndPort)
	if err != nil {
		return nil, err
	}
	st, err := state.New(cfg.StatePath)
	if err != nil {
		return nil, err
	}
	return &Monitor{
		scanner:  sc,
		store:    st,
		aleriter:  alerter,
		interval: cfg.Interval,
		stop:     make(chan struct{}),
	}, nil
}

// Start begins periodic scanning in the background.
func (m *Monitor) Start() {
	go func() {
		ticker := time.NewTicker(m.interval)
		defer ticker.Stop()
		m.scan()
		for {
			select {
			case <-ticker.C:
				m.scan()
			case <-m.stop:
				return
			}
		}
	}()
}

// Stop halts the monitor.
func (m *Monitor) Stop() {
	close(m.stop)
}

// ScanOnce performs a single scan synchronously and returns the number of open ports found.
func (m *Monitor) ScanOnce() int {
	open := m.scanner.Scan()
	m.processScan(open)
	return len(open)
}

func (m *Monitor) scan() {
	open := m.scanner.Scan()
	m.processScan(open)
}

func (m *Monitor) processScan(open []scanner.Port) {
	openSet := make(map[int]bool, len(open))
	for _, p := range open {
		openSet[p.Number] = true
	}

	// Detect newly opened ports.
	for _, p := range open {
		prev, known := m.store.Get(p.Number)
		if !known || !prev.Open {
			m.aleriter.PortOpened(p)
		}
		m.store.Set(state.PortState{Port: p.Number, Open: true, LastSeen: time.Now()})
	}

	// Detect newly closed ports.
	for _, prev := range m.store.All() {
		if prev.Open && !openSet[prev.Port] {
			m.aleriter.PortClosed(scanner.Port{Number: prev.Port})
			m.store.Set(state.PortState{Port: prev.Port, Open: false, LastSeen: time.Now()})
		}
	}
}
