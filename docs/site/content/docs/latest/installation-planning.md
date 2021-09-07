# Planning Your Installation
There are four main steps involved in deploying Tanzu Community Edition. The following section describes the main steps and how they are invoked:

1. Install Tanzu Community Edition.
   You will download this from GitHub and install it on your desktop machine. This installs the Tanzu CLI.
   For more information, see [Installing the Tanzu CLI](cli-installation).
2. Prepare to deploy a cluster. For more information, see [Preparing to Deploy a Cluster](prepare-deployment).
3. Create a cluster on your [target platform](installation-planning/#target-platform). There are two ways to approach this:
   * Use the Tanzu CLI to open the Tanzu Community Edition installer interface, create a [management cluster](installation-planning/#management-cluster-description), and then create a [workload cluster](installation-planning/#workload-cluster). <br>
   **or**  <br>

   * Use the Tanzu CLI to open the installer interface and create a [standalone cluster](installation-planning/#standalone-cluster).

   For more information, see [Deploying Clusters](clusters-deploy.md).
4. Install and configure packages using the Tanzu CLI. For more information, see [Packages Overview](packages-intro).


## Component Descriptions
The following section provides descriptions of the main components involved in a Tanzu Community Edition installation:

{{% include "/docs/assets/mgmt-desc.md" %}}

### Workload cluster
After you deploy the management cluster, you can deploy a workload cluster. The workload cluster is deployed by the management cluster. The workload cluster is used to run your application workloads. The workload clusters is deployed using the Tanzu CLI.

{{% include "/docs/assets/standalone-desc.md" %}}

### Bootstrap
The bootstrap (noun) machine is the laptop, host, or server on which you download and run the Tanzu CLI. This is where the initial bootstrapping (verb) of a management or standalone cluster occurs before it is pushed to the platform where it will run. You run tanzu, kubectl and other commands on the bootstrap machine.
Using the Tanzu CLI to deploy a cluster to a target platform is often referred to as bootstrapping (verb).


### Tanzu Community Edition installer interface
The installer interface is a graphical wizard that you start up by running the ``tanzu management-cluster create --ui`` command. The installer interface runs locally in a browser on the bootstrap machine and provides a user interface to guide you through the process of deploying a management or standalone cluster.

The installer interface launches in a browser and takes you through steps to configure the management or standalone cluster.

### Target Platform
The target platform is the cloud provider or local Docker where you will deploy your cluster.
There are four available target platforms:

- Amazon EC2
- Microsoft Azure
- Docker
- vSphere
### Package
{{% include "/docs/assets/package-description.md" %}}

### Kind cluster
During the deployment of the management cluster, either from the installer interface or the CLI, Tanzu Kubernetes Grid creates a temporary management cluster using a [Kubernetes in Docker](https://kind.sigs.k8s.io/), `kind`, cluster on the bootstrap machine. Then, Tanzu Kubernetes Grid uses it to provision the final management cluster on the platform of your choice, depending on whether you are deploying to vSphere, Amazon EC2, or Docker. After the deployment of the management cluster finishes successfully, Tanzu deletes the temporary `kind` cluster.

### Tanzu CLI
Tanzu CLI provides commands that facilitate many of the operations that you can perform with your management cluster. However, for certain operations, you still need to use `kubectl`.