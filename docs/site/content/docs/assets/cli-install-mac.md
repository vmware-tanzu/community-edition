## Installation Procedure

1. Make sure you have the [Homebrew package manager installed](https://brew.sh/)

1. You must download and install the latest version of `kubectl`. For more information, see [Install and Set Up kubectl on MacOS](https://kubernetes.io/docs/tasks/tools/install-kubectl-macos/) in the Kubernetes documentation.

1. You must download and install the latest version of `docker`. For more information, see [Install Docker Desktop on MacOS](https://docs.docker.com/desktop/mac/install/) in the Docker documentation.

1. Run the following in your terminal:

    ```sh
    brew tap vmware-tanzu/tanzu
    brew install tanzu-community-edition
    ```

1. Run the post install configuration script. Note the output of the `brew install` step for the correct location of the configure script:

    ```sh
    {HOMEBREW-INSTALL-LOCATION}/{{< release_latest >}}/libexec/configure-tce.sh
    ```

    > This puts all the Tanzu plugins in the correct location.
    > The first time you run the `tanzu` command the installed plugins and plugin repositories are initialized. This action might take a minute.
