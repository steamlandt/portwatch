package dedup_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/dedup"
)

func TestFirstEventIsNotDuplicate(t *testing.T) {
	d := dedup.New()
	if d.IsDuplicate("tcp", 8080, "open") {
		t.Fatal("expected first event to not be a duplicate")
	}
}

func TestSameStateIsDuplicate(t *testing.T) {
	d := dedup.New()
	d.IsDuplicate("tcp", 8080, "open")
	if !d.IsDuplicate("tcp", 8080, "open") {
		t.Fatal("expected repeated same-state event to be a duplicate")
	}
}

func TestDifferentStateIsNotDuplicate(t *testing.T) {
	d := dedup.New()
	d.IsDuplicate("tcp", 8080, "open")
	if d.IsDuplicate("tcp", 8080, "closed") {
		t.Fatal("expected state change to not be a duplicate")
	}
}

func TestDifferentPortIsNotDuplicate(t *testing.T) {
	d := dedup.New()
	d.IsDuplicate("tcp", 8080, "open")
	if d.IsDuplicate("tcp", 9090, "open") {
		t.Fatal("expected different port to not be a duplicate")
	}
}

func TestResetClearsState(t *testing.T) {
	d := dedup.New()
	d.IsDuplicate("tcp", 8080, "open")
	d.Reset("tcp", 8080)
	if d.IsDuplicate("tcp", 8080, "open") {
		t.Fatal("expected reset port to not be a duplicate")
	}
}

func TestResetAllClearsAllStates(t *testing.T) {
	d := dedup.New()
	d.IsDuplicate("tcp", 8080, "open")
	d.IsDuplicate("tcp", 9090, "open")
	d.ResetAll()
	if d.IsDuplicate("tcp", 8080, "open") {
		t.Fatal("expected all states cleared after ResetAll")
	}
	if d.IsDuplicate("tcp", 9090, "open") {
		t.Fatal("expected all states cleared after ResetAll")
	}
}

func TestDifferentProtoIsNotDuplicate(t *testing.T) {
	d := dedup.New()
	d.IsDuplicate("tcp", 8080, "open")
	if d.IsDuplicate("udp", 8080, "open") {
		t.Fatal("expected different proto to not be a duplicate")
	}
}

func TestStateTransitionsAreTrackedIndependently(t *testing.T) {
	d := dedup.New()

	// Establish initial state for two ports
	d.IsDuplicate("tcp", 8080, "open")
	d.IsDuplicate("tcp", 9090, "open")

	// Change state on one port should not affect the other
	if d.IsDuplicate("tcp", 8080, "closed") {
		t.Fatal("expected state change on port 8080 to not be a duplicate")
	}
	if !d.IsDuplicate("tcp", 9090, "open") {
		t.Fatal("expected port 9090 state to remain independently tracked as duplicate")
	}
}
