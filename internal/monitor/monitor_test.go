package monitor

import (
	"bytes"
	"net"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func tempState(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-mon-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	os.Remove(f.Name())
	return f.Name()
}

func startListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func TestMonitorDetectsOpenPort(t *testing.T) {
	port, stop := startListener(t)
	defer stop()

	var buf bytes.Buffer
	alerter := alert.New(&buf)

	m, err := New(Config{
		StartPort: port,
		EndPort:   port,
		Interval:  50 * time.Millisecond,
		StatePath: tempState(t),
	}, alerter)
	if err != nil {
		t.Fatal(err)
	}

	m.Start()
	time.Sleep(150 * time.Millisecond)
	m.Stop()

	if buf.Len() == 0 {
		t.Error("expected alert output, got none")
	}
}

func TestMonitorDetectsClosedPort(t *testing.T) {
	port, stop := startListener(t)

	var buf bytes.Buffer
	alerter := alert.New(&buf)

	m, err := New(Config{
		StartPort: port,
		EndPort:   port,
		Interval:  50 * time.Millisecond,
		StatePath: tempState(t),
	}, alerter)
	if err != nil {
		t.Fatal(err)
	}

	m.Start()
	time.Sleep(100 * time.Millisecond)
	stop() // close the port
	time.Sleep(150 * time.Millisecond)
	m.Stop()

	out := buf.String()
	if len(out) == 0 {
		t.Error("expected alert output after port closed")
	}
}
