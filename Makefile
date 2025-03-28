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
BIN_PATH ?= dist/$(PLATFORM)/$(BIN_NAME)

GO_MODULE_DIRS ?= $(shell go list -m -f "{{ .Dir }}" | grep -v mod-vendor)

.PHONY: version
version:
	@cat internal/version/VERSION

.PHONY: dev
dev:
	CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags="-X github.com/hashicorp/consul-telemetery-collector/internal/version.GitCommit=${GITHUB_SHA::8}" -o $(BIN_NAME) ./cmd/$(BIN_NAME)

.PHONY: build
build:
	CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags="-X github.com/hashicorp/consul-telemetery-collector/internal/version.GitCommit=${GITHUB_SHA::8}" -o $(BIN_PATH) ./cmd/$(BIN_NAME)

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

build/docker:
	DOCKER_BUILDKIT=1 docker build -t consul-telemetry-collector --build-arg BIN_NAME=consul-telemetry-collector .

.PHONY: deps
deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.0
