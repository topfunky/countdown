.PHONY: default build test format lint install-deps clean run snapshot release-local ci-local

VERSION := $(shell git describe --tags --always)
LDFLAGS := -X 'main.version=$(VERSION)'

# Default task
default: format lint build

# Build the application
build: clean test
	go build -ldflags "$(LDFLAGS)" -o countdown .

# Run tests
test: clean
	go test -v ./...

# Format code
format:
	gofumpt -w .

# Lint code
lint:
	golangci-lint run
	goreleaser check

# Install development dependencies
install-deps:
	go mod download
	go mod tidy
	go install mvdan.cc/gofumpt@latest
	go install github.com/goreleaser/goreleaser/v2@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Clean build artifacts
clean:
	rm -f countdown
	rm -rf dist/
	go clean

# Run the application
run: build
	./countdown

# Build snapshot release (for testing, no publish)
snapshot:
	goreleaser release --snapshot --clean

# Build release locally (no publish)
release-local:
	goreleaser release --skip=publish --clean

# Run CI locally using act
ci-local:
	sudo -i bash -c "cd /home/dev/projects/countdown && act --job lint && act --job test"
