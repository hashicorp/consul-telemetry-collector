SHELL := /usr/bin/env bash -euo pipefail -c

REPO_NAME    ?= $(shell basename "$(CURDIR)")
PRODUCT_NAME ?= $(REPO_NAME)
BIN_NAME     ?= $(PRODUCT_NAME)

# Get local ARCH; on Intel Mac, 'uname -m' returns x86_64 which we turn into amd64.
# Not using 'go env GOOS/GOARCH' here so 'make docker' will work without local Go install.
ARCH     = $(shell A=$$(uname -m); [ $$A = x86_64 ] && A=amd64; echo $$A)
OS       = $(shell uname | tr [[:upper:]] [[:lower:]])
PLATFORM = $(OS)/$(ARCH)
DIST     = dist/$(PLATFORM)
BIN      = $(DIST)/$(BIN_NAME)

VERSION = $(shell ./build-scripts/version.sh pkg/version/version.go)

GIT_COMMIT?=$(shell git rev-parse --short HEAD)
GIT_DIRTY?=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GOLDFLAGS=-X github.com/hashicorp/consul-telemetry-collector/pkg/version.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)

# Get latest revision (no dirty check for now).
REVISION = $(shell git rev-parse HEAD)
GOLANGCI_CONFIG_DIR ?= $(CURDIR)

.PHONY: goversion
goversion:
	@go version

.PHONY: version
version:
	@echo $(VERSION)

dist:
	mkdir -p $(DIST)
	echo '*' > dist/.gitignore

.PHONY: bin
bin: goversion dist
	GOARCH=$(ARCH) GOOS=$(OS) CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags="$(GOLDFLAGS)" -o $(BIN) ./cmd/$(BIN_NAME)

.PHONY: dev
dev: bin
	cp $(BIN) $(GOBIN)/$(BIN_NAME)

.PHONY: tests
tests: goversion
	go test -timeout 10s ./...

.PHONY: lint
lint:
	golangci-lint run --config $(GOLANGCI_CONFIG_DIR)/.golangci.yml

.PHONY: deps
deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1


