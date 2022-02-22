# Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# Add tooling binaries here and in hack/tools/Makefile
TOOLS_DIR := ../../../../hack/tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin
GOLANGCI_LINT := $(TOOLS_BIN_DIR)/golangci-lint

# Directories
UI_DIR := web/tanzu-ui
GOBINDATA := $(TOOLS_BIN_DIR)/gobindata

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
	echo "N/A: No unit tests for ui"

e2e-test: ## Run e2e testing suite
	echo "N/A: No e2e tests for ui"

build: ## Build the executable
	echo "N/A: Implement building"

ui-dependencies: ## install UI dependencies (node modules)
	cd $(UI_DIR); npm ci

ui-build: ui-dependencies ## Install dependencies, then compile client UI for production
	cd $(UI_DIR); npm run build:prod
	$(MAKE) generate-ui-bindata

generate-ui-bindata: $(GOBINDATA) ## Generate go-bindata for ui files
	$(GOBINDATA) -mode=420 -modtime=1 -o=web/zz_generated.bindata.go -pkg=main $(UI_DIR)/build/...