#!/usr/bin/env bash
# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o nounset
set -o pipefail

CHECK_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPO_DIR="$(dirname "${CHECK_DIR}")/.."

docker run --rm -t cytopia/yamllint --version
CONTAINER_NAME="tce_yamllint_$RANDOM"
docker run --name ${CONTAINER_NAME} -t -v "${REPO_DIR}":/community-edition:ro cytopia/yamllint -s -c /community-edition/hack/check/.yamllintconfig.yaml /community-edition
EXIT_CODE=$(docker inspect ${CONTAINER_NAME} --format='{{.State.ExitCode}}')
docker rm -f ${CONTAINER_NAME} &> /dev/null

if [[ ${EXIT_CODE} == "0" ]]; then
  echo "yamllint passed!"
else
  echo "yamllint exit code ${EXIT_CODE}: YAML linting failed!"
  echo "Please fix the listed yamllint errors and verify using 'make yamllint'"
  exit "${EXIT_CODE}"
fi
