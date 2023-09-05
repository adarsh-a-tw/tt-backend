# Makefile

# Variables
MIGRATION_PATH = ./db/migrations

# Default target when you just run 'make' without specifying a target
.DEFAULT_GOAL := help

# Help command
help:
	@echo "Available targets:"
	@echo "  test       : Run Go tests"
	@echo "  build      : Build the Go application"
	@echo "  up         : Apply up migrations"
	@echo "  down       : Rollback down migrations"

# Run Go tests
test:
	go test -v ./...

# Build the Go application
build:
	go build -o your-app-name ./cmd/your-app-name

# Apply up migrations
up:
	migrate -path $(MIGRATION_PATH) -database ${DB_DRIVER} up

# Rollback down migrations
down:
	migrate -path $(MIGRATION_PATH) -database ${DB_DRIVER} down
