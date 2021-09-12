#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

# Change directories to the parent directory of the one in which this
# script is located.
ROOT_REPO_DIR="$(dirname "${BASH_SOURCE[0]}")/.."
cd "${ROOT_REPO_DIR}" || exit 1

if [[ -z "${BUILD_VERSION}" ]]; then
    echo "BUILD_VERSION is not set"
    exit 1
fi
if [[ -z "${FRAMEWORK_BUILD_VERSION}" ]]; then
    echo "FRAMEWORK_BUILD_VERSION is not set"
    exit 1
fi
if [[ -z "${TCE_RELEASE_DIR}" ]]; then
    echo "TCE_RELEASE_DIR is not set"
    exit 1
fi

# Required directories
BUILD_ROOT_DIR="${ROOT_REPO_DIR}/build"
FRAMEWORK_BUILD_VERSION="${FRAMEWORK_BUILD_VERSION:-latest}"
TCE_BUILD_VERSION="${BUILD_VERSION:-latest}"

DEP_BUILD_DIR="${TCE_RELEASE_DIR}"
ROOT_FRAMEWORK_DIR="${DEP_BUILD_DIR}/tanzu-framework"

PACKAGE_LINUX_AMD64_DIR="${BUILD_ROOT_DIR}/tce-linux-amd64-${BUILD_VERSION}"
PACKAGE_DARWIN_AMD64_DIR="${BUILD_ROOT_DIR}/tce-darwin-amd64-${BUILD_VERSION}"
PACKAGE_DARWIN_ARM64_DIR="${BUILD_ROOT_DIR}/tce-darwin-arm64-${BUILD_VERSION}"
PACKAGE_WINDOWS_AMD64_DIR="${BUILD_ROOT_DIR}/tce-windows-amd64-${BUILD_VERSION}"

# delete and create everything
rm -rf "${BUILD_ROOT_DIR}"
rm -rf "${PACKAGE_LINUX_AMD64_DIR}"
rm -rf "${PACKAGE_DARWIN_AMD64_DIR}"
rm -rf "${PACKAGE_DARWIN_ARM64_DIR}"
rm -rf "${PACKAGE_WINDOWS_AMD64_DIR}"
mkdir -p "${BUILD_ROOT_DIR}"
mkdir -p "${PACKAGE_LINUX_AMD64_DIR}/bin"
mkdir -p "${PACKAGE_DARWIN_AMD64_DIR}/bin"
mkdir -p "${PACKAGE_DARWIN_ARM64_DIR}/bin"
mkdir -p "${PACKAGE_WINDOWS_AMD64_DIR}/bin"

# Common directories
ROOT_FRAMEWORK_ARTFACTS_DIR="${ROOT_FRAMEWORK_DIR}/artifacts"
ROOT_FRAMEWORK_ARTFACTS_ADMIN_DIR="${ROOT_FRAMEWORK_DIR}/artifacts-admin"
ROOT_TCE_ARTIFACTS_DIR="${ROOT_REPO_DIR}/artifacts"

# copy tanzu cli bits Linux AMD64
# Tanzu bits
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/linux/amd64/cli/core/${FRAMEWORK_BUILD_VERSION}/tanzu-core-linux_amd64" "${PACKAGE_LINUX_AMD64_DIR}/bin/tanzu"
# cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/linux/amd64/cli/alpha/${FRAMEWORK_BUILD_VERSION}/tanzu-alpha-linux_amd64" "${PACKAGE_LINUX_AMD64_DIR}/bin/tanzu-plugin-alpha"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/linux/amd64/cli/cluster/${FRAMEWORK_BUILD_VERSION}/tanzu-cluster-linux_amd64" "${PACKAGE_LINUX_AMD64_DIR}/bin/tanzu-plugin-cluster"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/linux/amd64/cli/kubernetes-release/${FRAMEWORK_BUILD_VERSION}/tanzu-kubernetes-release-linux_amd64" "${PACKAGE_LINUX_AMD64_DIR}/bin/tanzu-plugin-kubernetes-release"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/linux/amd64/cli/login/${FRAMEWORK_BUILD_VERSION}/tanzu-login-linux_amd64" "${PACKAGE_LINUX_AMD64_DIR}/bin/tanzu-plugin-login"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/linux/amd64/cli/package/${FRAMEWORK_BUILD_VERSION}/tanzu-package-linux_amd64" "${PACKAGE_LINUX_AMD64_DIR}/bin/tanzu-plugin-package"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/linux/amd64/cli/pinniped-auth/${FRAMEWORK_BUILD_VERSION}/tanzu-pinniped-auth-linux_amd64" "${PACKAGE_LINUX_AMD64_DIR}/bin/tanzu-plugin-pinniped-auth"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/linux/amd64/cli/management-cluster/${FRAMEWORK_BUILD_VERSION}/tanzu-management-cluster-linux_amd64" "${PACKAGE_LINUX_AMD64_DIR}/bin/tanzu-plugin-management-cluster"

