// Package portmeta aggregates metadata about a port, combining label, group,
// rank, and baseline deviation into a single enriched description.
package portmeta

import (
	"fmt"

	"github.com/example/portwatch/internal/baseline"
	"github.com/example/portwatch/internal/portgroup"
	"github.com/example/portwatch/internal/portlabel"
	"github.com/example/portwatch/internal/portrank"
	"github.com/example/portwatch/internal/scanner"
)

// Meta holds enriched metadata for a single port.
type Meta struct {
	Port      scanner.Port
	Label     string
	Group     string
	Rank      string
	Deviation string // "added", "removed", or "" if within baseline
}

// String returns a human-readable summary of the port metadata.
func (m Meta) String() string {
	label := m.Label
	if label == "" {
		label = fmt.Sprintf("%d", m.Port.Number)
	}
	s := fmt.Sprintf("port=%s group=%s rank=%s", label, m.Group, m.Rank)
	if m.Deviation != "" {
		s += fmt.Sprintf(" deviation=%s", m.Deviation)
	}
	return s
}

// Enricher builds Meta values for ports.
type Enricher struct {
	labeler  *portlabel.Labeler
	grouper  *portgroup.Grouper
	ranker   *portrank.Ranker
	baseline *baseline.Baseline
}

// New returns an Enricher using the provided sub-components.
// baseline may be nil if deviation tracking is not required.
func New(l *portlabel.Labeler, g *portgroup.Grouper, r *portrank.Ranker, b *baseline.Baseline) *Enricher {
	return &Enricher{labeler: l, grouper: g, ranker: r, baseline: b}
}

// Enrich returns a Meta for the given port.
func (e *Enricher) Enrich(p scanner.Port) Meta {
	m := Meta{
		Port:  p,
		Label: e.labeler.Label(p.Number),
		Group: e.grouper.Categorize(p.Number),
		Rank:  e.ranker.Score(p.Number).String(),
	}
	if e.baseline != nil {
		deviations := e.baseline.Deviations([]scanner.Port{p})
		for _, d := range deviations {
			if d.Port.Number == p.Number {
				m.Deviation = d.Kind
				break
			}
		}
	}
	return m
}

// EnrichAll returns a Meta slice for every port in the provided list.
func (e *Enricher) EnrichAll(ports []scanner.Port) []Meta {
	out := make([]Meta, len(ports))
	for i, p := range ports {
		out[i] = e.Enrich(p)
	}
	return out
}
