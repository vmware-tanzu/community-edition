:: Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
:: SPDX-License-Identifier: Apache-2.0

@echo off

set STARTTIME=%TIME%
SET /A errno=0

:: start delete tanzu cli
SET TANZU_CLI_DIR=%ProgramFiles%\tanzu
if exist "%TANZU_CLI_DIR%" (
  rmdir /Q /S "%TANZU_CLI_DIR%"
)
if exist "%TANZU_CLI_DIR%" (
  SET /A errno=1
)

:: start delete plugins
SET PLUGIN_DIR="%LocalAppData%\tanzu-cli"
if exist "%PLUGIN_DIR%" (
  rmdir /Q /S %PLUGIN_DIR%
)
if exist "%PLUGIN_DIR%" (
  SET /A errno=1
)

:: start delete tanzu configuration
SET TANZU_CONFIG_DIR=%USERPROFILE%\.config\tanzu
if exist "%TANZU_CONFIG_DIR%" (
  rmdir /Q /S %TANZU_CONFIG_DIR%
)
if exist "%TANZU_CONFIG_DIR%" (
  SET /A errno=1
)

:: start delete tanzu cache
SET TANZU_CACHE_DIR=%USERPROFILE%\.cache\tanzu
if exist "%TANZU_CACHE_DIR%" (
  rmdir /Q /S %TANZU_CACHE_DIR%
)
if exist "%TANZU_CACHE_DIR%" (
  SET /A errno=1
)

:: start delete tce
SET TCE_DIR="%LocalAppData%\tce"
if exist "%TCE_DIR%" (
  rmdir /Q /S %TCE_DIR%
)
if exist "%TCE_DIR%" (
  SET /A errno=1
)

:: Get elapsed time:
set ENDTIME=%TIME%

:: Change formatting for the start and end times, adjusting for midnight crossover
    for /F "tokens=1-4 delims=:.," %%a in ("%STARTTIME%") do (
       set /A "start=(((%%a*60)+1%%b %% 100)*60+1%%c %% 100)*100+1%%d %% 100"
    )

    for /F "tokens=1-4 delims=:.," %%a in ("%ENDTIME%") do (
       IF %ENDTIME% GTR %STARTTIME% set /A "end=(((%%a*60)+1%%b %% 100)*60+1%%c %% 100)*100+1%%d %% 100"
       IF %ENDTIME% LSS %STARTTIME% set /A "end=((((%%a+24)*60)+1%%b %% 100)*60+1%%c %% 100)*100+1%%d %% 100"
    )

:: Calculate the elapsed time by subtracting values
    set /A elapsed=end-start

:: Format the results for output
    set /A hh=elapsed/(60*60*100), rest=elapsed%%(60*60*100), mm=rest/(60*100), rest%%=60*100, ss=rest/100, cc=rest%%100
    if %hh% lss 10 set hh=0%hh%
    if %mm% lss 10 set mm=0%mm%
    if %ss% lss 10 set ss=0%ss%
    if %cc% lss 10 set cc=0%cc%

    set DURATION=%hh%:%mm%:%ss%,%cc%

echo Start    : %STARTTIME%
echo Finish   : %ENDTIME%
echo          ---------------
echo Duration : %DURATION%
IF %errno%==0 (echo "Uninstall complete!") ELSE (echo "Uninstall incomplete!")
IF %errno%==0 (echo "Remove %TANZU_CLI_DIR% from your system's PATH.")
EXIT /B %errno%
