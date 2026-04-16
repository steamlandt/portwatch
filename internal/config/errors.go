package config

import "errors"

// Sentinel errors for config validation.
var (
	ErrInvalidPortRange = errors.New("config: port range must be between 1-65535 with start <= end")
	ErrIntervalTooShort = errors.New("config: interval must be at least 1 second")
)
