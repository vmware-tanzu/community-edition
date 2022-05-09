## Installation Procedure

You can install the Tanzu CLI using Chocolatey **or** you can download the binary and install.
### Option 1: Use the Chocolatey installation method

1. Make sure you have the [Chocolatey](https://chocolatey.org/install) package manager installed.

1. Open Windows PowerShell **as an administrator** and run the following:

    ```sh
    choco install tanzu-community-edition
    ```

1. The `tanzu` command will be added to your `$PATH` variable automatically by Chocolatey.

### Option 2: Use the Binary download/installation method

1. Download the release for [Windows](https://github.com/vmware-tanzu/community-edition/releases/download/v0.12.0/tce-windows-amd64-v0.12.0.zip).

1. Open Windows PowerShell **as an administrator**, change to the download directory and unpack the release, for example,

    ```sh
    cd <DOWNLOAD-DIR>
    Expand-Archive -Path 'tce-windows-amd64-v0.12.0.zip'
    ```

1. Change to the extracted directory and run `.\install.bat`.

    ```sh
    cd tce-windows-amd64-v0.12.0\tce-windows-amd64-v0.12.0
    .\install.bat
    ```

1. Complete the following steps to add `C:\Program Files\tanzu` to your PATH.
      1. Search for  **env**, and choose **Edit the system environment variables**.
      2. Select **Environment Variables**, under **System Variable**, select **Path**, and  create a new path for `C:\Program Files\tanzu`.
