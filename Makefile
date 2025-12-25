.PHONY: build test format lint install-deps clean run snapshot release-local

# Build the application
build:
	go build -o countdown .

# Run tests
test:
	go test -v ./...

# Format code
format:
	go fmt ./...
	gofumpt -l -w .

# Lint code
lint:
	golangci-lint run

# Install development dependencies
install-deps:
	go mod download
	go mod tidy
	go install mvdan.cc/gofumpt@latest
	go install github.com/goreleaser/goreleaser/v2@latest

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
