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
BUILD_SHA ?= $$(git rev-parse --short HEAD)
BUILD_DATE ?= $$(date -u +"%Y-%m-%d")
ifndef BUILD_VERSION
BUILD_VERSION ?= $$(git describe --tags --abbrev=0)
endif
CONFIG_VERSION ?= $$(echo "$(BUILD_VERSION)" | cut -d "-" -f1)
BUILD_EDITION=tce-standalone

ifeq ($(strip $(BUILD_VERSION)),)
BUILD_VERSION = dev
endif
ifndef IS_OFFICIAL_BUILD
IS_OFFICIAL_BUILD = ""
endif

FRAMEWORK_BUILD_VERSION=$$(cat "./hack/FRAMEWORK_BUILD_VERSION")
TANZU_FRAMEWORK_REPO_HASH ?= f7e4fbdac6478244ce5dbf25288d212f12df9e28

LD_FLAGS = -s -w
LD_FLAGS += -X "github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli.BuildDate=$(BUILD_DATE)"
LD_FLAGS += -X "github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli.BuildSHA=$(BUILD_SHA)"
LD_FLAGS += -X "github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli.BuildVersion=$(BUILD_VERSION)"
LD_FLAGS += -X 'main.BuildEdition=$(BUILD_EDITION)'
LD_FLAGS += -X 'github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/buildinfo.IsOfficialBuild=$(IS_OFFICIAL_BUILD)'

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

lint: tools
	@printf "\n===> Linting standalone plugin\n"
	@cd cli/cmd/plugin/standalone-cluster && $(GOLANGCI_LINT) run -v --timeout=5m
	@printf "\n===> Linting hack packages\n"
	@cd hack/asset && $(GOLANGCI_LINT) run -v --timeout=5m
	@cd hack/packages && $(GOLANGCI_LINT) run -v --timeout=5m
	@cd hack/release && $(GOLANGCI_LINT) run -v --timeout=5m

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
		-e GH_ACCESS_TOKEN=${GH_ACCESS_TOKEN} \
		-e GITLAB_CI_BUILD=true \
		-w /go/src/tce \
		-v ${PWD}:/go/src/tce \
		-v /tmp:/tmp \
		golang:1.16.2 \
		sh -c "cd /go/src/tce &&\
			./hack/fix-for-ci-build.sh &&\
			make release"

clean: clean-release clean-plugin clean-framework

release-env-check:
ifndef GH_ACCESS_TOKEN
	$(error GH_ACCESS_TOKEN is undefined)
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
	cd ./hack/asset && go run ./asset.go -tag $(BUILD_VERSION) && cd ../..
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
		$(GO) mod download;\
		cd $${working_dir};\
	done

.PHONY: build-cli-plugins
build-cli-plugins: prep-build-cli
	@cd ./hack/builder/ && \
		$(GO) run github.com/vmware-tanzu/tanzu-framework/cmd/cli/plugin-admin/builder cli compile --version $(BUILD_VERSION) \
			--ldflags "$(LD_FLAGS)" --path ../../cli/cmd/plugin --artifacts ../../${ARTIFACTS_DIR}

.PHONY: install-cli-plugins
install-cli-plugins: build-cli-plugins
	@cd ./hack/builder/ && \
		TANZU_CLI_NO_INIT=true $(GO) run -ldflags "$(LD_FLAGS)" github.com/vmware-tanzu/tanzu-framework/cmd/cli/tanzu \
			plugin install all --local ../../$(ARTIFACTS_DIR)

test-plugins: ## run tests on TCE plugins
	# TODO(joshrosso): update once we get our testing strategy in place
	@echo "No tests to run."

.PHONY: clean-plugin
clean-plugin:
	rm -rf ${ARTIFACTS_DIR}
# PLUGINS

# MISC
.PHONY: create-package
create-package: # Stub out new package directories and manifests. Usage: make create-package NAME=foobar
	@hack/create-package-dir.sh $(NAME)

.PHONY: create-channel
create-channel: # Stub out new channel values file. Usage: make create-channel NAME=foobar
	@hack/create-channel.sh $(NAME)
# MISC

##### BUILD TARGETS #####

##### PACKAGE OPERATIONS #####

check-carvel:
	$(foreach exec,$(REQUIRED_BINARIES),\
		$(if $(shell which $(exec)),,$(error "'$(exec)' not found. Carvel toolset is required. See instructions at https://carvel.dev/#install")))

vendir-sync-package: check-carvel # Performs a `vendir sync` for a package. Usage: make vendir-package-sync PACKAGE=foobar VERSION=1.0.0
	@printf "\n===> syncing $${PACKAGE}/$${VERSION}\n";\
	cd addons/packages/$${PACKAGE}/$${VERSION}/bundle && vendir sync >> /dev/null;\

lock-package-images: check-carvel # Updates the image lock file for a package. Usage: make lock-package-images PACKAGE=foobar VERSION=1.0.0
	@printf "\n===> Updating image lockfile for package $${PACKAGE}/$${VERSION}\n";\
	cd addons/packages/$${PACKAGE}/$${VERSION} && kbld --file bundle --imgpkg-lock-output bundle/.imgpkg/images.yml >> /dev/null;\

push-package: check-carvel # Build and push a package template. Tag will default to `latest`. Usage: make push-package PACKAGE=foobar VERSION=1.0.0
	@printf "\n===> pushing $${PACKAGE}/$${VERSION}\n";\
	cd addons/packages/$${PACKAGE}/$${VERSION} && imgpkg push --bundle $(OCI_REGISTRY)/$${PACKAGE}:$${VERSION} --file bundle/;\

generate-package-repo: check-carvel # Generate and push the package repository. Usage: make generate-package-repo REPO=main
	cd ./hack/packages/ && go run generate-package-repository.go $${REPO}

get-package-config: # Extracts the package values.yaml file. Usage: make get-package-config PACKAGE=foo VERSION=1.0.0
	TEMP_DIR=`mktemp -d` \
	&& imgpkg pull --bundle ${OCI_REGISTRY}/$${PACKAGE}:$${VERSION} -o $${TEMP_DIR} \
	&& cp $${TEMP_DIR}/config/values.yaml ./$${PACKAGE}-$${VERSION}-values.yaml \
	&& rm -rf $${TEMP_DIR}

test-packages-unit: check-carvel
	$(GO) test -coverprofile cover.out -v `go list ./... | grep github.com/vmware-tanzu/tce/addons/packages | grep -v e2e`

create-repo: # Usage: make create-repo NAME=my-repo
	cp hack/packages/templates/repo.yaml addons/repos/${NAME}.yaml

##### PACKAGE OPERATIONS #####

generate-channel:
	@print "\nGenerating CHANNEL file:\n";\

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
