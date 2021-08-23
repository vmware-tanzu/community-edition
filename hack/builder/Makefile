.DEFAULT_GOAL:=help

BUILD_DATE ?= $$(date -u +"%Y-%m-%d")
BUILD_SHA ?= $$(git rev-parse --short HEAD)
BUILD_EDITION=tce-standalone

ifndef IS_OFFICIAL_BUILD
IS_OFFICIAL_BUILD = ""
endif

ifndef BUILD_VERSION
BUILD_VERSION ?= $$(git describe --tags --abbrev=0)
endif

ifeq ($(strip $(BUILD_VERSION)),)
BUILD_VERSION = dev
endif

LD_FLAGS = -s -w
LD_FLAGS += -X "github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli.BuildDate=$(BUILD_DATE)"
LD_FLAGS += -X "github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli.BuildSHA=$(BUILD_SHA)"
LD_FLAGS += -X "github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli.BuildVersion=$(BUILD_VERSION)"
LD_FLAGS += -X 'main.BuildEdition=$(BUILD_EDITION)'
LD_FLAGS += -X 'github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/buildinfo.IsOfficialBuild=$(IS_OFFICIAL_BUILD)'

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
	go mod download

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

compile: check ## Compiles the TCE plugins using the Tanzu Framework builder
	$(GO) run github.com/vmware-tanzu/tanzu-framework/cmd/cli/plugin-admin/builder cli compile --version $(BUILD_VERSION) \
		--ldflags "$(LD_FLAGS)" --path ../../cli/cmd/plugin --artifacts ../../${ARTIFACTS_DIR}

install-plugins: check ## Installs the compiled TCE plugins
	TANZU_CLI_NO_INIT=true $(GO) run -ldflags "$(LD_FLAGS)" github.com/vmware-tanzu/tanzu-framework/cmd/cli/tanzu \
		plugin install all --local ../../$(ARTIFACTS_DIR)

