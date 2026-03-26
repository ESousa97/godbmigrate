# godbmigrate

<div align="center">
  <h1>godbmigrate</h1>
  <p>A fast, flexible, and robust database migration tool for Go projects with advisory locks.</p>

  <img src="assets/github-go.png" alt="godbmigrate Banner" width="600px">

  <br>

[![CI](https://github.com/ESousa97/godbmigrate/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/ESousa97/godbmigrate/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ESousa97/godbmigrate?=v7)](https://goreportcard.com/report/github.com/ESousa97/godbmigrate?=v7)
[![Go Reference](https://pkg.go.dev/badge/github.com/ESousa97/godbmigrate.svg)](https://pkg.go.dev/github.com/ESousa97/godbmigrate)
[![License](https://img.shields.io/github/license/ESousa97/godbmigrate)](https://github.com/ESousa97/godbmigrate/blob/master/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ESousa97/godbmigrate)](https://github.com/ESousa97/godbmigrate)
[![Last Commit](https://img.shields.io/github/last-commit/ESousa97/godbmigrate)](https://github.com/ESousa97/godbmigrate/commits/master)

</div>

---

`godbmigrate` is a lightweight CLI tool and Go library designed to manage database schema evolutions with ease and safety. It leverages PostgreSQL advisory locks to prevent race conditions during concurrent migration attempts, ensuring your database remains consistent.

## Demonstration

### CLI Usage

```bash
# Create a new migration pair
./godbmigrate new add_users_table

# Apply all pending migrations
./godbmigrate up --dsn "postgres://user:pass@localhost:5432/dbname?sslmode=disable"

# Revert the last migration
./godbmigrate down --dsn "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
```

### Library Usage

```go
import "github.com/ESousa97/godbmigrate/internal/db"

// Connect to the database
store, err := db.Connect(dsn)
if err != nil {
    log.Fatal(err)
}
defer store.Close()

// Apply pending migrations
err = store.ApplyPending("migrations/")
```

## Tech Stack

| Technology | Role |
|---|---|
| Go 1.25 | Core language and logic |
| Cobra | CLI interface and command management |
| lib/pq | PostgreSQL driver |
| slog | Structured logging for observability |

## Prerequisites

- Go >= 1.25 (defined in `go.mod`)
- PostgreSQL instance (for migrations)
- `golangci-lint` (for development)

## Installation and Usage

### As a Binary

```bash
go install github.com/ESousa97/godbmigrate@latest
```

### From Source

```bash
git clone https://github.com/ESousa97/godbmigrate.git
cd godbmigrate
make build
./godbmigrate --help
```

## Makefile Targets

| Target | Description |
|---|---|
| `make build` | Compiles the project binary |
| `make db-up` | Starts a PostgreSQL container via Docker |
| `make db-down` | Stops and removes the PostgreSQL container |
| `make test-full` | Runs a complete end-to-end migration test |
| `make clean` | Removes build artifacts and test migrations |

## Architecture

The project adopts a modular architecture focused on safety and simplicity:

- **`cmd/`**: CLI command definitions using Cobra.
- **`internal/db/`**: Core database logic and migration engine.
- **`Advisory Locks`**: Uses PostgreSQL's `pg_try_advisory_lock` to ensure only one migration runs at a time.
- **`Transactions`**: Every migration is executed within a database transaction for atomicity.

## API Reference

Detailed documentation of the internal packages can be found at [pkg.go.dev](https://pkg.go.dev/github.com/ESousa97/godbmigrate).

## Configuration

The CLI accepts flags for database connection.

| Flag | Description | Default |
|---|---|---|
| `--dsn` | PostgreSQL Data Source Name | `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable` |
| `--debug` | Enable debug logging | `false` |
| `--all` | Revert all migrations (used with `down`) | `false` |

## Roadmap

- [x] Phase 1: Core Migration Engine (.up/.down SQL)
- [x] Phase 2: CLI Interface with Cobra
- [x] Phase 3: PostgreSQL Advisory Locking
- [x] Phase 4: Structured Logging with slog
- [x] Phase 5: Transactional Atomicity & Makefile Automation

## Contributing

Contributions are welcome! See the full guide at [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Distributed under the MIT license. See [LICENSE](LICENSE) for more details.

<div align="center">

## Author

**Enoque Sousa**

[![LinkedIn](https://img.shields.io/badge/LinkedIn-0077B5?style=flat&logo=linkedin&logoColor=white)](https://www.linkedin.com/in/enoque-sousa-bb89aa168/)
[![GitHub](https://img.shields.io/badge/GitHub-100000?style=flat&logo=github&logoColor=white)](https://github.com/ESousa97)
[![Portfolio](https://img.shields.io/badge/Portfolio-FF5722?style=flat&logo=target&logoColor=white)](https://enoquesousa.vercel.app)

**[⬆ Back to top](#godbmigrate)**

Made with ❤️ by [Enoque Sousa](https://github.com/ESousa97)

**Project Status:** Active — Study Project

</div>
