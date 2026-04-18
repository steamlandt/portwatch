package healthcheck

import (
	"testing"
)

func TestSetAndGet(t *testing.T) {
	h := New()
	h.Set("scanner", StatusOK, "running")
	c, ok := h.Get("scanner")
	if !ok {
		t.Fatal("expected check to exist")
	}
	if c.Status != StatusOK {
		t.Errorf("expected ok, got %s", c.Status)
	}
	if c.Message != "running" {
		t.Errorf("unexpected message: %s", c.Message)
	}
}

func TestGetMissing(t *testing.T) {
	h := New()
	_, ok := h.Get("missing")
	if ok {
		t.Fatal("expected missing check to return false")
	}
}

func TestOverallReturnsWorst(t *testing.T) {
	h := New()
	h.Set("a", StatusOK, "")
	h.Set("b", StatusDegraded, "slow")
	h.Set("c", StatusOK, "")
	if got := h.Overall(); got != StatusDegraded {
		t.Errorf("expected degraded, got %s", got)
	}
}

func TestOverallDownWinsOverDegraded(t *testing.T) {
	h := New()
	h.Set("a", StatusDegraded, "")
	h.Set("b", StatusDown, "crashed")
	if got := h.Overall(); got != StatusDown {
		t.Errorf("expected down, got %s", got)
	}
}

func TestOverallEmptyIsOK(t *testing.T) {
	h := New()
	if got := h.Overall(); got != StatusOK {
		t.Errorf("expected ok for empty healthcheck, got %s", got)
	}
}

func TestAll(t *testing.T) {
	h := New()
	h.Set("x", StatusOK, "")
	h.Set("y", StatusDown, "")
	all := h.All()
	if len(all) != 2 {
		t.Errorf("expected 2 checks, got %d", len(all))
	}
}

func TestStatusString(t *testing.T) {
	cases := map[Status]string{
		StatusOK:       "ok",
		StatusDegraded: "degraded",
		StatusDown:     "down",
	}
	for s, want := range cases {
		if got := s.String(); got != want {
			t.Errorf("Status(%d).String() = %q, want %q", s, got, want)
		}
	}
}
