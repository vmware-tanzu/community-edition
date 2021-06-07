1. Download the release for [Linux](https://github.com/vmware-tanzu/tce/releases/download/v0.4.0/tce-linux-amd64-v0.4.0.tar.gz).

1. Unpack the release.

    ```sh
    tar xzvf ~/Downloads/tce-linux-amd64-v0.5.0.tar.gz
    ```

1. Run the install script (make sure to use the appropriate directory for your platform).

    ```sh
    cd tce-linux-amd64-v0.5.0
    ./install.sh
    ```

    > This installs the `Tanzu` CLI and puts all the plugins in their proper location.
    > The first time you run the `tanzu` command the installed plugins and plugin repositories are initialized. This action might take a minute.

1. If you wish to run commands against any of the Kubernetes clusters that are created, you will need to download and install `kubectl`.

    ```sh
    curl -LO https://dl.k8s.io/release/v1.20.1/bin/linux/amd64/kubectl
    sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
    ```
