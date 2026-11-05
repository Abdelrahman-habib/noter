// Package env reads settings from environment variables.
package env

import (
	"os"
	"strconv"
)

// Default returns the value of the environment variable if set, otherwise fallback
func Default(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Bool parses the environment variable with strconv.ParseBool
// (1, t, T, TRUE, true, True); unset or unrecognized values return fallback
func Bool(key string, fallback bool) bool {
	value, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return value
}
