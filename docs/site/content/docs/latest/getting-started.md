# Getting Started with Managed Clusters

<!-- markdownlint-disable MD036 -->
<!-- markdownlint-disable MD024 -->

This guide walks you through standing up a management and workload cluster using
Tanzu CLI.

## Before You Begin

Review [Plan Your Deployment](planning).

{{% include "/docs/assets/tce-feedback.md" %}}

## Install Tanzu CLI

{{< tabs tabTotal="3" tabID="1" tabName1="Linux" tabName2="Mac" tabName3="Windows">}}
{{< tab tabNum="1" >}}

{{% include "/docs/assets/prereq-linux-short.md" %}}

### Package Manager

{{% include "/docs/assets/install-homebrew.md" %}}

### Direct Download

{{% include "/docs/assets/direct-download.md" %}}

{{< /tab >}}
{{< tab tabNum="2" >}}

{{% include "/docs/assets/prereq-mac-short.md" %}}

### Package Manager

{{% include "/docs/assets/install-homebrew.md" %}}

### Direct Download

{{% include "/docs/assets/direct-download.md" %}}

{{< /tab >}}
{{< tab tabNum="3" >}}

{{% include "/docs/assets/prereq-windows-short.md" %}}

### Package Manager

Install using [Chocolatey](https://chocolatey.org/install), in **Powershell, as an administrator**.

```sh
choco install tanzu-community-edition
```

### Direct Download

{{% include "/docs/assets/direct-download.md" %}}

{{< /tab >}}
{{< /tabs >}}

## Deploy Clusters

{{< tabs tabTotal="4" tabID="2" tabName1="AWS" tabName2="Azure" tabName3="Docker" tabName4="vSphere" >}}
{{< tab tabNum="1" >}}

{{% include "/docs/assets/aws-clusters.md" %}}

{{< /tab >}}
{{< tab tabNum="2" >}}

{{% include "/docs/assets/azure-clusters.md" %}}

{{< /tab >}}
{{< tab tabNum="3" >}}

### ⚠️ If bootstrapping docker-based clusters on Windows, [see our Windows guide](../ref-windows-capd)

{{% include "/docs/assets/capd-clusters.md" %}}

{{< /tab >}}
{{< tab tabNum="4" >}}

{{% include "/docs/assets/vsphere-clusters.md" %}}

{{< /tab >}}
{{< /tabs >}}

## Deploy a Package

{{% include "/docs/assets/package-installation.md" %}}

## Install a Local Dashboard (octant)

{{% include "/docs/assets/octant-install.md" %}}

## Delete Clusters

{{% include "/docs/assets/clean-up.md" %}}

## Next Steps

{{% include "/docs/assets/next-steps.md" %}}
