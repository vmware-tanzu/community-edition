# Docker-based Clusters on Windows

In order to run Docker-based clusters on Windows, multiple additional steps are
required. At this time, we don't recommend deploying Tanzu Community Edition clusters onto Docker
for Windows unless you're willing to tinker with lower level details around
Windows Subsystem for Linux. If you wish to continue, the following steps will
take you through deploying Docker-based clusters on Windows.

## Compile the WSL Kernel

⚠️ : These steps will have you use a custom-built kernel that will be used for
**all** your WSL-based VMs.

The CNI used by Tanzu Community Edition (antrea) requires specific configuration in the kernel that
is not enabled in the default WSL kernel. In future versions of antrea, this
kernel configuration will not be required (tracked in
[antrea#2635](https://github.com/antrea-io/antrea/issues/2635)). This section
covers compiling a kernel that will work with Antrea.

  > Thanks to [the kind
  > project](https://kind.sigs.k8s.io/docs/user/using-wsl2/) for hosting this
  > instructions, which we were able to build atop.

1. Run and enter an Ubuntu container to build the kernel

    ```txt
    docker run --name wsl-kernel-builder --rm -it ubuntu@sha256:9d6a8699fb5c9c39cf08a0871bd6219f0400981c570894cd8cbea30d3424a31f bash
    ```

1. From inside the container, run the following

    ```sh
    WSL_COMMIT_REF=linux-msft-5.4.72 # change this line to the version you want to build
    apt update
    apt install -y git build-essential flex bison libssl-dev libelf-dev bc

    mkdir src
    cd src
    git init
    git remote add origin https://github.com/microsoft/WSL2-Linux-Kernel.git
    git config --local gc.auto 0
    git -c protocol.version=2 fetch --no-tags --prune --progress --no-recurse-submodules --depth=1 origin +${WSL_COMMIT_REF}:refs/remotes/origin/build/linux-msft-wsl-5.4.y
    git checkout --progress --force -B build/linux-msft-wsl-5.4.y refs/remotes/origin/build/linux-msft-wsl-5.4.y

    # adds support for clientIP-based session affinity
    sed -i 's/# CONFIG_NETFILTER_XT_MATCH_RECENT is not set/CONFIG_NETFILTER_XT_MATCH_RECENT=y/' Microsoft/config-wsl

    # required module for antrea
    sed -i 's/# CONFIG_NETFILTER_XT_TARGET_CT is not set/CONFIG_NETFILTER_XT_TARGET_CT=y/' Microsoft/config-wsl

    # build the kernel
    make -j2 KCONFIG_CONFIG=Microsoft/config-wsl
    ```

1. Once the above completes, in a **new** powershell session, run the following

    ```sh
    docker cp wsl-kernel-builder:/src/arch/x86/boot/bzImage .
    ```

1. Set the `.wslconfig` file to point at the new kernel (`bzImage`).

    ```txt
    [wsl2]
    kernel=C:\\Users\\<your_user>\\bzImage
    ```

    > As seen above, you should escape the `\` by writing `\\`.
    > The above path may differ for you depending on where you compiled/saved
    > the kernel.

1. Shutdown WSL.

    ```sh
    wsl --shutdown
    ```

1. Restart WSL VMs.

    > This can be done via Docker desktop or using `wsl`.
    > You may need to restart Docker desktop even after restarting wsl.

1. Verify the kernel version run by WSL is consistent with what you compiled
   above.

    ```sh
    wsl uname -a
    Linux DESKTOP-4T1VL4L 5.4.72-microsoft-standard-WSL2+ #1 SMP Sat Sep 11 16:50:20 UTC 2021 x86_64 Linux
    ```

1. In a WSL VM with appropriate tools (e.g. Ubuntu) verify the kernel
   configuration required by antrea is present.

    ```sh
    wsl zgrep CONFIG_NETFILTER_XT_TARGET_CT /proc/config.gz

    CONFIG_NETFILTER_XT_TARGET_CT=y
    ```

## Create a Managed or Standalone CAPD Cluster

{{< tabs tabTotal="2" tabID="1" tabName1="Managed" tabName2="Standalone" >}}
{{< tab tabNum="1" >}}

{{% include "/docs/assets/capd-clusters-windows.md" %}}

{{< /tab >}}
{{< tab tabNum="2" >}}

{{% include "/docs/assets/capd-standalone-clusters-windows.md" %}}

{{< /tab >}}
{{< /tabs >}}

{{% include "/docs/assets/package-installation.md" %}}
{{% include "/docs/assets/octant-install.md" %}}
{{% include "/docs/assets/clean-up-standalone.md" %}}
