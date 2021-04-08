#!/bin/sh

# this shell script creates the prescribed directory structure for a tanzu package.

# set this value to your package name
NAME=$1

if [ -z "$NAME" ]
then
  # this name var comes a Makefile
  # kinda hacky, should figure out a better way
  echo "create package failed. must set NAME env variable!"
  exit 2
fi

# Handle differences in MacOS sed
SEDARGS="-i"
if [ "$(uname -s)" = "Darwin" ]; then
    SEDARGS="-e"
fi

ROOT_DIR="addons/packages"
REPO_DIR="addons/repos"
BUNDLE_DIR="bundle"
CONFIG_DIR="config"
OVERLAY_DIR="overlay"
UPSTREAM_DIR="upstream"
IMGPKG_DIR=".imgpkg"
DIR="${ROOT_DIR}/${NAME}"

# create directory structure for package
mkdir -vp "${DIR}/${BUNDLE_DIR}/${CONFIG_DIR}"
mkdir -v "${DIR}/${BUNDLE_DIR}/${IMGPKG_DIR}"
mkdir -v "${DIR}/${BUNDLE_DIR}/${CONFIG_DIR}/${OVERLAY_DIR}"
mkdir -v "${DIR}/${BUNDLE_DIR}/${CONFIG_DIR}/${UPSTREAM_DIR}"

# create README and fill with name of package
# cp docs/package-templates/readme.md ${DIR}/README.md
sed $SEDARGS "s/PACKAGE_NAME/${NAME}/g" docs/package-templates/readme.md > "${DIR}/README.md"

# create manifests and fill with name of package
sed $SEDARGS "s/PACKAGE_NAME/${NAME}/g" docs/package-templates/clusterrolebinding.yaml > "${DIR}/clusterrolebinding.yaml"
sed $SEDARGS "s/PACKAGE_NAME/${NAME}/g" docs/package-templates/installedpackage.yaml > "${DIR}/installedpackage.yaml"
sed $SEDARGS "s/PACKAGE_NAME/${NAME}/g" docs/package-templates/serviceaccount.yaml > "${DIR}/serviceaccount.yaml"
sed $SEDARGS "s/PACKAGE_NAME/${NAME}/g" docs/package-templates/values.yaml > "${DIR}/bundle/config/values.yaml"
sed $SEDARGS "s/PACKAGE_NAME/${NAME}/g" docs/package-templates/vendir.yml > "${DIR}/bundle/vendir.yml"
sed $SEDARGS "s/PACKAGE_NAME/${NAME}/g" docs/package-templates/package.yaml > "${REPO_DIR}/main/packages/${NAME}.yml"

echo
echo "package boostrapped at ${DIR}"
echo
