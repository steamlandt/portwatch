package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestRecordAndAll(t *testing.T) {
	h, err := New("")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := h.Record(8080, "tcp", "opened"); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries := h.All()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Port != 8080 || entries[0].Event != "opened" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestPersistence(t *testing.T) {
	p := tempPath(t)
	h, _ := New(p)
	_ = h.Record(443, "tcp", "closed")

	h2, err := New(p)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	entries := h2.All()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after reload, got %d", len(entries))
	}
	if entries[0].Port != 443 {
		t.Errorf("wrong port: %d", entries[0].Port)
	}
}

func TestInvalidJSONReturnsError(t *testing.T) {
	p := tempPath(t)
	_ = os.WriteFile(p, []byte("not json"), 0644)
	_, err := New(p)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestTimestampIsSet(t *testing.T) {
	h, _ := New("")
	_ = h.Record(22, "tcp", "opened")
	if h.All()[0].Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestFlushWritesValidJSON(t *testing.T) {
	p := tempPath(t)
	h, _ := New(p)
	_ = h.Record(80, "tcp", "opened")
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("invalid JSON on disk: %v", err)
	}
}
