# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

ROOT_DIR := $(shell git rev-parse --show-toplevel)
VERSION ?= 2.6.1
SRIOV_CNI_TEST_ROOT := "${ROOT_DIR}"/addons/packages/sriov-cni/"${VERSION}"/test/unit

.DEFAULT_GOAL := help

.PHONY: help
help: ## print help messages
	@awk 'BEGIN {FS = ":.*##"; printf "Useful Targets:\n"} /^[a-z]+/ { printf "	\033[0;34m%-8s\033[0m	%s\n", $$1, $$2} '  ${MAKEFILE_LIST}

# tests for sriov-cni

.PHONY: test
unittest: ## unit tests
	@printf "\nunit tests for sriov-cni.${VERSION}\n\n"; \
	cd "${SRIOV_CNI_TEST_ROOT}" && ginkgo --flakeAttempts=2 -v .

e2e-test: ## e2e tests
	@printf "\ne2e tests for sriov-cni.${VERSION}\n\n"; \
	echo "No E2E, because SRIOV needs a environment with SR-IOV enabled"

lint: ## lint check for tests files
	@printf "\nlint check for tests files\n\n"; \
	echo "TODO: Currently no lint"

get-deps: ## get go sources dependencies.
	@printf "\nget go sources dependencies\n\n"; \
	go mod download

build: ## build
	@echo "TODO: No build steps now"
