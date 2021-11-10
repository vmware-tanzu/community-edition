# Planning Your Installation

There are four main steps involved in deploying Tanzu Community Edition. The following section describes the main steps. Specific end to end installation/deployment documentation is available below for each target platform (infrastructure provider):

1. **Install Tanzu Community Edition.**

   You will download this from GitHub and install it on your desktop machine. This installs the Tanzu CLI. For information about the supported operating systems and prerequisites for your desktop machine, see the[Support Matrix](support-matrix/#local-client-bootstrap-machine-prerequisites). For information about the Tanzu Community Edition architecture, see [Architecture](architecture).

1. **Prepare to deploy clusters.**
   Choose the [target platform](installation-planning/#target-platform-infrastructure-provider) where you want to deploy clusters and ensure that the prerequisites are met for this platform.
   {{% include "/docs/assets/support-matrix.md" %}}

1. **Deploy a cluster to your target platform.**

   There are two ways to approach this:

   * Use the Tanzu CLI to launch the Tanzu Community Edition installer, deploy a [management cluster](installation-planning/#managed-clusters), and then deploy a [workload cluster](installation-planning/#workload-cluster).

     **or**

   * Use the Tanzu CLI to launch the installer and deploy a [standalone cluster](installation-planning/#standalone-clusters).

   **Note**: the installer is a web based interface, if you need to perform an installation on a machine that does not have a desktop environment, see [Headless Installation](headless-install).

1. **Install and configure packages.**

   Use the Tanzu CLI to install and configure [Packages](installation-planning/#package).

1. **Start here:**
   ||
   |:------------------------ |
   |**If your target platform is AWS start [here](aws-intro).**|
   |**If your target platform is Microsoft Azure start [here](azure-intro).**|
   |**If your target platform is Docker start [here](docker-intro).**|
   |**If your target platform is vSphere start [here](vsphere-intro).**|
