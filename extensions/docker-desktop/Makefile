# Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

REGISTRY?=ghcr.io/
ORG?=vmware-tanzu/
IMAGE_NAME=community-edition-extension-for-dockerdesktop
IMAGE?=$(REGISTRY)${ORG}${IMAGE_NAME}
EXTENSION_REGISTRY?=docker.io/
EXTENSION_IMAGE?=$(EXTENSION_REGISTRY)vmware/vmware-tanzu-community-edition-extension-for-docker-desktop
TAG?=main
DEV_UI_SOURCE?=http://localhost:3000
BUILD_ARGS=--progress=plain

BUILDER=buildx-multi-arch

INFO_COLOR = \033[0;36m
NO_COLOR   = \033[m

.DEFAULT_GOAL := extension

build-tanzucli: ## Build the tanzu cli container image
	docker buildx build $(BUILD_ARGS) --tag=$(IMAGE)-tanzu-cli:$(TAG) -f Dockerfile-tanzu-cli .
.PHONY: build-tanzucli

build-appsbin: ## Build the dashboard container image
	docker buildx build $(BUILD_ARGS) --tag=$(IMAGE)-apps-bin:$(TAG) -f Dockerfile-apps-bin .
.PHONY: build-appsbin

build-downloader: ## Build the binary downloader container image
	docker buildx build $(BUILD_ARGS) --tag=$(IMAGE)-downloader:$(TAG) -f Dockerfile-downloader .
.PHONY: build-downloader

build-client: ## Build the UI container image
	docker buildx build $(BUILD_ARGS) --tag=$(IMAGE)-client:$(TAG) -f Dockerfile-client .
.PHONY: build-client

build-deps: ## Build the required images used in the extension image
	make -j build-downloader build-tanzucli build-appsbin build-client
.PHONY: build-deps

build-extension: ## Build the extension container image
	docker build $(BUILD_ARGS) --tag=$(EXTENSION_IMAGE):$(TAG) --build-arg TAG=$(TAG) --build-arg REGISTRY=$(REGISTRY) --build-arg ORG=${ORG} --build-arg IMAGE_NAME=${IMAGE_NAME} .
.PHONY: build-extension

extension: ## Build the extension and its dependencies
	make build-deps build-extension
.PHONY: extension

install: ## Install the extension or update it if already exists
	docker extension install $(EXTENSION_IMAGE):$(TAG) || docker extension update $(EXTENSION_IMAGE):$(TAG)

update: ## Update the extension with a new image
	docker pull $(EXTENSION_IMAGE):$(TAG) && docker extension update $(EXTENSION_IMAGE):$(TAG)

debug: ## Enable debug in the extension
	docker extension dev debug $(EXTENSION_IMAGE):$(TAG)

source: ## Replace the ui source of the extension
	docker extension dev ui-source $(EXTENSION_IMAGE):$(TAG) $(DEV_UI_SOURCE)

validate: ## Validate the extension
	docker extension validate $(EXTENSION_IMAGE):$(TAG)

dev-reset: ## Reset development status of the extension
	docker extension dev reset $(EXTENSION_IMAGE):$(TAG)

delete: ## Remove the extension
	docker extension rm $(EXTENSION_IMAGE):$(TAG)

prepare-buildx: ## Create buildx builder for multi-arch build, if not exists
	docker buildx inspect $(BUILDER) || docker buildx create --name=$(BUILDER) --driver=docker-container --driver-opt=network=host

push-extension: prepare-buildx ## Build & upload extension image to a registry
	# Do not push if tag already exists: make push-extension tag=0.1
	docker pull $(EXTENSION_IMAGE):$(tag) && echo "Failure: Tag already exists" || docker buildx build --push --builder=$(BUILDER) --platform=linux/amd64,linux/arm64 --build-arg TAG=${tag)} --tag=$(EXTENSION_IMAGE):$(tag) .

help: ## Show this help
	@echo Please specify a build target. The choices are:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(INFO_COLOR)%-30s$(NO_COLOR) %s\n", $$1, $$2}'
