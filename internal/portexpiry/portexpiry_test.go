package portexpiry_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portexpiry"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(n int) scanner.Port {
	return scanner.Port{Number: n, Proto: "tcp"}
}

func TestObserveRecordsFirstSeen(t *testing.T) {
	e := portexpiry.New(time.Minute)
	p := makePort(80)
	e.Observe(p)
	age, ok := e.Age(p)
	if !ok {
		t.Fatal("expected port to be tracked")
	}
	if age < 0 {
		t.Fatalf("unexpected negative age: %v", age)
	}
}

func TestObserveIsIdempotent(t *testing.T) {
	now := time.Now()
	calls := 0
	e := portexpiry.New(time.Minute)
	// inject a clock that advances on each call
	_ = now
	p := makePort(443)
	e.Observe(p)
	age1, _ := e.Age(p)
	e.Observe(p) // second observe must not reset the timestamp
	age2, _ := e.Age(p)
	if age2 < age1 {
		t.Errorf("second Observe reset the first-seen timestamp (calls=%d)", calls)
	}
}

func TestRemoveClearsEntry(t *testing.T) {
	e := portexpiry.New(time.Minute)
	p := makePort(22)
	e.Observe(p)
	e.Remove(p)
	_, ok := e.Age(p)
	if ok {
		t.Fatal("expected port to be removed")
	}
}

func TestExpiredReturnsLongLivedPorts(t *testing.T) {
	maxAge := 10 * time.Millisecond
	e := portexpiry.New(maxAge)
	p80 := makePort(80)
	p443 := makePort(443)
	e.Observe(p80)
	e.Observe(p443)
	time.Sleep(20 * time.Millisecond)
	expired := e.Expired([]scanner.Port{p80, p443})
	if len(expired) != 2 {
		t.Fatalf("expected 2 expired ports, got %d", len(expired))
	}
}

func TestExpiredExcludesFreshPorts(t *testing.T) {
	e := portexpiry.New(time.Hour)
	p := makePort(8080)
	e.Observe(p)
	expired := e.Expired([]scanner.Port{p})
	if len(expired) != 0 {
		t.Fatalf("expected no expired ports, got %d", len(expired))
	}
}

func TestExpiredEntryHasDuration(t *testing.T) {
	maxAge := 5 * time.Millisecond
	e := portexpiry.New(maxAge)
	p := makePort(3306)
	e.Observe(p)
	time.Sleep(10 * time.Millisecond)
	expired := e.Expired([]scanner.Port{p})
	if len(expired) == 0 {
		t.Fatal("expected expired entry")
	}
	if expired[0].Duration < maxAge {
		t.Errorf("duration %v less than maxAge %v", expired[0].Duration, maxAge)
	}
}

func TestAgeUnknownPortReturnsFalse(t *testing.T) {
	e := portexpiry.New(time.Minute)
	_, ok := e.Age(makePort(9999))
	if ok {
		t.Fatal("expected false for untracked port")
	}
}
