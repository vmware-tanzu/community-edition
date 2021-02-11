# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

.DEFAULT_GOAL:=help

### TOOLS ###
TOOLS_DIR := hack/tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin

# Add tooling binaries here and in hack/tools/Makefile
GOLANGCI_LINT := $(TOOLS_BIN_DIR)/golangci-lint
GOLINT := $(TOOLS_BIN_DIR)/golint
TOOLING_BINARIES := $(GOLANGCI_LINT) $(GOLINT)

##### BUILD #####
ifndef ${BUILD_VERSION}
BUILD_VERSION ?= $$(git describe --tags --dirty=-dev --abbrev=0)
endif
BUILD_SHA ?= $$(git rev-parse --short HEAD)
BUILD_DATE ?= $$(date -u +"%Y-%m-%d")
CONFIG_VERSION ?= $$(echo "$(BUILD_VERSION)" | cut -d "-" -f1)

build_OS := $(shell uname 2>/dev/null || echo Unknown)

LD_FLAGS = -X 'github.com/vmware-tanzu-private/core/pkg/v1/cli.BuildDate=$(BUILD_DATE)'
LD_FLAGS += -X 'github.com/vmware-tanzu-private/core/pkg/v1/cli.BuildSHA=$(BUILD_SHA)'
LD_FLAGS += -X 'github.com/vmware-tanzu-private/core/pkg/v1/cli.BuildVersion=$(BUILD_VERSION)'

ARTIFACTS_CORE_DIR ?= ./artifacts-core
ARTIFACTS_DIR ?= ./artifacts

ifeq ($(build_OS), Linux)
XDG_DATA_HOME := ${HOME}/.local/share
endif
ifeq ($(build_OS), Darwin)
XDG_DATA_HOME := ${HOME}/Library/ApplicationSupport
endif

export XDG_DATA_HOME

PRIVATE_REPOS="github.com/vmware-tanzu-private,github.com/vmware-tanzu,github.com/dvonthenen"
GO := GOPRIVATE=${PRIVATE_REPOS} go
##### BUILD #####

##### IMAGE #####
OCI_REGISTRY := projects.registry.vmware.com/tce
EXTENSION_NAMESPACE := tanzu-extensions
##### IMAGE #####

##### COMMON TARGETS #####
help: ## display help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

fmt: ## Run go fmt
	$(GO) fmt ./...

vet: ## Run go vet
	$(GO) vet ./...

lint: tools ## Run linting checks
	$(GOLANGCI_LINT) run -v
	$(GOLINT) -set_exit_status ./...

test: ## Run go tests
	$(GO) test ./...

##### COMMON TARGETS #####

##### BUILD TARGETS #####
build: build-plugin

build-plugin: clean-plugin prep-build-cli build-cli-plugins install-cli-plugins copy-release tag-release

clean: clean-plugin

tools: $(TOOLING_BINARIES) ## Build tooling binaries
.PHONY: $(TOOLING_BINARIES)
$(TOOLING_BINARIES):
	make -C $(TOOLS_DIR) $(@F)

# RELEASE MANAGEMENT
PHONY: copy-release
copy-release:
	mkdir -p ${XDG_DATA_HOME}/tanzu-repository
	cp -f ./hack/config.yaml ${XDG_DATA_HOME}/tanzu-repository/config.yaml

.PHONY: tag-release
tag-release:
	@echo "BUILD_VERSION:" ${BUILD_VERSION}
	@echo "CONFIG_VERSION:" ${CONFIG_VERSION}
ifeq ($(shell expr $(BUILD_VERSION)), $(shell expr $(CONFIG_VERSION)))
	sed -i "s/version: latest/version: $(CONFIG_VERSION)/g" ${XDG_DATA_HOME}/tanzu-repository/config.yaml
endif
# RELEASE MANAGEMENT

# TANZU CLI
# .PHONY: build-cli
# build-cli:
# 	$(GO) run github.com/vmware-tanzu-private/core/cmd/cli/plugin-admin/builder cli compile --version $(BUILD_VERSION) --ldflags "$(LD_FLAGS)" --corepath ../../vmware-tanzu-private/core/cmd/cli/tanzu --path ../../vmware-tanzu-private/core/cmd/cli/plugin --artifacts ${ARTIFACTS_CORE_DIR}
# 	$(GO) run github.com/vmware-tanzu-private/core/cmd/cli/plugin-admin/builder cli compile --version $(BUILD_VERSION) --ldflags "$(LD_FLAGS)" --path ../../vmware-tanzu-private/core/cmd/cli/plugin-admin --artifacts ${ARTIFACTS_CORE_DIR}-admin

# .PHONY: install-cli
# install-cli:
# 	sudo cp -f ${ARTIFACTS_CORE_DIR}/core/latest/tanzu-core-linux_amd64 /usr/local/bin/tanzu
# 	TANZU_CLI_NO_INIT=true $(GO) run -ldflags "$(LD_FLAGS)" ../../vmware-tanzu-private/core/cmd/cli/tanzu \
# 		plugin install all --local $(ARTIFACTS_CORE_DIR) --local $(ARTIFACTS_CORE_DIR)-admin

