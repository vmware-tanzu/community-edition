# Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

param (
    # TCE release version argument
    [Parameter(Mandatory=$True)]
    [string]$version,

    # Path to the signtool
    [Parameter(Mandatory=$True)]
    [string]$signToolPath
)

$ErrorActionPreference = 'Stop';

if ((Test-Path env:GITHUB_TOKEN) -eq $False) {
  throw "GITHUB_TOKEN environment variable is not set"
}

$tempFolderPath = Join-Path $Env:Temp $(New-Guid)
New-Item -Type Directory -Path $tempFolderPath

$TCE_REPO_URL = "https://github.com/vmware-tanzu/community-edition"

gh release download $version --repo $TCE_REPO_URL --pattern "tce-windows-amd64-$version.zip" --dir $tempFolderPath

Expand-Archive -LiteralPath "$tempFolderPath\tce-windows-amd64-$version.zip" -Destination $tempFolderPath

# Check if the binaries are all signed
Get-ChildItem -Path "$tempFolderPath\tce-windows-amd64-$version\bin\tanzu-*" -File -Recurse | Foreach-Object {
  & $signToolPath verify /pa $_.FullName
  if ($LastExitCode -ne 0) {
    throw "Error verifying: " + $_.FullName
  }
}

Push-Location "$tempFolderPath\tce-windows-amd64-$version"

& ".\install.bat"

Pop-Location

$Env:Path += ";C:\Program Files\tanzu"

tanzu version

if ($LastExitCode -ne 0) {
  throw "Error verifying tanzu CLI using version command: " + $_.FullName
}

tanzu cluster version

if ($LastExitCode -ne 0) {
  throw "Error verifying tanzu cluster plugin using version command: " + $_.FullName
}

tanzu conformance version

if ($LastExitCode -ne 0) {
  throw "Error verifying tanzu conformance plugin using version command: " + $_.FullName
}

tanzu diagnostics version

if ($LastExitCode -ne 0) {
  throw "Error verifying tanzu diagnostics plugin using version command: " + $_.FullName
}

tanzu kubernetes-release version

if ($LastExitCode -ne 0) {
  throw "Error verifying tanzu kubernetes-release plugin using version command: " + $_.FullName
}

tanzu management-cluster version

if ($LastExitCode -ne 0) {
  throw "Error verifying tanzu management-cluster plugin using version command: " + $_.FullName
}

tanzu package version

if ($LastExitCode -ne 0) {
  throw "Error verifying tanzu package plugin using version command: " + $_.FullName
}

tanzu pinniped-auth version

if ($LastExitCode -ne 0) {
  throw "Error verifying tanzu pinniped-auth plugin using version command: " + $_.FullName
}

tanzu builder version

if ($LastExitCode -ne 0) {
  throw "Error verifying tanzu builder plugin using version command: " + $_.FullName
}

tanzu login version

if ($LastExitCode -ne 0) {
  throw "Error verifying tanzu login plugin using version command: " + $_.FullName
}

tanzu unmanaged-cluster version

if ($LastExitCode -ne 0) {
  throw "Error verifying tanzu unmanaged-cluster plugin using version command: " + $_.FullName
}
