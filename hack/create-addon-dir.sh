#!/bin/sh

# this shell script creates the prescribed directory structure for a tanzu add-on.

# set this value to your add-on name
EXT_NAME=$1

if [ -z "$EXT_NAME" ]
then
  # this name var comes a Makefile
  # kinda hacky, should figure out a better way
  echo "create add-on failed. must set NAME env variable!"
  exit 2
fi

# Handle differences in MacOS sed
SEDARGS="-i"
if [ "$(uname -s)" = "Darwin" ]; then
    SEDARGS="-i '' -e"
fi

EXT_ROOT_DIR="addons/packages"
EXT_BUNDLE_DIR="bundle"
EXT_CONFIG_DIR="config"
EXT_OVERLAY_DIR="overlay"
EXT_UPSTREAM_DIR="upstream"
EXT_IMGPKG_DIR=".imgpkg"
EXT_DIR="${EXT_ROOT_DIR}/${EXT_NAME}"

# create directory structure for extension
mkdir -vp "${EXT_DIR}/${EXT_BUNDLE_DIR}/${EXT_CONFIG_DIR}"
mkdir -v "${EXT_DIR}/${EXT_BUNDLE_DIR}/${EXT_IMGPKG_DIR}"
mkdir -v "${EXT_DIR}/${EXT_BUNDLE_DIR}/${EXT_CONFIG_DIR}/${EXT_OVERLAY_DIR}"
mkdir -v "${EXT_DIR}/${EXT_BUNDLE_DIR}/${EXT_CONFIG_DIR}/${EXT_UPSTREAM_DIR}"

# create README and fill with name of extension
cp docs/extension-readme-template.md "${EXT_DIR}/README.md"
sed "$SEDARGS" "s/EXT_NAME/${EXT_NAME}/g" "${EXT_DIR}/README.md"

echo
echo "add-on boostrapped at ${EXT_DIR}"
echo
