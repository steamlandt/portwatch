package portrank_test

import (
	"testing"

	"github.com/user/portwatch/internal/portrank"
)

func TestScoreKnownCriticalPort(t *testing.T) {
	r := portrank.New(nil)
	got := r.Score(23) // Telnet
	if got != portrank.ScoreCritical {
		t.Errorf("expected critical score for port 23, got %v", got)
	}
}

func TestScoreKnownHighPort(t *testing.T) {
	r := portrank.New(nil)
	got := r.Score(22) // SSH
	if got != portrank.ScoreHigh {
		t.Errorf("expected high score for port 22, got %v", got)
	}
}

func TestScoreUnknownPortReturnsNone(t *testing.T) {
	r := portrank.New(nil)
	got := r.Score(9999)
	if got != portrank.ScoreNone {
		t.Errorf("expected none score for unknown port, got %v", got)
	}
}

func TestCustomOverridesDefault(t *testing.T) {
	overrides := map[int]portrank.Score{
		22: portrank.ScoreLow, // downgrade SSH risk
	}
	r := portrank.New(overrides)
	got := r.Score(22)
	if got != portrank.ScoreLow {
		t.Errorf("expected overridden low score for port 22, got %v", got)
	}
}

func TestCustomAddsNewPort(t *testing.T) {
	overrides := map[int]portrank.Score{
		9999: portrank.ScoreMedium,
	}
	r := portrank.New(overrides)
	got := r.Score(9999)
	if got != portrank.ScoreMedium {
		t.Errorf("expected medium score for custom port 9999, got %v", got)
	}
}

func TestIsCritical(t *testing.T) {
	r := portrank.New(nil)
	if !r.IsCritical(3389) {
		t.Error("expected port 3389 (RDP) to be critical")
	}
	if r.IsCritical(80) {
		t.Error("expected port 80 (HTTP) not to be critical")
	}
}

func TestScoreLabelCritical(t *testing.T) {
	if portrank.ScoreCritical.Label() != "critical" {
		t.Errorf("unexpected label: %s", portrank.ScoreCritical.Label())
	}
}

func TestScoreLabelNone(t *testing.T) {
	if portrank.ScoreNone.Label() != "none" {
		t.Errorf("unexpected label: %s", portrank.ScoreNone.Label())
	}
}

func TestScoreString(t *testing.T) {
	s := portrank.ScoreHigh
	got := s.String()
	if got != "75 (high)" {
		t.Errorf("unexpected string: %s", got)
	}
}
