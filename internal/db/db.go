package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type MigrationStore struct {
	DB *sql.DB
}

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

func (s *MigrationStore) ApplyMigration(version int64, sqlContent string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}

	// Execute migration SQL
	if _, err := tx.Exec(sqlContent); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Record the migration
	if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record migration version: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *MigrationStore) RevertMigration(version int64, sqlContent string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}

	// Execute down migration SQL
	if _, err := tx.Exec(sqlContent); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute down migration: %w", err)
	}

	// Remove the migration record
	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", version); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration version: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *MigrationStore) GetAppliedVersions() ([]int64, error) {
	var versions []int64
	query := "SELECT version FROM schema_migrations ORDER BY version DESC"
	
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not query applied versions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var v int64
		if err := rows.Scan(&v); err != nil {
			return nil, fmt.Errorf("could not scan version: %w", err)
		}
		versions = append(versions, v)
	}
	
	return versions, nil
}

func (s *MigrationStore) Close() error {
	return s.DB.Close()
}
