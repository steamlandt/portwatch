package schedule

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestSchedulerCallsFn(t *testing.T) {
	var count int32
	s := New(20 * time.Millisecond)

	go s.Start(func() {
		atomic.AddInt32(&count, 1)
	})

	time.Sleep(75 * time.Millisecond)
	s.Stop()

	got := atomic.LoadInt32(&count)
	if got < 2 {
		t.Errorf("expected at least 2 calls, got %d", got)
	}
}

func TestSchedulerStops(t *testing.T) {
	var count int32
	s := New(20 * time.Millisecond)

	go s.Start(func() {
		atomic.AddInt32(&count, 1)
	})

	time.Sleep(50 * time.Millisecond)
	s.Stop()
	time.Sleep(10 * time.Millisecond)
	before := atomic.LoadInt32(&count)
	time.Sleep(50 * time.Millisecond)
	after := atomic.LoadInt32(&count)

	if before != after {
		t.Errorf("scheduler continued after Stop: before=%d after=%d", before, after)
	}
}

func TestInterval(t *testing.T) {
	d := 5 * time.Second
	s := New(d)
	if s.Interval() != d {
		t.Errorf("expected %v, got %v", d, s.Interval())
	}
}
