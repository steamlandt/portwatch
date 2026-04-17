package snapshot

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// DiffReport summarises the changes between two snapshots.
type DiffReport struct {
	Added     []scanner.Port
	Removed   []scanner.Port
	GeneratedAt time.Time
}

// NewDiffReport builds a DiffReport from two snapshots.
func NewDiffReport(prev, next Snapshot) DiffReport {
	added, removed := Diff(prev, next)
	return DiffReport{
		Added:       added,
		Removed:     removed,
		GeneratedAt: time.Now().UTC(),
	}
}

// HasChanges returns true when at least one port was added or removed.
func (r DiffReport) HasChanges() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0
}

// WriteTo formats the report as human-readable text and writes it to w.
func (r DiffReport) WriteTo(w io.Writer) (int64, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Snapshot diff at %s\n", r.GeneratedAt.Format(time.RFC3339)))

	if len(r.Added) == 0 && len(r.Removed) == 0 {
		sb.WriteString("  No changes detected.\n")
	}
	for _, p := range r.Added {
		sb.WriteString(fmt.Sprintf("  [+] %s\n", p.String()))
	}
	for _, p := range r.Removed {
		sb.WriteString(fmt.Sprintf("  [-] %s\n", p.String()))
	}

	n, err := io.WriteString(w, sb.String())
	return int64(n), err
}
