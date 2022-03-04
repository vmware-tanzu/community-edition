#!/usr/bin/env bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace


BUILD_OS=$(uname 2>/dev/null || echo Unknown)
TCE_REPO_PATH="$(git rev-parse --show-toplevel)"
E2E_PATH="${TCE_REPO_PATH}/test/e2e/testdata/ensure-e2e-deps"

CMD="aws"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
    unzip awscliv2.zip
    ./aws/install -b ~/bin/aws
    rm -rf ./aws
    rm awscliv2.zip
    ;;
  Darwin)
    export E2E_PATH
    curl "https://awscli.amazonaws.com/AWSCLIV2.pkg" -o "awscliv2.pkg"

    # Creating a temp file to substitute environment variable
    cat "${E2E_PATH}"/aws-cli-choices.xml > "${E2E_PATH}"/choices-template.xml
    envsubst < "${E2E_PATH}"/choices-template.xml > "${E2E_PATH}"/aws-cli-choices.xml

    # installing AWS Cli
    installer -pkg AWSCLIV2.pkg -target CurrentUserHomeDirectory -applyChoiceChangesXML "${E2E_PATH}"/aws-cli-choices.xml
    ln -s "${E2E_PATH}"/aws-cli/aws /usr/local/bin/aws  
    ln -s "${E2E_PATH}"/aws-cli/aws_completer /usr/local/bin/aws_completer 
    
    # Clean up
    rm awscliv2.pkg
    cp "${E2E_PATH}"/choices-template.xml "${E2E_PATH}"/aws-cli-choices.xml
    rm "${E2E_PATH}"/choices-template.xml
    ;;
esac
fi
