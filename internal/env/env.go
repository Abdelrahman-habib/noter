// Package env reads settings from environment variables.
package env

import (
	"os"
	"strconv"
)

// Default uses the environment variable as the flag default
func Default(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Bool parses a boolean environment variable for flag defaults
// (strconv.ParseBool accepts 1, t, T, TRUE, true, True.
// Unrecognized values are false).
func Bool(key string) bool {
	value, _ := strconv.ParseBool(Default(key, "false"))
	return value
}
