# Getting Started with Tanzu Community Edition

This guide walks you through standing up a management and workload cluster using
Tanzu Community Edition.

## Management Clusters

{{% include "/docs/assets/mgmt-desc.md" %}}

{{% include "/docs/assets/tce-feedback.md" %}}

## Tanzu Community Edition Installation

Tanzu Community Edition consists of the Tanzu CLI and a select set of plugins. You will install Tanzu Community Edition on your local machine and then use the Tanzu CLI on your local machine to deploy ([bootstrap](../glossary/#bootstrap)) a cluster to your chosen target platform.

Installing the Tanzu Community Edition extracts the binaries and configures the plugin repositories. The first time you run the `tanzu` command the installed plugins and plugin repositories are initialized. This action might take a minute.

{{< tabs tabTotal="3" tabID="1" tabName1="Linux" tabName2="Mac" tabName3="Windows">}}
{{< tab tabNum="1" >}}

{{% include "/docs/assets/prereq-linux.md" %}}
{{% include "/docs/assets/cli-install-linux.md" %}}

{{< /tab >}}
{{< tab tabNum="2" >}}

{{% include "/docs/assets/prereq-mac.md" %}}
{{% include "/docs/assets/cli-install-mac.md" %}}

{{< /tab >}}
{{< tab tabNum="3" >}}

{{% include "/docs/assets/prereq-windows.md" %}}
{{% include "/docs/assets/cli-install-windows.md" %}}

{{< /tab >}}
{{< /tabs >}}

## Creating Clusters

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

{{% include "/docs/assets/package-installation.md" %}}
{{% include "/docs/assets/octant-install.md" %}}
{{% include "/docs/assets/clean-up.md" %}}
