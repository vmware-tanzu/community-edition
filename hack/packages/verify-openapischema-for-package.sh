#!/bin/sh

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
BUNDLE_DIR="bundle"
CONFIG_DIR="config"
PACKAGE_DIR="${ROOT_DIR}/${PACKAGE}"
VERSION_DIR="${PACKAGE_DIR}/${VERSION}"
BUNDLE_DIR="${VERSION_DIR}/bundle"
CONFIG_DIR="${BUNDLE_DIR}/config"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

verify_openapischema_for_package() {
  cd "${CONFIG_DIR}"
	ytt -f schema.yaml --data-values-schema-inspect -o openapi-v3 > openapi-schema.yaml
	yq e '.components.schemas.dataValues' openapi-schema.yaml > file1
	yq e '.spec.valuesSchema.openAPIv3' ../../package.yaml > file2
	diffyaml file1 file2

	# Remove generated files if no diff exists
	rm file1 file2 openapi-schema.yaml
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
		echo -e "${RED}Error: Found diff in ${YELLOW}$file1 ${RED}and ${YELLOW}$file2${NC}"
		echo -e "${RED}Please examine both files${NC}"
		echo -e "${RED}Generate the openAPIv3 schema for the package before running make push-package , Run:${NC}"
		echo -e "${YELLOW}make generate-openapischema-package${NC}"
		exit 1
	fi
}

verify_openapischema_for_package



