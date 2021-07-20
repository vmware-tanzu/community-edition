:: Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
:: SPDX-License-Identifier: Apache-2.0

@echo off

:: start copy tanzu cli
SET TANZU_CLI_DIR=%ProgramFiles%\tanzu
mkdir "%TANZU_CLI_DIR%"
copy /B /Y bin\tanzu.exe "%TANZU_CLI_DIR%"

:: set cli path
set PATH=%PATH%;%TANZU_CLI_DIR%
:: setx /M path "%path%;%TANZU_CLI_DIR%"
:: end copy tanzu cli

:: start copy plugins
SET PLUGIN_DIR=%LocalAppData%\tanzu-cli
mkdir %PLUGIN_DIR%

:: core
copy /B /Y bin\tanzu-plugin-cluster.exe %PLUGIN_DIR%
copy /B /Y bin\tanzu-plugin-kubernetes-release.exe %PLUGIN_DIR%
copy /B /Y bin\tanzu-plugin-login.exe %PLUGIN_DIR%
copy /B /Y bin\tanzu-plugin-management-cluster.exe %PLUGIN_DIR%
copy /B /Y bin\tanzu-plugin-pinniped-auth.exe %PLUGIN_DIR%

:: tce
copy /B /Y bin\tanzu-plugin-package.exe %PLUGIN_DIR%
copy /B /Y bin\tanzu-plugin-standalone-cluster.exe %PLUGIN_DIR%
:: start copy plugins

:: explicit init of tanzu cli and add tce repo
tanzu plugin repo add --name tce --gcp-bucket-name tce-cli-plugins --gcp-root-path artifacts

echo "Installation complete!"
echo "Please add %TANZU_CLI_DIR% permanently into your system's PATH."
