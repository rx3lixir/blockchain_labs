# ============================================================================
# CONFIGURATION
# ============================================================================

BINARY_NAME=bc

# ============================================================================
# MAIN COMMANDS
# ============================================================================

.PHONY: all build run clean help

all: build ## Default: build the application

build: ## Build the binary
	@echo "Building..."
	@go build -o ./bin/$(BINARY_NAME) ./cmd/lab/main.go
	@echo "Build complete!"

run: build ## Build and run the application
	@echo "Running..."
	@./bin/$(BINARY_NAME)

clean: ## Clean binaries and blockchain data
	@echo "Cleaning..."
	@go clean
	@rm -f ./bin/$(BINARY_NAME)
	@rm -f blockchain.json
	@echo "Clean complete!"

help: ## Show help
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}'
