.PHONY: all build clean run test package

MAIN_BINARY := GoNetPad
BUILD_DIR := GoNetPad
CMD_DIR := cmd

GOOS ?= darwin
GOARCH ?= amd64
CGO_ENABLED ?= 1
LDFLAGS = -w -s

all: build

build:
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags="$(LDFLAGS)" \
		-trimpath \
		-o $(BUILD_DIR)/$(MAIN_BINARY) \
		./$(CMD_DIR)/GoNetPad/main.go

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(MAIN_BINARY).app

run: build
	@echo "Running application..."
	./$(BUILD_DIR)/$(MAIN_BINARY)

package:
	@command -v fyne >/dev/null 2>&1 || { echo >&2 "Fyne CLI not found. Please install using: go install fyne.io/fyne/v2/cmd/fyne@latest"; exit 1; }
	@echo "Packaging the application as a macOS app..."
	@fyne package -os darwin -name $(MAIN_BINARY) -executable $(BUILD_DIR)/$(MAIN_BINARY)
	@echo "Packaged app created ($(MAIN_BINARY).app)"