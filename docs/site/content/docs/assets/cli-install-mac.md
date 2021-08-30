1. Download the release for [macOS](https://github.com/vmware-tanzu/community-edition/releases/download/v0.7.0/tce-darwin-amd64-v0.7.0.tar.gz).

1. Unpack the release.

    ```sh
    tar xzvf ~/<DOWNLOAD-DIR>/tce-darwin-amd64-v0.7.0.tar.gz
    ```

1. Run the install script.

    ```sh
    cd tce-darwin-amd64-v0.7.0
    ./install.sh
    ```

    > This installs the `Tanzu` CLI and puts all the plugins in the correct location.
    > The first time you run the `tanzu` command the installed plugins and plugin repositories are initialized. This action might take a minute.

1. You must download and install the latest version of `kubectl`. For more information, see [Install and Set Up kubectl on macOS](https://kubernetes.io/docs/tasks/tools/install-kubectl-macos/) in the Kubernetes documentation.



