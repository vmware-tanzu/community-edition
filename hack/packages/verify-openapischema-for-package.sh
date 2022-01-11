#!/usr/bin/env bash

# Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# this shell script checks the existence of openAPIv3 schema in package before pushing the imgpkg bundle for package.

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
ARTIFACTS_DIR="${VERSION_DIR}/artifacts"
BUNDLE_DIR="${VERSION_DIR}/bundle"
CONFIG_DIR="${BUNDLE_DIR}/config"
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color
GREEN='\033[0;32m'

verify_openapischema_for_package() {
  schema_file="${CONFIG_DIR}/schema.yaml"
  if [ -f "${schema_file}" ]; then
    mkdir -p "${ARTIFACTS_DIR}"
    cd "${ARTIFACTS_DIR}" || exit
    ytt -f ../bundle/config/schema.yaml --data-values-schema-inspect -o openapi-v3 > generated-openapi-schema.yaml
    status=$?
    if [ $status -eq 0 ]; then
      yq e '.components.schemas.dataValues' generated-openapi-schema.yaml > schema-contents.yaml
      yq e '.spec.valuesSchema.openAPIv3' ../package.yaml > package-schema-contents.yaml
      diffyaml schema-contents.yaml package-schema-contents.yaml
      echo -e "${GREEN}===> OpenAPIv3 contents match successful for schema and package${NC}"
    else
      echo -e "${RED}===> ytt manifests could not be generated!!${NC}"
      exit 1
    fi
  fi
}

diffyaml() {
  local -r file1="$1"
  local -r file2="$2"

  # Use kapp to intelligently diff YAML files.
  # Get rid of all the output lines that say there is no difference.
  # Get rid of the output line that says "Succeeded".
  # Expect that there are no output lines leftover in the output.
  if eval "$(kapp tools diff --file "$file1" --file2 "$file2" --changes --json \
    | yq e '.Lines[]' --tojson - \
    | grep -v '@@ noop' \
    | grep -v Succeeded || true)"; then
		return 0
	else
		echo -e "${RED}Error: Found diff in ${YELLOW}${ARTIFACTS_DIR}/$file1 ${RED}and ${YELLOW}${ARTIFACTS_DIR}/$file2${NC}"
		echo -e "${RED}openAPIv3 schema in package doesn't match what is defined in schema.yaml${NC}"
		echo -e "${RED}Generate the openAPIv3 schema for the package before running make push-package , Run:${NC}"
		echo -e "${YELLOW}make generate-openapischema-package PACKAGE=${PACKAGE} VERSION=${VERSION} ${NC}"
		exit 1
	fi
}

verify_openapischema_for_package



