:: Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
:: SPDX-License-Identifier: Apache-2.0

@echo off

:: start delete tanzu cli
SET TANZU_CLI_DIR=%ProgramFiles%\tanzu
rmdir /Q /S "%TANZU_CLI_DIR%"

:: start delete plugins
SET PLUGIN_DIR=%LocalAppData%\tanzu-cli
rmdir /Q /S %PLUGIN_DIR%

:: start delete tce
SET TCE_DIR=%LocalAppData%\tce
rmdir /Q /S %TCE_DIR%

echo "Uninstall complete!"
