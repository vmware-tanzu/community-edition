## Installation Procedure

You can install the Tanzu CLI using Homebrew **or** you can download the binary and install.

### Option 1: Use the Homebrew install method

1. Make sure you have the [Homebrew](https://brew.sh/) package manager installed.

1. Run the following in your terminal:

    ```sh
    brew install vmware-tanzu/tanzu/tanzu-community-edition
    ```

1. Run the post install configuration script. Note the output of the `brew install` step for the correct location of the configure script:

    ```sh
    <HOMEBREW-INSTALL-LOCATION>/v0.10.0/libexec/configure-tce.sh
    ```

## Option 2: Use the binary download/install  method

1. Download the release for [macOS](https://github.com/vmware-tanzu/community-edition/releases/download/{{< release_latest >}}/tce-darwin-amd64-{{< release_latest >}}.tar.gz).

1. Change to the download directory and unpack the release. Run the following in your terminal:

    ```sh
    cd <DOWNLOAD-DIR>
    tar xzvf ~/<DOWNLOAD-DIR>/tce-darwin-amd64-{{< release_latest >}}.tar.gz

    ```

1. Change to the extracted directory and run the install script.

    ```sh
    cd tce-darwin-amd64-{{< release_latest >}}
    ./install.sh
    ```
