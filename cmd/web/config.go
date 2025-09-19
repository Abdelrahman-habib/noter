package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

const (
	envDevelopment = "development"
	envProduction  = "production"
	envTest        = "test"
)

type config struct {
	// app
	env       string
	debugMode bool

	// server
	addr string

	// tls
	tlsCert string
	tlsKey  string

	// db
	dsn string
}

func parseFlags() *config {
	addr := flag.String("addr", ":4000", "HTTP network address")
	debugMode := flag.Bool("debug", false, "enable debug mode")
	env := flag.String("env", "development", "Environment (development, production, test)")

	tlsCert := flag.String("tls-cert", "./tls/cert.pem", "Path to TLS certificate file")
	tlsKey := flag.String("tls-key", "./tls/key.pem", "Path to TLS key file")

	dsn := flag.String("dsn", "noter_web:pass@/noter?parseTime=true", "MySQL data source name")

	flag.Parse()

	// validate env
	if *env != envDevelopment && *env != envProduction && *env != envTest {
		log.Fatal("invalid environment")
		os.Exit(1)
	}

	// Get DSN from environment variable if available, otherwise use flag
	dsnValue := *dsn
	if envDSN := os.Getenv("DB_DSN"); envDSN != "" {
		dsnValue = envDSN
	}

	return &config{
		addr:      *addr,
		debugMode: *debugMode,
		env:       *env,

		tlsCert: *tlsCert,
		tlsKey:  *tlsKey,

		dsn: dsnValue,
	}
}

// buildDSN safely adds required parameters to the DSN
func buildDSN(dsn string) (string, error) {
	// Check if DSN already has query parameters
	separator := "?"
	if strings.Contains(dsn, "?") {
		separator = "&"
	}

	// Add parseTime parameter if not already present
	if !strings.Contains(dsn, "parseTime=") {
		dsn += separator + "parseTime=true"
	}

	return dsn, nil
}
