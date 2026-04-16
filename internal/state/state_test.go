package state

import (
	"os"
	"testing"
	"time"
)

func tempPath(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-state-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	os.Remove(f.Name()) // let Store create it fresh
	return f.Name()
}

func TestSetAndGet(t *testing.T) {
	s, err := New(tempPath(t))
	if err != nil {
		t.Fatal(err)
	}
	ps := PortState{Port: 8080, Open: true, LastSeen: time.Now()}
	if err := s.Set(ps); err != nil {
		t.Fatal(err)
	}
	got, ok := s.Get(8080)
	if !ok {
		t.Fatal("expected port 8080 to be present")
	}
	if got.Open != true {
		t.Errorf("expected Open=true, got %v", got.Open)
	}
}

func TestPersistence(t *testing.T) {
	path := tempPath(t)
	s1, err := New(path)
	if err != nil {
		t.Fatal(err)
	}
	ps := PortState{Port: 443, Open: true, LastSeen: time.Now()}
	if err := s1.Set(ps); err != nil {
		t.Fatal(err)
	}

	s2, err := New(path)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := s2.Get(443)
	if !ok {
		t.Fatal("expected port 443 after reload")
	}
	if got.Port != 443 {
		t.Errorf("expected port 443, got %d", got.Port)
	}
}

func TestGetMissing(t *testing.T) {
	s, err := New(tempPath(t))
	if err != nil {
		t.Fatal(err)
	}
	_, ok := s.Get(9999)
	if ok {
		t.Error("expected missing port to return ok=false")
	}
}

func TestAll(t *testing.T) {
	s, err := New(tempPath(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range []int{80, 443, 8080} {
		s.Set(PortState{Port: p, Open: true, LastSeen: time.Now()})
	}
	all := s.All()
	if len(all) != 3 {
		t.Errorf("expected 3 ports, got %d", len(all))
	}
}
