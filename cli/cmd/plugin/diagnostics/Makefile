# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

ROOT_DIR := $(shell git rev-parse --show-toplevel)

.DEFAULT_GOAL:=help

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[0-9A-Za-z_-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

get-deps: ## Download the necessary dependencies from go mod
	go mod tidy

lint: ## Run Go code linting
ifeq ($(origin GOLANGCI_LINT),undefined)
	@echo "Error! GOLANGCI_LINT env var not set"
else
	$(GOLANGCI_LINT) run -v --timeout=5m
endif

test: ## Run unit testing suite
	echo "N/A: No unit tests for hack/packages"

.PHONY: e2e-test
e2e-test: ## Run e2e testing suite
	cd e2e-test && BUILD_VERSION=$(BUILD_VERSION) ROOT_DIR=$(ROOT_DIR) go test -tags=e2e -timeout 180m

build: ## Build the executable
	echo "N/A: Implement building"
