# Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

$ErrorActionPreference = 'Stop';

New-Item -ItemType Directory -Force -Path $HOME\tce-pkg

$parentDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

Write-Host "${parentDir}" -ForegroundColor Cyan

$packAndPushScriptPath = "${parentDir}\..\packpush.ps1"

Unblock-File -LiteralPath $packAndPushScriptPath

$packAndPushScriptBlock  =  [scriptblock]::Create($packAndPushScriptPath)

Invoke-Command -ScriptBlock $packAndPushScriptBlock

choco install tanzu-community-edition --source $HOME\tce-pkg -y

tanzu

tanzu version

tanzu standalone-cluster version

tanzu conformance version

tanzu diagnostics version

choco uninstall tanzu-community-edition --source $HOME\tce-pkg -y
