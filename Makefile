SHELL := /usr/bin/env bash -euo pipefail -c

REPO_NAME    ?= $(shell basename "$(CURDIR)")
PRODUCT_NAME ?= $(REPO_NAME)
BIN_NAME     ?= $(PRODUCT_NAME)
GOPATH       ?= $(shell go env GOPATH)
GOBIN        ?= $(GOPATH)/bin

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

.PHONY: version
version:
	@echo $(VERSION)

dist:
	mkdir -p $(DIST)
	echo '*' > dist/.gitignore

.PHONY: bin
bin: dist
	GOARCH=$(ARCH) GOOS=$(OS) CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags="$(GOLDFLAGS)" -o $(BIN) ./cmd/$(BIN_NAME)

.PHONY: dev
dev: bin
	cp $(BIN) $(GOBIN)/$(BIN_NAME)

# Docker Stuff.
export DOCKER_BUILDKIT=1
BUILD_ARGS = BIN_NAME=$(BIN_NAME) PRODUCT_VERSION=$(VERSION) PRODUCT_REVISION=$(REVISION)
TAG        = $(PRODUCT_NAME)/$(TARGET):$(VERSION)
BA_FLAGS   = $(addprefix --build-arg=,$(BUILD_ARGS))
FLAGS      = --target $(TARGET) --platform $(PLATFORM) --tag $(TAG) $(BA_FLAGS)

# Set OS to linux for all docker/* targets.
docker/%: OS = linux

# DOCKER_TARGET is a macro that generates the build and run make targets
# for a given Dockerfile target.
# Args: 1) Dockerfile target name (required).
#       2) Build prerequisites (optional).
define DOCKER_TARGET
.PHONY: docker/$(1)
docker/$(1): TARGET=$(1)
docker/$(1): $(2)
	docker build $$(FLAGS) .
	@echo 'Image built; run "docker run --rm $$(TAG)" to try it out.'

.PHONY: docker/$(1)/run
docker/$(1)/run: TARGET=$(1)
docker/$(1)/run: docker/$(1)
	docker run --rm $$(TAG)
endef

# Create docker/<target>[/run] targets.
$(eval $(call DOCKER_TARGET,release-default,bin))
$(eval $(call DOCKER_TARGET,release-ubi,bin))

.PHONY: docker
docker: docker/release-default

.PHONY: tests
tests:
	go test ./...

.PHONY: lint
lint:
	@PATH=$$PATH:$(GOBIN) golangci-lint run

.PHONY: deps
deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1

.PHONY: changelog
changelog:
ifdef LAST_RELEASE_GIT_TAG
	@changelog-build \
		-last-release $(LAST_RELEASE_GIT_TAG) \
		-entries-dir .changelog/ \
		-changelog-template .changelog/changelog.tmpl \
		-note-template .changelog/note.tmpl \
		-this-release $(REVISION)
else
	$(error Cannot generate changelog without LAST_RELEASE_GIT_TAG)
endif

INTEGRATION_TESTS_SERVER_IMAGE    ?= hashicorppreview/consul:1.14-dev
INTEGRATION_TESTS_DATAPLANE_IMAGE ?= $(PRODUCT_NAME)/release-default:$(VERSION)