cp -f "${ROOT_FRAMEWORK_ARTFACTS_ADMIN_DIR}/linux/amd64/cli/builder/${FRAMEWORK_BUILD_VERSION}/tanzu-builder-linux_amd64" "${PACKAGE_LINUX_AMD64_DIR}/bin/tanzu-plugin-builder"

# TCE bits (New folder structure using tanzu-framework main)
cp -f "${ROOT_TCE_ARTIFACTS_DIR}/standalone-cluster/${TCE_BUILD_VERSION}/tanzu-standalone-cluster-linux_amd64" "${PACKAGE_LINUX_AMD64_DIR}/bin/tanzu-plugin-standalone-cluster"
cp -f "${ROOT_TCE_ARTIFACTS_DIR}/conformance/${TCE_BUILD_VERSION}/tanzu-conformance-linux_amd64" "${PACKAGE_LINUX_AMD64_DIR}/bin/tanzu-plugin-conformance"
cp -f "${ROOT_TCE_ARTIFACTS_DIR}/diagnostics/${TCE_BUILD_VERSION}/tanzu-diagnostics-linux_amd64" "${PACKAGE_LINUX_AMD64_DIR}/bin/tanzu-plugin-diagnostics"


# copy tanzu cli bits Darwin AMD64
# Tanzu bits
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/darwin/amd64/cli/core/${FRAMEWORK_BUILD_VERSION}/tanzu-core-darwin_amd64" "${PACKAGE_DARWIN_AMD64_DIR}/bin/tanzu"
# cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/darwin/amd64/cli/alpha/${FRAMEWORK_BUILD_VERSION}/tanzu-alpha-darwin_amd64" "${PACKAGE_DARWIN_AMD64_DIR}/bin/tanzu-plugin-alpha"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/darwin/amd64/cli/cluster/${FRAMEWORK_BUILD_VERSION}/tanzu-cluster-darwin_amd64" "${PACKAGE_DARWIN_AMD64_DIR}/bin/tanzu-plugin-cluster"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/darwin/amd64/cli/kubernetes-release/${FRAMEWORK_BUILD_VERSION}/tanzu-kubernetes-release-darwin_amd64" "${PACKAGE_DARWIN_AMD64_DIR}/bin/tanzu-plugin-kubernetes-release"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/darwin/amd64/cli/login/${FRAMEWORK_BUILD_VERSION}/tanzu-login-darwin_amd64" "${PACKAGE_DARWIN_AMD64_DIR}/bin/tanzu-plugin-login"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/darwin/amd64/cli/package/${FRAMEWORK_BUILD_VERSION}/tanzu-package-darwin_amd64" "${PACKAGE_DARWIN_AMD64_DIR}/bin/tanzu-plugin-package"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/darwin/amd64/cli/pinniped-auth/${FRAMEWORK_BUILD_VERSION}/tanzu-pinniped-auth-darwin_amd64" "${PACKAGE_DARWIN_AMD64_DIR}/bin/tanzu-plugin-pinniped-auth"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/darwin/amd64/cli/management-cluster/${FRAMEWORK_BUILD_VERSION}/tanzu-management-cluster-darwin_amd64" "${PACKAGE_DARWIN_AMD64_DIR}/bin/tanzu-plugin-management-cluster"

cp -f "${ROOT_FRAMEWORK_ARTFACTS_ADMIN_DIR}/darwin/amd64/cli/builder/${FRAMEWORK_BUILD_VERSION}/tanzu-builder-darwin_amd64" "${PACKAGE_DARWIN_AMD64_DIR}/bin/tanzu-plugin-builder"

# TCE bits (New folder structure using tanzu-framwork main)
cp -f "${ROOT_TCE_ARTIFACTS_DIR}/standalone-cluster/${TCE_BUILD_VERSION}/tanzu-standalone-cluster-darwin_amd64" "${PACKAGE_DARWIN_AMD64_DIR}/bin/tanzu-plugin-standalone-cluster"
cp -f "${ROOT_TCE_ARTIFACTS_DIR}/conformance/${TCE_BUILD_VERSION}/tanzu-conformance-darwin_amd64" "${PACKAGE_DARWIN_AMD64_DIR}/bin/tanzu-plugin-conformance"
cp -f "${ROOT_TCE_ARTIFACTS_DIR}/diagnostics/${TCE_BUILD_VERSION}/tanzu-diagnostics-darwin_amd64" "${PACKAGE_DARWIN_AMD64_DIR}/bin/tanzu-plugin-diagnostics"

