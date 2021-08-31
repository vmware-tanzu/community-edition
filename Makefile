# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# detect the build OS
ifeq ($(OS),Windows_NT)
	build_OS := Windows
	NUL = NUL
else
	build_OS := $(shell uname -s 2>/dev/null || echo Unknown)
	NUL = /dev/null
endif

REQUIRED_BINARIES := imgpkg kbld ytt

.DEFAULT_GOAL:=help

### GLOBAL ###
ROOT_DIR := $(shell git rev-parse --show-toplevel)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

TOOLS_DIR := $(ROOT_DIR)/hack/tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin

# Add tooling binaries here and in hack/tools/Makefile
GOLANGCI_LINT := $(TOOLS_BIN_DIR)/golangci-lint
TOOLING_BINARIES := $(GOLANGCI_LINT)

help: #### display help
	@awk 'BEGIN {FS = ":.*## "; printf "\nTargets:\n"} /^[a-zA-Z_-]+:.*?#### / { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@awk 'BEGIN {FS = ":.* ## "; printf "\n  Build targets \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*? ## / { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@awk 'BEGIN {FS = ":.* ### "; printf "\n  Release targets \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*? ### / { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
### GLOBAL ###

##### BUILD #####
ifndef BUILD_VERSION
BUILD_VERSION ?= $$(git describe --tags --abbrev=0)
endif
CONFIG_VERSION ?= $$(echo "$(BUILD_VERSION)" | cut -d "-" -f1)

ifeq ($(strip $(BUILD_VERSION)),)
BUILD_VERSION = dev
endif

FRAMEWORK_BUILD_VERSION=$$(cat "./hack/FRAMEWORK_BUILD_VERSION")
TANZU_FRAMEWORK_REPO_HASH ?= 25b04ec4069e946146226b0acdd79066ff501b6f

ARTIFACTS_DIR ?= ./artifacts

# this captures where the tanzu CLI will be installed (due to usage of go install)
# When GOBIN is set, this is where the tanzu binary is installed
# When GOBIN is not set, but GOPATH is, $GOPATH/bin is where the tanzu binary is installed
# When GOBIN is not set and GOPATH is not set, ${HOME}/go/bin is where the tanzu binary is installed
TANZU_CLI_INSTALL_PATH = "$${HOME}/go/bin/tanzu"
ifdef GOPATH
TANZU_CLI_INSTALL_PATH = "$${GOPATH}/bin/tanzu"
endif
ifdef GOBIN
TANZU_CLI_INSTALL_PATH = "$${GOBIN}/tanzu"
endif

#INSTALLED_CLI_DIR

ifeq ($(GITLAB_CI_BUILD), true)
XDG_DATA_HOME := /tmp/mylocal
SED := sed -i
endif
ifeq ($(XDG_DATA_HOME),)
ifeq ($(build_OS), Darwin)
XDG_DATA_HOME := "${HOME}/Library/Application Support"
SED := sed -i '' -e
endif
endif
ifeq ($(XDG_DATA_HOME),)
ifeq ($(build_OS), Linux)
XDG_DATA_HOME := ${HOME}/.local/share
SED := sed -i
endif
endif

export XDG_DATA_HOME
export GO
export GOLANGCI_LINT
export ARTIFACTS_DIR

PRIVATE_REPOS="github.com/vmware-tanzu/*,github.com/dvonthenen/*,github.com/joshrosso/*"
GO := GOPRIVATE=${PRIVATE_REPOS} go
##### BUILD #####

##### IMAGE #####
OCI_REGISTRY := projects.registry.vmware.com/tce
##### IMAGE #####

##### LINTING TARGETS #####
.PHONY: lint mdlint shellcheck check
check: ensure-deps lint mdlint shellcheck

.PHONY: ensure-deps
ensure-deps:
	hack/ensure-dependencies.sh

GO_MODULES=$(shell find . -path "*/go.mod" | xargs -I _ dirname _)

get-deps:
	@for i in $(GO_MODULES); do \
		echo "-- Getting deps for $$i --"; \
		(cd $$i; $(MAKE) get-deps); \
	done

lint: tools get-deps
	@for i in $(GO_MODULES); do \
		echo "-- Linting $$i --"; \
		(cd $$i; $(MAKE) lint); \
	done

mdlint:
	hack/check-mdlint.sh

shellcheck:
	hack/check-shell.sh
##### LINTING TARGETS #####

##### Tooling Binaries
tools: $(TOOLING_BINARIES)
.PHONY: $(TOOLING_BINARIES)
$(TOOLING_BINARIES):
	make -C $(TOOLS_DIR) $(@F)
##### Tooling Binaries

##### BUILD TARGETS #####
build: build-plugin

build-all: release-env-check version clean install-cli install-cli-plugins ## build all CLI plugins that are used in TCE
	@printf "\n[COMPLETE] installed plugins at $${XDG_DATA_HOME}/tanzu-cli/. "
	@printf "These plugins will be automatically detected by tanzu CLI.\n"
	@printf "\n[COMPLETE] installed tanzu CLI at $(TANZU_CLI_INSTALL_PATH). "
	@printf "Move this binary to a location in your path!\n"

build-plugin: version clean-plugin install-cli-plugins ## build only CLI plugins that live in the TCE repo
	@printf "\n[COMPLETE] installed TCE-specific plugins at $${XDG_DATA_HOME}/tanzu-cli/. "
	@printf "These plugins will be automatically detected by your tanzu CLI.\n"

release: build-all package-release ### builds and produces the release packaging/tarball for TCE in your local Go environment

release-docker: release-env-check ### builds and produces the release packaging/tarball for TCE in a containerized environment
	docker run --rm \
		-e HOME=/go \
		-e GITHUB_TOKEN=${GITHUB_TOKEN} \
		-e GITLAB_CI_BUILD=true \
		-w /go/src/community-edition \
		-v ${PWD}:/go/src/community-edition \
		-v /tmp:/tmp \
		golang:1.16.2 \
		sh -c "cd /go/src/community-edition &&\
			./hack/fix-for-ci-build.sh &&\
			make release"

clean: clean-release clean-plugin clean-framework

release-env-check:
ifndef GITHUB_TOKEN
	$(error GITHUB_TOKEN is undefined)
endif

# RELEASE MANAGEMENT
version:
	@echo "BUILD_VERSION:" ${BUILD_VERSION}
	@echo "CONFIG_VERSION:" ${CONFIG_VERSION}
	@echo "FRAMEWORK_BUILD_VERSION:" ${FRAMEWORK_BUILD_VERSION}
	@echo "XDG_DATA_HOME:" $(XDG_DATA_HOME)

.PHONY: package-release
package-release:
	FRAMEWORK_BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} BUILD_VERSION=${BUILD_VERSION} hack/package-release.sh

