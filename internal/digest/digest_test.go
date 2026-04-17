package digest

import (
	"testing"
)

func TestComputeReturnsDeterministicHash(t *testing.T) {
	d := New()
	ports := []Port{
		{Proto: "tcp", Number: 80, State: "open"},
		{Proto: "tcp", Number: 443, State: "open"},
	}

	h1, err := d.Compute(ports)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	h2, err := d.Compute(ports)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if h1 != h2 {
		t.Errorf("expected identical hashes, got %s and %s", h1, h2)
	}
}

func TestComputeEmptyPortsReturnsStableHash(t *testing.T) {
	d := New()
	h, err := d.Compute([]Port{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == "" {
		t.Error("expected non-empty hash for empty port list")
	}
}

func TestComputeDifferentPortsProduceDifferentHashes(t *testing.T) {
	d := New()
	a, _ := d.Compute([]Port{{Proto: "tcp", Number: 80, State: "open"}})
	b, _ := d.Compute([]Port{{Proto: "tcp", Number: 8080, State: "open"}})

	if a == b {
		t.Error("expected different hashes for different port sets")
	}
}

func TestEqualMatchingHashes(t *testing.T) {
	d := New()
	if !d.Equal("abc123", "abc123") {
		t.Error("expected equal hashes to match")
	}
}

func TestEqualMismatchedHashes(t *testing.T) {
	d := New()
	if d.Equal("abc123", "xyz789") {
		t.Error("expected different hashes to not match")
	}
}

func TestEqualEmptyHashReturnsFalse(t *testing.T) {
	d := New()
	if d.Equal("", "") {
		t.Error("expected empty hash comparison to return false")
	}
}
