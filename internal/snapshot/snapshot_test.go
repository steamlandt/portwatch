package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func tempPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "snapshot.json")
}

func makePort(number int, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func TestSaveAndLoad(t *testing.T) {
	store := snapshot.New(tempPath(t))
	ports := []scanner.Port{makePort(80, "tcp"), makePort(443, "tcp")}

	if err := store.Save(ports); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(snap.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(snap.Ports))
	}
	if snap.CapturedAt.IsZero() {
		t.Error("expected non-zero CapturedAt")
	}
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	store := snapshot.New(tempPath(t))
	snap, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap.Ports) != 0 {
		t.Errorf("expected empty ports, got %d", len(snap.Ports))
	}
}

func TestLoadInvalidJSONReturnsError(t *testing.T) {
	path := tempPath(t)
	_ = os.WriteFile(path, []byte("not json"), 0644)
	store := snapshot.New(path)
	_, err := store.Load()
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestDiffDetectsAddedAndRemoved(t *testing.T) {
	prev := snapshot.Snapshot{Ports: []scanner.Port{makePort(22, "tcp"), makePort(80, "tcp")}}
	next := snapshot.Snapshot{Ports: []scanner.Port{makePort(80, "tcp"), makePort(443, "tcp")}}

	added, removed := snapshot.Diff(prev, next)

	if len(added) != 1 || added[0].Number != 443 {
		t.Errorf("expected port 443 added, got %v", added)
	}
	if len(removed) != 1 || removed[0].Number != 22 {
		t.Errorf("expected port 22 removed, got %v", removed)
	}
}

func TestDiffNoDifference(t *testing.T) {
	ports := []scanner.Port{makePort(80, "tcp")}
	prev := snapshot.Snapshot{Ports: ports}
	next := snapshot.Snapshot{Ports: ports}

	added, removed := snapshot.Diff(prev, next)
	if len(added) != 0 || len(removed) != 0 {
		t.Errorf("expected no diff, got added=%v removed=%v", added, removed)
	}
}

func TestDiffProtocolDistinct(t *testing.T) {
	// Ports with the same number but different protocols should be treated as distinct.
	prev := snapshot.Snapshot{Ports: []scanner.Port{makePort(80, "tcp")}}
	next := snapshot.Snapshot{Ports: []scanner.Port{makePort(80, "udp")}}

	added, removed := snapshot.Diff(prev, next)

	if len(added) != 1 || added[0].Protocol != "udp" {
		t.Errorf("expected udp/80 added, got %v", added)
	}
	if len(removed) != 1 || removed[0].Protocol != "tcp" {
		t.Errorf("expected tcp/80 removed, got %v", removed)
	}
}
