# Copyright 2021-2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This Makefile exists because the root TCE Makefile calls each package's unit
# tests via "make -C path/to/package/dir test".

GO_TEST_FLAGS ?=

# Inspired by Cluster-API Makefile
# Any target that has '## ' append to it will be included in the help message
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[0-9A-Za-z_-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# We pass -count 1 to the unit tests so that the go test cache is ignored. Since
# the tests reference files outside of this package (i.e., they reference test
# fixtures in other packages), then we could end up with false positive test runs
# because our test fixtures is not part of the go test cache key.
test: ## Run the unit tests
	go test -count 1 $(GO_TEST_FLAGS) .

test-v0121: ## Run the unit tests for just v0.12.1
	go test -count 1 -run TestTemplate/0.12.1 $(GO_TEST_FLAGS) .

test-v012: ## Run the unit tests for just v0.12.0
	go test -count 1 -run TestTemplate/0.12.0 $(GO_TEST_FLAGS) .

test-v044: ## Run the unit tests for just v0.4.4
	go test -count 1 -run TestTemplate/0.4.4 $(GO_TEST_FLAGS) .

test-v043: ## Run the unit tests for just v0.4.3
	go test -count 1 -run TestTemplate/0.4.3 $(GO_TEST_FLAGS) .

update: ## Update all expected test fixtures with their actual value
	@echo "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"
	@echo "  WARNING: this is janky!"
	@echo "  WARNING: this might break if we change test code!"
	@echo "  WARNING: You should manually verify that it worked after running!"
	@echo "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"
	make test \
    | awk 'BEGIN { FS = "to update expected:" } /to update expected:/ {printf("%s\n", $$2)}' \
    | while read line; do echo "running: $$line"; $(SHELL) -c "$$line"; done

get-deps: ## Get all dependencies
	go mod tidy

lint: ## Run Go code linting
ifeq ($(origin GOLANGCI_LINT),undefined)
	@echo "Error! GOLANGCI_LINT env var not set"
else
	$(GOLANGCI_LINT) run -v --timeout=5m
endif
