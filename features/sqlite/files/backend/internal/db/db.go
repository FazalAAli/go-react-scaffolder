package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

// Pool is the database handle attached to the app. Aliased so callers depend on
// this package rather than the concrete driver type.
type Pool = sql.DB

// New opens a SQLite handle from DATABASE_URL and verifies it with a ping.
func New(ctx context.Context) (*Pool, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	if err := conn.PingContext(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}
	return conn, nil
}
