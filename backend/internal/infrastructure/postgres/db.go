package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for database/sql
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

var ErrNotFound = errors.New("record not found")

//go:embed migrations/*.sql
var migrations embed.FS

// Connect creates a pgxpool and runs goose migrations automatically.
func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	// Run migrations via database/sql (goose requirement)
	if err := runMigrations(dsn); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	// Main pool via pgxpool (high-performance)
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}
	cfg.MaxConns = 20

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return pool, nil
}

// runMigrations uses goose with embedded SQL files.
// Goose tracks applied migrations in the `goose_db_version` table.
// Adding a new *.sql file → goose applies only the new ones on next startup.
func runMigrations(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open db for migration: %w", err)
	}
	defer db.Close()

	goose.SetBaseFS(migrations)
	goose.SetLogger(goose.NopLogger()) // silence goose output in prod; remove for verbose

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	// "migrations" = the directory inside the embed.FS
	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}