# copy tanzu cli bits Darwin ARM64
# Tanzu bits
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/darwin/arm64/cli/core/${FRAMEWORK_BUILD_VERSION}/tanzu-core-darwin_arm64" "${PACKAGE_DARWIN_ARM64_DIR}/bin/tanzu"
# cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/darwin/arm64/cli/alpha/${FRAMEWORK_BUILD_VERSION}/tanzu-alpha-darwin_arm64" "${PACKAGE_DARWIN_ARM64_DIR}/bin/tanzu-plugin-alpha"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/darwin/arm64/cli/cluster/${FRAMEWORK_BUILD_VERSION}/tanzu-cluster-darwin_arm64" "${PACKAGE_DARWIN_ARM64_DIR}/bin/tanzu-plugin-cluster"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/darwin/arm64/cli/kubernetes-release/${FRAMEWORK_BUILD_VERSION}/tanzu-kubernetes-release-darwin_arm64" "${PACKAGE_DARWIN_ARM64_DIR}/bin/tanzu-plugin-kubernetes-release"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/darwin/arm64/cli/login/${FRAMEWORK_BUILD_VERSION}/tanzu-login-darwin_arm64" "${PACKAGE_DARWIN_ARM64_DIR}/bin/tanzu-plugin-login"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/darwin/arm64/cli/package/${FRAMEWORK_BUILD_VERSION}/tanzu-package-darwin_arm64" "${PACKAGE_DARWIN_ARM64_DIR}/bin/tanzu-plugin-package"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/darwin/arm64/cli/pinniped-auth/${FRAMEWORK_BUILD_VERSION}/tanzu-pinniped-auth-darwin_arm64" "${PACKAGE_DARWIN_ARM64_DIR}/bin/tanzu-plugin-pinniped-auth"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/darwin/arm64/cli/management-cluster/${FRAMEWORK_BUILD_VERSION}/tanzu-management-cluster-darwin_arm64" "${PACKAGE_DARWIN_ARM64_DIR}/bin/tanzu-plugin-management-cluster"

cp -f "${ROOT_FRAMEWORK_ARTFACTS_ADMIN_DIR}/darwin/arm64/cli/builder/${FRAMEWORK_BUILD_VERSION}/tanzu-builder-darwin_arm64" "${PACKAGE_DARWIN_ARM64_DIR}/bin/tanzu-plugin-builder"

# TCE bits (New folder structure using tanzu-framwork main)
cp -f "${ROOT_TCE_ARTIFACTS_DIR}/standalone-cluster/${TCE_BUILD_VERSION}/tanzu-standalone-cluster-darwin_arm64" "${PACKAGE_DARWIN_ARM64_DIR}/bin/tanzu-plugin-standalone-cluster"
cp -f "${ROOT_TCE_ARTIFACTS_DIR}/conformance/${TCE_BUILD_VERSION}/tanzu-conformance-darwin_arm64" "${PACKAGE_DARWIN_ARM64_DIR}/bin/tanzu-plugin-conformance"
cp -f "${ROOT_TCE_ARTIFACTS_DIR}/diagnostics/${TCE_BUILD_VERSION}/tanzu-diagnostics-darwin_amd64" "${PACKAGE_DARWIN_AMD64_DIR}/bin/tanzu-plugin-diagnostics"


# copy tanzu cli bits Windows AMD64
# Tanzu bits
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/windows/amd64/cli/core/${FRAMEWORK_BUILD_VERSION}/tanzu-core-windows_amd64.exe" "${PACKAGE_WINDOWS_AMD64_DIR}/bin/tanzu.exe"
# cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/windows/amd64/cli/alpha/${FRAMEWORK_BUILD_VERSION}/tanzu-alpha-windows_amd64.exe" "${PACKAGE_WINDOWS_AMD64_DIR}/bin/tanzu-plugin-alpha.exe"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/windows/amd64/cli/cluster/${FRAMEWORK_BUILD_VERSION}/tanzu-cluster-windows_amd64.exe" "${PACKAGE_WINDOWS_AMD64_DIR}/bin/tanzu-plugin-cluster.exe"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/windows/amd64/cli/kubernetes-release/${FRAMEWORK_BUILD_VERSION}/tanzu-kubernetes-release-windows_amd64.exe" "${PACKAGE_WINDOWS_AMD64_DIR}/bin/tanzu-plugin-kubernetes-release.exe"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/windows/amd64/cli/login/${FRAMEWORK_BUILD_VERSION}/tanzu-login-windows_amd64.exe" "${PACKAGE_WINDOWS_AMD64_DIR}/bin/tanzu-plugin-login.exe"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/windows/amd64/cli/package/${FRAMEWORK_BUILD_VERSION}/tanzu-package-windows_amd64.exe" "${PACKAGE_WINDOWS_AMD64_DIR}/bin/tanzu-plugin-package.exe"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/windows/amd64/cli/pinniped-auth/${FRAMEWORK_BUILD_VERSION}/tanzu-pinniped-auth-windows_amd64.exe" "${PACKAGE_WINDOWS_AMD64_DIR}/bin/tanzu-plugin-pinniped-auth.exe"
cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/windows/amd64/cli/management-cluster/${FRAMEWORK_BUILD_VERSION}/tanzu-management-cluster-windows_amd64.exe" "${PACKAGE_WINDOWS_AMD64_DIR}/bin/tanzu-plugin-management-cluster.exe"

