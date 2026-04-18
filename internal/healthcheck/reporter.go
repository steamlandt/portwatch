package healthcheck

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Report is a JSON-serialisable snapshot of all health checks.
type Report struct {
	Overall string            `json:"overall"`
	Checks  []CheckReport     `json:"checks"`
	At      time.Time         `json:"at"`
}

// CheckReport is a JSON-serialisable single check result.
type CheckReport struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
	At      time.Time `json:"at"`
}

// Reporter writes health reports to a writer.
type Reporter struct {
	hc  *HealthCheck
	out io.Writer
}

// NewReporter returns a Reporter writing to out. If out is nil, os.Stdout is used.
func NewReporter(hc *HealthCheck, out io.Writer) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{hc: hc, out: out}
}

// Write serialises the current health state as JSON to the writer.
func (r *Reporter) Write() error {
	all := r.hc.All()
	checks := make([]CheckReport, 0, len(all))
	for _, c := range all {
		checks = append(checks, CheckReport{
			Name:    c.Name,
			Status:  c.Status.String(),
			Message: c.Message,
			At:      c.At,
		})
	}
	rep := Report{
		Overall: r.hc.Overall().String(),
		Checks:  checks,
		At:      time.Now(),
	}
	b, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return fmt.Errorf("healthcheck: marshal report: %w", err)
	}
	_, err = fmt.Fprintln(r.out, string(b))
	return err
}
