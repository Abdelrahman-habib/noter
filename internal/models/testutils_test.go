package models

import (
	"database/sql"
	"flag"
	"os"
	"sync"
	"testing"

	schema "github.com/Abdelrahman-habib/noter/db/schema"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	// database/sql uses "pgx"; goose uses "postgres".
	sqlDriver      = "pgx"
	defaultTestDSN = "postgres://noter_test_web:test_pass@localhost:5433/noter_test?sslmode=disable"
)

type testDBConfig struct {
	dbDialect string
	testDSN   string
}

var (
	// Register flags before the test runner parses them.
	dbDialect = flag.String("db-dialect", "postgres", "Goose migration dialect (eg: postgres)")
	testDSN   = flag.String("test-dsn", envDefault("TEST_DB_DSN", defaultTestDSN), "PostgreSQL test database connection string")

	configOnce     sync.Once
	configInstance *testDBConfig
)

func newTestDB(t *testing.T) *sql.DB {
	config := getTestDBConfig()

	// Establish a sql.DB connection pool for our test database
	db, err := sql.Open(sqlDriver, config.testDSN)
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
			t.Errorf("failed to reset seeds: %v", err)
		}

		// Reset migrations
		if err := goose.Reset(db, "migrations", goose.WithNoVersioning()); err != nil {
			t.Errorf("failed to reset migrations: %v", err)
		}
	})

	// Return the database connection pool.
	return db
}

// envDefault uses the environment variable as the flag default
func envDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getTestDBConfig() *testDBConfig {
	configOnce.Do(func() {
		configInstance = &testDBConfig{
			dbDialect: *dbDialect,
			testDSN:   *testDSN,
		}
	})
	return configInstance
}
