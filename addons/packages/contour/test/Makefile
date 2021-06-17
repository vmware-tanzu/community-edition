# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

.DEFAULT_GOAL:=help

REPO_ROOT := $(shell git rev-parse --show-toplevel)

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

## --------------------------------------
## Testing
## --------------------------------------

.EXPORT_ALL_VARIABLES:
TOOLS_DIR ?= $(REPO_ROOT)/hack/tools
TOOLS_BIN_DIR ?= $(TOOLS_DIR)/bin
GINKGO := $(TOOLS_BIN_DIR)/ginkgo
CONTOUR_INTEGRATION_TESTER := $(TOOLS_BIN_DIR)/integration-tester
KIND := $(TOOLS_BIN_DIR)/kind
KUBECTL := $(TOOLS_BIN_DIR)/kubectl
TEST_TAGS ?= e2e
GINKGO_OPTIONS ?= $(GINKGO_FAILFAST)  # Other needed ginkgo options to be appended to the command.
CONFIG_YAML ?= config/e2e_config.yaml

GINKGO_VERSION=v1.16.4
CONTOUR_INTEGRATION_TESTER_VERSION=7.0.0
KIND_VERSION=0.11.1
KUBECTL_VERSION=1.21.1

export PATH := $(abspath $(TOOLS_BIN_DIR)):$(PATH)

# Host information.
OS ?= $(shell uname -s)
HOST_OS ?= $(shell go env GOOS)
HOST_ARCH ?= $(shell go env GOARCH)

## --------------------------------------
## Ginkgo
## --------------------------------------
ginkgo: $(GINKGO) ## Download ginkgo
$(GINKGO):
	GOBIN=$(TOOLS_BIN_DIR) go get -u github.com/onsi/ginkgo/ginkgo@$(GINKGO_VERSION)

## --------------------------------------
## Contour Integration Tester
## --------------------------------------
contour-integration-tester: $(CONTOUR_INTEGRATION_TESTER) ## Download contour integration tester
$(CONTOUR_INTEGRATION_TESTER):
	mkdir -p "$(TOOLS_DIR)/integration-tester" && \
	curl -sL "https://github.com/projectcontour/integration-tester/releases/download/v$(CONTOUR_INTEGRATION_TESTER_VERSION)/integration-tester_$(CONTOUR_INTEGRATION_TESTER_VERSION)_$(OS)_x86_64.tar.gz" -o "$(TOOLS_DIR)/integration-tester/integration-tester_$(CONTOUR_INTEGRATION_TESTER_VERSION)_$(OS)_x86_64.tar.gz" && \
	tar -xf "$(TOOLS_DIR)/integration-tester/integration-tester_$(CONTOUR_INTEGRATION_TESTER_VERSION)_$(OS)_x86_64.tar.gz" -C "$(TOOLS_DIR)/integration-tester" && \
	mv "$(TOOLS_DIR)/integration-tester/integration-tester" $(@) && \
	chmod +x $(@) && \
	rm -rf "$(TOOLS_DIR)/integration-tester"

## --------------------------------------
## Kind
## --------------------------------------
kind: $(KIND) ## Download kind
$(KIND):
	curl -sL https://github.com/kubernetes-sigs/kind/releases/download/v$(KIND_VERSION)/kind-$(HOST_OS)-$(HOST_ARCH) -o $(@) && \
	chmod +x $(@)

## --------------------------------------
## Kubectl
## --------------------------------------
kubectl: $(KUBECTL) ## Download kind
$(KUBECTL):
	curl -sL https://dl.k8s.io/release/v$(KUBECTL_VERSION)/bin/$(HOST_OS)/$(HOST_ARCH)/kubectl -o $(@) && \
    chmod a+x $(@)

.PHONY: e2e-test
e2e-test: $(GINKGO) $(CONTOUR_INTEGRATION_TESTER) ## Run the end-to-end tests
	git config --global url."git@github.com:".insteadOf https://github.com/
	cd e2e; $(GINKGO) $(GINKGO_OPTIONS) -tags="$(TEST_TAGS)" -timeout 1h  -race -v . -- \
	    -e2e.config="$(CONFIG_YAML)" -package-repo-imgpkg-bundle="$(PACKAGE_REPO_IMGPKG_BUNDLE)" \
	    -iaas="$(IAAS)" -cluster-type="$(CLUSTER_TYPE)" .

.PHONY: prepare-kind-cluster
prepare-kind-cluster: $(KIND) $(KUBECTL) ## Create kind cluster
	./e2e/scripts/kind-cluster.sh prepare_cluster

.PHONY: delete-kind-cluster
delete-kind-cluster: $(KIND) ## Delete kind cluster
	./e2e/scripts/kind-cluster.sh delete_cluster

.PHONY: kind-e2e-test
kind-e2e-test: IAAS=kind
kind-e2e-test: prepare-kind-cluster e2e-test ## Run the end-to-end tests on kind

.PHONY: clean
clean: delete-kind-cluster ## clean