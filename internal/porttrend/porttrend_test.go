package porttrend_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/porttrend"
)

var t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestObserveCreatesEntry(t *testing.T) {
	tr := porttrend.New()
	tr.Observe("80", t0)

	e, ok := tr.Get("80")
	if !ok {
		t.Fatal("expected entry for port 80")
	}
	if e.Seen != 1 {
		t.Fatalf("expected Seen=1, got %d", e.Seen)
	}
	if !e.FirstSeen.Equal(t0) {
		t.Fatalf("unexpected FirstSeen: %v", e.FirstSeen)
	}
}

func TestObserveIncrementsCount(t *testing.T) {
	tr := porttrend.New()
	t1 := t0.Add(time.Minute)
	tr.Observe("443", t0)
	tr.Observe("443", t1)

	e, _ := tr.Get("443")
	if e.Seen != 2 {
		t.Fatalf("expected Seen=2, got %d", e.Seen)
	}
	if !e.FirstSeen.Equal(t0) {
		t.Fatalf("FirstSeen should not change on second observe")
	}
	if !e.LastSeen.Equal(t1) {
		t.Fatalf("LastSeen should be updated")
	}
}

func TestGetMissingReturnsFalse(t *testing.T) {
	tr := porttrend.New()
	_, ok := tr.Get("9999")
	if ok {
		t.Fatal("expected ok=false for unseen port")
	}
}

func TestAllReturnsAllEntries(t *testing.T) {
	tr := porttrend.New()
	tr.Observe("22", t0)
	tr.Observe("80", t0)
	tr.Observe("443", t0)

	all := tr.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(all))
	}
}

func TestTransientFiltersLowCount(t *testing.T) {
	tr := porttrend.New()
	tr.Observe("22", t0)
	tr.Observe("22", t0)
	tr.Observe("22", t0)
	tr.Observe("8080", t0) // seen once — transient

	transient := tr.Transient(3)
	if len(transient) != 1 {
		t.Fatalf("expected 1 transient port, got %d", len(transient))
	}
	if transient[0].Port != "8080" {
		t.Fatalf("expected port 8080, got %s", transient[0].Port)
	}
}

func TestResetClearsAllEntries(t *testing.T) {
	tr := porttrend.New()
	tr.Observe("80", t0)
	tr.Reset()

	if len(tr.All()) != 0 {
		t.Fatal("expected empty tracker after Reset")
	}
}
