package store

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
	log  *slog.Logger
}

// New creates a new SQLite database connection with optimized pragmas
func New(dbPath string, logger *slog.Logger) (*DB, error) {
	conn, err := sql.Open("sqlite", "file:"+dbPath)
	if err != nil {
		return nil, fmt.Errorf("store: failed to open database: %w", err)
	}

	db := &DB{
		conn: conn,
		log:  logger,
	}

	// Configure SQLite for optimal performance
	if err := db.configurePragmas(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("store: failed to configure pragmas: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("store: failed to ping database: %w", err)
	}

	return db, nil
}

// configurePragmas sets up SQLite for performance and reliability
func (db *DB) configurePragmas() error {
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA cache_size = -128000;",
		"PRAGMA busy_timeout = 5000;",
		"PRAGMA temp_store = MEMORY;",
		"PRAGMA mmap_size = 268435456;",
		"PRAGMA auto_vacuum = INCREMENTAL;",
		"PRAGMA foreign_keys = ON;",
	}

	for _, pragma := range pragmas {
		if _, err := db.conn.Exec(pragma); err != nil {
			return fmt.Errorf("failed to execute pragma: %w", err)
		}
	}

	return nil
}

// Ping tests the database connection
func (db *DB) Ping() error {
	return db.conn.Ping()
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// Exec executes a statement without returning rows
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.conn.Exec(query, args...)
}

// Query executes a query that returns rows
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.conn.Query(query, args...)
}

// QueryRow executes a query that returns a single row
func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.conn.QueryRow(query, args...)
}

// BeginTx starts a new transaction
func (db *DB) BeginTx() (*sql.Tx, error) {
	return db.conn.Begin()
}

// Conn returns the underlying sql.DB connection for advanced operations
func (db *DB) Conn() *sql.DB {
	return db.conn
}
