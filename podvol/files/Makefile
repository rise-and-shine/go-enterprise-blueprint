# Makefile

# Go parameters
GOCMD=go
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOGENERATE=$(GOCMD) generate

# Tool versions (pin for reproducibility)
MOCKERY_VERSION=v2.46.0

# Directories
MOCK_DIR=./internal/.../mocks

.PHONY: all test clean mocks mocks-check tools lint

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## tools: Install development tools
tools:
	@echo "Installing mockery..."
	go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION)
	@echo "Tools installed successfully"

## mocks: Generate all mocks using .mockery.yaml config
mocks:
	@echo "Generating mocks..."
	mockery
	@echo "Mocks generated successfully"

## mocks-dry: Show what mocks would be generated (dry run)
mocks-dry:
	mockery --dry-run

## mocks-clean: Remove all generated mock files
mocks-clean:
	@echo "Cleaning generated mocks..."
	find . -type d -name 'mocks' -exec rm -rf {} + 2>/dev/null || true
	@echo "Mocks cleaned"

## mocks-check: Verify mocks are up to date (for CI)
mocks-check: mocks-clean mocks
	@echo "Checking if mocks are up to date..."
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "ERROR: Generated mocks are out of date. Run 'make mocks' and commit."; \
		git diff; \
		exit 1; \
	fi
	@echo "Mocks are up to date"

## generate: Run all go:generate directives
generate:
	$(GOGENERATE) ./...

## test: Run all tests
test:
	$(GOTEST) -v -race ./...

## test-coverage: Run tests with coverage
test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

## lint: Run linter
lint:
	golangci-lint run

## ci: Run all CI checks
ci: tools mocks-check lint test
