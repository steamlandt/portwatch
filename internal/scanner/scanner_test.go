package scanner

import (
	"net"
	"testing"
	"time"
)

func startTestListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return port, func() { ln.Close() }
}

func TestScanDetectsOpenPort(t *testing.T) {
	port, stop := startTestListener(t)
	defer stop()

	s := New("127.0.0.1", 500*time.Millisecond)
	ports, err := s.Scan(port, port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 1 || ports[0].Port != port {
		t.Errorf("expected port %d to be open, got %v", port, ports)
	}
}

func TestScanInvalidRange(t *testing.T) {
	s := New("127.0.0.1", 100*time.Millisecond)
	_, err := s.Scan(500, 100)
	if err == nil {
		t.Error("expected error for invalid range, got nil")
	}
}

func TestPortString(t *testing.T) {
	p := Port{Protocol: "tcp", Address: "127.0.0.1", Port: 8080}
	expected := "127.0.0.1:8080/tcp"
	if p.String() != expected {
		t.Errorf("expected %q, got %q", expected, p.String())
	}
}
