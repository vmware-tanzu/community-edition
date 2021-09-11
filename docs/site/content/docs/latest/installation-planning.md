# Planning Your Installation
There are four main steps involved in deploying Tanzu Community Edition. The following section describes the main steps. Specific end to end installation/deployment documentation is available below for each target platform (infrastructure provider):

1. **Install Tanzu Community Edition.**\
   You will download this from GitHub and install it on your desktop machine. This installs the Tanzu CLI. For information about the supported operating systems and prerequisites for your desktop machine, see the[Support Matrix](support-matrix/#local-client-machine). For information about the Tanzu Community Edition architecture, see [Architecture](architecture).

1. **Prepare to deploy clusters.**\
   Choose the [target platform](installation-planning/#target-platform) where you want to deploy clusters and ensure that the prerequisites are met for this platform.
   {{% include "/docs/assets/support-matrix.md" %}}
1. **Deploy a cluster to your target platform.**\
   There are two ways to approach this:
   * Use the Tanzu CLI to launch the Tanzu Community Edition installer, deploy a [management cluster](installation-planning/#management-cluster-description), and then deploy a [workload cluster](installation-planning/#workload-cluster).

      **or**


   * Use the Tanzu CLI to launch the installer and deploy a [standalone cluster](installation-planning/#standalone-cluster).

1. **Install and configure packages.**\
   Use the Tanzu CLI to install and configure [Packages](installation-planning/#package).

1. **Start here:**
   ||
   |:------------------------ |
   |**If your target platform is Amazon EC2 start [here](aws-intro):**|
   |**If your target platform is Microsoft Azure start [here](azure-intro):**|
   |**If your target platform is Docker start [here](docker-intro):**|
   |**If your target platform is vSphere start [here](vsphere-intro):**|

## Glossary

The following section provides a glossary of the main components and concepts involved in a Tanzu Community Edition deployment:

{{% include "/docs/assets/mgmt-desc.md" %}}

### Workload cluster
After you deploy the management cluster, you can deploy a workload cluster. The workload cluster is deployed by the management cluster. The workload cluster is used to run your application workloads. The workload clusters is deployed using the Tanzu CLI.

{{% include "/docs/assets/standalone-desc.md" %}}

### Bootstrap
The bootstrap (noun) machine is the laptop, host, or server on which you download and run the Tanzu CLI. This is where the initial bootstrapping (verb) of a management or standalone cluster occurs before it is pushed to the platform where it will run. You run tanzu, kubectl and other commands on the bootstrap machine.
Using the Tanzu CLI to deploy a cluster to a target platform is often referred to as bootstrapping (verb).


### Tanzu Community Edition installer
The Tanzu Community Edition installer (the installer) is a graphical wizard that you launch in your browser by running the ``tanzu management-cluster create --ui`` command. The installer runs locally in a browser on the bootstrap machine and provides a user interface to guide you through the process of deploying a management or standalone cluster.

### Target Platform
The target platform is the cloud provider or local Docker where you will deploy your cluster. This is also referred to as your infrastructure provider.
There are four available target platforms:

- Amazon EC2
- Microsoft Azure
- Docker
- vSphere
### Package
{{% include "/docs/assets/package-description.md" %}}

### Kind cluster
During the deployment of the management cluster, either from the installer interface or the CLI, Tanzu Kubernetes Grid creates a temporary management cluster using a [Kubernetes in Docker](https://kind.sigs.k8s.io/), `kind`, cluster on the bootstrap machine. Then, Tanzu Kubernetes Grid uses it to provision the final management cluster to the platform of your choice, depending on whether you are deploying to vSphere, Amazon EC2, or Docker. After the deployment of the management cluster finishes successfully, Tanzu deletes the temporary `kind` cluster.

### Tanzu CLI
Tanzu CLI provides commands that facilitate many of the operations that you can perform with your management cluster. However, for certain operations, you still need to use `kubectl`.