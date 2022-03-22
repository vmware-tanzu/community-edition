# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

ROOT_DIR := $(shell git rev-parse --show-toplevel)
VERSION ?= 3.3.2
SRIOV_DP_TEST_ROOT := "${ROOT_DIR}"/addons/packages/sriov-network-device-plugin/"${VERSION}"/test/unit

.DEFAULT_GOAL := help

.PHONY: help
help: ## print help messages
	@awk 'BEGIN {FS = ":.*##"; printf "Useful Targets:\n"} /^[a-z]+/ { printf "	\033[0;34m%-8s\033[0m	%s\n", $$1, $$2} '  ${MAKEFILE_LIST}

# tests for sriov-network-device-plugin

.PHONY: test
unittest: ## unit tests
	@printf "\nunit tests for sriov-network-device-plugin.${VERSION}\n\n"; \
	cd "${SRIOV_DP_TEST_ROOT}" && ginkgo --flakeAttempts=2 -v .

e2e-test: ## e2e tests
	@printf "\ne2e tests for sriov-network-device-plugin.${VERSION}\n\n"; \
	echo "No E2E, because SRIOV needs a environment with SR-IOV enabled"

all-test: unittest e2e-test ## run all tests

lint: ## lint check for tests files
ifeq ($(origin GOLANGCI_LINT),undefined)
	@echo "Error! GOLANGCI_LINT env var not set"
else
	$(GOLANGCI_LINT) run -v --timeout=5m
endif

get-deps: ## get go sources dependencies.
	go mod download

build: ## build
	@echo "TODO: No build steps now"
