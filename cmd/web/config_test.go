package main

import (
	"flag"
	"io"
	"os"
	"testing"
)

func TestParseFlagsDsnPrecedence(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		envDSN  string
		wantDsn string
	}{
		{"explicit flag beats env", []string{"-dsn=postgres://flag-dsn"}, "postgres://env-dsn", "postgres://flag-dsn"},
		{"explicit flag equal to default beats env", []string{"-dsn=" + defaultDSN}, "postgres://env-dsn", defaultDSN},
		{"explicit empty -dsn acts as unset", []string{"-dsn="}, "postgres://env-dsn", "postgres://env-dsn"},
		{"env beats default", []string{}, "postgres://env-dsn", "postgres://env-dsn"},
		{"default without env", []string{}, "", defaultDSN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("DB_DSN", tt.envDSN)

			flag.CommandLine = flag.NewFlagSet("noter", flag.ContinueOnError)
			flag.CommandLine.SetOutput(io.Discard)
			os.Args = append([]string{os.Args[0]}, tt.args...)

			config := parseFlags()
			if config.dsn != tt.wantDsn {
				t.Errorf("got %q; want %q", config.dsn, tt.wantDsn)
			}
		})
	}
}
