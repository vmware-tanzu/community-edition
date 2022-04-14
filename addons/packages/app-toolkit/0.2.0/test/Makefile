# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

.DEFAULT_GOAL:=help

help: ## Display this help message
	# Inspired by Cluster-API Makefile
	# Any target that has '## ' append to it will be included in the help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[0-9A-Za-z_-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

test: ## Run unit testing suite
	@echo "TODO: implement running unit tests"

e2e-test: ## Run e2e testing suite
	./app-toolkit-test.sh

build: ## Build the executable
	@echo "TODO: implement building"
