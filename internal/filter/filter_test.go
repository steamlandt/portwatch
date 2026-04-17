package filter

import (
	"testing"
)

func TestAllowPermitsUnignoredPort(t *testing.T) {
	f := New([]int{22, 80})
	if !f.Allow(443) {
		t.Error("expected port 443 to be allowed")
	}
}

func TestAllowBlocksIgnoredPort(t *testing.T) {
	f := New([]int{22, 80})
	if f.Allow(22) {
		t.Error("expected port 22 to be blocked")
	}
}

func TestFilterPortsRemovesIgnored(t *testing.T) {
	f := New([]int{22, 80})
	input := []int{22, 443, 80, 8080}
	result := f.FilterPorts(input)

	expected := []int{443, 8080}
	if len(result) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, result)
	}
	for i, p := range expected {
		if result[i] != p {
			t.Errorf("index %d: expected %d, got %d", i, p, result[i])
		}
	}
}

func TestFilterPortsEmptyIgnoreList(t *testing.T) {
	f := New([]int{})
	input := []int{22, 80, 443}
	result := f.FilterPorts(input)
	if len(result) != len(input) {
		t.Fatalf("expected all ports to pass, got %v", result)
	}
}

func TestFilterPortsAllIgnored(t *testing.T) {
	f := New([]int{22, 80})
	result := f.FilterPorts([]int{22, 80})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}
