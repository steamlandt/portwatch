package watchdog

import (
	"testing"
	"time"
)

func TestStatusUnknownBeforeFirstBeat(t *testing.T) {
	w := New(time.Second)
	if s := w.Status(); s != StatusUnknown {
		t.Fatalf("expected unknown, got %s", s)
	}
}

func TestStatusHealthyAfterBeat(t *testing.T) {
	w := New(time.Second)
	w.Beat()
	if s := w.Status(); s != StatusHealthy {
		t.Fatalf("expected healthy, got %s", s)
	}
}

func TestStatusStaleAfterTimeout(t *testing.T) {
	w := New(10 * time.Millisecond)
	w.Beat()
	time.Sleep(20 * time.Millisecond)
	if s := w.Status(); s != StatusStale {
		t.Fatalf("expected stale, got %s", s)
	}
}

func TestBeatResetsStale(t *testing.T) {
	w := New(10 * time.Millisecond)
	w.Beat()
	time.Sleep(20 * time.Millisecond)
	w.Beat()
	if s := w.Status(); s != StatusHealthy {
		t.Fatalf("expected healthy after re-beat, got %s", s)
	}
}

func TestResetReturnsUnknown(t *testing.T) {
	w := New(time.Second)
	w.Beat()
	w.Reset()
	if s := w.Status(); s != StatusUnknown {
		t.Fatalf("expected unknown after reset, got %s", s)
	}
}

func TestLastBeatUpdates(t *testing.T) {
	w := New(time.Second)
	before := time.Now()
	w.Beat()
	after := time.Now()
	lb := w.LastBeat()
	if lb.Before(before) || lb.After(after) {
		t.Fatalf("last beat %v not in expected range", lb)
	}
}
