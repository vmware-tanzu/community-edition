# Install the Tanzu CLI and Other Tools

This topic explains how to install and initialize the Tanzu command line interface (CLI) on a bootstrap machine.
The bootstrap machine is the laptop, host, or server that you deploy management and workload clusters from,
and that keeps the Tanzu and Kubernetes configuration files for your deployments.
The bootstrap machine is typically local, but it can also be a physical machine or VM that you access remotely.

Once the Tanzu CLI is installed, the second and last step to deploying Tanzu Kubernetes Grid is using the Tanzu CLI to create or designate a management cluster on each cloud provider that you use.
The Tanzu CLI then communicates with the management cluster to create and manage workload clusters on the cloud provider.

## <a id="prereqs"></a> Prerequisites

VMware provides Tanzu CLI binaries for Linux, macOS, and Windows systems.

The bootstrap machine on which you run the Tanzu CLI must meet the following requirements:

- A browser, if you intend to use the Tanzu Kubernetes Grid installer interface. You can use the Tanzu CLI without a browser, but for first deployments, it is **strongly recommended** to use the installer interface.
- A Linux, Windows, or macOS operating system with a minimum system configuration of 6 GB of RAM and a 2-core CPU.
- A Docker client installed and running on your bootstrap machine:
   - **Linux**: [Docker](https://docs.docker.com/install/)
   - **Windows**: [Docker Desktop](https://www.docker.com/products/docker-desktop)
   - **macOS**: [Docker Desktop](https://www.docker.com/products/docker-desktop)
- For Windows and macOS Docker clients, you must allocate at least 6 GB of memory in Docker Desktop to accommodate the `kind` container. See [Settings for Docker Desktop](https://kind.sigs.k8s.io/docs/user/quick-start/#settings-for-docker-desktop) in the `kind` documentation.
- System time is synchronized with a Network Time Protocol (NTP) server.
- On VMware Cloud on AWS and Azure VMware Solution, the bootstrap machine must be a cloud VM, not a local physical machine.  See [Prepare a vSphere Management as a Service Infrastructure](mgmt-clusters/prepare-maas.md) for setup instructions.
- If you intend to run the Tanzu CLI on a Linux machine, add your non-root user to the `docker` group. Create the group if it does not already exist. This enables the Tanzu CLI to access the Docker socket, which is owned by the `root` user. For more information, see the [Docker documentation](https://docs.docker.com/engine/install/linux-postinstall/#manage-docker-as-a-non-root-user).

## <a id="download"></a> Download and Unpack the Tanzu CLI and `kubectl`

The `tanzu` CLI ships with a compatible version of the `kubectl` CLI. To download and unpack both:

1. Go to [https://my.vmware.com](https://my.vmware.com) and log in with your My VMware credentials.

1. Visit [the Tanzu Kubernetes Grid downloads page](https://my.vmware.com/en/web/vmware/downloads/info/slug/infrastructure_operations_management/vmware_tanzu_kubernetes_grid/1_x)

1. In the **VMware Tanzu Kubernetes Grid** row, click **Go to Downloads**.

1. In the **Select Version** dropdown, select **1.4.0**.

1. Under **Product Downloads**, scroll to the section labeled **VMware Tanzu CLI 1.4.0 CLI**.
   - For macOS, locate **VMware Tanzu CLI for Mac** and click **Download Now**.
   - For Linux, locate **VMware Tanzu CLI for Linux** and click **Download Now**.
   - For Windows, locate **VMware Tanzu CLI for Windows** and click **Download Now**.

1. Navigate to the **Kubectl 1.20.5 for VMware Tanzu Kubernetes Grid 1.4.0** section of the download page.
   - For macOS, locate **kubectl cluster cli v1.20.5 for Mac** and click **Download Now**.
   - For Linux, locate **kubectl cluster cli v1.20.5 for Linux** and click **Download Now**.
   - For Windows, locate **kubectl cluster cli v1.20.5 for Windows** and click **Download Now**.

1. (Optional) Verify that your downloaded files are unaltered from the original. VMware provides a SHA-1, a SHA-256, and an MD5 checksum for each download. To obtain these checksums, click **Read More** under the entry that you want to download. For more information, see [Using Cryptographic Hashes](https://www.vmware.com/download/cryptographichashes.html).

1. On your system, create a new directory named `tanzu`. If you previously unpacked artifacts for previous releases to this folder, delete the folder's existing contents.

1. In the `tanzu` folder, unpack the Tanzu Kubernetes Grid CLI bundle file for your operating system. To unpack
the bundle file, use the extraction tool of your choice. For example, the `tar -xvf` command.

   - For macOS, unpack `tanzu-cli-bundle-v1.4.0-darwin-amd64.tar`.
   - For Linux, unpack `tanzu-cli-bundle-v1.4.0-linux-amd64.tar`.
   - For Windows, unpack `tanzu-cli-bundle-v1.4.0-windows-amd64.tar`.

    After you unpack the bundle file, in your `tanzu` folder, you will see a `cli` folder with multiple subfolders and files.

    The files in the `cli` directory, such as `ytt`, `kapp`, and `kbld`, are required by the Tanzu Kubernetes Grid extensions and add-ons. You will need these files later when you install the extensions and register add-ons.

1. Unpack the `kubectl` binary for your operating system:

   - For macOS, unpack `kubectl-mac-v1.20.5-vmware.1.gz`.
   - For Linux, unpack `kubectl-linux-v1.20.5-vmware.1.gz`.
   - For Windows, unpack `kubectl-windows-v1.20.5-vmware.1.exe.gz`.

## <a id="install-cli"></a> Install the Tanzu CLI

After you have downloaded and unpacked the Tanzu CLI on your bootstrap machine, you must make it available to the system.

1. Navigate to the `tanzu/cli` folder that you unpacked in the previous section.

1. Make the CLI available to the system:

   - For macOS:

       1. Install the binary to `/usr/local/bin`:
         <pre>sudo install core/v1.4.0/tanzu-core-darwin_amd64 /usr/local/bin/tanzu</pre>
       1. Confirm that the binary is executable by running the `ls` command.

    - For Linux:

       1. Install the binary to `/usr/local/bin`:
         <pre>sudo install core/v1.4.0/tanzu-core-linux_amd64 /usr/local/bin/tanzu</pre>
       1. Confirm that the binary is executable by running the `ls` command.

   - For Windows:

       1. Create a new `Program Files\tanzu` folder.
       1. In the unpacked `cli` folder, locate and copy the `core\v1.4.0\tanzu-core-windows_amd64.exe` into the new `Program Files\tanzu` folder.
       1. Rename `tanzu-core-windows_amd64.exe` to `tanzu.exe`.
       1. Right-click the `tanzu` folder, select **Properties** > **Security**, and make sure that your user account has the **Full Control** permission.
       1. Use Windows Search to search for `env`.
       1. Select **Edit the system environment variables** and click the **Environment Variables** button.
       1. Select the `Path` row under **System variables**, and click **Edit**.
       1. Click **New** to add a new row and enter the path to the `tanzu` CLI.

1. At the command line in a new terminal, run `tanzu version` to check that the correct version of the CLI is properly installed.

    If you are running on macOS, you might encounter the following error:

   ```
   "tanzu" cannot be opened because the developer cannot be verified.
   ```

   If this happens, you need to create a security exception for the `tanzu` executable. Locate the `tanzu` app in Finder, control-click the app, and select **Open**.

## <a id="install-plugins"></a> Install the Tanzu CLI Plugins

After you have installed the `tanzu` core executable, you must install the CLI plugins related to Tanzu Kubernetes cluster management and feature operations.

1. (Optional) Remove existing plugins from any previous CLI installations.

    ```
    tanzu plugin clean
    ```

1. Navigate to the `tanzu` folder that contains the `cli` folder.

1. Run the following command from the `tanzu` directory to install all the plugins for this release.

   ```
   tanzu plugin install --local cli all
   ```

1. Check plugin installation status.

   ```
   tanzu plugin list
   ```

   If successful, you should see a list of all installed plugins. For example:

   ```
   NAME                LATEST VERSION  DESCRIPTION                                                        REPOSITORY  VERSION  STATUS
   cluster             v1.4.0          Kubernetes cluster operations                                      core        v1.4.0   installed
   login               v1.4.0          Login to the platform                                              core        v1.4.0   installed
   pinniped-auth       v1.4.0          Pinniped authentication operations (usually not directly invoked)  core        v1.4.0   installed
   kubernetes-release  v1.4.0          Kubernetes release operations                                      core        v1.4.0   installed
   management-cluster  v1.4.0          Kubernetes management cluster operations                           core        v1.4.0   installed
   ```

## <a id="common-options"></a> Tanzu CLI Help

Run `tanzu --help` to see the list of commands that the Tanzu CLI provides.

You can view help text for any command group with the `--help` option to see information about that specific command or command group. For example, `tanzu login --help`, `tanzu management-cluster --help`, or `tanzu management-cluster create --help`.

For more information about the Tanzu CLI, see the [Tanzu CLI Command Reference](tanzu-cli-reference.md).

## <a id="install-kubectl"></a> Install `kubectl`

After you have downloaded and unpacked `kubectl` on your bootstrap machine, you must make it available to the system.

1. Navigate to the `kubectl` binary that you unpacked in [Download and Unpack the Tanzu CLI and kubectl](#download) above.

1. Make the CLI available to the system:

   - For macOS:

      1. Install the binary to `/usr/local/bin`:

         ```
         sudo install kubectl-mac-v1.20.5-vmware.1 /usr/local/bin/kubectl
         ```

      1. Confirm that the binary is executable by running the `ls` command.

   - For Linux:

      1. Install the binary to `/usr/local/bin`:

         ```
         sudo install kubectl-linux-v1.20.5-vmware.1 /usr/local/bin/kubectl
         ```

      1. Confirm that the binary is executable by running the `ls` command.

   - For Windows:

      1. Create a new `Program Files\kubectl` folder.
      1. Locate and copy the `kubectl-windows-v1.20.5-vmware.1.exe` file into the new `Program Files\kubectl` folder.
      1. Rename `kubectl-windows-v1.20.5-vmware.1.exe` to `kubectl.exe`.
      1. Right-click the `kubectl` folder, select **Properties** > **Security**, and make sure that your user account has the **Full Control** permission.
      1. Use Windows Search to search for `env`.
      1. Select **Edit the system environment variables** and click the **Environment Variables** button.
      1. Select the `Path` row under **System variables** and click **Edit**. 
      1. Click **New** to add a new row and enter the path to the `kubectl` CLI.

1. Run `kubectl version` to check that the correct version of the CLI is properly installed.

## <a id="install-carvel"></a> Install the Carvel Tools

The Tanzu Kubernetes Grid uses the following tools from the [Carvel open-source project](https://carvel.dev/):

- [`ytt`](https://carvel.dev/ytt/)
- [`kapp`](https://carvel.dev/kapp/)
- [`kbld`](https://carvel.dev/kbld/)
- [`imgpkg`](https://carvel.dev/imgpkg/)

Tanzu Kubernetes Grid provides signed binaries for `ytt`, `kapp`, `kbld`, `imgpkg`, and `vendir` that are bundled with the Tanzu CLI. The bundle also includes [`vendir`](https://carvel.dev/vendir/), a Kubernetes directory structure tool, that is not currently required by end users, but is provided for convenience.

1. Navigate to the location on your bootstrap environment machine where you unpacked the Tanzu CLI bundle tar file for your OS.

   For example, the `tanzu` folder, that you created in the previous procedure.

1. Open the `cli` folder.

   ```
   cd cli
   ```

### Install `ytt`

[`ytt`](https://carvel.dev/ytt/) is a YAML templating tool, that is used to deploy the Tanzu Kubernetes Grid extensions. You might need to install `ytt` if you need to troubleshoot the deployment of the Tanzu Kubernetes Grid extensions, or if you need to use overlays to customize your extensions or cluster templates.

MacOS and Linux:

1. Unpack the `ytt` binary and make it executable.

  Linux:

  ```
  gunzip ytt-linux-amd64-v0.31.0+vmware.1.gz
  chmod ugo+x ytt-linux-amd64-v0.31.0+vmware.1
  ```

  Mac OS:

  ```
  gunzip ytt-darwin-amd64-v0.31.0+vmware.1.gz
  chmod ugo+x ytt-darwin-amd64-v0.31.0+vmware.1
  ```

2. Move the binary to `/usr/local/bin` and rename it to `ytt`:

  Linux:

  ```
  mv ./ytt-linux-amd64-v0.31.0+vmware.1 /usr/local/bin/ytt
  ```

  Mac OS:

  ```
  mv ./ytt-darwin-amd64-v0.31.0+vmware.1 /usr/local/bin/ytt
  ```
1. Confirm that the binary is executable by running the `ls` command.


Windows:

1. Unpack the  the `ytt` binary.

  ```
  gunzip ytt-windows-amd64-v0.31.0+vmware.1.gz
  ```

1. Rename `ytt-windows-amd64-v0.30.0+vmware.1` to `ytt.exe`
1. Create a new `Program Files\ytt` folder and copy the `ytt.exe` file into it.
1. Right-click the `ytt` folder, select **Properties** > **Security**, and make sure that your user account has the **Full Control** permission.
1. Use Windows Search to search for `env`.
1. Select **Edit the system environment variables** and click the **Environment Variables** button.
1. Select the `Path` row under **System variables**, and click **Edit**.
1. Click **New** to add a new row and enter the path to the `ytt` tool.

At the command line in a new terminal, run `ytt version` to check that the correct version of `ytt` is properly installed.

### Install `kapp`

[`kapp`](https://carvel.dev/kapp/) is a Kubernetes applications CLI, that is used to manage the Tanzu Kubernetes Grid extensions. You might need to install `kapp` if you need to troubleshoot the deployment of the Tanzu Kubernetes Grid extensions, or if you need to use overlays to customize your extensions or cluster templates.

MacOS and Linux:

1. Unpack the `kapp` binary and make it executable.

  Linux:

  ```
  gunzip kapp-linux-amd64-v0.36.0+vmware.1.gz
  chmod ugo+x kapp-linux-amd64-v0.36.0+vmware.1
  ```

  Mac OS:

  ```
  gunzip kapp-darwin-amd64-v0.36.0+vmware.1.gz
  chmod ugo+x kapp-darwin-amd64-v0.36.0+vmware.1
  ```

1. Move the binary to `/usr/local/bin` and rename it to `kapp`:

  Linux:

  ```
  mv ./kapp-linux-amd64-v0.36.0+vmware.1 /usr/local/bin/kapp
  ```

  Mac OS:

  ```
  mv ./kapp-darwin-amd64-v0.36.0+vmware.1 /usr/local/bin/kapp
  ```
1. Confirm that the binary is executable by running the `ls` command.

Windows:

1. Unpack the  the `kapp` binary.

  ```
  gunzip kapp-windows-amd64-v0.36.0+vmware.1.gz
  ```

1. Rename `kapp-windows-amd64-v0.36.0+vmware.1` to `kapp.exe`
1. Create a new `Program Files\kapp` folder and copy the `kapp.exe` file into it.
1. Right-click the `kapp` folder, select **Properties** > **Security**, and make sure that your user account has the **Full Control** permission.
1. Use Windows Search to search for `env`.
1. Select **Edit the system environment variables** and click the **Environment Variables** button.
1. Select the `Path` row under **System variables**, and click **Edit**.
1. Click **New** to add a new row and enter the path to the `kapp` tool.

At the command line in a new terminal, run `kapp version` to check that the correct version of `kapp` is properly installed.

### Install `kbld`

[`kbld`](https://carvel.dev/kbld/) is a Kubernetes image builder. It is used by Tanzu Kubernetes Grid to build the Tanzu Kubernetes Grid extensions. You might need to install `kbld` if you need to troubleshoot the deployment of the Tanzu Kubernetes Grid extensions, or if you need to use overlays to customize your extensions or cluster templates.

MacOS and Linux:

1. Unpack the `kbld` binary and make it executable.

  Linux:

  ```
  gunzip kbld-linux-amd64-v0.28.0+vmware.1.gz
  chmod ugo+x kbld-linux-amd64-v0.28.0+vmware.1
  ```

  Mac OS:

  ```
  gunzip kbld-darwin-amd64-v0.28.0+vmware.1.gz
  chmod ugo+x kbld-darwin-amd64-v0.28.0+vmware.1
  ```

1. Move the binary to `/usr/local/bin` and rename it to `kbld`:

  Linux:

  ```
  mv ./kbld-linux-amd64-v0.28.0+vmware.1 /usr/local/bin/kbld
  ```

  Mac OS:

  ```
  mv ./kbld-darwin-amd64-v0.28.0+vmware.1 /usr/local/bin/kbld
  ```
1. Confirm that the binary is executable by running the `ls` command.

Windows:

1. Unpack the  the `kbld` binary.

  ```
  gunzip kbld-windows-amd64-v0.28.0+vmware.1.gz
  ```

1. Rename `kbld-windows-amd64-v0.28.0+vmware.1` to `kbld.exe`
1. Create a new `Program Files\kbld` folder and copy the `kbld.exe` file into it.
1. Right-click the `kbld` folder, select **Properties** > **Security**, and make sure that your user account has the **Full Control** permission.
1. Use Windows Search to search for `env`.
1. Select **Edit the system environment variables** and click the **Environment Variables** button.
1. Select the `Path` row under **System variables**, and click **Edit**.
1. Click **New** to add a new row and enter the path to the `kbld` tool.

At the command line in a new terminal, run `kbld version` to check that the correct version of `kbld` is properly installed.

### Install `imgpkg`

[`imgpkg`](https://carvel.dev/imgpkg/): Kubernetes image packaging tool, that is required to deploy Tanzu Kubernetes Grid in Internet-restricted environments and when building your own node images.

MacOS and Linux:

1. Unpack the `imgpkg` binary and make it executable.

  Linux:

  ```
  gunzip imgpkg-linux-amd64-v0.5.0+vmware.1.gz
  chmod ugo+x imgpkg-linux-amd64-v0.5.0+vmware.1
  ```

  Mac OS:

  ```
  gunzip imgpkg-darwin-amd64-v0.5.0+vmware.1.gz
  chmod ugo+x imgpkg-darwin-amd64-v0.5.0+vmware.1
  ```

1. Move the binary to `/usr/local/bin` and rename it to `imgpkg`:

  Linux:

  ```
  mv ./imgpkg-linux-amd64-v0.5.0+vmware.1 /usr/local/bin/imgpkg
  ```

  Mac OS:

  ```
  mv ./imgpkg-darwin-amd64-v0.5.0+vmware.1 /usr/local/bin/imgpkg
  ```
1. Confirm that the binary is executable by running the `ls` command.

Windows:

1. Unpack the  the `imgpkg` binary.

  ```
  gunzip imgpkg-windows-amd64-v0.5.0+vmware.1.gz
  ```

1. Rename `imgpkg-windows-amd64-v0.5.0+vmware.1` to `imgpkg.exe`
1. Create a new `Program Files\imgpkg` folder and copy the `imgpkg.exe` file into it.
1. Right-click the `imgpkg` folder, select **Properties** > **Security**, and make sure that your user account has the **Full Control** permission.
1. Use Windows Search to search for `env`.
1. Select **Edit the system environment variables** and click the **Environment Variables** button.
1. Select the `Path` row under **System variables**, and click **Edit**.
1. Click **New** to add a new row and enter the path to the `imgpkg` tool.

At the command line in a new terminal, run `imgpkg version` to check that the correct version of `imgpkg` is properly installed.

## <a id="what-next"></a> What to Do Next

With the Tanzu CLI and other tools installed, you can set up and use your bootstrap machine to deploy Kubernetes clusters to vSphere, Amazon EC2, and Microsoft Azure.

- For information about how to deploy management clusters to your chosen platform, see [Deploying Management Clusters](mgmt-clusters/deploy-management-clusters.md).
- If you have vSphere 7 and the vSphere with Tanzu feature is enabled, you can directly use the Tanzu CLI to deploy Tanzu Kubernetes clusters to vSphere with Tanzu, without deploying a management cluster. For information about how to connect the Tanzu CLI to a vSphere with Tanzu Supervisor Cluster, see [Use the Tanzu CLI with a vSphere with Tanzu Supervisor Cluster](tanzu-k8s-clusters/connect-vsphere7.md).
