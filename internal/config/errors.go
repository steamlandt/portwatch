package config

import "errors"

var (
	ErrNotFound       = errors.New("config file not found")
	ErrInvalidRange   = errors.New("invalid port range")
	ErrInvalidInterval = errors.New("interval must be >= 1")
)
