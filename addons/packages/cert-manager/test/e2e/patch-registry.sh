#!/usr/bin/env bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# shellcheck disable=SC2086
set -eoux pipefail

CUSTOM_REG_ROOT='projects-stg.registry.vmware.com/tkg'
CUSTOM_REG='projects-stg.registry.vmware.com/tkg/cert-manager-e2e-test'
NAMESPACE_SELECTOR='metadata.namespace!=kube-system,metadata.namespace!=tkg-system,metadata.namespace!=tkr-system'

# Patches registry of all pods filtered by namespace selector to use custom registry.
function patch_registry {
  while IFS='|' read -r kind name namespace; do
    for element in containers initContainers
    do
      while IFS='|' read -r containerImg containerName; do

        # only patch containers that aren't in CUSTOM_REG_ROOT
        if ! [[ $containerImg = $CUSTOM_REG_ROOT* ]]
        then
          # new image is CUSTOM_REG + name and tag of existing image
          IFS='/' read -a containerImgTokens <<< "$containerImg"
          containerImg="${CUSTOM_REG}/${containerImgTokens[${#containerImgTokens[@]}-1]}"
          # also update imagePullPolicy because cert-manager builds some containers from source and makes it available to kind but doesn't work with external clusters
          kubectl patch $kind/$name -n $namespace -p '{"spec":{"template":{"spec":{"'$element'":[{"name":"'$containerName'","image":"'$containerImg'", "imagePullPolicy": "IfNotPresent"}]}}}}'
        fi

      done < <(kubectl get $kind/$name -n $namespace -o json | jq '.spec.template.spec.'$element' |.[] | "\(.image)|\(.name)"' | tr -d \")

    done
  done < <(kubectl get daemonsets,deployments -A --field-selector=$NAMESPACE_SELECTOR -o json | jq '.items |.[] | "\(.kind)|\(.metadata.name)|\(.metadata.namespace)"' | tr -d \")
}



patch_registry
