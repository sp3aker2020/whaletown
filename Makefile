.PHONY: build install clean test generate

BINARY := wt
BUILD_DIR := .

# Get version info for ldflags
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -X github.com/speaker20/whaletown/internal/cmd.Version=$(VERSION) \
           -X github.com/speaker20/whaletown/internal/cmd.Commit=$(COMMIT) \
           -X github.com/speaker20/whaletown/internal/cmd.BuildTime=$(BUILD_TIME)

generate:
	go generate ./...

build: generate
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY) ./cmd/wt
ifeq ($(shell uname),Darwin)
	@codesign -s - -f $(BUILD_DIR)/$(BINARY) 2>/dev/null || true
	@echo "Signed $(BINARY) for macOS"
endif

install: build
	cp $(BUILD_DIR)/$(BINARY) ~/.local/bin/$(BINARY)
ifeq ($(shell uname),Darwin)
	@codesign -s - -f ~/.local/bin/$(BINARY) 2>/dev/null || true
endif

clean:
	rm -f $(BUILD_DIR)/$(BINARY)

test:
	go test ./...
