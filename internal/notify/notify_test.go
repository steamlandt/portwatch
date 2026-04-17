package notify

import (
	"bytes"
	"strings"
	"testing"
)

func TestSendStdoutWritesMessage(t *testing.T) {
	n := New(Config{Method: MethodStdout})
	var buf bytes.Buffer
	n.out = &buf

	if err := n.Send("port opened", "tcp/8080"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "port opened") {
		t.Errorf("expected subject in output, got: %s", got)
	}
	if !strings.Contains(got, "tcp/8080") {
		t.Errorf("expected body in output, got: %s", got)
	}
}

func TestSendUnknownMethodReturnsError(t *testing.T) {
	n := New(Config{Method: "sms"})
	err := n.Send("x", "y")
	if err == nil {
		t.Fatal("expected error for unknown method")
	}
	if !strings.Contains(err.Error(), "unknown notify method") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestSendExecRunsCommand(t *testing.T) {
	n := New(Config{Method: MethodExec, Target: "echo"})
	if err := n.Send("subject", "body"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendExecBadCommandReturnsError(t *testing.T) {
	n := New(Config{Method: MethodExec, Target: "/no/such/binary"})
	err := n.Send("subject", "body")
	if err == nil {
		t.Fatal("expected error for missing binary")
	}
}

func TestNewDefaultsMethod(t *testing.T) {
	n := New(Config{Method: MethodStdout})
	if n.cfg.Method != MethodStdout {
		t.Errorf("expected stdout method, got %s", n.cfg.Method)
	}
}
