:: Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
:: SPDX-License-Identifier: Apache-2.0

@echo off

:: start copy tanzu cli
SET TANZU_CLI_DIR=%ProgramFiles%\tanzu
mkdir "%TANZU_CLI_DIR%"
copy /B /Y tanzu.exe "%TANZU_CLI_DIR%"

:: set cli path
set PATH=%PATH%;%TANZU_CLI_DIR%
:: setx /M path "%path%;%TANZU_CLI_DIR%"
:: end copy tanzu cli

:: start copy plugins
SET PLUGIN_DIR="%USERPROFILE%\.config\tanzu-plugins"
SET TCE_DIR=%USERPROFILE%\tce"
SET TANZU_CACHE_DIR="%USERPROFILE%\.cache\tanzu"
mkdir %PLUGIN_DIR%
mkdir %TCE_DIR%
:: delete the plugin cache if it exists, before installing new plugins
rmdir /Q /S %TANZU_CACHE_DIR% 2>nul

:: Workaround!!!
:: For TF 0.17.0 or higher
:: tanzu plugin install all --local windows-amd64-default
:: For 0.11.2
:: setup
xcopy /Y /E /H /C /I default-local "%USERPROFILE%\.config\tanzu-plugins"
:: install plugins
tanzu plugin install builder
tanzu plugin install codegen
tanzu plugin install cluster
tanzu plugin install kubernetes-release
tanzu plugin install login
tanzu plugin install management-cluster
tanzu plugin install package
tanzu plugin install pinniped-auth
tanzu plugin install secret
tanzu plugin install conformance
tanzu plugin install diagnostics
tanzu plugin install unmanaged-cluster

:: copy uninstall.bat
copy /B /Y uninstall.bat %TCE_DIR%

:: explicit init of tanzu cli and add tce repo
:: For TF 0.17.0 or higher
:: tanzu init
tanzu plugin repo add --name tce --gcp-bucket-name tce-tanzu-cli-plugins --gcp-root-path artifacts
tanzu plugin repo add --name core-admin --gcp-bucket-name tce-tanzu-cli-framework-admin --gcp-root-path artifacts-admin

echo "Installation complete!"
echo "Please add %TANZU_CLI_DIR% permanently into your system's PATH."
