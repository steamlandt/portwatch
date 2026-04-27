package porttrend

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
	"time"
)

// Reporter writes a human-readable trend summary to a writer.
type Reporter struct {
	tracker *Tracker
	out     io.Writer
}

// NewReporter creates a Reporter backed by the given Tracker.
func NewReporter(tracker *Tracker, out io.Writer) *Reporter {
	return &Reporter{tracker: tracker, out: out}
}

// Report writes a sorted table of all tracked ports to the writer.
func (r *Reporter) Report() error {
	entries := r.tracker.All()
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Seen > entries[j].Seen
	})

	w := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PORT\tSEEN\tFIRST SEEN\tLAST SEEN")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\n",
			e.Port,
			e.Seen,
			e.FirstSeen.Format(time.RFC3339),
			e.LastSeen.Format(time.RFC3339),
		)
	}
	return w.Flush()
}

// ReportTransient writes only ports seen fewer than minSeen times.
func (r *Reporter) ReportTransient(minSeen int) error {
	entries := r.tracker.Transient(minSeen)
	if len(entries) == 0 {
		_, err := fmt.Fprintln(r.out, "no transient ports detected")
		return err
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Port < entries[j].Port
	})

	w := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PORT\tSEEN\tFIRST SEEN")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%d\t%s\n",
			e.Port,
			e.Seen,
			e.FirstSeen.Format(time.RFC3339),
		)
	}
	return w.Flush()
}
