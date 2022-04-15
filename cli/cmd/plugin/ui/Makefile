# Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0


UID := $(shell id -u)
GID := $(shell id -g)
SWAGGER=docker run --rm -it --user $(UID):$(GID) -e GOPATH=$(shell go env GOPATH):/go -v $$HOME:$$HOME -w $(PWD) quay.io/goswagger/swagger:v0.21.0

# Directories
UI_DIR := web/tanzu-ui

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

ui-dependencies: ## Install UI dependencies (node modules)
	cd $(UI_DIR); npm ci

ui-build: ui-dependencies ## Install dependencies, then compile client UI for production
	cd $(UI_DIR); npm run build:prod

ui-build-docker:
	docker run -it --rm \
		-v "$$(pwd)":/usr/src -w /usr/src/$(UI_DIR) \
		--user $(UID):$(GID) \
		node:16 \
		sh -c "npm ci; npm run build:prod"

.PHONY: generate-ui-swagger-api
generate-ui-swagger-api: ## Generate swagger files for UI backend
	rm -fr server/models server/restapi/operations
	${SWAGGER} generate server -A TanzuUI -t server -f api/spec.yaml --exclude-main
	git reset HEAD server/restapi/server.go
	git checkout HEAD server/restapi/server.go
