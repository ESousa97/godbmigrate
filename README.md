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

### Check Status
```bash
godbmigrate status --dsn "postgres://user:pass@host:5432/db?sslmode=disable"
```

## Roadmap

- [x] **Phase 1**: Initial CLI structure and local migration generation.
- [x] **Phase 2**: PostgreSQL integration and migration tracking table.
- [ ] **Phase 3**: Execution of migrations (Up/Down) and transaction support.

## Technologies
- Go (Golang)
- Cobra CLI
- PostgreSQL (lib/pq)
