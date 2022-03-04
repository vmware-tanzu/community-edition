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
GO := go
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOHOSTOS ?= $(shell go env GOHOSTOS)
GOHOSTARCH ?= $(shell go env GOHOSTARCH)
# Add supported OS-ARCHITECTURE combinations here
ENVS := linux-amd64 windows-amd64 darwin-amd64 darwin-arm64

TOOLS_DIR := $(ROOT_DIR)/hack/tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin

# Add tooling binaries here and in hack/tools/Makefile
GOLANGCI_LINT := $(TOOLS_BIN_DIR)/golangci-lint
TOOLING_BINARIES := $(GOLANGCI_LINT)

help: #### display help
	@awk 'BEGIN {FS = ":.*## "; printf "\nTargets:\n"} /^[a-zA-Z_-]+:.*?#### / { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@awk 'BEGIN {FS = ":.* ## "; printf "\n  \033[1;32mBuild targets\033[36m\033[0m\n  \033[0;37mTargets for building and/or installing CLI plugins on the system.\n  Append \"ENVS=<os-arch>\" to the end of these targets to limit the binaries built.\n  e.g.: make build-all-tanzu-cli-plugins ENVS=linux-amd64  \n  List available at https://github.com/golang/go/blob/master/src/go/build/syslist.go\033[36m\033[0m\n\n"} /^[a-zA-Z_-]+:.*? ## / { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@awk 'BEGIN {FS = ":.* ### "; printf "\n  \033[1;32mRelease targets\033[36m\033[0m\n\033[0;37m  Targets for producing a TCE release package.\033[36m\033[0m\n\n"} /^[a-zA-Z_-]+:.*? ### / { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
### GLOBAL ###

##### BUILD #####
ifndef PLUGINS
PLUGINS ?= "conformance diagnostics unmanaged-cluster"
endif
ifndef ifndef
# For TF 0.17.0 or higher
# DISCOVERY_NAME ?= "default"
# For 0.11.1
DISCOVERY_NAME ?= "standalone"
endif
ifndef BUILD_VERSION
BUILD_VERSION ?= $$(git describe --tags --abbrev=0)
endif
CONFIG_VERSION ?= $$(echo "$(BUILD_VERSION)" | cut -d "-" -f1)

ifeq ($(strip $(BUILD_VERSION)),)
BUILD_VERSION = dev
endif

# TANZU_FRAMEWORK_REPO override for being able to use your own fork
TANZU_FRAMEWORK_REPO ?= https://github.com/vmware-tanzu/tanzu-framework.git
# TANZU_FRAMEWORK_REPO_BRANCH sets a branch or tag to build Tanzu Framework
TANZU_FRAMEWORK_REPO_BRANCH ?= v0.11.1
# if the hash below is set, this overrides the value of TANZU_FRAMEWORK_REPO_BRANCH
TANZU_FRAMEWORK_REPO_HASH ?=
# TKG_DEFAULT_IMAGE_REPOSITORY override for using a different image repo
ifndef TKG_DEFAULT_IMAGE_REPOSITORY
TKG_DEFAULT_IMAGE_REPOSITORY ?= projects-stg.registry.vmware.com/tkg
endif
# TKG_DEFAULT_COMPATIBILITY_IMAGE_PATH override for using a different image path
ifndef TKG_DEFAULT_COMPATIBILITY_IMAGE_PATH
TKG_DEFAULT_COMPATIBILITY_IMAGE_PATH ?= framework-zshippable/tkg-compatibility
endif
FRAMEWORK_BUILD_VERSION=$$(cat "./hack/FRAMEWORK_BUILD_VERSION")

ARTIFACTS_DIR ?= ./artifacts
TCE_PLUGIN_BUILD_DIR ?= ./build
TCE_SCRATCH_DIR ?= /tmp/tce-scratch-space

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

ifeq ($(TCE_CI_BUILD), true)
XDG_DATA_HOME := /tmp/mylocal
XDG_CACHE_HOME := /tmp/mycache
XDG_CONFIG_HOME := /tmp/myconfig
SED := sed -i
endif
ifeq ($(XDG_DATA_HOME),)
ifeq ($(build_OS), Darwin)
XDG_DATA_HOME := "${HOME}/Library/Application Support"
XDG_CACHE_HOME := ${HOME}/.cache
XDG_CONFIG_HOME := ${HOME}/.config
SED := sed -i '' -e
endif
endif
ifeq ($(XDG_DATA_HOME),)
ifeq ($(build_OS), Linux)
XDG_DATA_HOME := ${HOME}/.local/share
XDG_CACHE_HOME := ${HOME}/.cache
XDG_CONFIG_HOME := ${HOME}/.config
SED := sed -i
endif
endif

