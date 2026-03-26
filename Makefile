# Database Configuration
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=postgres
POSTGRES_PORT=5432
DSN="postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"

# Commands
.PHONY: build db-up db-down test-full clean

build:
	go build -o godbmigrate.exe main.go

db-up:
	docker run --name godbmigrate-postgres -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -p $(POSTGRES_PORT):5432 -d postgres
	@echo "Waiting for Postgres to be ready..."
	@timeout /t 5 /nobreak > nul

db-down:
	docker stop godbmigrate-postgres || true
	docker rm godbmigrate-postgres || true

test-full: build
	@echo "Creating new migration..."
	./godbmigrate.exe new test_users
	@echo "Checking status (should be empty/no migrations)..."
	./godbmigrate.exe status --dsn $(DSN) || true
	@echo "Running migrations UP..."
	./godbmigrate.exe up --dsn $(DSN)
	@echo "Checking status (should show applied version)..."
	./godbmigrate.exe status --dsn $(DSN)
	@echo "Running migrations DOWN..."
	./godbmigrate.exe down --dsn $(DSN)
	@echo "Checking status (should be empty again)..."
	./godbmigrate.exe status --dsn $(DSN)
	@echo "Test completed successfully!"

clean:
	rm -f godbmigrate.exe
	rm -rf migrations/
