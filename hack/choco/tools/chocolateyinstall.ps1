# Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

$ErrorActionPreference = 'Stop';
$packageName = 'tanzu-community-edition'
$packageFullName = 'tce-windows-amd64-v0.8.0-rc.1'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
# This line is here for doing local testing
#$url64 = 'C:\Users\<user>\Downloads\tce-windows-amd64-v0.8.0-rc.1.zip'
$url64 = 'https://github.com/vmware-tanzu/community-edition/releases/download/v0.8.0-rc.1/tce-windows-amd64-v0.8.0-rc.1.tar.gz'
$checksum64 = '38ddc27423130469ba6e1b1fb6b8eed07fd08ccd92cda0dd9114e3211ec24d05'
$checksumType64= 'sha256'

$packageArgs = @{
  packageName   = $packageName
  unzipLocation = $toolsDir
  url64bit      = $url64

  softwareName  = 'tanzu-community-edition'

  checksum64    = $checksum64
  checksumType64= 'sha256'
}

function Setup-TanzuEnvironment {
    # important locations
    $PluginDir = "${env:LOCALAPPDATA}\tanzu-cli"
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

    ## end env clean up ##

    ## begin env setup ##

    # create the plugin directory for tanzu CLI
    New-Item -Path ${PluginDir} -ItemType directory -Force | Out-Null
    Write-Host "  - Created CLI plugin directory at ${pluginDir}" -ForegroundColor Cyan

    # for every plugin (syntax == "tanzu-*"), move it to ${XDG_DATA_HOME}/tanzu-cli
    # this is where tanzu CLI will lookup the plugin to wire into its command
    Get-ChildItem -Path "${toolsDir}\${packageFullName}\bin\tanzu-*" -Recurse | Move-Item -Destination ${PluginDir} -Force
    Write-Host "  - Moved CLI plugins to ${pluginDir}" -ForegroundColor Cyan


    # initialize CLI and add TCE plugin repo (bucket)
    Write-Host "  - Initializing tanzu CLI and plugin repository" -ForegroundColor Cyan
    & "${toolsDir}\${packageFullName}\bin\tanzu.exe" plugin repo add --name tce --gcp-bucket-name tce-tanzu-cli-plugins --gcp-root-path artifacts

    ## end env setup ##

    Write-Host "Completed tanzu CLI environment setup`n" -ForegroundColor Green
}

# this is a built-in function, read https://docs.chocolatey.org/en-us/create/functions/install-chocolateyzippackage
Install-ChocolateyZipPackage @packageArgs

Setup-TanzuEnvironment
