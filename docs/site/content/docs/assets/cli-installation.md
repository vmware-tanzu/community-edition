## CLI Installation

1. Download the release.

    Log into GitHub and go to the [TCE Releases](https://github.com/vmware-tanzu/tce/releases/tag/v0.3.0) page and download the Tanzu CLI for either

    * [Linux](https://github.com/vmware-tanzu/tce/releases/download/v0.3.0/tce-linux-amd64-v0.3.0.tar.gz)
    * [macOS](https://github.com/vmware-tanzu/tce/releases/download/v0.3.0/tce-darwin-amd64-v0.3.0.tar.gz)

1. Unpack the release.

    **Linux**

    ```sh
    tar xzvf ~/Downloads/tce-linux-amd64-v0.3.0.tar.gz
    ```

    **macOS**

    ```sh
    tar xzvf ~/Downloads/tce-darwin-amd64-v0.3.0.tar.gz
    ```

1. Run the install script (make sure to use the appropriate directory for your platform).

    **linux**

    ```sh
    cd tce-linux-amd64-v0.3.0
    ./install.sh
    ```

    **macOS**

    ```sh
    cd tce-darwin-amd64-v0.3.0
    ./install.sh
    ```

    > This installs the `Tanzu` CLI and puts all the plugins in their proper location.
    > The first time you run the `tanzu` command the installed plugins and plugin repositories are initialized. This action might take a minute.

1. If you wish to run commands against any of the Kubernetes clusters that are created, you will need to download and install `kubectl`.

    **linux**

    ```sh
    curl -LO https://dl.k8s.io/release/v1.20.1/bin/linux/amd64/kubectl
    sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
    ```

    **macOS**

    ```sh
    curl -LO https://dl.k8s.io/release/v1.20.1/bin/darwin/amd64/kubectl
    sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
    ```
