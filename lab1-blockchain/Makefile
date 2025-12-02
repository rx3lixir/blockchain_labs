# ============================================================================
# CONFIGURATION
# ============================================================================

# Binary name for compilation
BINARY_NAME=blockchain-app

# Auto-load .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

# ============================================================================
# DEFAULT TARGET
# ============================================================================

all: build ## Default build command

# ============================================================================
# BUILD & RUN
# ============================================================================

build: ## Build the binary
	@echo "Building..."
	go build -o ./bin/$(BINARY_NAME) .

run: build ## Build and run the application
	@echo "Running..."
	./bin/$(BINARY_NAME)

clean: ## Clean all binaries
	@echo "Cleaning..."
	go clean
	rm -f ./bin/$(BINARY_NAME)
	rm -f blockchain_data.json

# ============================================================================
# PROGRAM OPERATIONS
# ============================================================================

list: ## Show all blocks
	go run . -list

validate: ## Validate blockchain integrity
	go run . -validate

debts: ## Find students with debts
	go run . -debts

init-blockchain: ## Initialize new blockchain (delete existing)
	rm -f blockchain.json
	go run . -list

search-index: ## Search by index (usage: make search-index INDEX=1)
	@if [ -z "$(INDEX)" ]; then \
		echo "Ошибка: введите параметр INDEX"; \
		echo "Пример: make search-index INDEX=1"; \
		exit 1; \
	fi
	go run . -index $(INDEX)

search: ## Search by custom keyword (usage: make search KEYWORD="Петров")
	@if [ -z "$(KEYWORD)" ]; then \
		echo "Ошибка: введите параметр KEYWORD"; \
		echo "Пример: make search KEYWORD=Петров"; \
		exit 1; \
	fi
	go run . -keyword "$(KEYWORD)"

# Универсальная команда для добавления с параметрами
add: ## Add custom record (usage: make add COURSE=3 GROUP="ИВТ-301" NAME="Иванов" RECORD="202101" SUBJECT="Математика" GRADE=5)
	@if [ -z "$(COURSE)" ] || [ -z "$(GROUP)" ] || [ -z "$(NAME)" ] || [ -z "$(RECORD)" ] || [ -z "$(SUBJECT)" ] || [ -z "$(GRADE)" ]; then \
		echo "Ошибка: Параметры: COURSE, GROUP, NAME, RECORD, SUBJECT, GRADE"; \
		echo "Пример: make add COURSE=3 GROUP='ИВТ-301' NAME='Иванов И.И.' RECORD='202101' SUBJECT='Математика' GRADE=5"; \
		exit 1; \
	fi
	go run . -add -course $(COURSE) -group "$(GROUP)" -name "$(NAME)" -record "$(RECORD)" -subject "$(SUBJECT)" -grade $(GRADE)

# ============================================================================
# UTILITIES
# ============================================================================

help: ## Show help
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}'
