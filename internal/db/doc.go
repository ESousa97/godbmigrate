// Package db provides a robust and thread-safe implementation for managing
// database migrations in PostgreSQL.
//
// It handles connection management, schema version tracking, and ensures
// execution safety using advisory locks to prevent concurrent migration runs.
//
// The primary entry point is [Connect], which returns a [*MigrationStore]:
//
//	store, err := db.Connect(dsn)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer store.Close()
//
// You can then apply pending migrations using [MigrationStore.ApplyMigration]:
//
//	err := store.ApplyMigration(20231027120000, "CREATE TABLE users (...)")
package db
