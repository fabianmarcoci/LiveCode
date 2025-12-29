package database

import (
	"errors"
	"fmt"

	"livecode-api/middleware"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func RunMigrations(databaseURL string) error {
	middleware.Logger.Info("starting database migrations")

	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		middleware.Logger.Error("failed to create migration driver",
			zap.Error(err),
		)
		return fmt.Errorf("migration driver creation failed: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		middleware.Logger.Error("failed to create migrate instance",
			zap.Error(err),
		)
		return fmt.Errorf("migrate instance creation failed: %w", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			middleware.Logger.Info("no new migrations to apply - database schema is up to date")
			return nil
		}

		middleware.Logger.Error("migration execution failed",
			zap.Error(err),
		)
		return fmt.Errorf("migration failed: %w", err)
	}

	middleware.Logger.Info("database migrations completed successfully")
	return nil
}

func CheckMigrationsApplied() error {
	var version int
	err := DB.QueryRow("SELECT version FROM schema_migrations WHERE dirty = false LIMIT 1").Scan(&version)
	if err != nil {
		return fmt.Errorf("no migrations applied: %w", err)
	}
	return nil
}
