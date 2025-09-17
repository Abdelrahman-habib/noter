package models

import (
	"flag"
	"sync"
)

type config struct {
	dbDialect string
	testDSN   string
}

var (
	configOnce     sync.Once
	configInstance *config
)

func parseFlags() *config {
	configOnce.Do(func() {
		dbDialect := flag.String("db-dialect", "mysql", "Database Dialect (eg: mysql, postgres, etc.)")
		testDSN := flag.String("test-dsn", "test_web:test_pass@/test_snippetbox?parseTime=true", "MySQL test data source name")
		flag.Parse()

		configInstance = &config{
			dbDialect: *dbDialect,
			testDSN:   *testDSN,
		}
	})
	return configInstance
}
