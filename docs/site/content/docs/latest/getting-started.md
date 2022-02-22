# Getting Started with Managed Clusters

This guide walks you through standing up a management and workload cluster using
Tanzu Community Edition.

{{% include "/docs/latest/assets/unmanaged-cluster-note.md" %}}

## Management Clusters

{{% include "/docs/latest/assets/mgmt-desc.md" %}}

{{% include "/docs/latest/assets/tce-feedback.md" %}}

## Tanzu Community Edition Installation

Tanzu Community Edition consists of the Tanzu CLI and a select set of plugins. You will install Tanzu Community Edition on your local machine and then use the Tanzu CLI on your local machine to deploy ([bootstrap](../glossary/#bootstrap)) a cluster to your chosen target platform.

Installing the Tanzu Community Edition extracts the binaries and configures the plugin repositories. The first time you run the `tanzu` command the installed plugins and plugin repositories are initialized. This action might take a minute.

{{< tabs tabTotal="3" tabID="1" tabName1="Linux" tabName2="Mac" tabName3="Windows">}}
{{< tab tabNum="1" >}}

{{% include "/docs/latest/assets/prereq-linux.md" %}}
{{% include "/docs/latest/assets/cli-install-linux.md" %}}

{{< /tab >}}
{{< tab tabNum="2" >}}

{{% include "/docs/latest/assets/prereq-mac.md" %}}
{{% include "/docs/latest/assets/cli-install-mac.md" %}}

{{< /tab >}}
{{< tab tabNum="3" >}}

{{% include "/docs/latest/assets/prereq-windows.md" %}}
{{% include "/docs/latest/assets/cli-install-windows.md" %}}

{{< /tab >}}
{{< /tabs >}}

## Creating Clusters

{{< tabs tabTotal="4" tabID="2" tabName1="AWS" tabName2="Azure" tabName3="Docker" tabName4="vSphere" >}}
{{< tab tabNum="1" >}}

{{% include "/docs/latest/assets/aws-clusters.md" %}}

{{< /tab >}}
{{< tab tabNum="2" >}}

{{% include "/docs/latest/assets/azure-clusters.md" %}}

{{< /tab >}}
{{< tab tabNum="3" >}}

### ⚠️ If bootstrapping docker-based clusters on Windows, [see our Windows guide](../ref-windows-capd)

{{% include "/docs/latest/assets/capd-clusters.md" %}}

{{< /tab >}}
{{< tab tabNum="4" >}}

{{% include "/docs/latest/assets/vsphere-clusters.md" %}}

{{< /tab >}}
{{< /tabs >}}

{{% include "/docs/latest/assets/package-installation.md" %}}
{{% include "/docs/latest/assets/octant-install.md" %}}
{{% include "/docs/latest/assets/clean-up.md" %}}
{{% include "/docs/latest/assets/next-steps.md" %}}
