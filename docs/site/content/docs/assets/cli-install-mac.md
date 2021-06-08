1. Download the release for [macOS](https://github.com/vmware-tanzu/tce/releases/download/v0.5.0/tce-darwin-amd64-v0.5.0.tar.gz).

1. Unpack the release.

    ```sh
    tar xzvf ~/Downloads/tce-darwin-amd64-v0.5.0.tar.gz
    ```

1. Run the install script.

    ```sh
    cd tce-darwin-amd64-v0.5.0
    ./install.sh
    ```

    > This installs the `Tanzu` CLI and puts all the plugins in their proper location.
    > The first time you run the `tanzu` command the installed plugins and plugin repositories are initialized. This action might take a minute.

1. If you wish to run commands against any of the Kubernetes clusters that are created, you will need to download and install `kubectl`.

    ```sh
    curl -LO https://dl.k8s.io/release/v1.20.1/bin/darwin/amd64/kubectl
    sudo install -o root -g wheel -m 0755 kubectl /usr/local/bin/kubectl
    ```