# PHONY: clean-core
# clean-core: clean-cli-metadata
# 	rm -rf ${ARTIFACTS_CORE_DIR}
# 	rm -rf ${ARTIFACTS_CORE_DIR}-admin

# .PHONY: clean-cli-metadata
# clean-cli-metadata:
# 	- rm -rf ${XDG_DATA_HOME}/tanzu-cli/*
# TANZU CLI

# PLUGINS
.PHONY: prep-build-cli
prep-build-cli:
	$(GO) mod download

.PHONY: build-cli-plugins
build-cli-plugins: prep-build-cli
	tanzu builder cli compile --version $(BUILD_VERSION) --ldflags "$(LD_FLAGS)" --path ./cmd/plugin --artifacts artifacts

.PHONY: install-cli-plugins
install-cli-plugins:
	TANZU_CLI_NO_INIT=true $(GO) run -ldflags "$(LD_FLAGS)" ../../vmware-tanzu-private/core/cmd/cli/tanzu \
		plugin install all --local $(ARTIFACTS_DIR)

PHONY: clean-plugin
clean-plugin: clean-plugin-metadata
	rm -rf ${ARTIFACTS_DIR}
	rm -rf ${ARTIFACTS_DIR}-admin

.PHONY: clean-plugin-metadata
clean-plugin-metadata:
	- rm -rf ${XDG_DATA_HOME}/tanzu-repository/*
# PLUGINS

# MISC
.PHONY: prune
prune:
	find $(ARTIFACTS_DIR) -name "*.exe" -type f -delete
	# find $(ARTIFACTS_DIR) -name "*darwin*" -type f -delete
	# find ${ARTIFACTS_CORE_DIR} -name "*.exe" -type f -delete
	# find ${ARTIFACTS_CORE_DIR} -name "*darwin*" -type f -delete
	# find ${ARTIFACTS_CORE_DIR}-admin -name "*.exe" -type f -delete
	# find ${ARTIFACTS_CORE_DIR}-admin -name "*darwin*" -type f -delete
# MISC
##### BUILD TARGETS #####

##### IMAGE TARGETS #####
deploy-kapp-controller: ## deploys the latest version of kapp-controller
	kubectl create ns kapp-controller || true
	kubectl --namespace kapp-controller apply -f https://gist.githubusercontent.com/joshrosso/e6f73bee6ade35b1be5280be4b6cb1de/raw/b9f8570531857b75a90c1e961d0d134df13adcf1/kapp-controller-build.yaml

push-extensions: ## build and push extension templates
	imgpkg push --bundle $(OCI_REGISTRY)/velero-extension-templates:dev --file extensions/velero/bundle/
	imgpkg push --bundle $(OCI_REGISTRY)/contour-extension-templates:dev --file extensions/contour/bundle/
	imgpkg push --bundle $(OCI_REGISTRY)/gatekeeper-extension-templates:dev --file extensions/gatekeeper/bundle/
	imgpkg push --bundle $(OCI_REGISTRY)/cert-manager-extension-templates:dev --file extensions/cert-manager/bundle/

update-image-lockfiles: ## updates the ImageLock files in each extension
	kbld --file extensions/velero/bundle --imgpkg-lock-output extensions/velero/bundle/.imgpkg/images.yml
	kbld --file extensions/cert-manager/bundle --imgpkg-lock-output extensions/cert-manager/bundle/.imgpkg/images.yml
	kbld --file extensions/gatekeeper/bundle --imgpkg-lock-output extensions/gatekeeper/bundle/.imgpkg/images.yml
	kbld --file extensions/cert-manager/bundle --imgpkg-lock-output extensions/cert-manager/bundle/.imgpkg/images.yml

redeploy-velero: ## delete and redeploy the velero extension
	kubectl --namespace $(EXTENSION_NAMESPACE) --ignore-not-found=true delete app velero
	kubectl apply --filename extensions/velero/extension.yaml

redeploy-gatekeeper: ## delete and redeploy the velero extension
	kubectl -n tanzu-extensions delete app gatekeeper || true
	kubectl apply -f extensions/gatekeeper/extension.yaml

uninstall-contour:
	kubectl --ignore-not-found=true delete namespace projectcontour contour-operator
	kubectl --ignore-not-found=true --namespace $(EXTENSION_NAMESPACE) delete apps contour
	kubectl --ignore-not-found=true delete clusterRoleBinding contour-extension

deploy-contour:
	kubectl apply --filename extensions/contour/extension.yaml

redeploy-cert-manager: ## delete and redeploy the cert-manager extension
	kubectl --namespace tanzu-extensions delete app cert-manager
	kubectl apply --filename extensions/cert-manager/extension.yaml
##### IMAGE TARGETS #####
