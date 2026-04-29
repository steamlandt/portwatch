// Package portcorrelate identifies relationships between ports that tend to
// appear or disappear together, helping distinguish coordinated service
// restarts from unexpected port changes.
package portcorrelate

import (
	"sync"
	"time"
)

// Pair represents two port numbers that have been observed co-occurring.
type Pair struct {
	A        int
	B        int
	Count    int
	LastSeen time.Time
}

// Correlator tracks which ports open or close within the same scan cycle.
type Correlator struct {
	mu      sync.Mutex
	pairs   map[string]*Pair
	window  time.Duration
}

// New returns a Correlator that groups ports observed within window of each other.
func New(window time.Duration) *Correlator {
	return &Correlator{
		pairs:  make(map[string]*Pair),
		window: window,
	}
}

// Observe records that ports a and b were seen changing state in the same cycle.
func (c *Correlator) Observe(a, b int) {
	if a == b {
		return
	}
	if a > b {
		a, b = b, a
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	k := buildKey(a, b)
	p, ok := c.pairs[k]
	if !ok {
		p = &Pair{A: a, B: b}
		c.pairs[k] = p
	}
	p.Count++
	p.LastSeen = time.Now()
}

// All returns all recorded pairs.
func (c *Correlator) All() []Pair {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]Pair, 0, len(c.pairs))
	for _, p := range c.pairs {
		out = append(out, *p)
	}
	return out
}

// Strong returns pairs whose co-occurrence count meets or exceeds minCount.
func (c *Correlator) Strong(minCount int) []Pair {
	c.mu.Lock()
	defer c.mu.Unlock()
	var out []Pair
	for _, p := range c.pairs {
		if p.Count >= minCount {
			out = append(out, *p)
		}
	}
	return out
}

// Reset clears all recorded correlation data.
func (c *Correlator) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pairs = make(map[string]*Pair)
}

func buildKey(a, b int) string {
	return fmt.Sprintf("%d:%d", a, b)
}
