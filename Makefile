.PHONY: build run clean test

# Define variables
APP_NAME = littleScrapper
SRC = main.go
DIST = bin
LDFLAGS = -w -s

# Determine the OS
# Set OS-specific variables
ifeq ($(OS),Windows_NT)
	BINARY_NAME = $(DIST)/$(APP_NAME).exe
	CLEAN_CMD = rmdir /s /q
else
	BINARY_NAME = $(DIST)/$(APP_NAME)
	CLEAN_CMD = rm -rf
endif

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir $(DIST)
	@go mod vendor
	@go build -o $(BINARY_NAME) -ldflags "$(LDFLAGS)" $(SRC)

# Run the application
run:
	@go run $(SRC)

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@$(CLEAN_CMD) $(DIST)

# Run tests
test:
	@go test ./...

# Help
help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  build     - Build the application"
	@echo "  run       - Run the application"
	@echo "  test      - Run tests"
	@echo "  clean     - Clean up build artifacts"

# Default target
.DEFAULT_GOAL := help
