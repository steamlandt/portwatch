package baseline

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func makePort(proto string, number int) scanner.Port {
	return scanner.Port{Protocol: proto, Number: number}
}

func TestCaptureAndLoad(t *testing.T) {
	b := New(tempPath(t))
	ports := []scanner.Port{makePort("tcp", 80), makePort("tcp", 443)}
	if err := b.Capture(ports); err != nil {
		t.Fatalf("capture: %v", err)
	}
	b2 := New(b.path)
	if err := b2.Load(); err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(b2.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(b2.Ports))
	}
	if b2.CapturedAt.IsZero() {
		t.Error("expected non-zero CapturedAt")
	}
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	b := New(filepath.Join(t.TempDir(), "missing.json"))
	if err := b.Load(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(b.Ports) != 0 {
		t.Errorf("expected empty ports")
	}
}

func TestDeviationsDetectsAdded(t *testing.T) {
	b := New(tempPath(t))
	_ = b.Capture([]scanner.Port{makePort("tcp", 80)})
	added, removed := b.Deviations([]scanner.Port{makePort("tcp", 80), makePort("tcp", 8080)})
	if len(added) != 1 || added[0].Number != 8080 {
		t.Errorf("expected port 8080 added, got %v", added)
	}
	if len(removed) != 0 {
		t.Errorf("expected no removed ports, got %v", removed)
	}
}

func TestDeviationsDetectsRemoved(t *testing.T) {
	b := New(tempPath(t))
	_ = b.Capture([]scanner.Port{makePort("tcp", 80), makePort("tcp", 443)})
	added, removed := b.Deviations([]scanner.Port{makePort("tcp", 80)})
	if len(removed) != 1 || removed[0].Number != 443 {
		t.Errorf("expected port 443 removed, got %v", removed)
	}
	if len(added) != 0 {
		t.Errorf("expected no added ports")
	}
}

func TestLoadInvalidJSONReturnsError(t *testing.T) {
	p := tempPath(t)
	_ = os.WriteFile(p, []byte("not json"), 0644)
	b := New(p)
	if err := b.Load(); err == nil {
		t.Error("expected error for invalid JSON")
	}
}
