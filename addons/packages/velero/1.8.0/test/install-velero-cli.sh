#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -e

# Script to download asset file from tag release using GitHub API v3.
# Inspired by: https://github.com/vmware-tanzu/community-edition/blob/main/hack/get-tce-release.sh

# Should have two arguments, release-tag and file name
if [ $# -ne 2 ]; then
	echo "Usage: $0 [tag] [name]"
	echo "Example: ./install-velero-cli.sh v1.8.0 linux"
	exit 1
fi

# Check dependencies exist
echo "Validating dependencies ..."
type curl grep sed tr jq

# Define variables
tag=$1
name=$2

GH_API="https://api.github.com"
GH_REPO="$GH_API/repos/vmware-tanzu/velero"
GH_TAGS="$GH_REPO/releases/tags/$tag"
AUTH=""
CURL_ARGS="-LJO#"

# Validate GitHub access token
if [ -z "${GITHUB_TOKEN:-}" ]; then
	echo "Warning: No GITHUB_TOKEN variable defined - requests to the GitHub API may be rate limited"
else
    AUTH="Authorization: token $GITHUB_TOKEN"
    curl -o /dev/null -sH "$AUTH" $GH_REPO || { echo "Error: Unauthenticated token or network issue!";  exit 1; }
fi

# Read asset tags
assets=$(curl -sH "$AUTH" "$GH_TAGS")

if [[ $(echo "$assets" | jq ".message") = "\"Not Found\"" ]]; then
	echo "Error: release tag $tag not found! Please enter a valid release tag version. Ex: v1.8.0"
	echo "Available release tags include:"
	curl -sH "$AUTH" "${GH_REPO}/releases" | jq "[ .[] | { html_url }]"
	exit 1
fi

# Get ID of the asset based on given name.
id=$(echo "$assets" | jq ".assets[] | select( .browser_download_url | contains(\"${name}-amd64\"))" | jq '.id')
if [ -z "$id" ]; then
	echo "Error: Failed to get asset id for release containing substring $name"
	echo "Please provide a substring for the name of an assetfile. Ex: darwin, linux"
	echo "Available asset files include:"
	echo "$assets" | jq "[ .assets[] | { browser_download_url }]"
	exit 1
fi

GH_ASSET="$GH_REPO/releases/assets/$id"

# Download asset file
echo "Downloading asset ..."
curl $CURL_ARGS -H "$AUTH" -H 'Accept: application/octet-stream' "$GH_ASSET"
tar xvzf velero-"$tag"-"$name"-amd64.tar.gz
