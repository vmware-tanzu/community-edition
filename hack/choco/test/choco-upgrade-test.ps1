# Copyright 2021-2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

param (
    # TCE release version argument
    [Parameter(Mandatory=$True)]
    [string]$version
)

$ErrorActionPreference = 'Stop';

if ((Test-Path env:GITHUB_TOKEN) -eq $False) {
  throw "GITHUB_TOKEN environment variable is not set"
}

$parentDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$TCE_REPO = "https://github.com/vmware-tanzu/community-edition" 
$TCE_REPO_RELEASES_URL = "https://github.com/vmware-tanzu/community-edition/releases"
$TCE_WINDOWS_ZIP_FILE="tce-windows-amd64-${version}.zip"
$TCE_CHECKSUMS_FILE = "tce-checksums.txt"

Write-Host "${parentDir}" -ForegroundColor Cyan

# Testing for current release
& "${parentDir}\e2e-test.ps1"

Write-Host "Checking if the necessary files exist for the TCE $version release"

invoke-webrequest "${TCE_REPO_RELEASES_URL}/download/${version}/${TCE_WINDOWS_ZIP_FILE}" -DisableKeepAlive -UseBasicParsing -Method head
invoke-webrequest "${TCE_REPO_RELEASES_URL}/download/${version}/${TCE_CHECKSUMS_FILE}" -OutFile "${parentDir}/tce-checksums.txt"

$Checksum64 = ((Select-String -Path "./test/tce-checksums.txt" -Pattern "tce-windows-amd64-${version}.zip").Line.Split(" "))[0]

# Updating the version in tanzu-community-edition-temp.nuspec file
$textnuspec = Get-Content .\tanzu-community-edition.nuspec -Raw
$temptextnuspec = Get-Content .\tanzu-community-edition.nuspec -Raw 
$Regex = [Regex]::new("(?<=<version>)(.*)(?=<\/version>)")
$oldVersion = $Regex.Match($textnuspec)
$textnuspec = $textnuspec.Replace( $oldVersion.value  , $version.Substring(1) )
Set-Content -Path .\tanzu-community-edition.nuspec -Value $textnuspec


# Updating the version in chocolateyinstall.ps1 file
$textchocoinstall = Get-Content .\tools\chocolateyinstall.ps1 -Raw 
$temptextchocoinstall = Get-Content .\tools\chocolateyinstall.ps1 -Raw 
$Regex = [Regex]::new("(?<=releaseVersion = ')(.*)(?=')")
$oldVersion = $Regex.Match($textchocoinstall)
$textchocoinstall = $textchocoinstall.Replace( $oldVersion.value  , $version )

# Updating the Checksum64 in chocolateyinstall.ps1 file
$Regex = [Regex]::new("(?<=checksum64 = ')(.*)(?=')")
$oldChecksum64 = $Regex.Match($textchocoinstall)
$textchocoinstall = $textchocoinstall.Replace( $oldChecksum64.value  , $Checksum64 )

Set-Content -Path .\tools\chocolateyinstall.ps1 -Value $textchocoinstall

# Testing for latest release
& "${parentDir}\e2e-test.ps1"

# Updating files for current version
Set-Content -Path .\tanzu-community-edition.nuspec -NoNewline -Value $temptextnuspec
Set-Content -Path .\tools\chocolateyinstall.ps1 -NoNewline -Value $temptextchocoinstall

Remove-Item "${parentDir}/tce-checksums.txt"
