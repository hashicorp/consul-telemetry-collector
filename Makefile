SHELL := /usr/bin/env bash -euo pipefail -c

BIN_NAME     ?= consul-telemetry-collector
GIT_COMMIT?=$(shell git rev-parse --short HEAD)
GIT_DIRTY?=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GOLDFLAGS=-X github.com/hashicorp/consul-telemetry-collector/internal/version.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)
GOLANGCI_CONFIG_DIR ?= $(CURDIR)
ARCH := $(shell uname -m)
ifeq ($(ARCH),x86_64)
    ARCH := amd64
endif
ifeq ($(ARCH),aarch64)
    ARCH := amd64
endif
OS       = $(shell uname | tr [[:upper:]] [[:lower:]])
PLATFORM = $(OS)/$(ARCH)

GO_MODULE_DIRS ?= $(shell go list -m -f "{{ .Dir }}" | grep -v mod-vendor)

.PHONY: version
version:
	@bin/$(PLATFORM)/$(BIN_NAME) --version

.PHONY: bin
bin:
	@ GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags="$(GOLDFLAGS)" -o bin/linux/amd64/$(BIN_NAME) ./cmd/$(BIN_NAME)
	@ GOARCH=amd64 GOOS=darwin CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags="$(GOLDFLAGS)" -o bin/darwin/x86_64/$(BIN_NAME) ./cmd/$(BIN_NAME)
	@ GOARCH=arm64 GOOS=darwin CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags="$(GOLDFLAGS)" -o bin/darwin/arm64/$(BIN_NAME) ./cmd/$(BIN_NAME)

go/test:
	@ for mod in $(GO_MODULE_DIRS) ; do \
		cd $$mod > /dev/null; \
		echo "testing $$mod"; \
		go test -timeout 10s ./... ;\
		cd - > /dev/null; \
	done

go/mod:
	@ for mod in $(GO_MODULE_DIRS); do \
		(	cd $$mod > /dev/null; \
		echo "go mod tidy $$mod"; \
		go mod tidy; \
		); \
	done

go/lint:
	@ for mod in $(GO_MODULE_DIRS) ; do \
		cd $$mod > /dev/null; \
		echo "linting $$mod"; \
		golangci-lint run --timeout 5m --config $(GOLANGCI_CONFIG_DIR)/.golangci.yml ;\
		cd - > /dev/null; \
	done

go/fix:
	@ for mod in $(GO_MODULE_DIRS) ; do \
		cd $$mod > /dev/null; \
		echo "linting $$mod"; \
		golangci-lint run --timeout 5m --config $(GOLANGCI_CONFIG_DIR)/.golangci.yml --fix ;\
		cd - > /dev/null; \
	done

.PHONY: deps
deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2
