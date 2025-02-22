.PHONY: all build clean run

MAIN_BINARY := GoNetPad
BUILD_DIR := GoNetPad
CMD_DIR := cmd

GOOS ?= darwin
GOARCH ?= amd64
CGO_ENABLED ?= 1
LDFLAGS = -w -s

all: clean build run

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
	@rm -f $(MAIN_BINARY).app

run: 
	@echo "Running application..."
	./$(BUILD_DIR)/$(MAIN_BINARY)
