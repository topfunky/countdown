.PHONY: build test format lint install-deps clean run

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

# Clean build artifacts
clean:
	rm -f countdown
	go clean

# Run the application
run: build
	./countdown
