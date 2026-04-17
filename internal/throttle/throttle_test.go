package throttle_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

func TestFirstCallIsAllowed(t *testing.T) {
	th := throttle.New(100 * time.Millisecond)
	if !th.Allow("tcp:8080") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestSecondCallWithinIntervalIsBlocked(t *testing.T) {
	th := throttle.New(200 * time.Millisecond)
	th.Allow("tcp:8080")
	if th.Allow("tcp:8080") {
		t.Fatal("expected second call within interval to be blocked")
	}
}

func TestCallAfterIntervalIsAllowed(t *testing.T) {
	th := throttle.New(20 * time.Millisecond)
	th.Allow("tcp:8080")
	time.Sleep(30 * time.Millisecond)
	if !th.Allow("tcp:8080") {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestDifferentKeysAreIndependent(t *testing.T) {
	th := throttle.New(200 * time.Millisecond)
	th.Allow("tcp:8080")
	if !th.Allow("tcp:9090") {
		t.Fatal("expected different key to be allowed")
	}
}

func TestResetClearsKey(t *testing.T) {
	th := throttle.New(200 * time.Millisecond)
	th.Allow("tcp:8080")
	th.Reset("tcp:8080")
	if !th.Allow("tcp:8080") {
		t.Fatal("expected allow after reset")
	}
}

func TestResetAllClearsAllKeys(t *testing.T) {
	th := throttle.New(200 * time.Millisecond)
	th.Allow("tcp:8080")
	th.Allow("tcp:9090")
	th.ResetAll()
	if !th.Allow("tcp:8080") || !th.Allow("tcp:9090") {
		t.Fatal("expected all keys cleared after ResetAll")
	}
}

func TestInterval(t *testing.T) {
	d := 500 * time.Millisecond
	th := throttle.New(d)
	if th.Interval() != d {
		t.Fatalf("expected interval %v, got %v", d, th.Interval())
	}
}
