# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

ifeq ($(OS),Windows_NT)
	build_OS := Windows
	NUL = NUL
else
	build_OS := $(shell uname -s 2>/dev/null || echo Unknown)
	NUL = /dev/null
endif

.DEFAULT_GOAL:=help

### GLOBAL ###
ROOT_DIR := $(shell git rev-parse --show-toplevel)

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

help: ## display help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
### GLOBAL ###

##### BUILD #####
BUILD_VERSION ?= $$(git describe --tags --abbrev=0)
BUILD_SHA ?= $$(git rev-parse --short HEAD)
BUILD_DATE ?= $$(date -u +"%Y-%m-%d")
CONFIG_VERSION ?= $$(echo "$(BUILD_VERSION)" | cut -d "-" -f1)

ifeq ($(strip $(BUILD_VERSION)),)
BUILD_VERSION = dev
endif
CORE_BUILD_VERSION=$$(cat "./hack/CORE_BUILD_VERSION")

LD_FLAGS = -X "github.com/vmware-tanzu-private/core/pkg/v1/cli.BuildDate=$(BUILD_DATE)"
LD_FLAGS += -X "github.com/vmware-tanzu-private/core/pkg/v1/cli.BuildSHA=$(BUILD_SHA)"
LD_FLAGS += -X "github.com/vmware-tanzu-private/core/pkg/v1/cli.BuildVersion=$(BUILD_VERSION)"

ARTIFACTS_DIR ?= ./artifacts

ifeq ($(build_OS), Linux)
XDG_DATA_HOME := ${HOME}/.local/share
endif
ifeq ($(build_OS), Darwin)
XDG_DATA_HOME := ${HOME}/Library/ApplicationSupport
endif

export XDG_DATA_HOME

PRIVATE_REPOS="github.com/vmware-tanzu-private/*,github.com/vmware-tanzu/*"
GO := GOPRIVATE=${PRIVATE_REPOS} go
##### BUILD #####

##### IMAGE #####
OCI_REGISTRY := projects.registry.vmware.com/tce
EXTENSION_NAMESPACE := tanzu-extensions
##### IMAGE #####

##### LINTING TARGETS #####
.PHONY: fmt vet lint mdlint shellcheck staticcheck check
check: fmt lint mdlint shellcheck staticcheck vet

fmt:
	hack/check-format.sh

lint:
	hack/check-lint.sh

mdlint:
	hack/check-mdlint.sh

shellcheck:
	hack/check-shell.sh

staticcheck:
	hack/check-staticcheck.sh

vet:
	hack/check-vet.sh
##### LINTING TARGETS #####

##### BUILD TARGETS #####
build: build-plugin

build-all: version clean copy-release tag-release install-cli install-cli-plugins
build-plugin: version clean-plugin copy-release tag-release install-cli-plugins

release: build-all gen-metadata package-release

clean: clean-release clean-plugin clean-core

# RELEASE MANAGEMENT
version:
	@echo "BUILD_VERSION:" ${BUILD_VERSION}
	@echo "CONFIG_VERSION:" ${CONFIG_VERSION}
	@echo "CORE_BUILD_VERSION:" ${CORE_BUILD_VERSION}

PHONY: gen-metadata
gen-metadata:
	go run ./hack/release/release.go
ifeq ($(shell expr $(BUILD_VERSION)), $(shell expr $(CONFIG_VERSION)))
	go run ./hack/metadata/metadata.go -tag $(CONFIG_VERSION) -release
else
	go run ./hack/metadata/metadata.go -tag $(CONFIG_VERSION)
endif

PHONY: copy-release
copy-release:
	mkdir -p ${XDG_DATA_HOME}/tanzu-repository
	cp -f ./hack/config.yaml ${XDG_DATA_HOME}/tanzu-repository/config.yaml

.PHONY: tag-release
tag-release:
ifeq ($(shell expr $(BUILD_VERSION)), $(shell expr $(CONFIG_VERSION)))
	sed -i "s/version: latest/version: $(CONFIG_VERSION)/g" ./hack/config.yaml
	sed -i "s/version: latest/version: $(CONFIG_VERSION)/g" ${XDG_DATA_HOME}/tanzu-repository/config.yaml
endif

.PHONY: package-release
package-release:
	CORE_BUILD_VERSION=${CORE_BUILD_VERSION} BUILD_VERSION=${BUILD_VERSION} hack/package-release.sh

clean-release:
	rm -rf ./build
	rm -rf ./metadata
	rm -rf ./offline
