#!/usr/bin/env bash

# Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# this shell script checks the existence os a sample values file for a Tanzu package and verifies if ytt templates are generated.

# set this value to your package name
PACKAGE=$1

# set this value to the version of your package
VERSION=$2

if [ -z "$PACKAGE" ]
then
  echo "check sample values for package failed. must set NAME"
  exit 2
fi

if [ -z "$VERSION" ]
then
  echo "check sample values for package failed. must set VERSION"
  exit 2
fi

ROOT_DIR="addons/packages"
BUNDLE_DIR="bundle"
CONFIG_DIR="config"
PACKAGE_DIR="${ROOT_DIR}/${PACKAGE}"
VERSION_DIR="${PACKAGE_DIR}/${VERSION}"
BUNDLE_DIR="${VERSION_DIR}/bundle"
CONFIG_DIR="${BUNDLE_DIR}/config"

check_sample_values_existence_for_package() {
  sample_values_dir="${BUNDLE_DIR}/sample-values"
  echo "$sample_values_dir"
  yttCmd="ytt -f ."
  if [ -d "${sample_values_dir}" ]
  then
    yttCmd="${yttCmd} -f ../sample-values/*.yaml"
  fi
  cd "${CONFIG_DIR}" || exit
	${yttCmd}
}

check_sample_values_existence_for_package



