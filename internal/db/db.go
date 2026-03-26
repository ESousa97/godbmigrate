package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

const migrationLockID = 123456789 // Unique ID for our advisory lock

// MigrationStore manages the database connection and the schema_migrations table.
// It uses a standard [*sql.DB] handle for all operations.
type MigrationStore struct {
	// DB is the underlying database connection.
	DB *sql.DB
}

// Connect establishes a connection to a PostgreSQL database using the provided DSN.
// It verifies the connection with a ping and ensures the schema_migrations table exists.
//
// Returns a [*MigrationStore] for managing migrations.
func Connect(dsn string) (*MigrationStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not connect to postgres: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping postgres: %w", err)
	}

	store := &MigrationStore{DB: db}
	if err := store.EnsureSchemaTable(); err != nil {
		return nil, err
	}

	return store, nil
}

// EnsureSchemaTable creates the schema_migrations table if it doesn't already exist.
// This table tracks applied migration versions and their timestamps.
func (s *MigrationStore) EnsureSchemaTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version BIGINT PRIMARY KEY,
		applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := s.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("could not ensure schema_migrations table: %w", err)
	}
	return nil
}

// AcquireLock attempts to obtain a PostgreSQL advisory lock.
// This prevents multiple migration processes from running concurrently.
//
// Returns an error if the lock cannot be acquired or if another process holds it.
func (s *MigrationStore) AcquireLock() error {
	slog.Debug("Acquiring advisory lock", "lock_id", migrationLockID)
	var locked bool
	query := "SELECT pg_try_advisory_lock($1)"
	err := s.DB.QueryRow(query, migrationLockID).Scan(&locked)
	if err != nil {
		return fmt.Errorf("could not acquire advisory lock: %w", err)
	}
	if !locked {
		return fmt.Errorf("could not acquire advisory lock: another migration process is currently running")
	}
	slog.Info("Advisory lock acquired")
	return nil
}

// ReleaseLock releases the previously acquired advisory lock.
func (s *MigrationStore) ReleaseLock() error {
	slog.Debug("Releasing advisory lock", "lock_id", migrationLockID)
	query := "SELECT pg_advisory_unlock($1)"
	_, err := s.DB.Exec(query, migrationLockID)
	if err != nil {
		return fmt.Errorf("could not release advisory lock: %w", err)
	}
	slog.Info("Advisory lock released")
	return nil
}

// GetLatestVersion returns the version number of the most recently applied migration.
// If no migrations have been applied, it returns 0 without an error.
func (s *MigrationStore) GetLatestVersion() (int64, error) {
	var version int64
	query := "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1"

	err := s.DB.QueryRow(query).Scan(&version)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("could not query latest version: %w", err)
	}

	return version, nil
}

// ApplyMigration executes a migration's SQL content and records its version in the schema table.
// It runs the entire operation within a single transaction.
func (s *MigrationStore) ApplyMigration(version int64, sqlContent string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}

	// Execute migration SQL
	if _, err := tx.Exec(sqlContent); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Record the migration
	if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to record migration version: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// RevertMigration executes a 'down' migration's SQL content and removes its record from the schema table.
// It runs the entire operation within a single transaction.
func (s *MigrationStore) RevertMigration(version int64, sqlContent string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}

	// Execute down migration SQL
	if _, err := tx.Exec(sqlContent); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to execute down migration: %w", err)
	}

	// Remove the migration record
	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", version); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to remove migration version: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetAppliedVersions retrieves all migration versions from the schema table, sorted in descending order.
func (s *MigrationStore) GetAppliedVersions() ([]int64, error) {
	var versions []int64
	query := "SELECT version FROM schema_migrations ORDER BY version DESC"

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not query applied versions: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var v int64
		if err := rows.Scan(&v); err != nil {
			return nil, fmt.Errorf("could not scan version: %w", err)
		}
		versions = append(versions, v)
	}

	return versions, rows.Err()
}

// Close closes the underlying database connection.
func (s *MigrationStore) Close() error {
	return s.DB.Close()
}
