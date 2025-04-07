# Makefile for Tapir Project

# Variables
APP_NAME = tapir
SRC_DIR = ./cmd/web
BIN_DIR = ./bin
GO_FILES = $(wildcard $(SRC_DIR)/*.go)

# Default target
all: build

# Build the application
build: $(GO_FILES)
	@echo "Building the application..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(GO_FILES)

# Run the application
run: build
	@echo "Running the application..."
	@PORT=4000 $(BIN_DIR)/$(APP_NAME)

# Clean the build
clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)

# Help message
help:
	@echo "Makefile commands:"
	@echo "  make all      - Build the application"
	@echo "  make run      - Build and run the application"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make help     - Show this help message"

.PHONY: all build run clean help