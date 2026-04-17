package config

import (
	"encoding/json"
	"os"
)

// Config holds all portwatch configuration.
type Config struct {
	StartPort    int    `json:"start_port"`
	EndPort      int    `json:"end_port"`
	IntervalSecs int    `json:"interval_secs"`
	StatePath    string `json:"state_path"`
	IgnoredPorts []int  `json:"ignored_ports"`
}

// Default returns a Config with sensible defaults.
func Default() *Config {
	return &Config{
		StartPort:    1,
		EndPort:      1024,
		IntervalSecs: 60,
		StatePath:    "/tmp/portwatch_state.json",
		IgnoredPorts: []int{},
	}
}

// Load reads a JSON config file and merges it over the defaults.
func Load(path string) (*Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrConfigNotFound
		}
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.StartPort < 1 || c.EndPort > 65535 || c.StartPort > c.EndPort {
		return ErrInvalidPortRange
	}
	if c.IntervalSecs < 1 {
		return ErrInvalidInterval
	}
	return nil
}
