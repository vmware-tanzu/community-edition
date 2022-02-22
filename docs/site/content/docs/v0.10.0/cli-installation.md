# Install Tanzu Community Edition

Tanzu Community Edition consists of the Tanzu CLI and a select set of plugins. You will install Tanzu Community Edition on your local machine and then use the Tanzu CLI on your local machine to deploy ([bootstrap](../glossary/#bootstrap)) a cluster to your chosen target platform.

Installing the Tanzu Community Edition extracts the binaries and configures the plugin repositories. The first time you run the `tanzu` command the installed plugins and plugin repositories are initialized. This action might take a minute.

{{% include "/docs/v0.10.0/assets/unmanaged-cluster-note.md" %}}

## Procedure

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
