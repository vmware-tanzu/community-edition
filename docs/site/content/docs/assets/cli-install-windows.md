## Installation Procedure

1. Make sure you have the [Chocolatey package manager installed](https://chocolatey.org/install).

1. You must download and install the latest version of `kubectl`. For more information, see [Install and Set Up kubectl on Windows](https://kubernetes.io/docs/tasks/tools/install-kubectl-windows/) in the Kubernetes documentation.

1. You must download and install the latest version of `docker`. For more information, see [Install Docker Desktop on Windows](https://docs.docker.com/desktop/windows/install/) in the Docker documentation.

1. Open PowerShell **as an administrator** and run the following:

    ```sh
    choco install tanzu-community-edition --version={{< choco_release_latest  >}}

    ```

    > Both `docker` and `kubectl` are required to be present on the system, but are not explicit Chocolatey dependencies.
    > Installing the Tanzu Community Edition package will extract the binaries and configure the plugin repositories. This might take a minute.
    > Using an explicit version is required for now.

1. The `tanzu` command will be added to your `$PATH` variable automatically by Chocolatey.

