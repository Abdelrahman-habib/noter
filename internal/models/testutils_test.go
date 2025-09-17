package models

import (
	"database/sql"
	"testing"

	schema "github.com/Abdelrahman-habib/snippetbox/db/schema"
	"github.com/pressly/goose/v3"
)

func newTestDB(t *testing.T) *sql.DB {
	config := parseFlags()
	// Establish a sql.DB connection pool for our test database
	db, err := sql.Open(config.dbDialect, config.testDSN)
	if err != nil {
		t.Fatal(err)
	}

	goose.SetBaseFS(schema.EmbedMigrations)

	// Set the dialect for goose
	if err := goose.SetDialect(config.dbDialect); err != nil {
		db.Close()
		t.Fatalf("failed to set goose dialect: %v", err)
	}

	// Run migrations using goose library
	if err := goose.Up(db, "migrations", goose.WithNoVersioning()); err != nil {
		db.Close()
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Run seed data using goose library (no versioning for seeds)
	if err := goose.Up(db, "seed", goose.WithNoVersioning()); err != nil {
		db.Close()
		t.Fatalf("failed to run seeds: %v", err)
	}

	// Use t.Cleanup() to register a function which will automatically be
	// called by Go when the current test (or sub-test) which calls newTestdb()
	// has finished*. In this function we reset the database and close the connection.
	t.Cleanup(func() {
		defer db.Close()

		// Reset seed data (no versioning for seeds)
		if err := goose.Reset(db, "seed", goose.WithNoVersioning()); err != nil {
			t.Logf("warning: failed to reset seeds: %v", err)
		}

		// Reset migrations
		if err := goose.Reset(db, "migrations", goose.WithNoVersioning()); err != nil {
			t.Logf("warning: failed to reset migrations: %v", err)
		}
	})

	// Return the database connection pool.
	return db
}
