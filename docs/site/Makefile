# Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

prepare:
	$(CURDIR)/load-docs.sh

production-build: prepare
	hugo --gc --minify

preview-build: prepare
	hugo \
		--baseURL $(DEPLOY_PRIME_URL) \
		--buildDrafts \
		--buildFuture 
