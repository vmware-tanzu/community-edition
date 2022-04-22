#!/bin/bash

# Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -eo pipefail

existingPackageRepo='projects.registry.vmware.com-tce-main-v0.11.0'
packageRepoUrl=$1

registryServer=$(grep 'registry.server:' app-toolkit-values.yaml | awk '{print $NF}')
registryUser=$(grep 'registry.username:' app-toolkit-values.yaml | awk '{print $NF}')
registryPass=$(grep 'registry.password:' app-toolkit-values.yaml | awk '{print $NF}')
developerNamespace=$(grep 'developer_namespace:' app-toolkit-values.yaml | awk '{print $NF}')
workloadURL="http://tanzu-simple-web-app.${developerNamespace}.127-0-0-1.sslip.io/"

function main() {
	echo -e "=== APP TOOLKIT TEST - START ===\n"

	deleteExistingCluster
	createCluster
	checkExecutables
	updatePackageRepository
	setupSecrets
	installPackage
	createWorkload
	checkWorkload

	echo -e "\n=== APP TOOLKIT TEST - PASSED! ===\n"
}

function deleteExistingCluster {
	validateCommand "tanzu uc" "unmanaged-cluster"
  tanzu uc list | grep -q app-toolkit-test
	retcode=$?

	if [ $retcode -eq 0 ]; then
		echo "Existing 'app-toolkit-test' cluster found"
		tanzu uc delete app-toolkit-test
		echo "'app-toolkit-test' cluster deleted"
	fi
}

function createCluster {
	tanzu uc create app-toolkit-test -p 80:80 -p 443:443
}

function checkExecutables() {
	echo -e "\n--- Executables Check : Start ---\n"

	validateCommand "tanzu" "Tanzu CLI"
	validateCommand "tanzu apps" "Applications on Kubernetes"
	validateCommand "tanzu secret" "Tanzu secret management"
	validateCommand "tanzu package" "Tanzu package management"
	validateCommand "kubectl" "kubectl controls the Kubernetes cluster manager"
	validateCommand "docker" "A self-sufficient runtime for containers"

	echo -e "\n--- Executables Check : OK! ---\n"
}

function updatePackageRepository() {
	if [ "$packageRepoUrl" != "" ]; then
		echo "Updating '$existingPackageRepo' to use '$packageRepoUrl'"
		tanzu package repository update "$existingPackageRepo" -n tanzu-package-repo-global --url "$packageRepoUrl"
	else
		echo "Using standard PackageRepo found in $existingPackageRepo"
	fi
}

function setupSecrets() {
	echo -e "\n--- Setting Up Secrets : Start ---\n"

	tanzu package install secretgen-controller --package-name secretgen-controller.community.tanzu.vmware.com --version 0.8.0
	tanzu secret registry add registry-credentials --server "$registryServer" --username "$registryUser" --password "$registryPass" --export-to-all-namespaces --yes

	validateCommand "tanzu secret registry list" "registry-credentials"

	echo -e "\n--- Setting Up Secrets : OK! ---\n"
}

function installPackage() {
	echo -e "\n--- Installing App Toolkit : Start ---\n"

	tanzu package install app-toolkit -p app-toolkit.community.tanzu.vmware.com -v 0.2.0 -n tanzu-package-repo-global -f app-toolkit-values.yaml --verbose 3
	validateCommand "tanzu package installed get app-toolkit -n tanzu-package-repo-global" "ReconcileSucceeded"

	echo -e "\n--- Installing App Toolkit : OK! ---\n"
}

function createWorkload(){
	echo -e "\n--- Creating the Workload : Start ---\n"

	tanzu apps workload create tanzu-simple-web-app --git-repo https://github.com/cgsamp/tanzu-simple-web-app --git-branch main --type=web --app tanzu-simple-web-app --yes -n "$developerNamespace"

	echo -e "\n--- Creating the Workload : OK! ---\n"
}

function checkWorkload(){
	echo -e "\n--- Checking the Workload : Start ---\n"
	
	pollCommand "tanzu apps workload list -n ${developerNamespace}" "Ready" 5
	pollCommand "curl $workloadURL" "Hello" 1

	echo -e "\n--- Checking the Workload : OK! ---\n"
}

function validateCommand() {
	cmd=$1
	match=$2
	echo "Validating '$cmd'"
	output=$($cmd 2>&1)
	echo "$output" | grep -q "${match}"
	retcode=$?
	
	if [ $retcode -ne 0 ]; then
		fail "'$match' not found after executing '$cmd'"
	fi
}

function pollCommand() {
	cmd=$1
	match=$2
	timeout=$3
	duration=5
	count=0
	flag=1
	echo "Polling '$cmd' until it contains '$match'"

	while [ $flag -ne 0 ] ; do
		set +e
		output=$($cmd 2>&1)
		echo "${output}" | grep "${match}"
		flag=$?
		set -e
		minutes=$(( count / 60 ))
		if [[ "$minutes" -ge "$timeout" ]]; then
			fail "Timeout exceeded polling for '$cmd' to return expected result"
		fi
		sleep $duration
		count=$((count+duration))
	done
	minutes=$(( count / 60 ))
	seconds=$(( count % 60 ))
	echo "Result returned after ${minutes}m${seconds}s"
}

function fail() {
	echo -e "\n=== APP TOOLKIT TEST - FAILED! ===\n"
	echo "$1"
	exit 1
}

main
exit 0
