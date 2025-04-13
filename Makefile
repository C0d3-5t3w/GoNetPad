.PHONY: help all build clean run

MAIN_BINARY := GoNetPad
BUILD_DIR := GoNetPad
CMD_DIR := cmd/app

GOOS ?= darwin
GOARCH ?= arm64
CGO_ENABLED ?= 1
LDFLAGS = -w -s

help:
	@echo "Makefile for $(MAIN_BINARY)"
	@echo "Usage:"
	@echo "  make all       - Build and run the application"
	@echo "  make build     - Build the application"
	@echo "  make clean     - Clean up build artifacts"
	@echo "  make run       - Run the application"
	@echo "  make deps      - Install dependencies"
	@echo "  make help      - Show this help message"

all: clean deps ts sass build run

build:
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags="$(LDFLAGS)" \
		-trimpath \
		-o $(BUILD_DIR)/$(MAIN_BINARY) \
		./$(CMD_DIR)/main.go

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -rf pkg/website/assets/js
	@rm -rf pkg/website/assets/css

deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download
	@go mod vendor
	@go mod verify

ts:
	@echo "Building TypeScript files..."
	@tsc --outDir pkg/website/assets/js

sass:
	@echo "Building SASS files..."
	@sass --style compressed pkg/website/static/sass/style.scss:pkg/website/assets/css/style.css

run: 
	@echo "Running application..."
	./$(BUILD_DIR)/$(MAIN_BINARY)
