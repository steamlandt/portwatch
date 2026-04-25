package portrank_test

import (
	"testing"

	"github.com/user/portwatch/internal/portrank"
)

// TestDefaultRanksCoverCommonRiskyPorts verifies that the built-in table
// covers a representative set of commonly exploited ports.
func TestDefaultRanksCoverCommonRiskyPorts(t *testing.T) {
	r := portrank.New(nil)

	tests := []struct {
		port     int
		minScore portrank.Score
		desc     string
	}{
		{21, portrank.ScoreHigh, "FTP"},
		{22, portrank.ScoreHigh, "SSH"},
		{23, portrank.ScoreHigh, "Telnet"},
		{445, portrank.ScoreHigh, "SMB"},
		{3389, portrank.ScoreHigh, "RDP"},
		{3306, portrank.ScoreMedium, "MySQL"},
		{5432, portrank.ScoreMedium, "PostgreSQL"},
		{6379, portrank.ScoreMedium, "Redis"},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := r.Score(tc.port)
			if got < tc.minScore {
				t.Errorf("port %d (%s): expected score >= %v, got %v",
					tc.port, tc.desc, tc.minScore, got)
			}
		})
	}
}

// TestNilOverridesUsesDefaults ensures passing nil does not panic.
func TestNilOverridesUsesDefaults(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("New(nil) panicked: %v", r)
		}
	}()
	r := portrank.New(nil)
	if r == nil {
		t.Fatal("expected non-nil Ranker")
	}
}

// TestOverrideDoesNotMutateOtherInstance verifies that two rankers with
// different overrides are independent.
func TestOverrideDoesNotMutateOtherInstance(t *testing.T) {
	r1 := portrank.New(map[int]portrank.Score{22: portrank.ScoreLow})
	r2 := portrank.New(nil)

	if r1.Score(22) == r2.Score(22) {
		t.Error("expected r1 and r2 to have different scores for port 22")
	}
}
