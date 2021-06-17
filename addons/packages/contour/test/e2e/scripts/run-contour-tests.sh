#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

################################################################################
##                               GLOBAL VARIABLES                             ##
################################################################################
ROOT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../../../.." >/dev/null 2>&1 && pwd )"
INFRA_PROVIDER="${1:-none}"
CONTOUR_VERSION="${CONTOUR_VERSION:-1.15.1}"
ADDRESS="${ADDRESS:-127.0.0.1}"
HTTP_PORT="${HTTP_PORT:-9080}"
HTTPS_PORT="${HTTPS_PORT:-9443}"
CHECK_TIMEOUT="${TIMEOUT:-360s}"
BASIC_E2E="${BASIC_E2E:-false}"
BASIC_E2E_TESTS=("httpproxy/003-path-condition-match.yaml" \
	         "httpproxy/004-https-sni-enforcement.yaml" \
		 "httpproxy/008-tcproute-https-termination.yaml")

################################################################################
##                           GET CONTOUR REPO                                 ##
################################################################################
rm -rf "${ROOT_DIR}"/.build/contour
mkdir -p "${ROOT_DIR}"/.build/contour
git clone https://github.com/projectcontour/contour "${ROOT_DIR}"/.build/contour
cd "${ROOT_DIR}"/.build/contour
git checkout v"${CONTOUR_VERSION}"
cd "${ROOT_DIR}"/.build/contour/_integration/testsuite

# Fix the test cases which failed due to DockerHub rate limit issues.
if [ "${INFRA_PROVIDER}" == "vc" ] || [ "${INFRA_PROVIDER}" == "none" ] || [ "${INFRA_PROVIDER}" == "kind" ]; then
  echo "Use internal proxy cache for DockerHub"
  sed -i -e 's|image: docker.io/kennethreitz/httpbin|image: harbor-repo.vmware.com/dockerhub-proxy-cache/kennethreitz/httpbin|g' fixtures/httpbin.yaml
  sed -i -e 's|image: tsaarni/echoserver|image: harbor-repo.vmware.com/dockerhub-proxy-cache/tsaarni/echoserver|g' fixtures/ingress-conformance-echo.yaml
  sed -i -e 's|image: docker.io/projectcontour/contour-authserver|image: harbor-repo.vmware.com/dockerhub-proxy-cache/projectcontour/contour-authserver|g' httpproxy/014-auth-basic-testserver.yaml
fi

################################################################################
##                           RUN CONTOUR E2E                                  ##
################################################################################
export ADDRESS
export HTTP_PORT
export HTTPS_PORT

if [ "${BASIC_E2E}" == "false" ]; then
  rm httpproxy/001-required-field-validation.yaml # This case failed in k8s v1.20.1
  rm httpproxy/018-external-name-service.yaml # This case was flaky as of Contour 1.12, has since been fixed (https://github.com/projectcontour/contour/pull/3342)
  rm httpproxy/020-global-rate-limiting.yaml
  ./run-test-case.sh httpproxy/*.yaml --check-timeout="${CHECK_TIMEOUT}"
else
  ./run-test-case.sh "${BASIC_E2E_TESTS[@]}" --check-timeout="${CHECK_TIMEOUT}"
fi