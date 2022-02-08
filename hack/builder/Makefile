# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

.DEFAULT_GOAL:=help

BUILD_DATE ?= $$(date -u +"%Y-%m-%d")
BUILD_SHA ?= $$(git rev-parse --short HEAD)
BUILD_EDITION=tce-standalone
OS ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)

ifndef TKG_DEFAULT_IMAGE_REPOSITORY
TKG_DEFAULT_IMAGE_REPOSITORY = "projects-stg.registry.vmware.com/tkg"
endif
ifndef TKG_DEFAULT_COMPATIBILITY_IMAGE_PATH
TKG_DEFAULT_COMPATIBILITY_IMAGE_PATH = "framework-zshippable/tkg-compatibility"
endif

ifndef IS_OFFICIAL_BUILD
IS_OFFICIAL_BUILD = ""
endif

ifndef TANZU_PLUGIN_UNSTABLE_VERSIONS
TANZU_PLUGIN_UNSTABLE_VERSIONS = "experimental"
endif

ifeq ($(strip $(BUILD_VERSION)),)
BUILD_VERSION = $(shell git describe --tags --abbrev=0)
endif

# TODO: Push this down into the conformance plugin's own Makefile
# The > /dev/null redirects are so we don't get the output mixed with the version number.
SONO_VERSION := $(shell working_dir=`pwd`; cd ../../cli/cmd/plugin/conformance; $(GO) list -f '{{.Version}}' -m github.com/vmware-tanzu/sonobuoy; cd $$working_dir;)

LD_FLAGS = -s -w
LD_FLAGS += -X "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgconfigpaths.TKGDefaultImageRepo=$(TKG_DEFAULT_IMAGE_REPOSITORY)"
LD_FLAGS += -X "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgconfigpaths.TKGDefaultCompatibilityImagePath=$(TKG_DEFAULT_COMPATIBILITY_IMAGE_PATH)"
LD_FLAGS += -X "github.com/vmware-tanzu/tanzu-framework/pkg/v1/buildinfo.Date=$(BUILD_DATE)"
LD_FLAGS += -X "github.com/vmware-tanzu/tanzu-framework/pkg/v1/buildinfo.SHA=$(BUILD_SHA)"
LD_FLAGS += -X "github.com/vmware-tanzu/tanzu-framework/pkg/v1/buildinfo.Version=$(BUILD_VERSION)"
LD_FLAGS += -X 'main.BuildEdition=$(BUILD_EDITION)'
LD_FLAGS += -X 'github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/buildinfo.IsOfficialBuild=$(IS_OFFICIAL_BUILD)'
LD_FLAGS += -X 'github.com/vmware-tanzu/community-edition/cli/cmd/plugin.Date=$(BUILD_DATE)'
LD_FLAGS += -X 'github.com/vmware-tanzu/community-edition/cli/cmd/plugin.SHA=$(BUILD_SHA)'
LD_FLAGS += -X 'github.com/vmware-tanzu/community-edition/cli/cmd/plugin.Version=$(BUILD_VERSION)'

ifneq ($(strip $(TANZU_CORE_BUCKET)),)
LD_FLAGS += -X 'github.com/vmware-tanzu/tanzu-framework/pkg/v1/config.CoreBucketName=$(TANZU_CORE_BUCKET)'
endif

# TODO: Push this down into the conformance plugin's own Makefile
LD_FLAGS += -X "github.com/vmware-tanzu/sonobuoy/pkg/buildinfo.Version=$(SONO_VERSION)"

