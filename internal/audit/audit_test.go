package audit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/audit"
)

func tempPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "audit.jsonl")
}

func TestRecordAndAll(t *testing.T) {
	path := tempPath(t)
	log := audit.New(path)

	if err := log.Record("opened", "tcp", 80, "http"); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := log.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Event != "opened" || e.Port != 80 || e.Protocol != "tcp" || e.Label != "http" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestAllMissingFileReturnsEmpty(t *testing.T) {
	path := tempPath(t)
	log := audit.New(path)

	entries, err := log.All()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestMultipleRecordsPreserveOrder(t *testing.T) {
	path := tempPath(t)
	log := audit.New(path)

	events := []struct {
		event string
		port  uint16
	}{
		{"opened", 22},
		{"opened", 443},
		{"closed", 22},
	}
	for _, ev := range events {
		if err := log.Record(ev.event, "tcp", ev.port, ""); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}

	entries, err := log.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Port != 22 || entries[1].Port != 443 || entries[2].Port != 22 {
		t.Errorf("order mismatch: %+v", entries)
	}
}

func TestTimestampIsSet(t *testing.T) {
	path := tempPath(t)
	log := audit.New(path)

	before := time.Now().UTC()
	_ = log.Record("opened", "tcp", 8080, "")
	after := time.Now().UTC()

	entries, _ := log.All()
	if len(entries) == 0 {
		t.Fatal("no entries")
	}
	ts := entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v out of expected range [%v, %v]", ts, before, after)
	}
}

func TestRecordInvalidPathReturnsError(t *testing.T) {
	log := audit.New("/nonexistent-dir/audit.jsonl")
	err := log.Record("opened", "tcp", 80, "")
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
	_ = os.Remove("/nonexistent-dir/audit.jsonl")
}
