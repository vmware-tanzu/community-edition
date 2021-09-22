## Installation Procedure

1. Download the release for [Windows](https://github.com/vmware-tanzu/community-edition/releases/download/v0.8.0/tce-windows-amd64-v0.8.0.zip).

1. Open PowerShell **as an administrator**, change to the download directory and unpack the release, for example,

    ```sh
    cd <DOWNLOAD-DIR>
    Expand-Archive -Path 'tce-windows-amd64-v0.8.0.zip'
    ```

1. Change to the extracted directory and run `install.bat`.

   ```sh
   cd <INSTALLATION-DIR>
   install.bat
   ```

   > This installs the `Tanzu` CLI and puts all the plugins in the correct location.
   > The first time you run the `tanzu` command the installed plugins and plugin repositories are initialized. This action might take a minute.
1. Add `Program Files\tanzu` to your PATH.
1. You must download and install the latest version of `kubectl`. For more information, see [Install and Set Up kubectl on Windows](https://kubernetes.io/docs/tasks/tools/install-kubectl-windows/) in the Kubernetes documentation.
