# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

ROOT_DIR := $(shell git rev-parse --show-toplevel)
PLATFORM := $(shell uname | tr '[:upper:]' '[:lower:]')
ifndef GOPATH
GOPATH := "~/go"
endif
ifndef PROVIDER
CLOUD_PROVIDER := "minio"
endif

e2e-test:
	sh -c "rm -rf /tmp/velero-e2e-test &&\
	mktemp -d /tmp/velero-e2e-test &&\
	git clone https://github.com/vmware-tanzu/velero /tmp/velero-e2e-test/velero &&\
	cd /tmp/velero-e2e-test &&\
	$(ROOT_DIR)/addons/packages/velero/1.5.2/test/install-velero-cli.sh v1.5.2 $(PLATFORM) &&\
	GOPATH=$(GOPATH) CLOUD_PROVIDER=$(CLOUD_PROVIDER) \
	OBJECT_STORE_PROVIDER=minio \
	CREDS_FILE=/tmp/credential BSL_BUCKET=velero \
	GINKGO_FOCUS=Basic INSTALL_VELERO=false VELERO_CLI=/tmp/velero-e2e-test/velero-v1.5.2-$(PLATFORM)-amd64/velero  \
	make -C /tmp/velero-e2e-test/velero/test/e2e run"

