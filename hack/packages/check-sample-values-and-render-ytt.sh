#!/usr/bin/env bash

# Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# this shell script checks the existence of a sample values file for a Tanzu package and verifies if ytt templates are generated.

# set this value to your package name
PACKAGE=$1

# set this value to the version of your package
VERSION=$2

if [ -z "$PACKAGE" ]
then
  echo "check sample values for package failed. must set PACKAGE"
  exit 2
fi

if [ -z "$VERSION" ]
then
  echo "check sample values for package failed. must set VERSION"
  exit 2
fi

ROOT_DIR="addons/packages"
PACKAGE_DIR="${ROOT_DIR}/${PACKAGE}"
VERSION_DIR="${PACKAGE_DIR}/${VERSION}"
BUNDLE_DIR="${VERSION_DIR}/bundle"
CONFIG_DIR="${BUNDLE_DIR}/config"
NC='\033[0m' # No Color
GREEN='\033[0;32m'
RED='\033[0;31m'

check_sample_values_and_render_ytt() {
  sample_values_dir="${VERSION_DIR}/sample-values"

  yttCmd="ytt -f ."
  if [ -d "${sample_values_dir}" ]
  then
    yttCmd="${yttCmd} -f ../../sample-values/*.yaml"
  fi
  cd "${CONFIG_DIR}" || exit
	${yttCmd} > /dev/null
	status=$?

	if [ $status -eq 0 ]; then
	  echo -e "${GREEN}===> ytt manifests successfully rendered for ${PACKAGE}/${VERSION}${NC}"
	else
	  echo -e "${RED}===> $yttCmd failed. ytt manifests could not be generated!!${NC}"
	  exit 1
	fi

}

check_sample_values_and_render_ytt