help: ## Display this help message
	@printf "The builder is used to pull in necessary build dependencies, primarily, Tanzu Framework"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[0-9A-Za-z_-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

lint:
	echo "N/A: Skipping linting builder"

test:
	echo "N/A: No unit tests for hack/builder"

e2e-test:
	echo "N/A: No e2e tests for hack/builder"

build:
	echo "N/A: No binary output for hack/builder"

get-deps: ## Download the necessary dependencies for the builder
	go mod tidy

check-go:
ifeq ($(origin GO),undefined)
	$(error Error! GO env var not set)
endif

check-build-version:
ifeq ($(origin BUILD_VERSION),undefined)
	$(error Error! BUILD_VERSION env var not set)
endif

check-ld-flags:
ifeq ($(origin LD_FLAGS),undefined)
	$(error Error! LD_FLAGS env var not set)
endif

check-artifacts:
ifeq ($(origin ARTIFACTS_DIR),undefined)
	$(error Error! ARTIFACTS_DIR env var not set)
endif

check: check-go check-build-version check-ld-flags check-artifacts

framework-set-unstable-versions:
	TANZU_CLI_NO_INIT=true $(GO) run -ldflags "$(LD_FLAGS)" github.com/vmware-tanzu/tanzu-framework/cmd/cli/tanzu config set unstable-versions $(TANZU_PLUGIN_UNSTABLE_VERSIONS)

# TODO: Post-MVP, this needs to be able to compile individual plugins with their own Makefiles.
# As it stands, we cannot customize LDFLAGS per plugin
compile: check ## Compiles the TCE plugins using the Tanzu Framework builder
	$(GO) run github.com/vmware-tanzu/tanzu-framework/cmd/cli/plugin-admin/builder cli compile \
	--version $(BUILD_VERSION) --ldflags "$(LD_FLAGS)" --path ../../cli/cmd/plugin \
	--ldflags "$(LD_FLAGS) -X 'github.com/vmware-tanzu/tanzu-framework/pkg/v1/config.DefaultStandaloneDiscoveryType=local'" \
	--artifacts ../../${ARTIFACTS_DIR} --target local

install-plugins: check framework-set-unstable-versions ## Installs the compiled TCE plugins
	TANZU_CLI_NO_INIT=true $(GO) run -ldflags "$(LD_FLAGS)" github.com/vmware-tanzu/tanzu-framework/cmd/cli/tanzu \
		plugin install all --local ../../$(ARTIFACTS_DIR)



# .PHONY: install-cli-plugins-without-discovery
# install-cli-plugins-without-discovery: set-unstable-versions set-context-aware-cli-for-plugins ## Install Tanzu CLI plugins when context-aware discovery is disabled
# 	TANZU_CLI_NO_INIT=true $(GO) run -ldflags "$(LD_FLAGS)" ./cmd/cli/tanzu/main.go \
# 		plugin install all --local $(ARTIFACTS_DIR)/$(GOHOSTOS)/$(GOHOSTARCH)/cli
		
# .PHONY: build-cli-local
# build-cli-local: configure-buildtags-embedproviders build-cli-local-${GOHOSTOS}-${GOHOSTARCH} publish-plugins-local ## Build Tanzu CLI with local standalone discovery. cluster and management-cluster plugins are built with embedded providers.
# 	$(GO) run ./cmd/cli/plugin-admin/builder/main.go cli compile --version $(BUILD_VERSION) --ldflags "$(LD_FLAGS) -X 'github.com/vmware-tanzu/tanzu-framework/pkg/v1/config.DefaultStandaloneDiscoveryType=local'" --tags "${BUILD_TAGS}" --path ./cmd/cli/plugin-admin --artifacts artifacts-admin/${GOHOSTOS}/${GOHOSTARCH}/cli --target local
# 	$(MAKE) publish-admin-plugins-local

# .PHONY: publish-admin-plugins-local
# publish-admin-plugins-local: ## Publish CLI admin plugins locally for current host os-arch only
# 	$(GO) run ./cmd/cli/plugin-admin/builder/main.go publish --type local --plugins "$(ADMIN_PLUGINS)" --version $(BUILD_VERSION) --os-arch "${GOHOSTOS}-${GOHOSTARCH}" --local-output-discovery-dir "$(TANZU_PLUGIN_PUBLISH_PATH)/discovery/admin" --local-output-distribution-dir "$(TANZU_PLUGIN_PUBLISH_PATH)/distribution" --input-artifact-dir $(ARTIFACTS_ADMIN_DIR)