export XDG_DATA_HOME
export XDG_CACHE_HOME
export XDG_CONFIG_HOME
export GO
export GOLANGCI_LINT
export ARTIFACTS_DIR
##### BUILD #####

##### IMAGE #####
OCI_REGISTRY := projects.registry.vmware.com/tce
##### IMAGE #####

##### LINTING TARGETS #####
.PHONY: lint mdlint shellcheck check yamllint misspell actionlint urllint imagelint
check: ensure-deps lint mdlint shellcheck yamllint misspell actionlint urllint imagelint

.PHONY: check-deps-minimum-build
check-deps-minimum-build:
	hack/ensure-deps/check-zip.sh

.PHONY: ensure-deps
ensure-deps:
	hack/ensure-deps/ensure-dependencies.sh

GO_MODULES=$(shell find . -path "*/go.mod" | xargs -I _ dirname _)

get-deps:
	@for i in $(GO_MODULES); do \
		echo "-- Getting deps for $$i --"; \
		working_dir=`pwd`; \
		if [ "$${i}" = "." ]; then \
			go mod tidy; \
		else \
			cd $${i}; \
			$(MAKE) get-deps || exit 1; \
			cd $$working_dir; \
		fi; \
	done

# Verify if go.mod and go.sum Go module files are out of sync
verify-modules: get-deps
	@for i in $(GO_MODULES); do \
		echo "-- Verifying modules for $$i --"; \
		working_dir=`pwd`; \
		cd $${i}; \
		if [ "`git diff --name-only HEAD -- go.sum go.mod`" != "" ]; then \
			echo "go module files in $$i directory are out of date, run 'go mod tidy' and commit the changes"; exit 1; \
		fi; \
		cd $$working_dir; \
	done

lint: tools verify-modules
	@for i in $(GO_MODULES); do \
		echo "-- Linting $$i --"; \
		working_dir=`pwd`; \
		if [ "$${i}" = "." ]; then \
			$(GOLANGCI_LINT) run -v --timeout=5m; \
		else \
			cd $${i}; \
			echo $(MAKE) lint || exit 1; \
			cd $$working_dir; \
		fi; \
	done

mdlint:
	# mdlint rules with common errors and possible fixes can be found here:
	# https://github.com/DavidAnson/markdownlint/blob/main/doc/Rules.md
	hack/check/check-mdlint.sh

shellcheck:
	hack/check/check-shell.sh

yamllint:
	hack/check/check-yaml.sh

misspell:
	hack/check/check-misspell.sh

actionlint:
	go install github.com/rhysd/actionlint/cmd/actionlint@latest
	actionlint -shellcheck=

urllint:
	cd ./hack/urllinter && go build -o urllinter main.go
	hack/urllinter/urllinter --path=./ --config=hack/check/.urllintconfig.yaml --summary=true --details=Fail

imagelint:
	cd ./hack/imagelinter && go build -o imagelinter main.go
	hack/imagelinter/imagelinter --path=./ --config=hack/check/.imagelintconfig.yaml --summary=true --details=all

##### LINTING TARGETS #####

##### Tooling Binaries
tools: $(TOOLING_BINARIES)
.PHONY: $(TOOLING_BINARIES)
$(TOOLING_BINARIES):
	make -C $(TOOLS_DIR) $(@F)
##### Tooling Binaries

##### BUILD TARGETS #####
build-tce-cli-plugins: version build-cli-plugins ## builds the CLI plugins that live in the TCE repo into the artifacts directory
	@printf "\n[COMPLETE] built TCE-specific plugins at $(ARTIFACTS_DIR)\n"
	@printf "To install these plugins, run \`make install-tce-cli-plugins\`\n"

install-tce-cli-plugins: version build-cli-plugins install-plugins ## builds and installs CLI plugins found in artifacts directory
	@printf "\n[COMPLETE] built and installed TCE-specific plugins at $${XDG_DATA_HOME}/tanzu-cli/. "
	@printf "These plugins will be automatically detected by your tanzu CLI.\n"	

build-all-tanzu-cli-plugins: version clean build-cli build-cli-plugins ## builds the Tanzu CLI and all CLI plugins that are used in TCE
	@printf "\n[COMPLETE] built plugins at $(ARTIFACTS_DIR)\n"
	@printf "These plugins will be automatically detected by tanzu CLI.\n"
	@printf "\n[COMPLETE] built tanzu CLI at $(TCE_SCRATCH_DIR). "
	@printf "Move this binary to a location in your path!\n"

