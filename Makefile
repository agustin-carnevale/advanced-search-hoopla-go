# ======================================
# Project configuration
# ======================================

BINARY_NAME := hoopla
MAIN_PATH := ./main.go
DOCS_DIR := ./docs
DATA_DIR := ./data
ENV_FILE := .env
ENV_EXAMPLE := .env.example
MOVIES_FILE := $(DATA_DIR)/movies.json
DATASET_URL := https://storage.googleapis.com/qvault-webapp-dynamic-assets/course_assets/course-rag-movies.json


# ======================================
# Phony targets
# ======================================

.PHONY: help prepare dataset build run install test fmt clean docs check 

# ======================================
# Help
# ======================================

## Show available make targets
help:
	@echo ""
	@echo "Available targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //'
	@echo ""

# ======================================
# Environment preparation
# ======================================

## Prepare local environment (.env, data folder, dataset hints)
prepare:
	@echo "🔧 Preparing local environment..."
	@echo ""

	@if [ ! -f $(ENV_FILE) ]; then \
		if [ -f $(ENV_EXAMPLE) ]; then \
			echo "📄 Creating $(ENV_FILE) from $(ENV_EXAMPLE)"; \
			cp $(ENV_EXAMPLE) $(ENV_FILE); \
			echo "⚠️  Please edit $(ENV_FILE) and add your API keys"; \
		else \
			echo "❌ $(ENV_EXAMPLE) not found"; \
			echo "👉 Create $(ENV_EXAMPLE) with required environment variables"; \
		fi \
	else \
		echo "✅ $(ENV_FILE) already exists"; \
	fi

	@echo ""

	@if [ ! -d $(DATA_DIR) ]; then \
		echo "📁 Creating $(DATA_DIR)/ directory"; \
		mkdir -p $(DATA_DIR); \
	else \
		echo "✅ $(DATA_DIR)/ directory already exists"; \
	fi

	@echo ""

	@if [ ! -f $(MOVIES_FILE) ]; then \
		echo "🎬 Dataset not found: $(MOVIES_FILE)"; \
		echo "👉 Download the movies dataset and place it at:"; \
		echo "   $(MOVIES_FILE)"; \
	else \
		echo "✅ Dataset found: $(MOVIES_FILE)"; \
	fi

	@echo ""
	@echo "Next steps:"
	@echo "  1. Edit $(ENV_FILE) and add your API keys"
	@echo "  2. Ensure movies.json exists in $(DATA_DIR)/"
	@echo "  3. Run: make build"
	@echo ""

# ======================================
# Dataset download
# ======================================


## Download movies dataset
dataset:
	@echo "⬇️  Downloading movies dataset..."
	@mkdir -p $(DATA_DIR)

	@if command -v curl >/dev/null 2>&1; then \
		curl -L $(DATASET_URL) -o $(MOVIES_FILE); \
	elif command -v wget >/dev/null 2>&1; then \
		wget $(DATASET_URL) -O $(MOVIES_FILE); \
	else \
		echo "❌ Neither curl nor wget is installed"; \
		exit 1; \
	fi

	@echo "✅ Dataset downloaded to $(MOVIES_FILE)"

# ======================================
# Build & Run
# ======================================

## Build the CLI binary
build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

## Run the CLI locally (without installing)
run:
	go run $(MAIN_PATH)

## Install the CLI to GOPATH/bin or GOBIN
install:
	go install .

# ======================================
# Quality & Maintenance
# ======================================

## Run tests
test:
	go test ./...

## Format Go code
fmt:
	go fmt ./...

## Remove build artifacts
clean:
	rm -f $(BINARY_NAME)

# ======================================
# Documentation
# ======================================

## Generate CLI documentation from Cobra
docs:
	@echo "📚 Generating CLI documentation..."
	go run tools/gen-docs/main.go

# ======================================
# Validation
# ======================================

## Check required files without modifying anything
check:
	@echo "🔍 Checking environment..."
	@echo ""

	@if [ ! -f $(ENV_FILE) ]; then \
		echo "❌ Missing $(ENV_FILE)"; \
	else \
		echo "✅ $(ENV_FILE) exists"; \
	fi

	@if [ ! -f $(MOVIES_FILE) ]; then \
		echo "❌ Missing dataset: $(MOVIES_FILE)"; \
	else \
		echo "✅ Dataset exists: $(MOVIES_FILE)"; \
	fi

	@echo ""
