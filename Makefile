# Copyright 2019 The Caicloud Authors.
#
# The old school Makefile, following are required targets. The Makefile is written
# to allow building multiple binaries. You are free to add more targets or change
# existing implementations, as long as the semantics are preserved.
#
#   make              - default to 'build' target
#   make lint         - code analysis with golangci-lint
#   make test         - run unit test
#   make build        - alias for `build-local` target
#   make build-local  - build local binary
#   make build-linux  - build amd64 Linux binary
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


# It's necessary to set the errexit flags for the bash shell.
export SHELLOPTS := errexit

# This will force go to use the vendor files instead of using the `$GOPATH/pkg/mod`. (vendor mode)
# more info: https://github.com/golang/go/wiki/Modules#how-do-i-use-vendoring-with-modules-is-vendoring-going-away
export GOFLAGS := -mod=vendor

# this is not a public registry; change it to your own
REGISTRY ?= caicloud/
BASE_REGISTRY ?=

ARCH ?=
GO_VERSION ?= 1.13

CPUS ?= $(shell /bin/bash hack/read_cpus_available.sh)

# Track code version with Docker Label.
DOCKER_LABELS ?= git-describe="$(shell date -u +v%Y%m%d)-$(shell git describe --tags --always --dirty)"

GOPATH ?= $(shell go env GOPATH)
BIN_DIR := $(GOPATH)/bin
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint
CMD_DIR := ./cmd
OUTPUT_DIR := ./bin
BUILD_DIR := ./build
.PHONY: lint test build build-local build-linux container push clean

build: build-local

build-local: clean
	@echo ">> building binaries"``
	@GOOS=$(shell uname -s | tr A-Z a-z) GOARCH=$(ARCH) CGO_ENABLED=0	 \
	go build -i -v -o $(OUTPUT_DIR)/event_exporter -p $(CPUS)			\
		 -ldflags "-s -w 										        \
	  -X $(ROOT)/pkg/version.Version=${VERSION}							 \
	  -X $(ROOT)/pkg/version.Branch=${BRANCH}							 \
	  -X $(ROOT)/pkg/version.Commit=${COMMIT}							 \
	  -X $(ROOT)/pkg/version.BuildDate=${BUILD_DATE}"					 \
	 $(CMD_DIR)/

build-linux:
	@docker run --rm -t                                                                \
	  -v $(PWD):/go/src/$(ROOT)                                                        \
	  -w /go/src/$(ROOT)                                                               \
	  -e GOOS=linux                                                                    \
	  -e GOARCH=amd64                                                                  \
	  -e GOPATH=/go                                                                    \
	  -e CGO_ENABLED=0																              \
	  -e GOFLAGS=$(GOFLAGS)   	                                                       \
	  -e SHELLOPTS=$(SHELLOPTS)                                                        \
	  $(BASE_REGISTRY)golang:$(GO_VERSION)                                            \
	    /bin/bash -c '                                    								\
	      	go build -i -v -o $(OUTPUT_DIR)/event_exporter -p $(CPUS)			\
          		 -ldflags "-s -w 										        \
          	  -X $(ROOT)/pkg/version.Version=${VERSION}							 \
          	  -X $(ROOT)/pkg/version.Branch=${BRANCH}							 \
          	  -X $(ROOT)/pkg/version.Commit=${COMMIT}							 \
          	  -X $(ROOT)/pkg/version.BuildDate=${BUILD_DATE}"					 \
          	 $(CMD_DIR)/'                                                    			\

container: build-linux
	@echo ">> building image"
	@docker build -t $(REGISTRY)event-exporter:$(VERSION) --label $(DOCKER_LABELS)  -f $(BUILD_DIR)/Dockerfile .

push: container
	@echo ">> pushing image"
	@docker push $(REGISTRY)event-exporter:$(VERSION)

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
	@-rm -vrf ${OUTPUT_DIR}
	@rm -f coverage.out
