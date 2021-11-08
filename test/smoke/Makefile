# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

ROOT_DIR := $(shell git rev-parse --show-toplevel)
MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

smoke-test: install-ginkgo  cert-manager contour fluent-bit local-path-storage external-dns multus-cni 

install-ginkgo:
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega/...

contour:
	./packages/contour/1.18.2/contour-test.sh

cert-manager:
	./packages/cert-manager/1.5.3/cert-manager-test.sh

fluent-bit:
	./packages/fluent-bit/1.7.5/fluent-bit-test.sh

local-path-storage:
	./packages/local-path-storage/0.0.20/local-path-storage-test.sh

external-dns:
	cd "${ROOT_DIR}"/test/e2e && ginkgo -v -- --packages="external-dns" --version="0.8.0" --guest-cluster-name="tce.public"

multus-cni:
	cd "${ROOT_DIR}"/test/e2e && ginkgo -v -- --packages="multus-cni" --version="3.7.1" --guest-cluster-name="tce.public"