# IMPORTANT: This should only ever be called CI/github-action
.PHONY: cut-release
cut-release: version
	BUILD_VERSION=$(BUILD_VERSION) FAKE_RELEASE=$(shell expr $(BUILD_VERSION) | grep fake) hack/cut-release.sh
	echo "$(BUILD_VERSION)" | tee -a ./cayman_trigger.txt

.PHONY: upload-signed-assets
upload-signed-assets: release-env-check
	cd ./hack/asset && $(MAKE) run && cd ../..
# IMPORTANT: This should only ever be called CI/github-action

clean-release:
	rm -rf ./build
	rm -f ./hack/NEW_BUILD_VERSION
# RELEASE MANAGEMENT

# TANZU CLI
.PHONY: build-cli
build-cli: install-cli

.PHONY: install-cli
install-cli:
	TKG_DEFAULT_COMPATIBILITY_IMAGE_PATH="tkg-compatibility" \
	TANZU_FRAMEWORK_REPO_HASH=$(TANZU_FRAMEWORK_REPO_HASH) BUILD_EDITION=tce TCE_BUILD_VERSION=$(BUILD_VERSION) \
	FRAMEWORK_BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} hack/build-tanzu.sh

.PHONY: clean-framework
clean-framework:
	rm -rf /tmp/tce-release
	rm -rf ${XDG_DATA_HOME}/tanzu-cli
	mkdir -p ${XDG_DATA_HOME}/tanzu-cli
# TANZU CLI

# PLUGINS
.PHONY: prep-build-cli
prep-build-cli:
	@cd ./cli/cmd/plugin/ && for plugin in *; do\
		printf "===> Preparing $${plugin}\n";\
		working_dir=`pwd`;\
		cd $${plugin};\
		$(MAKE) get-deps;\
		cd $${working_dir};\
	done

