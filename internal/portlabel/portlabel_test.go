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
		t.Fatalf("expected empty, got %q", got)
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
	l := portlabel.New(map[uint16]string{80: "my-app"})
	if got := l.Label(80); got != "my-app" {
		t.Fatalf("expected my-app, got %q", got)
	}
}

func TestCustomAddsNewLabel(t *testing.T) {
	l := portlabel.New(map[uint16]string{9200: "elasticsearch"})
	if got := l.Label(9200); got != "elasticsearch" {
		t.Fatalf("expected elasticsearch, got %q", got)
	}
}
