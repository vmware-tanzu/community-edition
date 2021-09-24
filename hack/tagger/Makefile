# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

REGISTRY ?= projects.registry.vmware.com/tce/tagger
TAG ?= dev

get-deps:
	go mod tidy

build:
	DOCKER_BUILDKIT=1 docker build . -t ${REGISTRY}:${TAG}

build-and-push: build
	docker push ${REGISTRY}:${TAG}

run:
	go run main.go

test:
	echo "N/A: No unit tests for hack/packages"

e2e-test:
	echo "N/A: No e2e tests for hack/packages"

lint:
	echo "N/A: No linters for hack/packages"

