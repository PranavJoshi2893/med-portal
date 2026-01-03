package database

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "github.com/lib/pq"

	"github.com/PranavJoshi2893/boilerplate/internal/config"
)

func NewPostgres(cfg *config.Config) (*sql.DB, error) {
	// 1. Build a Connection URI to safely handle special characters
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.DBUser, cfg.DBPassword),
		Host:   fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort),
		Path:   cfg.DBName,
	}

	// Add query parameters (e.g., sslmode)
	q := dsn.Query()
	q.Set("sslmode", cfg.DBSSLMode)
	dsn.RawQuery = q.Encode()

	// 2. Open the connection pool
	db, err := sql.Open("postgres", dsn.String())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 3. Set production-ready connection pool limits
	db.SetMaxOpenConns(25)                 // Max active connections
	db.SetMaxIdleConns(25)                 // Max idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Prevent "stale" connection errors

	// 4. Verify the connection (sql.Open is lazy and doesn't check connectivity)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
