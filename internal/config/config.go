package config

import (
	"encoding/json"
	"os"

	"github.com/user/portwatch/internal/notify"
)

// Config holds all portwatch runtime configuration.
type Config struct {
	PortRange  [2]int          `json:"port_range"`
	Interval   int             `json:"interval_seconds"`
	StatePath  string          `json:"state_path"`
	HistoryPath string         `json:"history_path"`
	IgnorePorts []int          `json:"ignore_ports"`
	Notify     NotifyConfig    `json:"notify"`
}

// NotifyConfig mirrors notify.Config for JSON unmarshalling.
type NotifyConfig struct {
	Method  notify.Method     `json:"method"`
	Target  string            `json:"target"`
	Headers map[string]string `json:"headers"`
}

// Default returns a Config populated with sensible defaults.
func Default() Config {
	return Config{
		PortRange:   [2]int{1, 65535},
		Interval:    60,
		StatePath:   "/var/lib/portwatch/state.json",
		HistoryPath: "/var/lib/portwatch/history.json",
		Notify: NotifyConfig{
			Method: notify.MethodStdout,
		},
	}
}

// Load reads a JSON config file, overlaying values onto defaults.
func Load(path string) (Config, error) {
	cfg := Default()
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, ErrNotFound
		}
		return cfg, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return cfg, err
	}
	return cfg, validate(cfg)
}

func validate(cfg Config) error {
	lo, hi := cfg.PortRange[0], cfg.PortRange[1]
	if lo < 1 || hi > 65535 || lo > hi {
		return ErrInvalidRange
	}
	if cfg.Interval < 1 {
		return ErrInvalidInterval
	}
	return nil
}
