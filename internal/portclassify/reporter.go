package portclassify

import (
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"
)

// Reporter prints a human-readable classification summary.
type Reporter struct {
	classifier *Classifier
	out        io.Writer
}

// NewReporter returns a Reporter that writes to w.
// If w is nil, os.Stdout is used.
func NewReporter(c *Classifier, w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{classifier: c, out: w}
}

// Report writes a tabulated classification report for the provided ports.
func (r *Reporter) Report(ports []Port) error {
	results := r.classifier.ClassifyAll(ports)

	// Collect and sort by port number for deterministic output.
	numbers := make([]int, 0, len(results))
	for n := range results {
		numbers = append(numbers, n)
	}
	sort.Ints(numbers)

	tw := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PORT\tPROTOCOL\tCLASS\tLABEL")
	for _, n := range numbers {
		res := results[n]
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\n",
			res.Port.Number, res.Port.Protocol, res.Class, res.Label)
	}
	return tw.Flush()
}

// Summary returns a count of ports per class.
func (r *Reporter) Summary(ports []Port) map[Class]int {
	counts := make(map[Class]int)
	for _, p := range ports {
		res := r.classifier.Classify(p)
		counts[res.Class]++
	}
	return counts
}
