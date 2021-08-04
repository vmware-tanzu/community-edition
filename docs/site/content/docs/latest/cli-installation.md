# CLI Installation
To use Tanzu Community Edition, you must download and run the Tanzu CLI on a local system, commonly referred to as the bootstrap environment. The bootstrap environment is the laptop, host, or server on which the initial bootstrapping of a management cluster or standalone cluster is performed. This is where you run Tanzu  CLI commands.

## Before you begin

* Ensure you install the Tanzu CLI on either Linux or Mac OS.
Windows is not directly supported right now, if you are  bootstrapping from a Windows desktop, you must run Linux from your Windows machine using one of the following methods:
    * Create a Linux VM, for example use VMware Workstation Pro
    * Install Windows Subsystem for Linux, for more information, see [Windows Subsystem for Linux Installation Guide for Windows 10](https://docs.microsoft.com/en-us/windows/wsl/install-win10)
* Ensure your bootstrap machine has the following prerequisites:
    * 6 GB of RAM and a 2-core CPU

## Procedure

{{< tabs tabTotal="2" tabID="1" tabName1="Mac" tabName2="Linux" >}}
{{< tab tabNum="1" >}}

{{% include "/docs/assets/cli-install-mac.md" %}}

{{< /tab >}}
{{< tab tabNum="2" >}}

{{% include "/docs/assets/cli-install-linux.md" %}}

{{< /tab >}}
{{< /tabs >}}