1. Download the release for [macOS](https://github.com/vmware-tanzu/tce/releases/download/v0.6.0/tce-darwin-amd64-v0.6.0.tar.gz).

1. Unpack the release.

    ```sh
    tar xzvf ~/Downloads/tce-darwin-amd64-v0.6.0.tar.gz
    ```

1. Run the install script.

    ```sh
    cd tce-darwin-amd64-v0.6.0
    ./install.sh
    ```

    > This installs the `Tanzu` CLI and puts all the plugins in their proper location.
    > The first time you run the `tanzu` command the installed plugins and plugin repositories are initialized. This action might take a minute.


