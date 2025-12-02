# ============================================================================
# CONFIGURATION
# ============================================================================

BINARY_NAME=bc-km

# Auto-load .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

# ============================================================================
# MAIN COMMANDS
# ============================================================================

.PHONY: all build run clean help

all: build ## Default: build the application

build: ## Build the binary
	@echo "Building..."
	@go build -o ./bin/$(BINARY_NAME) .
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

init: ## Initialize fresh blockchain (deletes existing data)
	@echo "Initializing fresh blockchain..."
	@rm -f blockchain.json
	@go run . -list
	@echo "Fresh blockchain initialized!"

help: ## Show this help message
	@echo "Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "For operations like add, search, validate, use:"
	@echo "   ./bin/$(BINARY_NAME) --help"
	@echo "   OR"
	@echo "   go run . --help"
