# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

ROOT_DIR := $(shell git rev-parse --show-toplevel)
VERSION ?= 0.5.1
USE_CONF_FILE ?= true
WHEREABOUTS_E2E_TEST_ROOT := "${ROOT_DIR}"/addons/packages/whereabouts/"${VERSION}"/test/e2e

.DEFAULT_GOAL := help

.PHONY: help
help: ## print help messages
	@awk 'BEGIN {FS = ":.*##"; printf "Useful Targets:\n"} /^[a-z]+/ { printf "	\033[0;34m%-8s\033[0m	%s\n", $$1, $$2} '  ${MAKEFILE_LIST}

# tests for whereabouts

.PHONY: test
test: ## unit tests
	@echo "TODO: No unit tests for ${VERSION} now"

e2e-test: ## e2e tests
	@printf "\ne2e tests for whereabouts.${VERSION}\n\n"; \
	cd "${WHEREABOUTS_E2E_TEST_ROOT}" && ginkgo --flakeAttempts=2 -v . -- --version=${VERSION} --whereabouts-use-conf=${USE_CONF_FILE}

lint: ## lint check for tests files
ifeq ($(origin GOLANGCI_LINT),undefined)
	@echo "Error! GOLANGCI_LINT env var not set"
else
	$(GOLANGCI_LINT) run -v --timeout=5m
endif

get-deps: ## get go sources dependencies.
	go mod download

build: ## build
	@echo "TODO: No build steps"
