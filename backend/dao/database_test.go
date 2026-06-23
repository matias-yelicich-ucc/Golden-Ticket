package dao

import (
	"errors"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestLoadDBConfigFromEnvUsesDefaults(t *testing.T) {
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "")

	cfg := loadDBConfigFromEnv()
	if cfg.host != "localhost" || cfg.port != "3306" || cfg.user != "root" || cfg.name != "golden_ticket" {
		t.Fatalf("unexpected default config: %+v", cfg)
	}
}

func TestLoadDBConfigFromEnvUsesProvidedValues(t *testing.T) {
	t.Setenv("DB_HOST", "db")
	t.Setenv("DB_PORT", "3307")
	t.Setenv("DB_USER", "mati")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "tickets")

	cfg := loadDBConfigFromEnv()
	if cfg.host != "db" || cfg.port != "3307" || cfg.user != "mati" || cfg.pass != "secret" || cfg.name != "tickets" {
		t.Fatalf("unexpected config from env: %+v", cfg)
	}
}

func TestBuildDSN(t *testing.T) {
	dsn := buildDSN(dbConfig{
		host: "db",
		port: "3306",
		user: "mati",
		pass: "secret",
		name: "golden_ticket",
	})

	expectedParts := []string{"mati:secret@tcp(db:3306)/golden_ticket", "charset=utf8mb4", "parseTime=True", "loc=Local"}
	for _, part := range expectedParts {
		if !strings.Contains(dsn, part) {
			t.Fatalf("expected DSN to contain %q, got %q", part, dsn)
		}
	}
}

func TestSetupConnectionPoolAndMigrationHelpers(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:database-helper?mode=memory&cache=shared&_foreign_keys=1"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql db: %v", err)
	}

	setupConnectionPool(sqlDB)
	if sqlDB.Stats().MaxOpenConnections != 100 {
		t.Fatalf("expected max open connections to be 100, got %d", sqlDB.Stats().MaxOpenConnections)
	}

	if err := migrateAndPrepareDatabase(db); err != nil {
		t.Fatalf("expected migrateAndPrepareDatabase to succeed, got %v", err)
	}

	if err := initializeDB("ignored", func(string) (*gorm.DB, error) {
		return gorm.Open(sqlite.Open("file:database-init?mode=memory&cache=shared&_foreign_keys=1"), &gorm.Config{})
	}); err != nil {
		t.Fatalf("expected initializeDB to succeed, got %v", err)
	}
}

func TestOpenDatabaseFailsWithInvalidTarget(t *testing.T) {
	_, err := openDatabase("invalid dsn")
	if err == nil {
		t.Fatalf("expected openDatabase to fail for invalid dsn")
	}
}

func TestInitializeDBPropagatesOpenErrors(t *testing.T) {
	expectedErr := errors.New("open failed")
	err := initializeDB("ignored", func(string) (*gorm.DB, error) {
		return nil, expectedErr
	})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected open error to be returned, got %v", err)
	}
}
