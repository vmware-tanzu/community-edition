#!/bin/sh

# this shell script creates the prescribed directory structure for a Tanzu package.

# set this value to your package name
NAME=$1

# set this value to the version of your package
VERSION=$2

if [ -z "$NAME" ]
then
  echo "create package failed. must set NAME"
  exit 2
fi

if [ -z "$VERSION" ]
then
  echo "create package failed. must set VERSION"
  exit 2
fi

# Handle differences in MacOS sed
SEDARGS=""
if [ "$(uname -s)" = "Darwin" ]; then
    SEDARGS="-e"
fi

ROOT_DIR="addons/packages"
BUNDLE_DIR="bundle"
CONFIG_DIR="config"
OVERLAY_DIR="overlay"
UPSTREAM_DIR="upstream"
IMGPKG_DIR=".imgpkg"
PACKAGE_DIR="${ROOT_DIR}/${NAME}"
VERSION_DIR="${PACKAGE_DIR}/${VERSION}"

# create directory structure for package
mkdir -vp "${VERSION_DIR}/${BUNDLE_DIR}/${CONFIG_DIR}"
mkdir -v "${VERSION_DIR}/${BUNDLE_DIR}/${IMGPKG_DIR}"
mkdir -v "${VERSION_DIR}/${BUNDLE_DIR}/${CONFIG_DIR}/${OVERLAY_DIR}"
mkdir -v "${VERSION_DIR}/${BUNDLE_DIR}/${CONFIG_DIR}/${UPSTREAM_DIR}"

# create README and fill with name of package
sed $SEDARGS "s/PACKAGE_NAME/${NAME}/g" hack/packages/templates/readme.md > "${VERSION_DIR}/README.md"

# create manifests and fill with name of package
sed $SEDARGS "s/PACKAGE_NAME/${NAME}/g" hack/packages/templates/metadata.yaml > "${PACKAGE_DIR}/metadata.yaml"
sed $SEDARGS "s/PACKAGE_NAME/${NAME}/g" hack/packages/templates/package.yaml > "${VERSION_DIR}/package_a.yaml"
sed $SEDARGS "s/VERSION/${VERSION}/g" "${VERSION_DIR}/package_a.yaml" > "${VERSION_DIR}/package.yaml"
sed $SEDARGS "s/PACKAGE_NAME/${NAME}/g" hack/packages/templates/values.yaml > "${VERSION_DIR}/bundle/config/values.yaml"
sed $SEDARGS "s/PACKAGE_NAME/${NAME}/g" hack/packages/templates/vendir.yml > "${VERSION_DIR}/bundle/vendir.yml"

rm "${VERSION_DIR}/package_a.yaml"

echo
echo "package bootstrapped at ${VERSION_DIR}"
echo
