# Contributing to godbmigrate

Thank you for your interest in contributing to godbmigrate! We welcome all contributions that help improve this tool.

## Development Environment

- **Go Version**: 1.25.0 or later
- **PostgreSQL**: Required for running integration tests
- **Docker**: Recommended for local database setup

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/ESousa97/godbmigrate.git`
3. Install dependencies: `go mod download`
4. Create a new branch: `git checkout -b feature/your-feature-name`

## Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Run `go fmt ./...` before committing
- Ensure all exported items have [Go doc comments](https://go.dev/doc/comment)

## Testing

Always run tests before submitting a Pull Request:

```bash
# Run unit tests
go test ./...

# Run full integration test (requires Docker)
make test-full
```

## Pull Request Process

1. Ensure your code follows the coding standards and includes tests.
2. Update the `README.md` or `docs/` if you've added new features or changed existing ones.
3. Update the `CHANGELOG.md` following the [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) format.
4. Open a PR with a clear description of the changes and why they are necessary.

---

By contributing, you agree that your contributions will be licensed under its MIT License.
