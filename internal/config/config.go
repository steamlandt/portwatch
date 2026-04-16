package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the runtime configuration for portwatch.
type Config struct {
	PortRange   PortRange     `json:"port_range"`
	Interval    time.Duration `json:"interval"`
	StatePath   string        `json:"state_path"`
	AlertOutput string        `json:"alert_output"`
}

// PortRange defines the inclusive start/end ports to scan.
type PortRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// Default returns a Config with sensible defaults.
func Default() *Config {
	return &Config{
		PortRange:   PortRange{Start: 1, End: 65535},
		Interval:    30 * time.Second,
		StatePath:   "/var/lib/portwatch/state.json",
		AlertOutput: "",
	}
}

// Load reads a JSON config file from path, overlaying values onto defaults.
func Load(path string) (*Config, error) {
	cfg := Default()
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate returns an error if the config contains invalid values.
func (c *Config) Validate() error {
	if c.PortRange.Start < 1 || c.PortRange.End > 65535 || c.PortRange.Start > c.PortRange.End {
		return ErrInvalidPortRange
	}
	if c.Interval < time.Second {
		return ErrIntervalTooShort
	}
	return nil
}
