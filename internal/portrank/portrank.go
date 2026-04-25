// Package portrank assigns a risk score to ports based on their exposure profile.
// Higher scores indicate ports that are more commonly targeted or sensitive.
package portrank

import "fmt"

// Score represents a risk level from 0 (unknown/neutral) to 100 (critical).
type Score int

const (
	ScoreNone     Score = 0
	ScoreLow      Score = 25
	ScoreMedium   Score = 50
	ScoreHigh     Score = 75
	ScoreCritical Score = 100
)

// Label returns a human-readable risk label for the score.
func (s Score) Label() string {
	switch {
	case s >= ScoreCritical:
		return "critical"
	case s >= ScoreHigh:
		return "high"
	case s >= ScoreMedium:
		return "medium"
	case s >= ScoreLow:
		return "low"
	default:
		return "none"
	}
}

// String returns a formatted representation of the score.
func (s Score) String() string {
	return fmt.Sprintf("%d (%s)", int(s), s.Label())
}

// defaultRanks maps well-known port numbers to their risk scores.
var defaultRanks = map[int]Score{
	21:   ScoreCritical, // FTP
	22:   ScoreHigh,     // SSH
	23:   ScoreCritical, // Telnet
	25:   ScoreHigh,     // SMTP
	80:   ScoreLow,      // HTTP
	443:  ScoreLow,      // HTTPS
	445:  ScoreCritical, // SMB
	1433: ScoreCritical, // MSSQL
	1521: ScoreHigh,     // Oracle DB
	3306: ScoreHigh,     // MySQL
	3389: ScoreCritical, // RDP
	5432: ScoreHigh,     // PostgreSQL
	5900: ScoreHigh,     // VNC
	6379: ScoreHigh,     // Redis
	8080: ScoreLow,      // HTTP Alt
	8443: ScoreLow,      // HTTPS Alt
	27017: ScoreHigh,    // MongoDB
}

// Ranker assigns risk scores to ports.
type Ranker struct {
	ranks map[int]Score
}

// New creates a Ranker with optional overrides merged on top of defaults.
func New(overrides map[int]Score) *Ranker {
	ranks := make(map[int]Score, len(defaultRanks))
	for k, v := range defaultRanks {
		ranks[k] = v
	}
	for k, v := range overrides {
		ranks[k] = v
	}
	return &Ranker{ranks: ranks}
}

// Score returns the risk score for the given port number.
// Unknown ports receive ScoreNone.
func (r *Ranker) Score(port int) Score {
	if s, ok := r.ranks[port]; ok {
		return s
	}
	return ScoreNone
}

// IsCritical reports whether the port is ranked at the critical level.
func (r *Ranker) IsCritical(port int) bool {
	return r.Score(port) >= ScoreCritical
}
