#!/usr/bin/env bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This is a helper script that makes it easier to push
# a compatibility file to the VMware TCE registry, primarily for dev purposes

# Expects imgpkg to be installed locally

# Must provide 2 arguments
if [ $# -ne 2 ]; then
	echo "No arguments provided. Please give path to compatibility file and tag"
	exit 1
fi

localFile=$1
vmwareRegistry="projects.registry.vmware.com"
compatPath="/tce/compatibility"
tag=$2
fullImage="${vmwareRegistry}${compatPath}:${tag}"

docker login $vmwareRegistry

echo
echo "-----"
echo "Pushing file $localFile to $fullImage"
echo "-----"
echo

imgpkg push -i "$fullImage" -f "$localFile"