.PHONY: build-cli-plugins
build-cli-plugins: prep-build-cli
	@cd ./hack/builder/ && $(MAKE) compile

.PHONY: install-cli-plugins
install-cli-plugins: build-cli-plugins
	@cd ./hack/builder/ && $(MAKE) install-plugins

test-plugins: ## run tests on TCE plugins
	# TODO(joshrosso): update once we get our testing strategy in place
	@echo "No tests to run."

.PHONY: clean-plugin
clean-plugin:
	rm -rf ${ARTIFACTS_DIR}
# PLUGINS

##### BUILD TARGETS #####

##### PACKAGE OPERATIONS #####

check-carvel:
	$(foreach exec,$(REQUIRED_BINARIES),\
		$(if $(shell which $(exec)),,$(error "'$(exec)' not found. Carvel toolset is required. See instructions at https://carvel.dev/#install")))

.PHONY: create-package
create-package: # Stub out new package directories and manifests. Usage: make create-package NAME=foobar VERSION=10.0.0
	@hack/packages/create-package.sh $(NAME) $(VERSION)

vendir-sync-package: check-carvel # Performs a `vendir sync` for a package. Usage: make vendir-package-sync PACKAGE=foobar VERSION=1.0.0
	@printf "\n===> syncing $${PACKAGE}/$${VERSION}\n";\
	cd addons/packages/$${PACKAGE}/$${VERSION}/bundle && vendir sync >> /dev/null;\

lock-package-images: check-carvel # Updates the image lock file for a package. Usage: make lock-package-images PACKAGE=foobar VERSION=1.0.0
	@printf "\n===> Updating image lockfile for package $${PACKAGE}/$${VERSION}\n";\
	cd addons/packages/$${PACKAGE}/$${VERSION} && kbld --file bundle --imgpkg-lock-output bundle/.imgpkg/images.yml >> /dev/null;\

push-package: check-carvel # Build and push a package template. Tag will default to `latest`. Usage: make push-package PACKAGE=foobar VERSION=1.0.0
	@printf "\n===> pushing $${PACKAGE}/$${VERSION}\n";\
	cd addons/packages/$${PACKAGE}/$${VERSION} && imgpkg push --bundle $(OCI_REGISTRY)/$${PACKAGE}:$${VERSION} --file bundle/;\

export REPO
generate-package-repo: check-carvel # Generate and push the package repository. Usage: make generate-package-repo REPO=main
	cd ./hack/packages/ && $(MAKE) run

get-package-config: # Extracts the package values.yaml file. Usage: make get-package-config PACKAGE=foo VERSION=1.0.0
	TEMP_DIR=`mktemp -d` \
	&& imgpkg pull --bundle ${OCI_REGISTRY}/$${PACKAGE}:$${VERSION} -o $${TEMP_DIR} \
	&& cp $${TEMP_DIR}/config/values.yaml ./$${PACKAGE}-$${VERSION}-values.yaml \
	&& rm -rf $${TEMP_DIR}

test-packages-unit: check-carvel
	$(GO) test -coverprofile cover.out -v `go list ./... | grep github.com/vmware-tanzu/community-edition/addons/packages | grep -v e2e`

create-repo: # Usage: make create-repo NAME=my-repo
	cp hack/packages/templates/repo.yaml addons/repos/${NAME}.yaml

##### PACKAGE OPERATIONS #####

##### NESTED MAKEFILE SUPPORT #####

makefile:
	@cat "./hack/makefile-template";

##### NESTED MAKEFILE SUPPORT #####

##### E2E TESTS #####

##### BUILD TARGETS #####

e2e-test:
	test/aws/deploy-tce.sh

# TCE Docker Standalone Cluster E2E Test
tce-docker-standalone-cluster-e2e-test:
	test/docker/run-tce-docker-standalone-cluster.sh

# TCE Docker Managed Cluster E2E Test
tce-docker-managed-cluster-e2e-test:
	test/docker/run-tce-docker-managed-cluster.sh

# TCE vSphere Standalone Cluster E2E Test
tce-vsphere-standalone-cluster-e2e-test:
	test/vsphere/run-tce-vsphere-standalone-cluster.sh

##### E2E TESTS #####
