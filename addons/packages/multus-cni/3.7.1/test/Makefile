# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

ROOT_DIR := $(shell git rev-parse --show-toplevel)
VERSION ?= 3.7.1
MULTUS_E2E_TEST_ROOT := "${ROOT_DIR}"/addons/packages/multus-cni/"${VERSION}"/test/e2e

.DEFAULT_GOAL := help

.PHONY: help
help: ## print help messages
	@awk 'BEGIN {FS = ":.*##"; printf "Useful Targets:\n"} /^[a-z]+/ { printf "	\033[0;34m%-8s\033[0m	%s\n", $$1, $$2} '  ${MAKEFILE_LIST}

# tests for multus-cni

.PHONY: test
test: ## unit tests
	@echo "TODO: No unit tests for ${VERSION} now"

e2e-test: ## e2e tests
	@printf "\ne2e tests for multus.${VERSION}\n\n"; \
	cd "${MULTUS_E2E_TEST_ROOT}" && ginkgo . -- --version ${VERSION}

lint: ## lint check for tests files
	@printf "\nlint check for tests files\n\n"; \
	echo "TODO: Currently no lint"

get-deps: ## get go sources dependencies.
	@printf "\nget go sources dependencies\n\n"; \
	go mod download

build: ## build
	@echo "TODO: No build steps now"
