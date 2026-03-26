# godbmigrate

A simple and efficient database migration engine for Go projects.

## Installation

```bash
go get github.com/ESousa97/godbmigrate
```

## Usage

### Create a New Migration
```bash
godbmigrate new <migration_name>
```

This will generate two files in the `migrations/` directory:
- `YYYYMMDDHHMMSS_<name>.up.sql`
- `YYYYMMDDHHMMSS_<name>.down.sql`

### List Migrations
```bash
godbmigrate list
```

### Apply Pending Migrations
```bash
godbmigrate up --dsn "postgres://user:pass@host:5432/db?sslmode=disable"
```

### Revert Migrations
```bash
# Revert the last applied migration
godbmigrate down --dsn "postgres://user:pass@host:5432/db?sslmode=disable"

# Revert all applied migrations
godbmigrate down --all --dsn "postgres://user:pass@host:5432/db?sslmode=disable"
```

### Check Status
```bash
godbmigrate status --dsn "postgres://user:pass@host:5432/db?sslmode=disable"
```

## Features
- **PostgreSQL Support**: Native integration using `database/sql` and `lib/pq`.
- **Transaction Safety**: Every migration runs in its own transaction.
- **Concurrency Control**: Uses `pg_advisory_lock` to prevent multiple simultaneous migration processes.
- **Structured Logging**: Professional logs using Go's `slog`.
- **Roadmap Complete**: All 5 phases of development have been finished.

## Roadmap

- [x] **Phase 1**: Initial CLI structure and local migration generation.
- [x] **Phase 2**: PostgreSQL integration and migration tracking table.
- [x] **Phase 3**: Execution of migrations (Up) and transaction support.
- [x] **Phase 4**: Reversion of migrations (Down) and rollback support.
- [x] **Phase 5**: Advisory locks, professional logging, and Docker testing environment.

## Technologies
- Go (Golang)
- Cobra CLI
- PostgreSQL (lib/pq)
- Docker (for testing)
