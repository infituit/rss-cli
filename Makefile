.PHONY: build platform-builds build-matrix npm-packages compress-binary check-upx fmt lint test integration-test clean pre

GOLANGCI_LINT_TIMEOUT ?= 2m

BUILDTIME := $(shell date +%Y-%m-%dT%T%z)
GITTAG := $(shell git describe --tags --always --dirty 2>/dev/null || printf 'v0.0.0')
GITHASH := $(shell git rev-parse --short HEAD 2>/dev/null || printf 'unknown')
VERSION ?= $(GITTAG)-$(shell date +%y.%m.%d)
ROOT_DIR := $(CURDIR)
OUTPUT_DIR := $(ROOT_DIR)/build/$(VERSION)
APP_NAME := rss-cli
OSES := linux darwin windows
ARCHES := amd64 arm64
UPX_ENABLED ?=
UPX_LEVEL ?= 9
UPX_ARGS ?= -$(UPX_LEVEL)
UPX_SUPPORTED_PLATFORMS := linux/amd64 linux/arm64 windows/amd64

VERSION_LDFLAGS := -X github.com/infituit/rss-cli/pkg/version.VERSION=$(VERSION) \
	-X github.com/infituit/rss-cli/pkg/version.BUILDTIME=$(BUILDTIME) \
	-X github.com/infituit/rss-cli/pkg/version.GITHASH=$(GITHASH)

build: | pre
	CGO_ENABLED=0 go build -trimpath -ldflags "$(VERSION_LDFLAGS)" -o $(OUTPUT_DIR)/$(APP_NAME) ./cmd/rss-cli
	@if [ -n "$(UPX_ENABLED)" ]; then \
		$(MAKE) compress-binary BINARY=$(OUTPUT_DIR)/$(APP_NAME); \
	fi

platform-builds: | pre
	@set -e; \
	for os in $(OSES); do \
		for arch in $(ARCHES); do \
			ext=""; \
			if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
			binary="$(OUTPUT_DIR)/$(APP_NAME)_$${os}_$${arch}$$ext"; \
			echo "Building $(APP_NAME) $(VERSION) for $$os/$$arch..."; \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -trimpath -ldflags "$(VERSION_LDFLAGS)" -o "$$binary" ./cmd/rss-cli; \
			echo "Built: $$binary"; \
			if [ -n "$(UPX_ENABLED)" ]; then \
				$(MAKE) compress-binary BINARY="$$binary" PLATFORM=$$os/$$arch; \
			fi; \
		done; \
	done

build-matrix: | pre
	@if [ -z "$(GOOS)" ] || [ -z "$(GOARCH)" ]; then \
		echo "Error: GOOS and GOARCH must be specified."; \
		echo "Example: make build-matrix GOOS=linux GOARCH=amd64"; \
		exit 1; \
	fi
	@ext=""; \
	if [ "$(GOOS)" = "windows" ]; then ext=".exe"; fi; \
	binary="$(OUTPUT_DIR)/$(APP_NAME)_$(GOOS)_$(GOARCH)$$ext"; \
	echo "Building $(APP_NAME) $(VERSION) for $(GOOS)/$(GOARCH)..."; \
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -trimpath -ldflags "$(VERSION_LDFLAGS)" -o "$$binary" ./cmd/rss-cli; \
	echo "Built: $$binary"; \
	if [ -n "$(UPX_ENABLED)" ]; then \
		$(MAKE) compress-binary BINARY="$$binary" PLATFORM=$(GOOS)/$(GOARCH); \
	fi

npm-packages:
	node npm/prepare.cjs --validate "$(VERSION)"
	$(MAKE) platform-builds VERSION="$(VERSION)"
	node npm/prepare.cjs "$(VERSION)"

pre: check-upx
	@mkdir -p $(OUTPUT_DIR)

check-upx:
	@if [ -n "$(UPX_ENABLED)" ] && ! command -v upx >/dev/null 2>&1; then \
		echo "Error: UPX compression requested but upx was not found in PATH"; \
		echo "Run 'mise install' or set UPX_ENABLED="; \
		exit 1; \
	fi

compress-binary:
	@if [ -z "$(BINARY)" ]; then \
		echo "Error: BINARY path not specified"; \
		echo "Usage: make compress-binary BINARY=path/to/binary"; \
		exit 1; \
	fi
	@if [ ! -f "$(BINARY)" ]; then \
		echo "Error: Binary file not found: $(BINARY)"; \
		exit 1; \
	fi
	@if [ -n "$(PLATFORM)" ] && ! echo " $(UPX_SUPPORTED_PLATFORMS) " | grep -q " $(PLATFORM) "; then \
		echo "Skipping UPX compression for unsupported platform $(PLATFORM)"; \
	else \
		echo "Compressing $(BINARY) with UPX $(UPX_ARGS)..."; \
		upx $(UPX_ARGS) "$(BINARY)"; \
		echo "Compression complete"; \
	fi

clean:
	rm -rf $(ROOT_DIR)/build

fmt:
	gofmt -w ./cmd ./internal ./pkg ./testsuite

lint:
	golangci-lint run --config .golangci.yml --timeout $(GOLANGCI_LINT_TIMEOUT)

test:
	go test ./...

integration-test:
	go test -tags=integration ./testsuite/...
