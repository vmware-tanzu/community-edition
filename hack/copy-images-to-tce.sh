#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script is responsible for copying all images referenced in a BOM from the
# https://projects.registry.vmware.com/tkr repository into
# https://projects.registry.vmware.com/tce.

# This script was created as a workaround for the fact that TKG only supports changing the root
# registry for *all* images and not a specific image.

TKG_REPO=projects.registry.vmware.com/tkg
TCE_REPO=projects.registry.vmware.com/tce
images=( "$(grep -E 'imagePath: .*$' -A1 ~/.tanzu/tkg/bom/tkg-bom-v1.3.0.yaml | sed -n "s/^.*imagePath: \s*\(\S*\).*$/\1/p")" )
tags=( "$(grep -E 'imagePath: .*$' -A1 ~/.tanzu/tkg/bom/tkg-bom-v1.3.0.yaml | sed -n "s/^.*tag: \s*\(\S*\).*$/\1/p")" )
for i in "${!tags[@]}"; do 
  src=$(printf "%s/%s:%s\n" "${TKG_REPO}" "${images[$i]}" "${tags[$i]}")
  dst=$(printf "%s/%s:%s\n" "${TCE_REPO}" "${images[$i]}" "${tags[$i]}")
  printf "\n\n====== copying %s:%s ======\n\n" "${images[$i]}" "${tags[$i]}"
  crane cp "${src}" "${dst}"
done
