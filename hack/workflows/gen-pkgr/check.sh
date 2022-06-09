#!/usr/bin/env bash

# Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

check_if_pkgr_should_be_generated() {
  cd hack/workflows/gen-pkgr && go run check-for-um-package.go "$1"
  status=$?
  if [ $status -eq 0 ]; then
    echo "generate"
    exit 0
  else
    echo "donotgenerate"
    exit 1
  fi
}

check_if_pkgr_should_be_generated "$1"
