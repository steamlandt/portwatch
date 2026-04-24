package portlabel

import "testing"

func TestLabelKnownPort(t *testing.T) {
	l := New(nil)
	if got := l.Label(22); got != "ssh" {
		t.Errorf("expected ssh, got %q", got)
	}
}

func TestLabelUnknownPortReturnsEmpty(t *testing.T) {
	l := New(nil)
	if got := l.Label(9999); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestLabelOrPortFallsBackToNumber(t *testing.T) {
	l := New(nil)
	if got := l.LabelOrPort(9999); got != "9999" {
		t.Errorf("expected \"9999\", got %q", got)
	}
}

func TestLabelOrPortReturnsName(t *testing.T) {
	l := New(nil)
	if got := l.LabelOrPort(443); got != "https" {
		t.Errorf("expected https, got %q", got)
	}
}

func TestCustomOverridesDefault(t *testing.T) {
	l := New(map[int]string{80: "my-app"})
	if got := l.Label(80); got != "my-app" {
		t.Errorf("expected my-app, got %q", got)
	}
}

func TestCustomPortNotInDefaults(t *testing.T) {
	l := New(map[int]string{12345: "custom-svc"})
	if got := l.Label(12345); got != "custom-svc" {
		t.Errorf("expected custom-svc, got %q", got)
	}
	if got := l.Label(22); got != "ssh" {
		t.Errorf("built-in should still work, got %q", got)
	}
}
