package database

import (
	"database/sql"
	"fmt"
	"livecode-api/middleware"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

var DB *sql.DB

func Connect(databaseURL string) error {
	var err error
	DB, err = sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)
	DB.SetConnMaxIdleTime(1 * time.Minute)

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	middleware.Logger.Info("connected to PostgreSQL",
		zap.String("host", "postgres"),
		zap.Int("max_conns", 25),
	)
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
