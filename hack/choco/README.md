# TCE Choco Package

## Local Development / Validation

1. Clone this repository

1. Move into this directory

1. Download the relevant Windows release you'd like to test from [https://github.com/vmware-tanzu/community-edition/releases](https://github.com/vmware-tanzu/community-edition/releases)

1. Unpack and repack the file as `.zip` (if the bundle wasn't already a `.zip`)

    > Moving to `.zip` is covered in [https://github.com/vmware-tanzu/community-edition/issues/1662](https://github.com/vmware-tanzu/community-edition/pull/1661/files#diff-9508030a41d1dcf13745795c27d65381fef43477b98a526545eeba25f5b9f457)

1. Store the `sha256` of the newly created zip, for example:

    ```sh
    josh@DESKTOP-4T1VL4L:/mnt/c/Users/joshr/Downloads$ sha256sum tce-windows-amd64-v0.8.0-rc.1.zip
    38ddc27423130469ba6e1b1fb6b8eed07fd08ccd92cda0dd9114e3211ec24d05  tce-windows-amd64-v0.8.0-rc.1.zip
    ```

1. Alter `tools/chocolateyinstall.ps1` to point at the local `.zip` bundle and update the SHA, for example:

    ```diff
    $ErrorActionPreference = 'Stop';
    $packageName = 'tanzu-community-edition'
    $packageFullName = 'tce-windows-amd64-v0.8.0-rc.1'
    $toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
    #This line is here for doing local testing
    +$url64 = 'C:\Users\<your_user>\Downloads\tce-windows-amd64-v0.8.0-rc.1.zip'
    -$url64 = 'https://github.com/vmware-tanzu/community-edition/releases/download/v0.8.0-rc.1/tce-windows-amd64-v0.8.0-rc.1.zip'
    +$checksum64 = '38ddc27423130469ba6e1b1fb6b8eed07fd08ccd92cda0dd9114e3211ec24d05'
    -$checksum64 = '2131332321321469ba6e1b1fb6b8eed07fd08ccd92cda0dd9114e3211ec24d05'
    $checksumType64= 'sha256'
    ```

1. Create an arbitrary directory on your computer to upload the package to, for example

    ```sh
    mkdir $HOME\tce-pkg
    ```

1. Create a package and upload it to this directory.

    ```sh
    choco pack; choco push --source $HOME\tce-pkg
    ```

1. Install the package by pointing to this directory

    ```sh
    choco install tanzu-community-edition --source $HOME\tce-choco
    ```

1. Verify the package installs correctly

    ```sh
    Attempt to use original download file name failed for 'C:\Users\joshr\Downloads\tce-windows-amd64-v0.8.0-rc.1.zip'.
    Copying tanzu-community-edition
    from 'C:\Users\joshr\Downloads\tce-windows-amd64-v0.8.0-rc.1.zip'
    Hashes match.
    Extracting C:\Users\joshr\AppData\Local\Temp\chocolatey\tanzu-community-edition\0.7.0\tanzu-community-editionInstall.zip to C:\ProgramData\chocolatey\lib\tanzu-community-edition\tools...
    C:\ProgramData\chocolatey\lib\tanzu-community-edition\tools

    Started tanzu CLI environment setup
    - Removed existing CLI config at C:\Users\joshr\.config\tanzu\config.yaml
    - Removed existing tanzu plugin cache file at C:\Users\joshr\.cache\tanzu\catalog.yaml
    - Created CLI plugin directory at C:\Users\joshr\AppData\Local\tanzu-cli
    - Moved CLI plugins to C:\Users\joshr\AppData\Local\tanzu-cli
    - Initializing tanzu CLI and plugin repository
    Completed tanzu CLI environment setup

    ShimGen has successfully created a shim for tanzu.exe
    The install of tanzu-community-edition was successful.
    Software installed to 'C:\ProgramData\chocolatey\lib\tanzu-community-edition\tools'

    Chocolatey installed 1/1 packages.
    See the log for details (C:\ProgramData\chocolatey\logs\chocolatey.log)
    ```

    ```sh
    PS C:\Users\joshr\community-edition\hack\choco> tanzu

    Tanzu CLI

    Usage:
      tanzu [command]

        Available command groups:

        Admin
            builder                 Build Tanzu components

        Run
            cluster                 Kubernetes cluster operations
            conformance             Run Sonobuoy conformance tests against clusters
            diagnostics             Cluster diagnostics
            kubernetes-release      Kubernetes release operations
            management-cluster      Kubernetes management cluster operations
            package                 Tanzu package managemen
            standalone-cluster      Create clusters without a dedicated management cluster

        System
            completion              Output shell completion code
            config                  Configuration for the CLI
            init                    Initialize the CLI
            login                   Login to the platform
            plugin                  Manage CLI plugins
            update                  Update the CLI
            version                 Version information


        Flags:
        -h, --help   help for tanzu

        Use "tanzu [command] --help" for more information about a command.

        Not logged in
    ```

1. When done, uninstall the package.

    ```sh
    PS C:\Users\joshr\community-edition\hack\choco> choco uninstall tanzu-community-edition --source $HOME\tce-pkg

    Chocolatey v0.10.15
    Uninstalling the following packages:
    tanzu-community-edition

    tanzu-community-edition v0.7.0
    Skipping auto uninstaller - No registry snapshot.
    tanzu-community-edition has been successfully uninstalled.

    Chocolatey uninstalled 1/1 packages.
    See the log for details (C:\ProgramData\chocolatey\logs\chocolatey.log).
    ```

1. The `tanzu` command should no longer be accessible.
