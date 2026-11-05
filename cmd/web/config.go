package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

const (
	envDevelopment = "development"
	envProduction  = "production"
	envTest        = "test"
)

// envDefault uses the environment variable as the flag default
func envDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// envBool parses a boolean environment variable for flag defaults
// (strconv.ParseBool accepts 1, t, T, TRUE, true, True).
func envBool(key string) bool {
	value, _ := strconv.ParseBool(envDefault(key, "false"))
	return value
}

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
	addr := flag.String("addr", ":"+envDefault("PORT", "4000"), "HTTP network address")
	debugMode := flag.Bool("debug", envBool("DEBUG"), "enable debug mode")
	env := flag.String("env", envDefault("ENVIROMENT", "development"), "Environment (development, production, test)")

	tlsCert := flag.String("tls-cert", envDefault("TLS_CERT", "./tls/cert.pem"), "Path to TLS certificate file")
	tlsKey := flag.String("tls-key", envDefault("TLS_KEY", "./tls/key.pem"), "Path to TLS key file")

	dsn := flag.String("dsn", envDefault("DB_DSN", "postgres://noter_web:pass@localhost:5433/noter?sslmode=disable"), "PostgreSQL connection string (libpq URL)")

	flag.Parse()

	// validate env
	if *env != envDevelopment && *env != envProduction && *env != envTest {
		log.Fatal("invalid environment")
		os.Exit(1)
	}

	return &config{
		addr:      *addr,
		debugMode: *debugMode,
		env:       *env,

		tlsCert: *tlsCert,
		tlsKey:  *tlsKey,

		dsn: *dsn,
	}
}
