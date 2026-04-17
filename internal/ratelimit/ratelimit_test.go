package ratelimit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
)

func TestAllowFirstCallPermits(t *testing.T) {
	l := ratelimit.New(1 * time.Second)
	if !l.Allow("port:8080:opened") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllowSecondCallWithinCooldownBlocks(t *testing.T) {
	l := ratelimit.New(1 * time.Second)
	l.Allow("port:8080:opened")
	if l.Allow("port:8080:opened") {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestAllowAfterCooldownPermits(t *testing.T) {
	l := ratelimit.New(10 * time.Millisecond)
	l.Allow("port:9090:closed")
	time.Sleep(20 * time.Millisecond)
	if !l.Allow("port:9090:closed") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestResetClearsKey(t *testing.T) {
	l := ratelimit.New(1 * time.Second)
	l.Allow("port:443:opened")
	l.Reset("port:443:opened")
	if !l.Allow("port:443:opened") {
		t.Fatal("expected allow after)
	}
}

func TestResetAllClearsAllKeys(t *testing.T) {
	l := ratelimit.New(1 * time.Second)
	l.Allow("port:80:opened")
	l.Allow("port:443:opened")
	l.ResetAll()
	if !l.Allow("port:80:opened") {
		t.Fatal("expected port:80 to be allowed after ResetAll")
	}
	if !l.Allow("port:443:opened") {
		t.Fatal("expected port:443 to be allowed after ResetAll")
	}
}

func TestDifferentKeysAreIndependent(t *testing.T) {
	l := ratelimit.New(1 * time.Second)
	l.Allow("port:22:opened")
	if !l.Allow("port:23:opened") {
		t.Fatal("expected different key to be allowed independently")
	}
}
