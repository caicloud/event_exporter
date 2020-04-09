# Copyright 2019 The Caicloud Authors.
#
# The old school Makefile, following are required targets. The Makefile is written
# to allow building multiple binaries. You are free to add more targets or change
# existing implementations, as long as the semantics are preserved.
#
#   make              - default to 'build' target
#   make lint         - code analysis with golangci-lint
#   make test         - run unit test
#   make build        - build binary in a Golang container
#   make build-local  - build local binary
#   make container    - build container
#   make push         - push container
#   make clean        - clean up
#
# The makefile is also responsible to populate project version information.
#

ROOT := github.com/caicloud/event_exporter
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT = $(shell git rev-parse HEAD)
BRANCH = $(shell git branch | grep \* | cut -d ' ' -f2)
BUILD_DATE = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# this is not a public registry; change it to your own
REGISTRY ?= cargo.dev.caicloud.xyz/release

ARCH ?= amd64
GO_VERSION = 1.13

CPUS ?= $(shell /bin/bash hack/read_cpus_available.sh)
GOPATH ?= $(shell go env GOPATH)
BIN_DIR := $(GOPATH)/bin
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint

.PHONY: lint test build build-local container push clean

build:
	@docker run --rm -t													\
	  -v "${PWD}:/go/src/github.com/caicloud/event_exporter"			\
	  -w /go/src/github.com/caicloud/event_exporter						\
	  golang:${GO_VERSION} make build-local

build-local: clean
	@echo ">> building binaries"
	@GOOS=$(shell uname -s | tr A-Z a-z) GOARCH=$(ARCH) CGO_ENABLED=0	\
	go build -mod=vendor -ldflags "-s -w								\
	  -X $(ROOT)/pkg/version.Version=${VERSION}							\
	  -X $(ROOT)/pkg/version.Branch=${BRANCH}							\
	  -X $(ROOT)/pkg/version.Commit=${COMMIT}							\
	  -X $(ROOT)/pkg/version.BuildDate=${BUILD_DATE}"					\
	-o event_exporter

container: build
	@echo ">> building image"
	@docker build -t $(REGISTRY)/event_exporter:$(VERSION) -f ./Dockerfile .

push: container
	@echo ">> pushing image"
	@docker push $(REGISTRY)/event_exporter:$(VERSION)

lint: $(GOLANGCI_LINT)
	@echo ">> running golangci-lint"
	@$(GOLANGCI_LINT) run

$(GOLANGCI_LINT):
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BIN_DIR) v1.23.6

test:
	@echo ">> running tests"
	@go test -p $(CPUS) $$(go list ./... | grep -v /vendor) -coverprofile=coverage.out
	@go tool cover -func coverage.out | tail -n 1 | awk '{ print "Total coverage: " $$3 }'

clean:
	@echo ">> cleaning up"
	@rm -f event_exporter coverage.out
