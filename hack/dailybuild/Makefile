# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

.DEFAULT_GOAL:=help

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[0-9A-Za-z_-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

lint: ## Run Go code linting. Requires the $GOLANGCI_LINT env var to be set.
ifeq ($(origin GOLANGCI_LINT),undefined)
	@echo "Error! GOLANGCI_LINT env var not set"
else
	$(GOLANGCI_LINT) run -v --timeout=5m
endif

get-deps: ## Get all go dependencies
	go mod tidy

test:
	@echo "N/A: No unit tests for hack/dailybuild"

e2e-test:
	@echo "N/A: No e2e tests for hack/dailybuild"

build: ## Build the executable
	go build -o dailybuild dailybuild.go

run: ## Run the dailybuild program to upload dailybuild
	@echo "N/A: No run target for hack/dailybuild"

