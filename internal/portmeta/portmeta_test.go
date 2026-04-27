package portmeta_test

import (
	"strings"
	"testing"

	"github.com/example/portwatch/internal/portgroup"
	"github.com/example/portwatch/internal/portlabel"
	"github.com/example/portwatch/internal/portmeta"
	"github.com/example/portwatch/internal/portrank"
	"github.com/example/portwatch/internal/scanner"
)

func makeEnricher() *portmeta.Enricher {
	l := portlabel.New(nil)
	g := portgroup.New(nil)
	r := portrank.New(nil)
	return portmeta.New(l, g, r, nil)
}

func makePort(n int) scanner.Port {
	return scanner.Port{Number: n, Protocol: "tcp"}
}

func TestEnrichKnownPort(t *testing.T) {
	e := makeEnricher()
	m := e.Enrich(makePort(80))

	if m.Port.Number != 80 {
		t.Fatalf("expected port 80, got %d", m.Port.Number)
	}
	if m.Label == "" {
		t.Error("expected non-empty label for port 80")
	}
	if m.Group == "" {
		t.Error("expected non-empty group for port 80")
	}
	if m.Rank == "" {
		t.Error("expected non-empty rank for port 80")
	}
}

func TestEnrichUnknownPortHasEmptyLabel(t *testing.T) {
	e := makeEnricher()
	m := e.Enrich(makePort(19999))

	if m.Label != "" {
		t.Errorf("expected empty label for unknown port, got %q", m.Label)
	}
}

func TestEnrichNoDeviationWhenBaselineNil(t *testing.T) {
	e := makeEnricher()
	m := e.Enrich(makePort(443))

	if m.Deviation != "" {
		t.Errorf("expected no deviation without baseline, got %q", m.Deviation)
	}
}

func TestEnrichAllReturnsSameLength(t *testing.T) {
	e := makeEnricher()
	ports := []scanner.Port{makePort(22), makePort(80), makePort(443)}
	metas := e.EnrichAll(ports)

	if len(metas) != len(ports) {
		t.Fatalf("expected %d metas, got %d", len(ports), len(metas))
	}
}

func TestMetaStringContainsFields(t *testing.T) {
	e := makeEnricher()
	m := e.Enrich(makePort(22))
	s := m.String()

	for _, key := range []string{"port=", "group=", "rank="} {
		if !strings.Contains(s, key) {
			t.Errorf("Meta.String() missing %q: %s", key, s)
		}
	}
}

func TestMetaStringFallsBackToPortNumber(t *testing.T) {
	e := makeEnricher()
	m := e.Enrich(makePort(19999))
	s := m.String()

	if !strings.Contains(s, "19999") {
		t.Errorf("expected port number in string for unknown port, got: %s", s)
	}
}
