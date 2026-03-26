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

func (s *MigrationStore) Close() error {
	return s.DB.Close()
}
