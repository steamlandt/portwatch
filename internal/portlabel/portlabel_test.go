package portlabel_test

import (
	"testing"

	"github.com/user/portwatch/internal/portlabel"
)

func TestLabelKnownPort(t *testing.T) {
	l := portlabel.New(nil)
	if got := l.Label(22); got != "ssh" {
		t.Fatalf("expected ssh, got %q", got)
	}
}

func TestLabelUnknownPortReturnsEmpty(t *testing.T) {
	l := portlabel.New(nil)
	if got := l.Label(9999); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestLabelOrPortFallsBackToNumber(t *testing.T) {
	l := portlabel.New(nil)
	if got := l.LabelOrPort(9999); got != "9999" {
		t.Fatalf("expected \"9999\", got %q", got)
	}
}

func TestLabelOrPortReturnsName(t *testing.T) {
	l := portlabel.New(nil)
	if got := l.LabelOrPort(443); got != "https" {
		t.Fatalf("expected https, got %q", got)
	}
}

func TestCustomOverridesDefault(t *testing.T) {
	custom := map[int]string{80: "my-http"}
	l := portlabel.New(custom)
	if got := l.Label(80); got != "my-http" {
		t.Fatalf("expected my-http, got %q", got)
	}
}

func TestCustomAddsNewEntry(t *testing.T) {
	custom := map[int]string{12345: "my-service"}
	l := portlabel.New(custom)
	if got := l.Label(12345); got != "my-service" {
		t.Fatalf("expected my-service, got %q", got)
	}
	// Built-in defaults still present.
	if got := l.Label(22); got != "ssh" {
		t.Fatalf("expected ssh, got %q", got)
	}
}
