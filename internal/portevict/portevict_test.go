package portevict_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portevict"
	"github.com/user/portwatch/internal/scanner"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "evict.json")
}

func makePort(n int) scanner.Port {
	return scanner.Port{Number: n, Protocol: "tcp"}
}

func TestRecordAndAll(t *testing.T) {
	tr, err := portevict.New("")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	now := time.Now()
	if err := tr.Record(makePort(80), now.Add(-5*time.Minute), now); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries := tr.All()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Port.Number != 80 {
		t.Errorf("expected port 80, got %d", entries[0].Port.Number)
	}
}

func TestDurationIsSet(t *testing.T) {
	tr, _ := portevict.New("")
	opened := time.Now().Add(-10 * time.Minute)
	closed := time.Now()
	_ = tr.Record(makePort(443), opened, closed)
	entry := tr.All()[0]
	if entry.Duration == "" {
		t.Error("expected non-empty duration")
	}
}

func TestPersistence(t *testing.T) {
	p := tempPath(t)
	tr, _ := portevict.New(p)
	now := time.Now()
	_ = tr.Record(makePort(22), now.Add(-1*time.Hour), now)

	tr2, err := portevict.New(p)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(tr2.All()) != 1 {
		t.Errorf("expected 1 persisted entry, got %d", len(tr2.All()))
	}
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	p := filepath.Join(t.TempDir(), "missing.json")
	tr, err := portevict.New(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tr.All()) != 0 {
		t.Error("expected empty entries for missing file")
	}
}

func TestResetClearsEntries(t *testing.T) {
	p := tempPath(t)
	tr, _ := portevict.New(p)
	now := time.Now()
	_ = tr.Record(makePort(8080), now.Add(-2*time.Minute), now)
	if err := tr.Reset(); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	if len(tr.All()) != 0 {
		t.Error("expected empty after reset")
	}
	data, _ := os.ReadFile(p)
	if string(data) != "null" && string(data) != "[]" {
		t.Errorf("unexpected file content after reset: %s", data)
	}
}

func TestMultipleRecordsPreserveOrder(t *testing.T) {
	tr, _ := portevict.New("")
	now := time.Now()
	ports := []int{80, 443, 8080}
	for _, n := range ports {
		_ = tr.Record(makePort(n), now.Add(-1*time.Minute), now)
	}
	entries := tr.All()
	for i, n := range ports {
		if entries[i].Port.Number != n {
			t.Errorf("index %d: expected port %d, got %d", i, n, entries[i].Port.Number)
		}
	}
}
