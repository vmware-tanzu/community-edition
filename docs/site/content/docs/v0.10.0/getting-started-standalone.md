# Getting Started with Tanzu Community Edition

This guide walks you through creating a standalone cluster using Tanzu
Community Edition.

## Standalone Clusters

{{% include "/docs/v0.10.0/assets/standalone-desc.md" %}}
{{% include "/docs/v0.10.0/assets/standalone-warning.md" %}}
{{% include "/docs/v0.10.0/assets/tce-feedback.md" %}}

## Tanzu Community Edition Installation

Tanzu Community Edition consists of the Tanzu CLI and a select set of plugins. You will install Tanzu Community Edition on your local machine and then use the Tanzu CLI on your local machine to deploy a cluster to your chosen target platform.

{{< tabs tabTotal="3" tabID="1" tabName1="Linux" tabName2="Mac" tabName3="Windows">}}
{{< tab tabNum="1" >}}

{{% include "/docs/v0.10.0/assets/prereq-linux.md" %}}
{{% include "/docs/v0.10.0/assets/cli-install-linux.md" %}}

{{< /tab >}}
{{< tab tabNum="2" >}}

{{% include "/docs/v0.10.0/assets/prereq-mac.md" %}}
{{% include "/docs/v0.10.0/assets/cli-install-mac.md" %}}

{{< /tab >}}
{{< tab tabNum="3" >}}

{{% include "/docs/v0.10.0/assets/prereq-windows.md" %}}
{{% include "/docs/v0.10.0/assets/cli-install-windows.md" %}}

{{< /tab >}}
{{< /tabs >}}

## Creating Clusters

{{< tabs tabTotal="4" tabID="2" tabName1="AWS" tabName2="Azure" tabName3="Docker" tabName4="vSphere" >}}
{{< tab tabNum="1" >}}

{{% include "/docs/v0.10.0/assets/aws-standalone-clusters.md" %}}

{{< /tab >}}
{{< tab tabNum="2" >}}

{{% include "/docs/v0.10.0/assets/azure-standalone-clusters.md" %}}

{{< /tab >}}
{{< tab tabNum="3" >}}

### ⚠️ If bootstrapping docker-based clusters on Windows, [see our Windows guide](../ref-windows-capd)

{{% include "/docs/v0.10.0/assets/capd-standalone-clusters.md" %}}

{{< /tab >}}
{{< tab tabNum="4" >}}

{{% include "/docs/v0.10.0/assets/vsphere-standalone-clusters.md" %}}

{{< /tab >}}
{{< /tabs >}}

{{% include "/docs/v0.10.0/assets/package-installation.md" %}}
{{% include "/docs/v0.10.0/assets/octant-install.md" %}}
{{% include "/docs/v0.10.0/assets/clean-up-standalone.md" %}}
