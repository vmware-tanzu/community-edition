# Getting Started with Docker Desktop

<!-- markdownlint-disable MD036 -->
<!-- markdownlint-disable MD024 -->

Docker Desktop offers a Tanzu Community Edition extension. This extension spins
up a local cluster that you can interact with via `kubectl` and the `tanzu` CLI.
Under the hood, the Docker Desktop extension utilizes [unmanaged
clusters](../ref-unmanaged-cluster), which can also be created exclusively using
the Tanzu CLI. This guide walks you through standing up a local cluster using
Docker Desktop.

## Before You Begin

Ensure you have Docker Desktop 4.8.0 or above installed.

You can download the Docker Desktop installer from the following locations:

* [Mac](https://docs.docker.com/desktop/mac/release-notes/)
* [Windows](https://docs.docker.com/desktop/windows/release-notes/)
* [Linux](https://docs.docker.com/desktop/linux/release-notes/)
  * The extension is not currently being validated against the Linux version
    of Docker Desktop. At this time we cannot guarentee support.

## Install Tanzu CLI

The Tanzu CLI is used for deploying and managing Tanzu Community Edition.
Choose your operating system below for guidance on installation.

{{< tabs tabTotal="3" tabID="1" tabName1="Linux" tabName2="Mac" tabName3="Windows">}}
{{< tab tabNum="1" >}}

### Linux System Requirements

{{% include "/docs/assets/prereq-unmanaged-linux.md" %}}

### Package Manager

**Homebrew**

{{% include "/docs/assets/install-homebrew.md" %}}

### Direct Download

{{% include "/docs/assets/direct-download.md" %}}

{{< /tab >}}
{{< tab tabNum="2" >}}

### Mac System Requirements

{{% include "/docs/assets/prereq-unmanaged-mac.md" %}}

### Package Manager

**Homebrew**

{{% include "/docs/assets/install-homebrew.md" %}}

### Direct Download

{{% include "/docs/assets/direct-download.md" %}}

{{< /tab >}}
{{< tab tabNum="3" >}}

### Windows System Requirements

{{% include "/docs/assets/prereq-unmanaged-windows.md" %}}

### Package Manager

**Chocolatey**

1. Install using [chocolatey](https://chocolatey.org/install), in **Powershell, as an administrator**.

    ```sh
    choco install tanzu-community-edition
    ```

### Direct Download

{{% include "/docs/assets/direct-download.md" %}}

{{< /tab >}}
{{< /tabs >}}

## Install the Extension

1. Launch Docker Desktop.

1. Click Add Extensions in the left navigation.

    ![Docker Desktop Add Extensions](../img/dd-add-extensions.png)

1. Locate the VMware Tanzu Community Edition extension.

    ![Docker Desktop TCE Extension](../img/dd-install-tce.png)

1. Click Install.

## Create a Cluster

1. In Docker Desktop, click Tanzu Community Edition from the left navigation.

1. Click Create cluster.

    ![Docker Desktop Create Cluster](../img/dd-create-cluster.png)

1. Wait for the cluster creation to complete.

    ![Docker Desktop Cluster Created](../img/dd-cluster-created.png)

Upon successful creation, `tanzu-community-edition` is added to your kubeconfig
and set as the default context. All commands run with `kubectl` and `tanzu` will
be against this cluster unless you switch contexts.

Ports `80` and `443` are forwarded by default. Consider this when installing
ingress controllers and other services to ensure you can route traffic into the
virtual machine.

> Today this is not configurable, but should be in the future.

## Install a Package

{{% include "/docs/assets/install-package-to-cluster.md" %}}

## Delete a Cluster

To delete a cluster, click Delete in Docker Desktop.

![Docker Desktop Delete Cluster](../img/dd-delete-cluster.png)

This is the equivalent of running a `docker kill` and `docker rm` against the
container. The image will remain on your system for quicker cluster startup in
the future.

## Next Steps

* [Deploy a Test Workload to an Unmanaged Cluster](sample)
