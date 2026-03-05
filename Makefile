.PHONY: build test clean dev lint help

# Build the project
build:
	go build -o bin/simple-gallery .
	@echo "Build complete: bin/simple-gallery"

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Run in development mode
dev: build
	./bin/simple-gallery -config config.toml

# Run linter
lint:
	golangci-lint run ./...

# Help
help:
	@echo "Available commands:"
	@echo "  make build   - Build the binary"
	@echo "  make test    - Run tests"
	@echo "  make clean   - Remove build artifacts"
	@echo "  make dev     - Build and run in development mode"
	@echo "  make lint    - Run linter"
	@echo "  make help    - Show this help message"