# RELEASE MANAGEMENT

# TANZU CLI
.PHONY: build-cli
build-cli: install-cli

.PHONY: install-cli
install-cli:
	CORE_BUILD_VERSION=${CORE_BUILD_VERSION} hack/build-tanzu.sh

PHONY: clean-core
clean-core: clean-cli-metadata
	rm -rf ./tkg-cli
	rm -rf ./tkg-providers
	rm -rf ./core
	rm -rf ./tanzu-cli-tkg-plugins

.PHONY: clean-cli-metadata
clean-cli-metadata:
	- rm -rf ${XDG_DATA_HOME}/tanzu-cli/*
# TANZU CLI

# PLUGINS
.PHONY: prep-build-cli
prep-build-cli:
	$(GO) mod download

.PHONY: build-cli-plugins
build-cli-plugins: prep-build-cli
	$(GO) run github.com/vmware-tanzu-private/core/cmd/cli/plugin-admin/builder cli compile --version $(BUILD_VERSION) \
		--ldflags "$(LD_FLAGS)" --path ./cmd/plugin --artifacts ${ARTIFACTS_DIR}
	# $(GO) run ../../vmware-tanzu-private/core/cmd/cli/plugin-admin/builder/main.go cli compile --version $(BUILD_VERSION) \
	# 	--ldflags "$(LD_FLAGS)" --corepath "cmd/cli/tanzu" --target linux_amd64
	# $(GO) run ../../vmware-tanzu-private/core/cmd/cli/plugin-admin/builder/main.go cli compile --version $(BUILD_VERSION) \
	# 	--ldflags "$(LD_FLAGS)" --corepath "cmd/cli/tanzu" --target windows_amd64
	# $(GO) run ../../vmware-tanzu-private/core/cmd/cli/plugin-admin/builder/main.go cli compile --version $(BUILD_VERSION) \
	# 	--ldflags "$(LD_FLAGS)" --corepath "cmd/cli/tanzu" --target darwin_amd64

.PHONY: install-cli-plugins
install-cli-plugins: build-cli-plugins
	TANZU_CLI_NO_INIT=true $(GO) run -ldflags "$(LD_FLAGS)" github.com/vmware-tanzu-private/core/cmd/cli/tanzu \
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

.PHONY: create-addon
create-addon: ## create the directory structure from a new add-on
	hack/create-addon-dir.sh $(NAME)

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
	imgpkg push --bundle $(OCI_REGISTRY)/knative-serving-extension-templates:dev --file extensions/knative-serving/bundle/

update-image-lockfiles: ## updates the ImageLock files in each extension
	kbld --file extensions/velero/bundle --imgpkg-lock-output extensions/velero/bundle/.imgpkg/images.yml
	kbld --file extensions/cert-manager/bundle --imgpkg-lock-output extensions/cert-manager/bundle/.imgpkg/images.yml
	kbld --file extensions/gatekeeper/bundle --imgpkg-lock-output extensions/gatekeeper/bundle/.imgpkg/images.yml
	kbld --file extensions/cert-manager/bundle --imgpkg-lock-output extensions/cert-manager/bundle/.imgpkg/images.yml
	kbld --file extensions/knative-serving/bundle --imgpkg-lock-output extensions/knative-serving/bundle/.imgpkg/images.yml

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

uninstall-knative-serving:
	kubectl --ignore-not-found=true --namespace $(EXTENSION_NAMESPACE) delete apps knative-serving
	kubectl --ignore-not-found=true delete clusterRoleBinding knative-serving-extension
	kubectl --ignore-not-found=true delete service knative-serving-extension

deploy-knative-serving:
	kubectl apply --filename extensions/knative-serving/serviceaccount.yaml
	kubectl apply --filename extensions/knative-serving/clusterrolebinding.yaml
	kubectl apply --filename extensions/knative-serving/extension.yaml

update-knative-serving: ## updates the ImageLock files in each extension
	kbld --file extensions/knative-serving/bundle --imgpkg-lock-output extensions/knative-serving/bundle/.imgpkg/images.yml
	imgpkg push --bundle $(OCI_REGISTRY)/knative-serving-extension-templates:dev --file extensions/knative-serving/bundle/

redeploy-cert-manager: ## delete and redeploy the cert-manager extension
	kubectl --namespace tanzu-extensions delete app cert-manager
	kubectl apply --filename extensions/cert-manager/extension.yaml
##### IMAGE TARGETS #####
