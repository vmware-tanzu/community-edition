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

## Glossary

The following section provides a glossary of the main components and concepts involved in a Tanzu Community Edition deployment:

{{% include "/docs/assets/mgmt-desc.md" %}}

### Workload cluster

After you deploy the management cluster, you can deploy a workload cluster. The workload cluster is deployed by the management cluster. The workload cluster is used to run your application workloads. The workload cluster is deployed using the Tanzu CLI.

{{% include "/docs/assets/standalone-desc.md" %}}

### Bootstrap

The bootstrap (noun) machine is the laptop, host, or server on which you download and run the Tanzu CLI. This is where the initial bootstrapping (verb) of a management or standalone cluster occurs before it is pushed to the platform where it will run. You run tanzu, kubectl and other commands on the bootstrap machine.

Using the Tanzu CLI to deploy a cluster to a target platform is often referred to as bootstrapping (verb).

### Tanzu Community Edition installer

The Tanzu Community Edition installer (the installer) is a graphical wizard that you launch in your browser by running the ``tanzu management-cluster create --ui`` command. The installer runs locally in a browser on the bootstrap machine and provides a user interface to guide you through the process of deploying a management or standalone cluster.

### Target Platform (Infrastructure Provider)

The target platform is the cloud provider or local Docker where you will deploy your cluster. This is also referred to as your infrastructure provider.
There are four available target platforms:

* AWS
* Microsoft Azure
* Docker
* vSphere

### Package

{{% include "/docs/assets/package-description.md" %}}

### Add-ons

Same as packages (see above).

### Extensions

Same as packages (see above).

### Kind cluster

During the deployment of the management or standalone cluster, either from the installer interface or the CLI, Tanzu Kubernetes Grid creates a temporary management cluster using a [Kubernetes in Docker](https://kind.sigs.k8s.io/), `kind`, cluster on the bootstrap machine. Then, Tanzu Community Edition uses it to provision the final management cluster to the platform of your choice, depending on whether you are deploying to vSphere, Amazon EC2, Azure, or Docker. After the deployment of the management cluster finishes successfully, the temporary `kind` cluster is deleted.

### Tanzu CLI

Tanzu CLI provides commands that facilitate many of the operations that you can perform with your clusters. However, for certain operations, you still need to use `kubectl`.
