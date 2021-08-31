#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

function aws-nuke-tear-down {
	echo "$1"
	envsubst < "${TCE_REPO_PATH}"/test/aws/nuke-config-template.yml > "${TCE_REPO_PATH}"/test/aws/nuke-config.yml
	aws-nuke -c "${TCE_REPO_PATH}"/test/aws/nuke-config.yml --access-key-id "$AWS_ACCESS_KEY_ID" --secret-access-key "$AWS_SECRET_ACCESS_KEY" --force --no-dry-run || { error "$2 CLUSTER DELETION FAILED!!!"; rm -rf "${TCE_REPO_PATH}"/test/aws/nuke-config.yml; exit 1; }
	rm -rf "${TCE_REPO_PATH}"/test/aws/nuke-config.yml
	echo "$2 DELETED using aws-nuke!"
}