install-all-tanzu-cli-plugins: version clean build-cli install-cli build-cli-plugins install-plugins ## installs the Tanzu CLI and all CLI plugins that are used in TCE
	@printf "\n[COMPLETE] built and installed TCE-specific plugins at $${XDG_DATA_HOME}/tanzu-cli/.\n"
	@printf "These plugins will be automatically detected by your tanzu CLI.\n"
	@printf "\n[COMPLETE] built and installed tanzu CLI at $(TANZU_CLI_INSTALL_PATH). "
	@printf "Move this binary to a location in your path!\n"

release: build-all-tanzu-cli-plugins package-release ### builds and produces the release packaging/tarball for TCE in your local Go environment

release-docker: ### builds and produces the release packaging/tarball for TCE in a containerized environment
	docker run --rm \
		-e HOME=/go \
		-e TCE_CI_BUILD=true \
		-w /go/src/community-edition \
		-v ${PWD}:/go/src/community-edition \
		-v /tmp:/tmp \
		golang:1.17.6 \
		sh -c "cd /go/src/community-edition &&\
			./hack/release/fix-for-ci-build.sh &&\
			make release"

clean: clean-release clean-plugin clean-framework

# RELEASE MANAGEMENT
version:
	@echo "BUILD_VERSION:" ${BUILD_VERSION}
	@echo "CONFIG_VERSION:" ${CONFIG_VERSION}
	@echo "FRAMEWORK_BUILD_VERSION:" ${FRAMEWORK_BUILD_VERSION}
	@echo "XDG_DATA_HOME:" $(XDG_DATA_HOME)

.PHONY: upload-daily-build
upload-daily-build:
	BUILD_VERSION=$(BUILD_VERSION) ./hack/dailybuild/publish-daily-build.sh

.PHONY: package-release
package-release:
	TCE_SCRATCH_DIR=${TCE_SCRATCH_DIR} FRAMEWORK_BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} BUILD_VERSION=${BUILD_VERSION} \
	DISCOVERY_NAME=${DISCOVERY_NAME} ENVS="${ENVS}" hack/release/package-release.sh

# IMPORTANT: This should only ever be called CI/github-action
.PHONY: cut-release
cut-release: version
	TCE_SCRATCH_DIR=${TCE_SCRATCH_DIR} FRAMEWORK_BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} \
	BUILD_VERSION=$(BUILD_VERSION) FAKE_RELEASE=$(shell expr $(BUILD_VERSION) | grep fake) \
	hack/release/cut-release.sh
	echo "$(BUILD_VERSION)" | tee -a ./cayman_trigger.txt

# This target creates the directory structure needed for the GCP update buckets. When the OCI functionality
# is implemented, this target along with all associated scripts, github actions, and etc can be deleted
.PHONY: prep-gcp-tanzu-bucket
prep-gcp-tanzu-bucket:
	TCE_SCRATCH_DIR=${TCE_SCRATCH_DIR} \
	TANZU_FRAMEWORK_REPO=${TANZU_FRAMEWORK_REPO} \
	TKG_DEFAULT_IMAGE_REPOSITORY=${TKG_DEFAULT_IMAGE_REPOSITORY} \
	TKG_DEFAULT_COMPATIBILITY_IMAGE_PATH=${TKG_DEFAULT_COMPATIBILITY_IMAGE_PATH} \
	TANZU_FRAMEWORK_REPO_BRANCH=${TANZU_FRAMEWORK_REPO_BRANCH} \
	TANZU_FRAMEWORK_REPO_HASH=${TANZU_FRAMEWORK_REPO_HASH} \
	BUILD_EDITION=tce TCE_BUILD_VERSION=$(BUILD_VERSION) \
	FRAMEWORK_BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} ENVS="${ENVS}" hack/builder/prep-gcp-tanzu.sh

# Please see above
.PHONY: prep-gcp-tce-bucket
prep-gcp-tce-bucket:
	hack/builder/prep-gcp-tce.sh

# Please see above
.PHONY: prune-buckets
prune-buckets:
	TCE_SCRATCH_DIR=${TCE_SCRATCH_DIR} hack/release/prune-buckets.sh

# The main target for GCP buckets. Please see above
.PHONY: release-buckets
release-buckets: version prep-gcp-tanzu-bucket prep-gcp-tce-bucket build-cli-plugins-nopublish

.PHONY: upload-signed-assets
upload-signed-assets:
	@cd ./hack/asset && $(MAKE) run

release-gate:
	./hack/ensure-deps/ensure-gh-cli.sh
	./hack/release/trigger-release-gate-pipelines.sh

# IMPORTANT: This should only ever be called CI/github-action

