package watchdog

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Reporter periodically writes the watchdog status to a writer.
type Reporter struct {
	dog    *Watchdog
	out    io.Writer
	ticker *time.Ticker
	stop   chan struct{}
}

// NewReporter creates a Reporter that logs status every interval.
func NewReporter(dog *Watchdog, interval time.Duration) *Reporter {
	return &Reporter{
		dog:    dog,
		out:    os.Stdout,
		ticker: time.NewTicker(interval),
		stop:   make(chan struct{}),
	}
}

// WithWriter overrides the output writer.
func (r *Reporter) WithWriter(w io.Writer) *Reporter {
	r.out = w
	return r
}

// Start begins the reporting loop in a goroutine.
func (r *Reporter) Start() {
	go func() {
		for {
			select {
			case <-r.ticker.C:
				r.report()
			case <-r.stop:
				return
			}
		}
	}()
}

// Stop halts the reporting loop.
func (r *Reporter) Stop() {
	r.ticker.Stop()
	close(r.stop)
}

func (r *Reporter) report() {
	status := r.dog.Status()
	last := r.dog.LastBeat()
	var age string
	if last.IsZero() {
		age = "never"
	} else {
		age = fmt.Sprintf("%.1fs ago", time.Since(last).Seconds())
	}
	fmt.Fprintf(r.out, "[watchdog] status=%s last_beat=%s\n", status, age)
}
