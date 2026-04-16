package config_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func writeTempConfig(t *testing.T, v any) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-cfg-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(v); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestDefaultValues(t *testing.T) {
	cfg := config.Default()
	if cfg.PortRange.Start != 1 || cfg.PortRange.End != 65535 {
		t.Errorf("unexpected default port range: %+v", cfg.PortRange)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("unexpected default interval: %v", cfg.Interval)
	}
}

func TestLoadOverridesDefaults(t *testing.T) {
	path := writeTempConfig(t, map[string]any{
		"port_range": map[string]int{"start": 1024, "end": 9000interval":   int(10 * time.Second),
	})
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PortRange.Start != 1024 || cfg.PortRange.End != 9000 {
		t.Errorf("port range not overridden: %+v", cfg.PortRange)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestValidateRejectsInvalidRange(t *testing.T) {
	cfg := config.Default()
	cfg.PortRange.Start = 9000
	cfg.PortRange.End = 1000
	if err := cfg.Validate(); err != config.ErrInvalidPortRange {
		t.Errorf("expected ErrInvalidPortRange, got %v", err)
	}
}

func TestValidateRejectsShortInterval(t *testing.T) {
	cfg := config.Default()
	cfg.Interval = 500 * time.Millisecond
	if err := cfg.Validate(); err != config.ErrIntervalTooShort {
		t.Errorf("expected ErrIntervalTooShort, got %v", err)
	}
}

func TestValidateAcceptsValid(t *testing.T) {
	if err := config.Default().Validate(); err != nil {
		t.Errorf("default config should be valid: %v", err)
	}
}