clean-release:
	rm -rf ./release
	rm -f ./hack/NEW_BUILD_VERSION
# RELEASE MANAGEMENT

# TANZU CLI
# this forces the reinstallation of all plugins found in the TF ./build directory
.PHONY: build-cli-force
build-cli-force:
export FORCE_UPDATE_PLUGIN=true

# if we are (re)building, then we want to explicitly force reinstall of plugins
.PHONY: build-cli
build-cli: check-deps-minimum-build build-cli-force
	TCE_SCRATCH_DIR=${TCE_SCRATCH_DIR} \
	TANZU_FRAMEWORK_REPO=${TANZU_FRAMEWORK_REPO} \
	TKG_DEFAULT_IMAGE_REPOSITORY=${TKG_DEFAULT_IMAGE_REPOSITORY} \
	TKG_DEFAULT_COMPATIBILITY_IMAGE_PATH=${TKG_DEFAULT_COMPATIBILITY_IMAGE_PATH} \
	TANZU_FRAMEWORK_REPO_BRANCH=${TANZU_FRAMEWORK_REPO_BRANCH} \
	TANZU_FRAMEWORK_REPO_HASH=${TANZU_FRAMEWORK_REPO_HASH} \
	BUILD_EDITION=tce TCE_BUILD_VERSION=$(BUILD_VERSION) \
	DISCOVERY_NAME=$(DISCOVERY_NAME) \
	FRAMEWORK_BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} ENVS="${ENVS}" hack/builder/build-tanzu.sh

# since we are installing off a previously built TF, just install the Tanzu CLI and only
# reinstall plugins if FORCE_UPDATE_PLUGIN=true is explicitly stated
.PHONY: install-cli
install-cli:
	TCE_SCRATCH_DIR=${TCE_SCRATCH_DIR} \
	DISCOVERY_NAME=$(DISCOVERY_NAME) \
	TCE_BUILD_VERSION=$(BUILD_VERSION) \
	FRAMEWORK_BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} ENVS="${ENVS}" hack/builder/install-tanzu.sh

.PHONY: clean-framework
clean-framework:
	rm -rf ${TCE_SCRATCH_DIR}
	rm -rf ${XDG_DATA_HOME}/tanzu-cli
	rm -rf ${XDG_CONFIG_HOME}/tanzu
	rm -rf ${XDG_CONFIG_HOME}/tanzu-plugins
	rm -rf ${XDG_CACHE_HOME}/tanzu
	mkdir -p ${XDG_DATA_HOME}/tanzu-cli
	mkdir -p ${XDG_CONFIG_HOME}/tanzu
	mkdir -p ${XDG_CONFIG_HOME}/tanzu-plugins
	mkdir -p ${XDG_CACHE_HOME}/tanzu

# TANZU CLI

# PLUGINS
# Dynamically generate OS-ARCH targets to allow for parallel execution
PLUGIN_JOBS_WITHOUT_PUBLISH := $(addprefix build-cli-plugins-,${ENVS})
PLUGIN_JOBS := $(addsuffix -publish,${PLUGIN_JOBS_WITHOUT_PUBLISH})

.PHONY: prep-build-cli
prep-build-cli:
	@cd ./cli/cmd/plugin/ && for plugin in */; do\
		printf "===> Preparing $${plugin}\n";\
		working_dir=`pwd`;\
		cd $${plugin};\
		$(MAKE) get-deps;\
		cd $${working_dir};\
	done

# we must call clean-plugin in order to not collide bits using 0.11.1. non-issue in 0.17.0
.PHONY: build-cli-plugins
build-cli-plugins: clean-plugin ${PLUGIN_JOBS}

# This builds all plugins but does not go through the publish step required for install
# The GCP update buckets use the old directory structure which errors out when calling "publish"
# Please see "prep-gcp-tanzu-bucket" for more info
.PHONY: build-cli-plugins-%
build-cli-plugins-%: prep-build-cli
	$(eval ARCH = $(word 2,$(subst -, ,$*)))
	$(eval OS = $(word 1,$(subst -, ,$*)))

	@printf "===> Building with ${OS}-${ARCH}\n";
	@cd ./hack/builder/ && $(MAKE) compile OS=${OS} ARCH=${ARCH} PLUGINS=${PLUGINS} DISCOVERY_NAME=${DISCOVERY_NAME} TANZU_CORE_BUCKET="tce-tanzu-cli-framework" TKG_DEFAULT_IMAGE_REPOSITORY=${TKG_DEFAULT_IMAGE_REPOSITORY} TKG_DEFAULT_COMPATIBILITY_IMAGE_PATH=${TKG_DEFAULT_COMPATIBILITY_IMAGE_PATH}