cp -f "${ROOT_FRAMEWORK_ARTFACTS_ADMIN_DIR}/windows/amd64/cli/builder/${FRAMEWORK_BUILD_VERSION}/tanzu-builder-windows_amd64.exe" "${PACKAGE_WINDOWS_AMD64_DIR}/bin/tanzu-plugin-builder.exe"

# TCE bits (New folder structure using tanzu-framwork main)
cp -f "${ROOT_TCE_ARTIFACTS_DIR}/standalone-cluster/${TCE_BUILD_VERSION}/tanzu-standalone-cluster-windows_amd64.exe" "${PACKAGE_WINDOWS_AMD64_DIR}/bin/tanzu-plugin-standalone-cluster.exe"
cp -f "${ROOT_TCE_ARTIFACTS_DIR}/conformance/${TCE_BUILD_VERSION}/tanzu-conformance-windows_amd64.exe" "${PACKAGE_WINDOWS_AMD64_DIR}/bin/tanzu-plugin-conformance.exe"
cp -f "${ROOT_TCE_ARTIFACTS_DIR}/diagnostics/${TCE_BUILD_VERSION}/tanzu-diagnostics-windows_amd64.exe" "${PACKAGE_WINDOWS_AMD64_DIR}/bin/tanzu-plugin-diagnostics.exe"

# change settings
chmod +x "${ROOT_REPO_DIR}/hack/install.sh"
cp -f "${ROOT_REPO_DIR}/hack/install.sh" "${PACKAGE_LINUX_AMD64_DIR}"
cp -f "${ROOT_REPO_DIR}/hack/uninstall.sh" "${PACKAGE_LINUX_AMD64_DIR}"
cp -f "${ROOT_REPO_DIR}/hack/install.sh" "${PACKAGE_DARWIN_AMD64_DIR}"
cp -f "${ROOT_REPO_DIR}/hack/uninstall.sh" "${PACKAGE_DARWIN_AMD64_DIR}"
cp -f "${ROOT_REPO_DIR}/hack/install.sh" "${PACKAGE_DARWIN_ARM64_DIR}"
cp -f "${ROOT_REPO_DIR}/hack/uninstall.sh" "${PACKAGE_DARWIN_ARM64_DIR}"
cp -f "${ROOT_REPO_DIR}/hack/install.bat" "${PACKAGE_WINDOWS_AMD64_DIR}"
cp -f "${ROOT_REPO_DIR}/hack/uninstall.bat" "${PACKAGE_WINDOWS_AMD64_DIR}"
chown -R "$(id -u -n)":"$(id -g -n)" "${PACKAGE_LINUX_AMD64_DIR}"
chown -R "$(id -u -n)":"$(id -g -n)" "${PACKAGE_DARWIN_AMD64_DIR}"
chown -R "$(id -u -n)":"$(id -g -n)" "${PACKAGE_DARWIN_ARM64_DIR}"
chown -R "$(id -u -n)":"$(id -g -n)" "${PACKAGE_WINDOWS_AMD64_DIR}"

# packaging
rm -f tce-linux-amd64-*.tar.gz
rm -f tce-darwin-amd64-*.tar.gz
rm -f tce-darwin-arm64-*.tar.gz
rm -f tce-windows-amd64-*.tar.gz
pushd "${BUILD_ROOT_DIR}" || exit 1
tar -czvf "tce-linux-amd64-${TCE_BUILD_VERSION}.tar.gz" "tce-linux-amd64-${BUILD_VERSION}"
tar -czvf "tce-darwin-amd64-${TCE_BUILD_VERSION}.tar.gz" "tce-darwin-amd64-${BUILD_VERSION}"
tar -czvf "tce-darwin-arm64-${TCE_BUILD_VERSION}.tar.gz" "tce-darwin-arm64-${BUILD_VERSION}"
zip -r "tce-windows-amd64-${TCE_BUILD_VERSION}.zip" "tce-windows-amd64-${BUILD_VERSION}"
popd || exit 1
