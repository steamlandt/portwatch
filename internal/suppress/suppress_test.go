package suppress_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/suppress"
)

func TestFirstCallIsNotSuppressed(t *testing.T) {
	s := suppress.New(5 * time.Second)
	if s.IsSuppressed("port:8080:open") {
		t.Fatal("expected first call to not be suppressed")
	}
}

func TestSecondCallWithinWindowIsSuppressed(t *testing.T) {
	s := suppress.New(5 * time.Second)
	s.IsSuppressed("port:8080:open")
	if !s.IsSuppressed("port:8080:open") {
		t.Fatal("expected second call within window to be suppressed")
	}
}

func TestCallAfterWindowIsNotSuppressed(t *testing.T) {
	s := suppress.New(10 * time.Millisecond)
	s.IsSuppressed("k")
	time.Sleep(20 * time.Millisecond)
	if s.IsSuppressed("k") {
		t.Fatal("expected call after window to not be suppressed")
	}
}

func TestDifferentKeysAreIndependent(t *testing.T) {
	s := suppress.New(5 * time.Second)
	s.IsSuppressed("port:80:open")
	if s.IsSuppressed("port:443:open") {
		t.Fatal("expected different key to not be suppressed")
	}
}

func TestResetClearsKey(t *testing.T) {
	s := suppress.New(5 * time.Second)
	s.IsSuppressed("k")
	s.Reset("k")
	if s.IsSuppressed("k") {
		t.Fatal("expected key to be cleared after Reset")
	}
}

func TestResetAllClearsAllKeys(t *testing.T) {
	s := suppress.New(5 * time.Second)
	s.IsSuppressed("a")
	s.IsSuppressed("b")
	s.ResetAll()
	if s.IsSuppressed("a") || s.IsSuppressed("b") {
		t.Fatal("expected all keys cleared after ResetAll")
	}
}

func TestResetUnknownKeyIsNoop(t *testing.T) {
	// Resetting a key that was never seen should not panic or affect other keys.
	s := suppress.New(5 * time.Second)
	s.IsSuppressed("existing")
	s.Reset("nonexistent") // should be a no-op
	if !s.IsSuppressed("existing") {
		t.Fatal("expected existing key to still be suppressed after resetting an unknown key")
	}
}
