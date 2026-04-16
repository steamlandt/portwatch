package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto string, number int) scanner.Port {
	return scanner.Port{Proto: proto, Number: number}
}

func TestPortOpenedWritesAlertLevel(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)
	n.PortOpened(makePort("tcp", 8080))

	out := buf.String()
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", out)
	}
}

func TestPortClosedWritesWarnLevel(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)
	n.PortClosed(makePort("tcp", 443))

	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN in output, got: %s", out)
	}
	if !strings.Contains(out, "443") {
		t.Errorf("expected port 443 in output, got: %s", out)
	}
}

func TestNewDefaultsToStdout(t *testing.T) {
	// Should not panic when nil writer is passed.
	n := alert.New(nil)
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}

func TestOutputContainsTimestamp(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)
	n.PortOpened(makePort("udp", 53))

	out := buf.String()
	// RFC3339 timestamps contain 'T' between date and time.
	if !strings.Contains(out, "T") {
		t.Errorf("expected RFC3339 timestamp in output, got: %s", out)
	}
}
