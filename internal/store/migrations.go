package store

import (
	"database/sql"
	"fmt"
	"log/slog"
)

// Migrator handles database schema migrations
type Migrator struct {
	db  *DB
	log *slog.Logger
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *DB, logger *slog.Logger) *Migrator {
	return &Migrator{
		db:  db,
		log: logger,
	}
}

// ApplyMigrations applies all pending migrations to the database
func (m *Migrator) ApplyMigrations() error {
	// Create schema_version table if it doesn't exist
	if err := m.createSchemaVersionTable(); err != nil {
		return fmt.Errorf("migrations: failed to create schema_version table: %w", err)
	}

	// Get current schema version
	currentVersion, err := m.getCurrentVersion()
	if err != nil {
		return fmt.Errorf("migrations: failed to get current version: %w", err)
	}

	// Define all migrations
	migrations := []struct {
		version int
		sql     string
		name    string
	}{
		{
			version: 1,
			name:    "Initial schema",
			sql:     schemav1,
		},
	}

	// Apply pending migrations
	for _, migration := range migrations {
		if migration.version > currentVersion {
			if err := m.applyMigration(migration.version, migration.sql, migration.name); err != nil {
				return err
			}
		}
	}

	m.log.Info("migrations applied successfully", slog.Int("current_version", currentVersion))
	return nil
}

// createSchemaVersionTable creates the schema_version table if it doesn't exist
func (m *Migrator) createSchemaVersionTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_version (
			version INTEGER PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := m.db.Exec(query)
	return err
}

// getCurrentVersion gets the current schema version
func (m *Migrator) getCurrentVersion() (int, error) {
	var version int
	err := m.db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_version").Scan(&version)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return version, nil
}

// applyMigration applies a single migration
func (m *Migrator) applyMigration(version int, sql string, name string) error {
	m.log.Info("applying migration", slog.Int("version", version), slog.String("name", name))

	tx, err := m.db.BeginTx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(sql); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	if _, err := tx.Exec("INSERT INTO schema_version (version) VALUES (?)", version); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// schemav1 defines the initial database schema
const schemav1 = `
-- Users table
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY,
	email TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	role TEXT NOT NULL DEFAULT 'user',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Sessions table
CREATE TABLE IF NOT EXISTS sessions (
	id TEXT PRIMARY KEY,
	user_id INTEGER NOT NULL,
	expires_at DATETIME NOT NULL,
	last_accessed_at DATETIME,
	FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Media table
CREATE TABLE IF NOT EXISTS media (
	id TEXT PRIMARY KEY,
	user_id INTEGER NOT NULL,
	file_path TEXT UNIQUE NOT NULL,
	file_name TEXT NOT NULL,
	media_type TEXT NOT NULL,
	mime_type TEXT NOT NULL,
	file_size INTEGER NOT NULL,
	width INTEGER,
	height INTEGER,
	duration INTEGER,
	lat REAL,
	lon REAL,
	taken_at DATETIME,
	blurhash TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Albums table
CREATE TABLE IF NOT EXISTS albums (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT,
	created_by INTEGER NOT NULL,
	cover_media_id TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (created_by) REFERENCES users(id),
	FOREIGN KEY (cover_media_id) REFERENCES media(id)
);

-- Album-Media junction table
CREATE TABLE IF NOT EXISTS album_media (
	album_id INTEGER NOT NULL,
	media_id TEXT NOT NULL,
	sort_order INTEGER DEFAULT 0,
	PRIMARY KEY (album_id, media_id),
	FOREIGN KEY (album_id) REFERENCES albums(id),
	FOREIGN KEY (media_id) REFERENCES media(id)
);

-- Shared links table
CREATE TABLE IF NOT EXISTS shared_links (
	token TEXT PRIMARY KEY,
	target_type TEXT NOT NULL,
	target_id TEXT NOT NULL,
	password_hash TEXT,
	max_access_count INTEGER,
	access_count INTEGER DEFAULT 0,
	expires_at DATETIME NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- FTS5 virtual table for full-text search
CREATE VIRTUAL TABLE IF NOT EXISTS media_search USING fts5(
	media_id UNINDEXED,
	file_name,
	tags,
	album_names,
	tokenize='unicode61'
);

-- Create indices for better query performance
CREATE INDEX IF NOT EXISTS idx_media_user_id ON media(user_id);
CREATE INDEX IF NOT EXISTS idx_media_taken_at ON media(taken_at);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_albums_created_by ON albums(created_by);
CREATE INDEX IF NOT EXISTS idx_album_media_media_id ON album_media(media_id);
CREATE INDEX IF NOT EXISTS idx_shared_links_expires_at ON shared_links(expires_at);
`