# This builds and publishes all plugins so we can install them
.PHONY: build-cli-plugins-%-publish
build-cli-plugins-%-publish: prep-build-cli
	$(eval ARCH = $(word 2,$(subst -, ,$*)))
	$(eval OS = $(word 1,$(subst -, ,$*)))

	@printf "===> Building with ${OS}-${ARCH}\n";
	@cd ./hack/builder/ && $(MAKE) compile publish OS=${OS} ARCH=${ARCH} PLUGINS=${PLUGINS} DISCOVERY_NAME=${DISCOVERY_NAME} TANZU_CORE_BUCKET="tce-tanzu-cli-framework" TKG_DEFAULT_IMAGE_REPOSITORY=${TKG_DEFAULT_IMAGE_REPOSITORY} TKG_DEFAULT_COMPATIBILITY_IMAGE_PATH=${TKG_DEFAULT_COMPATIBILITY_IMAGE_PATH}

# we must call clean-plugin in order to not collide bits using 0.11.1. non-issue in 0.17.0
.PHONY: build-cli-plugins-nopublish
build-cli-plugins-nopublish: clean-plugin ${PLUGIN_JOBS_WITHOUT_PUBLISH}

# install-plugins depends on install-cli for 2 reasons:
# 1. to install the plugins, we leverage the "tanzu plugin install" command which means the CLI is installed
# 2. if we are running a non GA build, this is what calls "set-unstable-versions" in framework for non-GA semvers
.PHONY: install-plugins
install-plugins: install-cli
	@cd ./hack/builder/ && $(MAKE) install-plugins DISCOVERY_NAME=$(DISCOVERY_NAME)

.PHONY: build-install-plugins
build-install-plugins: build-cli-plugins install-plugins
	@printf "CLI plugins built and installed\n"

test-plugins: ## run tests on TCE plugins
	# TODO(joshrosso): update once we get our testing strategy in place
	@echo "No tests to run."

.PHONY: clean-plugin
clean-plugin:
	rm -rf ${ARTIFACTS_DIR}
	rm -rf ./build
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

push-package: check-carvel # Verify openAPIv3 schema in package before build and push a package template. Tag will default to `latest`. Usage: make push-package PACKAGE=foobar VERSION=1.0.0
	@printf "\n===> pushing $${PACKAGE}/$${VERSION}\n";\
	./hack/packages/verify-openapischema-for-package.sh $(PACKAGE) $(VERSION) \
	&& cd addons/packages/$${PACKAGE}/$${VERSION} && imgpkg push --bundle $(OCI_REGISTRY)/$${PACKAGE}:$${VERSION} --file bundle/;\

generate-openapischema-package: #Generate package with OpenAPI v3 schema
	@printf "\n===> generating OpenAPIv3 schema for $${PACKAGE}/$${VERSION}\n";\
	./hack/packages/check-sample-values-and-render-ytt.sh $(PACKAGE) $(VERSION) \
	&& cd addons/packages/$${PACKAGE}/$${VERSION} \
	&& mkdir -p ${ARTIFACTS_DIR} \
	&& cd ${ARTIFACTS_DIR} \
	&& ytt -f ../bundle/config/schema.yaml --data-values-schema-inspect -o openapi-v3 > openapi-schema.yaml \
	&& ytt -f ../package.yaml -f ../../../package-overlay/package-overlay.yaml --data-value-file openapi=openapi-schema.yaml > generated-package.yaml \
	&& mv generated-package.yaml ../package.yaml
	@printf "===> package.yaml has been updated with openAPIv3 schema in its valuesSchema field for $${PACKAGE}/$${VERSION}\n";

export CHANNEL
generate-package-repo: check-carvel # Generate and push the package repository. Usage: make generate-package-repo CHANNEL=main
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

# AWS Management + Workload Cluster E2E Test
aws-management-and-workload-cluster-e2e-test:
	test/aws/deploy-tce-managed.sh

# Azure Management + Workload Cluster E2E Test
azure-management-and-workload-cluster-e2e-test:
	test/azure/deploy-management-and-workload-cluster.sh

# Docker Management + Workload Cluster E2E Test
docker-management-and-cluster-e2e-test:
	test/docker/run-tce-docker-managed-cluster.sh

# vSphere Management + Workload Cluster E2E Test
vsphere-management-and-workload-cluster-e2e-test:
	BUILD_VERSION=$(BUILD_VERSION) test/vsphere/run-tce-vsphere-management-and-workload-cluster.sh

##### E2E TESTS #####
