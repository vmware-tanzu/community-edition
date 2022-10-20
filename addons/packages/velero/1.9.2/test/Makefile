# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

ROOT_DIR := $(shell git rev-parse --show-toplevel)
PLATFORM := $(shell uname | tr '[:upper:]' '[:lower:]')
VERSION := 1.9.2
GOPATH ?= "~/go"
GINKGO_FOCUS ?= Basic
# Note: Because we are not asking the Velero e2e tests to install Velero, the CLOUD_PROVIDER variable is only meant to indicate
# if the tests are being run in an environment that supports taking volume snapshots, like AWS, Azure or vSphere. If the
# tests are running on Docker, set this variable to "kind". Otherwise, set it to anything else.
CLOUD_PROVIDER ?= notkind
REGISTRY_CREDENTIAL_FILE ?=

e2e-test:
	sh -c "rm -rf /tmp/velero-e2e-test &&\
	mktemp -d /tmp/velero-e2e-test.XXX &&\
	git clone https://github.com/vmware-tanzu/velero /tmp/velero-e2e-test/velero &&\
	cd /tmp/velero-e2e-test &&\
	$(ROOT_DIR)/addons/packages/velero/${VERSION}/test/install-velero-cli.sh v${VERSION} $(PLATFORM) &&\
	ACK_GINKGO_DEPRECATIONS=1.16.4 \
	CREDS_FILE=dummyvalue BSL_BUCKET=dummyvalue \
	GOPATH=$(GOPATH) CLOUD_PROVIDER=$(CLOUD_PROVIDER) REGISTRY_CREDENTIAL_FILE=$(REGISTRY_CREDENTIAL_FILE)  \
	GINKGO_FOCUS=$(GINKGO_FOCUS) INSTALL_VELERO=false VELERO_CLI=/tmp/velero-e2e-test/velero-v${VERSION}-$(PLATFORM)-amd64/velero  \
	make -C /tmp/velero-e2e-test/velero/test/e2e run"
