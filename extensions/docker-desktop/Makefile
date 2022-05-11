# Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

PLATFORMS=linux/amd64,linux/arm64
REGISTRY?=ghcr.io/
ORG?=vmware-tanzu/
IMAGE_NAME=community-edition
IMAGE?=$(REGISTRY)${ORG}${IMAGE_NAME}
EXTENSION_REGISTRY?=docker.io/
EXTENSION_IMAGE?=$(EXTENSION_REGISTRY)vmware/vmware-tanzu-community-edition-extension-for-docker-desktop
ifndef TAG
TAG=dev
endif
DEV_UI_SOURCE?=http://localhost:3000
BUILD_ARGS=--progress=plain

BUILDER=buildx-multi-arch

INFO_COLOR = \033[0;36m
NO_COLOR   = \033[m

.DEFAULT_GOAL:=help

help:  # Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[0-9A-Za-z_-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^\$$\([0-9A-Za-z_-]+\):.*?##/ { gsub("_","-", $$1); printf "  \033[36m%-45s\033[0m %s\n", tolower(substr($$1, 3, length($$1)-7)), $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##### BUILD-DEPS
.PHONY: build-tanzu-cli
build-tanzu-cli: ## Build the tanzu cli container image for local arch
	docker buildx build $(BUILD_ARGS) --tag=$(IMAGE)-tanzu-cli:$(TAG) -f Dockerfile-tanzu-cli .

.PHONY: build-push-tanzu-cli-all-arch
build-push-tanzu-cli-all-arch: ## Build and push the tanzu cli container image for all archs
	docker buildx inspect $(BUILDER)-tanzu-cli || docker buildx create --name=$(BUILDER)-tanzu-cli --driver=docker-container --driver-opt=network=host
	docker buildx build --push $(BUILD_ARGS) --push --builder=$(BUILDER)-tanzu-cli --tag=$(IMAGE)-tanzu-cli:$(TAG) -f Dockerfile-tanzu-cli --platform=$(PLATFORMS) .
	docker buildx stop $(BUILDER)-tanzu-cli

.PHONY: build-apps-bin
build-apps-bin: ## Build the dashboard container image for local arch
	docker buildx build $(BUILD_ARGS) --tag=$(IMAGE)-apps-bin:$(TAG) -f Dockerfile-apps-bin .

.PHONY: build-push-apps-bin-all-arch
build-push-apps-bin-all-arch: ## Build and push the dashboard container image for all archs
	docker buildx inspect $(BUILDER)-apps-bin || docker buildx create --name=$(BUILDER)-apps-bin --driver=docker-container --driver-opt=network=host
	docker buildx build --push $(BUILD_ARGS) --push --builder=$(BUILDER)-apps-bin --tag=$(IMAGE)-apps-bin:$(TAG) -f Dockerfile-apps-bin --platform=$(PLATFORMS) .
	docker buildx stop $(BUILDER)-apps-bin

.PHONY: build-downloader
build-downloader: ## Build the binary downloader container image for local arch
	docker buildx build $(BUILD_ARGS) --tag=$(IMAGE)-downloader:$(TAG) -f Dockerfile-downloader .

.PHONY: build-push-downloader-all-arch
build-push-downloader-all-arch: ## Build and push the binary downloader container image for all archs
	docker buildx inspect $(BUILDER)-downloader || docker buildx create --name=$(BUILDER)-downloader --driver=docker-container --driver-opt=network=host
	docker buildx build --push $(BUILD_ARGS) --push --builder=$(BUILDER)-downloader --tag=$(IMAGE)-downloader:$(TAG) -f Dockerfile-downloader --platform=$(PLATFORMS) .
	docker buildx stop $(BUILDER)-downloader

.PHONY: build-client
build-client: ## Build the UI container image for local arch
	docker buildx build $(BUILD_ARGS) --tag=$(IMAGE)-client:$(TAG) -f Dockerfile-client .

.PHONY: build-push-client-all-arch
build-push-client-all-arch: ## Build and push the UI container image for all archs
	docker buildx inspect $(BUILDER)-client || docker buildx create --name=$(BUILDER)-client --driver=docker-container --driver-opt=network=host
	docker buildx build --push $(BUILD_ARGS) --push --builder=$(BUILDER)-client --tag=$(IMAGE)-client:$(TAG) -f Dockerfile-client --platform=$(PLATFORMS) .
	docker buildx stop $(BUILDER)-client
##### BUILD-DEPS

##### BUILD
.PHONY: build-deps
build-deps: ## Build the required images used in the extension image for local arch
	make -j build-downloader build-tanzu-cli build-apps-bin build-client

.PHONY: build-extension
build-extension: build-deps ## Build the extension image for local architecture
	docker buildx build $(BUILD_ARGS) --tag=$(EXTENSION_IMAGE):$(TAG) --build-arg TAG=$(TAG) --build-arg REGISTRY=$(REGISTRY) --build-arg ORG=${ORG} --build-arg IMAGE_NAME=${IMAGE_NAME} .
##### BUILD

##### DOCKER
.PHONY: install
install: ## Install the extension or update it if already exists
	docker extension install $(EXTENSION_IMAGE):$(TAG) || docker extension update $(EXTENSION_IMAGE):$(TAG)

.PHONY: build-install 
build-install: build-extension install ## Build extension images for local arch and install

.PHONY: update
update: ## Pull and update the extension with a new image
	docker pull $(EXTENSION_IMAGE):$(TAG) && docker extension update $(EXTENSION_IMAGE):$(TAG)

.PHONY: debug
debug: ## Enable debug in the extension
	docker extension dev debug $(EXTENSION_IMAGE):$(TAG)

.PHONY: source
source: ## Replace the UI source of the extension
	docker extension dev ui-source $(EXTENSION_IMAGE):$(TAG) $(DEV_UI_SOURCE)

.PHONY: validate
validate: ## Validate the extension
	docker extension validate $(EXTENSION_IMAGE):$(TAG)

.PHONY: dev-reset
dev-reset: ## Reset development status of the extension
	docker extension dev reset $(EXTENSION_IMAGE):$(TAG)

.PHONY: delete
delete: ## Remove the extension
	docker extension rm $(EXTENSION_IMAGE):$(TAG)
##### DOCKER

##### PUBLISH
.PHONY: check-ci
check-ci:
	@if [[ "${TCE_CI_BUILD}" != "true" ]]; then echo "WARNING: Publishing the Docker Desktop extension is only meant to be run within GitHub Actions CI or manually with extreme caution" && exit 1; else echo "Passed CI check"; fi

.PHONY: prepare-buildx
prepare-buildx: 
	docker buildx inspect $(BUILDER) || docker buildx create --name=$(BUILDER) --driver=docker-container --driver-opt=network=host

.PHONY: stop-buildx
stop-buildx:
	docker buildx stop $(BUILDER)

.PHONY: push-extension
push-extension:
	# Do not push if tag already exists: TAG=0.1.0 make push-extension
	docker manifest inspect $(EXTENSION_IMAGE):$(TAG) && echo "Failure: Tag ${TAG} already exists" || docker buildx build --push --builder=$(BUILDER) --platform=$(PLATFORMS) --build-arg TAG=$(TAG) --tag=$(EXTENSION_IMAGE):$(TAG) .

.PHONY: build-push-deps-all-arch
build-push-deps-all-arch: check-ci ## Build and push the required images used in the extension image for a multi-arch build
	make -j build-push-downloader-all-arch build-push-tanzu-cli-all-arch build-push-apps-bin-all-arch build-push-client-all-arch

.PHONY: build-push-extension
build-push-extension: check-ci prepare-buildx push-extension stop-buildx ## For all architectures, build the extension and upload to the registry

.PHONY: build-push-everything
build-push-everything: check-ci prepare-buildx build-push-deps-all-arch push-extension stop-buildx ## For all architectures, build dependent containers, build extension, and upload all image to the registries
##### PUBLISH
