.PHONY: test bench lint fmt vet clean examples check release-check

# Default target
all: fmt vet test

# Run all tests
test:
	@echo "Running tests..."
	@go test -v -race -cover $(shell go list ./... | grep -v /examples/)

# Run tests with coverage report
cover:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out $(shell go list ./... | grep -v /examples/)
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Lint code
lint:
	@echo "Running linters..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1; }
	@golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting code..."
	@gofmt -s -w .
	@goimports -w . 2>/dev/null || true

# Vet code
vet:
	@echo "Vetting code..."
	@go vet ./...

# Check everything (pre-commit)
check: fmt vet lint test
	@echo "✓ All checks passed"

# Pre-release validation
release-check: check bench
	@echo "Checking go.mod..."
	@go mod tidy
	@git diff --exit-code go.mod go.sum || { echo "go.mod or go.sum has uncommitted changes"; exit 1; }
	@echo "Checking for uncommitted changes..."
	@git diff --exit-code || { echo "Working directory has uncommitted changes"; exit 1; }
	@echo "✓ Ready for release"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f coverage.out coverage.html
	@go clean -cache -testcache


# Help
help:
	@echo "Available targets:"
	@echo "  make test           - Run all tests"
	@echo "  make cover          - Run tests with coverage report"
	@echo "  make bench          - Run benchmarks"
	@echo "  make lint           - Run linters"
	@echo "  make fmt            - Format code"
	@echo "  make vet            - Vet code"
	@echo "  make check          - Run all checks (fmt, vet, lint, test)"
	@echo "  make release-check  - Validate everything before release"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make examples       - Run basic example"
	@echo "  make build-examples - Build all examples"
	@echo "  make help           - Show this help"
