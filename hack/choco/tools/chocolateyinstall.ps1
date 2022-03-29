# Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

$ErrorActionPreference = 'Stop';
$releaseVersion = 'v0.11.0'
$packageName = 'tanzu-community-edition'
$packageFullName = "tce-windows-amd64-$releaseVersion"
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
# This line is for local testing
#$url64 = "C:\Users\...\tce-windows-amd64-${releaseVersion}.zip"
$url64 = "https://github.com/vmware-tanzu/community-edition/releases/download/${releaseVersion}/tce-windows-amd64-${releaseVersion}.zip"
$checksum64 = '2a859d12f65fee5d1f6cced16484a6398e3e3190064185608eb59b6529ce054b'
$checksumType64 = 'sha256'

$packageArgs = @{
    packageName    = $packageName
    unzipLocation  = $toolsDir
    url64bit       = $url64

    softwareName   = 'tanzu-community-edition'

    checksum64     = $checksum64
    checksumType64 = $checksumType64
}

function Test-Prereqs {
    # Since Windows users can install docker and kubectl a lot of different ways,
    # just see if the executables are there.
    Write-Host "`nChecking Prerequisites" -ForegroundColor Green

    if ( -not (Get-Command docker -ErrorAction SilentlyContinue)) {
        Write-Host -ForegroundColor Red "  - Docker CLI not present! This is required to create bootstrap clusters`n"
    }
    else {
        Write-Host "  - Docker CLI found, proceeding" -ForegroundColor Cyan
    }

    if ( -not (Get-Command kubectl -ErrorAction SilentlyContinue)) {
        Write-Host -ForegroundColor Red "  - kubectl not present! This is required to create bootstrap clusters`n"
    }
    else {
        Write-Host "  - kubectl found, proceeding`n" -ForegroundColor Cyan
    }
}

function Install-TanzuEnvironment {
    # important locations
    # XDG_DATA_HOME -> LOCALAPPDATA on Windows
    $PluginDir = "${HOME}\.config\tanzu-plugins"
    $CacheLocation = "${HOME}\.cache\tanzu\catalog.yaml"
    $CLIConfigLocation = "${HOME}\.config\tanzu\config.yaml"
    $CompatabilityLocation = "${HOME}\.config\tanzu\tkg\compatibility\tkg-compatibility.yaml"

    Write-Host "`nStarted tanzu CLI environment setup" -ForegroundColor Green

    ## begin env clean up ##

    # if an existing compatibility file exists, remove it; the cli will redownload it
    if (Test-Path -Path $CompatabilityLocation -PathType Leaf) {
        Remove-Item -Path ${CompatabilityLocation} -Force
        Write-Host "  - Removed stale compatibility file at ${CompatabilityLocation}" -ForegroundColor Cyan
    }

    # if an existing config file exists, remove it in favor of a newly initialized one
    if (Test-Path -Path $CLIConfigLocation -PathType Leaf) {
        Remove-Item -Path ${CLIConfigLocation} -Force
        Write-Host "  - Removed existing CLI config at ${CLIConfigLocation}" -ForegroundColor Cyan      
    }
    
    # if plugin cache exists, remove it; this ensures stale commands don't show up when running tanzu
    if (Test-Path -Path $CacheLocation -PathType Leaf) {
        Remove-Item -Path ${CacheLocation} -Force
        Write-Host "  - Removed existing tanzu plugin cache file at ${CacheLocation}" -ForegroundColor Cyan
    }

    # if PluginDir exists, we should remove it entirely as stale files could cause issues when we run tanzu init
    if (Test-Path -Path $PluginDir) {
        Remove-Item -Path ${PluginDir} -Force
        Write-Host "  - Removed existing tanzu plugin directory ${PluginDir}" -ForegroundColor Cyan
    }

    ## end env clean up ##

    ## begin env setup ##

    # create the plugin directory for tanzu CLI
    New-Item -Path ${PluginDir} -ItemType directory -Force | Out-Null
    Write-Host "  - Created CLI plugin directory at ${pluginDir}" -ForegroundColor Cyan

    # for every plugin (syntax == "tanzu-*"), move it to ${XDG_DATA_HOME}/tanzu-cli
    # this is where tanzu CLI will lookup the plugin to wire into its command
    Get-ChildItem -Path "${toolsDir}\${packageFullName}\default-local" -Recurse | Move-Item -Destination ${PluginDir} -Force
    Write-Host "  - Moved CLI plugins to ${pluginDir}" -ForegroundColor Cyan

    # initialize CLI and add TCE plugin repo (bucket)
    # Note that we use the toolsDir path because chocolatey doesn't put
    # binaries on the $PATH until _after_ the install script runs.
    $tanzuExe = "${toolsDir}\${packageFullName}\tanzu.exe"

    # The & allows execution of a binary stored in a variable.
    Write-Host "  - Initializing Tanzu configuration" -ForegroundColor Cyan
    # This is turned on because in framework v0.11.x, we report errors as logs for
    # installing plugins, when ErrorActionPreference is set to Stop, this fails
    # the install. If this is fixed in the future in framework, we should remove this
    # setting.
    $ErrorActionPreference = 'SilentlyContinue';
    & $tanzuExe init
    $ErrorActionPreference = 'Stop';

}
Test-Prereqs

# this is a built-in function, read https://docs.chocolatey.org/en-us/create/functions/install-chocolateyzippackage
Install-ChocolateyZipPackage @packageArgs

Install-TanzuEnvironment
