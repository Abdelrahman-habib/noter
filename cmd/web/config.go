package main

import (
	"flag"
	"log"
	"os"

	"github.com/Abdelrahman-habib/noter/internal/env"
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
	addr := flag.String("addr", ":"+env.Default("PORT", "4000"), "HTTP network address")
	debugMode := flag.Bool("debug", env.Bool("DEBUG", false), "enable debug mode")
	envFlag := flag.String("env", env.Default("ENVIROMENT", "development"), "Environment (development, production, test)")

	tlsCert := flag.String("tls-cert", env.Default("TLS_CERT", "./tls/cert.pem"), "Path to TLS certificate file")
	tlsKey := flag.String("tls-key", env.Default("TLS_KEY", "./tls/key.pem"), "Path to TLS key file")

	dsn := flag.String("dsn", env.Default("DB_DSN", "postgres://noter_web:pass@localhost:5433/noter?sslmode=disable"), "PostgreSQL connection string (libpq URL)")

	flag.Parse()

	// validate env
	if *envFlag != envDevelopment && *envFlag != envProduction && *envFlag != envTest {
		log.Fatal("invalid environment")
		os.Exit(1)
	}

	return &config{
		addr:      *addr,
		debugMode: *debugMode,
		env:       *envFlag,

		tlsCert: *tlsCert,
		tlsKey:  *tlsKey,

		dsn: *dsn,
	}
}
