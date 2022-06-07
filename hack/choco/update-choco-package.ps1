# Copyright 2021-2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

$ErrorActionPreference = 'Stop';

if ((Test-Path env:BUILD_VERSION) -eq $False) {
  throw "BUILD_VERSION environment variable is not set"
}
if ((Test-Path env:TCE_CI_BUILD) -eq $False) {
  throw "TCE_CI_BUILD environment variable is not set"
}
if ((Test-Path env:CHOCO_API_KEY) -eq $False) {
  throw "CHOCO_API_KEY environment variable is not set"
}

# copy the modified nuspec file from a previous stage and stage it to push
if ($env:BUILD_VERSION -notlike '*-*') {
  Copy-Item tanzu-community-edition-release.nuspec -Destination .\tanzu-community-edition.nuspec
} else {
  Copy-Item tanzu-community-edition-unstable.nuspec -Destination .\tanzu-community-edition.nuspec
}

Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# push choco
choco apikey -k $env:CHOCO_API_KEY -s https://push.chocolatey.org/
choco pack .\tanzu-community-edition.nuspec
choco push --source https://push.chocolatey.org/ --api-key $env:CHOCO_API_KEY
