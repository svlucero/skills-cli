.PHONY: build install test clean run help fmt lint

# Variables
BINARY_NAME=skills
BUILD_DIR=bin
MAIN_PATH=cmd/skill/main.go
INSTALL_PATH=$(GOPATH)/bin

# Default target
help:
	@echo "Skill CLI - Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  install    - Install the binary to \$$GOPATH/bin"
	@echo "  run        - Run the CLI without building"
	@echo "  test       - Run tests"
	@echo "  fmt        - Format code with gofmt"
	@echo "  lint       - Run golangci-lint"
	@echo "  clean      - Clean generated files"
	@echo "  help       - Show this help"

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Binary built at: $(BUILD_DIR)/$(BINARY_NAME)"

# Install to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "$(BINARY_NAME) successfully installed"
	@echo "Make sure \$$GOPATH/bin is in your PATH"

# Run without building
run:
	@go run $(MAIN_PATH) $(ARGS)

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean generated files
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Build with version information
build-release:
	@echo "Building release of $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Release binary built at: $(BUILD_DIR)/$(BINARY_NAME)"

# Format code
fmt:
	@echo "Formatting code..."
	@gofmt -w -s .
	@echo "Code formatted"

# Run linter
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run
	@echo "Lint check complete"
