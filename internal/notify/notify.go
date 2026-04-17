package notify

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Method represents a notification delivery method.
type Method string

const (
	MethodStdout  Method = "stdout"
	MethodWebhook Method = "webhook"
	MethodExec    Method = "exec"
)

// Config holds configuration for a notifier.
type Config struct {
	Method  Method
	Target  string // URL for webhook, command for exec
	Headers map[string]string
}

// Notifier sends notifications via a configured method.
type Notifier struct {
	cfg Config
	out io.Writer
}

// New creates a new Notifier.
func New(cfg Config) *Notifier {
	return &Notifier{cfg: cfg, out: os.Stdout}
}

// Send delivers a message using the configured method.
func (n *Notifier) Send(subject, body string) error {
	switch n.cfg.Method {
	case MethodStdout:
		return n.sendStdout(subject, body)
	case MethodWebhook:
		return n.sendWebhook(subject, body)
	case MethodExec:
		return n.sendExec(subject, body)
	default:
		return fmt.Errorf("unknown notify method: %s", n.cfg.Method)
	}
}

func (n *Notifier) sendStdout(subject, body string) error {
	_, err := fmt.Fprintf(n.out, "[notify] %s: %s\n", subject, body)
	return err
}

func (n *Notifier) sendWebhook(subject, body string) error {
	payload := fmt.Sprintf(`{"subject":%q,"body":%q}`, subject, body)
	cmd := exec.Command("curl", "-s", "-X", "POST",
		"-H", "Content-Type: application/json",
		"-d", payload,
		n.cfg.Target)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("webhook failed: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func (n *Notifier) sendExec(subject, body string) error {
	cmd := exec.Command(n.cfg.Target, subject, body)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("exec notify failed: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}
