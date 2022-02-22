# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# Build tooling that is used for the main project

# Ensure Make is run with bash shell as some syntax below is bash-specific
SHELL := /usr/bin/env bash

.DEFAULT_GOAL := help

# Use GOPROXY environment variable if set
GOPROXY := $(shell go env GOPROXY)
ifeq (,$(strip $(GOPROXY)))
GOPROXY := https://proxy.golang.org
endif
export GOPROXY

# Active module mode, as we use go modules to manage dependencies
export GO111MODULE := on

# Directories.
SRCS := go.mod go.sum
GOBINDATA := $(BIN_DIR)/gobindata

# Host information.
HOST_OS=$(shell go env GOOS)
HOST_ARCH=$(shell go env GOARCH)


lint:
	echo "N/A: Skipping tools linting"

get-deps: ## Get all dependencies
	go mod tidy

test:
	echo "N/A: No unit tests for hack/tools"

e2e-test:
	echo "N/A: No e2e tests for hack/tools"

build:
	echo "N/A: No binary output for hack/tools"

## --------------------------------------
## Help
## --------------------------------------

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

## --------------------------------------
## Binaries
## --------------------------------------

golangci-lint:
	go build -tags=tools -o $(GOLANGCI_LINT) github.com/golangci/golangci-lint/cmd/golangci-lint

gobindata: $(GOBINDATA)
$(GOBINDATA): go.mod # Build go-bindata
	mkdir -p $(BIN_DIR)
	go build -tags=tools -o $(BIN_DIR) github.com/shuLhan/go-bindata/...; mv $(BIN_DIR)/go-bindata $(GOBINDATA)

gobindata: $(GOBINDATA)
$(GOBINDATA): go.mod # Build go-bindata
	mkdir -p $(BIN_DIR)
	go build -tags=tools -o $(BIN_DIR) github.com/shuLhan/go-bindata/...; mv $(BIN_DIR)/go-bindata $(GOBINDATA)

## --------------------------------------
## Cleanup
## --------------------------------------

.PHONY: clean
clean: ## Remove all generated binaries
	rm -rf bin
