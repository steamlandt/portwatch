package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event describes a change in port state.
type Event struct {
	Level     Level
	Port      scanner.Port
	Message   string
	Timestamp time.Time
}

// Notifier sends alert events to a destination.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier writing to the given writer.
// Pass nil to default to os.Stdout.
func New(out io.Writer) *Notifier {
	if out == nil {
		out = os.Stdout
	}
	return &Notifier{out: out}
}

// PortOpened fires an ALERT event for a newly opened port.
func (n *Notifier) PortOpened(p scanner.Port) {
	n.send(Event{
		Level:     LevelAlert,
		Port:      p,
		Message:   "unexpected port opened",
		Timestamp: time.Now(),
	})
}

// PortClosed fires a WARN event for a port that has closed.
func (n *Notifier) PortClosed(p scanner.Port) {
	n.send(Event{
		Level:     LevelWarn,
		Port:      p,
		Message:   "previously open port closed",
		Timestamp: time.Now(),
	})
}

func (n *Notifier) send(e Event) {
	fmt.Fprintf(n.out, "%s [%s] %s — %s\n",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.Port.String(),
		e.Message,
	)
}
