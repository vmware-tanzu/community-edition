# Getting Started with Unmanaged Clusters

<!-- markdownlint-disable MD036 -->
<!-- markdownlint-disable MD024 -->

This guide walks you through standing up an unmanaged cluster using
Tanzu CLI.

## Before You Begin

Review [Plan Your Deployment](planning).

## Install Tanzu CLI

The `tanzu` CLI is used for deploying and managing Tanzu Community Edition.
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

## Deploy a Cluster

{{% include "/docs/assets/create-unmanaged-cluster.md" %}}

## Deploy a Package

{{% include "/docs/assets/install-package-to-cluster.md" %}}

## Delete a Cluster

{{% include "/docs/assets/delete-unmanaged-cluster.md" %}}

## Next Steps

* [Unmanaged Clusters Reference Documentation](ref-unmanaged-cluster)
* [Deploy a Test Workload to an Unmanaged Cluster](sample)
