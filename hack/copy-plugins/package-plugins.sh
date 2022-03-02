#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

ValuesFile=""

TanzuFrameworkCliPluginSourceRepo=""
TanzuFrameworkCliPluginDestinationVersion=""
TceCliPluginReleaseManifestPath=""

Main()
{
  EnsureDependencies
  ParseConfig
  CopyTfPluginsIntoTceOciRepo
  CreatePluginReleaseManifest
  PushRelease
}

ParseConfig()
{
  echo ""
  echo "Source Path where the Tanzu Framework CLI Plugins will be copied from:"
  tanzuFrameworkCliPluginSourceRepoLine=$(grep tanzuFrameworkCliPluginSourceRepo "${ValuesFile}")
  tanzuFrameworkCliPluginSourceRepo=$(echo "${tanzuFrameworkCliPluginSourceRepoLine}" | rev | cut -d: -f2 | rev | tr -d '[:space:]')
  tanzuFrameworkCliPluginSourceRepoVersion=$(echo "${tanzuFrameworkCliPluginSourceRepoLine}" | rev | cut -d: -f1 | rev )
  TanzuFrameworkCliPluginSourceRepo="${tanzuFrameworkCliPluginSourceRepo}:${tanzuFrameworkCliPluginSourceRepoVersion}"
  echo "TanzuFrameworkCliPluginSourceRepo: ${TanzuFrameworkCliPluginSourceRepo}"
  echo "tanzuFrameworkCliPluginSourceRepoVersion: ${tanzuFrameworkCliPluginSourceRepoVersion}"

  echo ""
  echo "Path where the Tanzu Framework CLI Plugins will be copied to in the TCE OCI Registry:"
  tanzuFrameworkCliPluginDestinationInTCERepoLine=$(grep tanzuFrameworkCliPluginDestinationInTCERepo "${ValuesFile}")
  TanzuFrameworkCliPluginDestinationInTCERepo=$(echo "${tanzuFrameworkCliPluginDestinationInTCERepoLine}" | cut -d: -f2 | tr -d '[:space:]')
  TanzuFrameworkCliPluginDestinationVersion=$(echo "${tanzuFrameworkCliPluginSourceRepoVersion}" | cut -d- -f1)
  TanzuFrameworkCliPluginDestinationInTCERepo="${TanzuFrameworkCliPluginDestinationInTCERepo}${TanzuFrameworkCliPluginDestinationVersion}"
  echo "TanzuFrameworkCliPluginDestinationInTCERepo: ${TanzuFrameworkCliPluginDestinationInTCERepo}"
  echo "TanzuFrameworkCliPluginDestinationVersion: ${TanzuFrameworkCliPluginDestinationVersion}"

  echo ""
  echo "The, ummmm, other path where the Tanzu Framework CLI Plugins will be copied to in the TCE OCI Registry:"
  tanzuFrameworkCliPluginDestinationInTCERepoLine=$(grep tanzuFrameworkCliPluginDestinationInTCERepo "${ValuesFile}")
  tanzuFrameworkCliPluginDestinationInTCERepo=$(echo "${tanzuFrameworkCliPluginDestinationInTCERepoLine}" | rev | cut -d: -f1 | rev | tr -d '[:space:]')
  echo "tanzuFrameworkCliPluginDestinationInTCERepo: ${tanzuFrameworkCliPluginDestinationInTCERepo}"

  echo ""
  echo "Path in the TCE OCI Registry where the final TCE CLI Plugin Release Manifest will be pushed to:"
  tceCliPluginReleaseRepo=$(grep tceCliReleaseRepo "${ValuesFile}" | cut -d: -f2 | tr -d '[:space:]')
  tceReleaseVersion=$(grep tceReleaseVersion "${ValuesFile}" | cut -d: -f2 | tr -d '[:space:]')
  TceCliPluginReleaseManifestPath="${tceCliPluginReleaseRepo}:${tceReleaseVersion}"
  echo "TceCliPluginReleaseManifestPath: ${TceCliPluginReleaseManifestPath}"
  echo "tceReleaseVersion: ${tceReleaseVersion}"
}

CopyTfPluginsIntoTceOciRepo()
{
  echo ""
  echo "CopyTfPluginsIntoTceOciRepo"
  echo "imgpkg copy --bundle ${TanzuFrameworkCliPluginSourceRepo} --to-repo ${TanzuFrameworkCliPluginDestinationInTCERepo}"
  # imgpkg copy --bundle ${TanzuFrameworkCliPluginSourceRepo} --to-repo ${TanzuFrameworkCliPluginDestinationInTCERepo}
}

CreatePluginReleaseManifest()
{
  echo ""
  echo "Create Plugin Release Manifest"
  rm -rf release
  mkdir -p release/plugins/.imgpkg

  echo "Generate lock-config.yaml..."
  ytt -f overlays/lock-config.yaml -f release-config.yaml > release/lock-config.yaml

  echo "Generate CLIPlugin Manifests"
  ytt -f overlays/cli-plugins.yaml -f release-config.yaml > release/plugins/cli-plugins.yaml

  echo "Lock images with kbld..."
  kbld -f release/lock-config.yaml --file release/plugins --imgpkg-lock-output release/plugins/.imgpkg/images.yml 1>> /dev/null
}

PushRelease()
{
  echo ""
  echo "PushRelease"
  echo "imgpkg push --bundle ${TanzuFrameworkCliPluginSourceRepo} --to-repo ${TanzuFrameworkCliPluginDestinationInTCERepo}"
  #imgpkg push --bundle ${TceCliPluginReleaseManifestPath} --file release/
}

EnsureDependencies()
{
  if [[ -z "$(command -v imgpkg)" ]]; then
    echo "The required application imgpkg is missing. Please install it from https://carvel.dev/imgpkg."
    exit 1
  fi
  if [[ -z "$(command -v kbld)" ]]; then
    echo "The required application kbld is missing. Please install it from https://carvel.dev/kbld."
    exit 1
  fi
  if [[ -z "$(command -v ytt)" ]]; then
    echo "The required application ytt is missing. Please install it from https://carvel.dev/ytt."
    exit 1
  fi
}

Help()
{
  echo "Copies a plugin package to the TCE OCI registry. It requires the imgpkg application."
  echo
  echo "Usage: package-plugins.sh --values-file values.yaml"
  echo "options:"
  echo "-h,--help         Help"
  echo "-f,--values-file Path to the configuration values file"
  echo
}

while (( "$#" )); do
  case "$1" in
    -h|--help)
      Help
      exit 1
      ;;
    -f|--values-file)
      if [ -n "$2" ] && [ "${2:0:1}" != "-" ]; then
        ValuesFile=$2
        shift 2
      else
        echo "Error: Argument for $1 is missing" >&2
        exit 1
      fi
      ;;
    --*=) # unsupported flags
      echo "Error: Unsupported flag $1" >&2
      exit 1
      ;;
    -*) # unsupported flags
      echo "Error: Unsupported flag $1" >&2
      exit 1
      ;;
  esac
done

Main
