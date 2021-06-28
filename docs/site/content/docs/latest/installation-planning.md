There are four main steps involved in deploying Tanzu Community Edition. The following section describes the main steps and how they are invoked:

1. Install the Tanzu CLI.  
   You will download this from GitHub and install it on your desktop machine.
   For more information, see [Installing the Tanzu CLI](cli-installation).
2. Prepare to create a cluster. For more information, see [Preparing to Deploy a Cluster](prepare-deployment).
2. Create a cluster on your infrastructure provider. There are two ways to approach this:
   * Create a management cluster and then create a workload cluster. First, create the management cluster using the Tanzu Kubernetes Grid Installer. This installer is initiated from the Tanzu CLI. Then, create a workload cluster using the Tanzu CLI.
   or  
      
   * Create a standalone cluster using Tanzu Kubernetes Grid Installer.

   There are three infrastructure providers:   

    * vSphere
    * Amazon EC2
    * Docker

   For more information, see [Deploying Clusters](clusters-deploy.md).
4. Install and configure packages using the Tanzu CLI. For more information, see [Packages Overview](packages-intro).


This section provides descriptions of the components you deploy, and the elements required in the deployment.


## Management cluster description
The management cluster provides management and operations for your instance. It runs [Cluster-API](https://cluster-api.sigs.k8s.io/) which is used to create workload clusters, as well as creating shared services for all the clusters within the instance.  The management cluster is not intended to be used for application workloads. A management cluster is deployed using the Tanzu Kubernetes Grid Installer.

When you create a management cluster, a bootstrap cluster is created on your local machine. This is a [Kind](https://kind.sigs.k8s.io/)  based cluster -  a cluster in a container.  This bootstrap cluster then creates a cluster on your specified provider. The Cluster APIs then pivots this cluster into a management cluster. 
At this point, the local bootstrap cluster is deleted.  The management cluster can now instantiate more workload clusters. 

## Workload cluster description
After you deploy the management cluster, you can deploy a workload cluster. The workload cluster is deployed by the management cluster. The workload cluster is used to run your application workloads. The workload clusters is deployed using the Tanzu CLI.

{{% include "/docs/assets/standalone-desc.md" %}}

## Bootstrap Machine
The bootstrap machine is the laptop, host, or server on which you download and run the Tanzu CLI. This is where the initial bootstrapping of a management or stand-alone cluster occurs before it is pushed to the platform where it will run. You run tanzu, kubectl and other commands on the bootstrap machine. The bootstrap machine can be a local physical machine or a VM that you access via a console window or client shell.


## Tanzu Kubernetes Grid Installer
The Tanzu Kubernetes Grid installer is a graphical wizard that you start up by running the ``tanzu management-cluster create --ui`` command. The installer wizard runs locally in a browser on the bootstrap machine and provides a user interface to guide you through the process of deploying a management or stand-alone cluster. 

The installer interface launches in a browser and takes you through steps to configure the management or standalone cluster.



## Package
{{% include "/docs/assets/package-description.md" %}}

## Kind cluster
During the deployment of the management cluster, either from the installer interface or the CLI, Tanzu Kubernetes Grid creates a temporary management cluster using a [Kubernetes in Docker](https://kind.sigs.k8s.io/), `kind`, cluster on the bootstrap machine. Then, Tanzu Kubernetes Grid uses it to provision the final management cluster on the platform of your choice, depending on whether you are deploying to vSphere, Amazon EC2, or Docker. After the deployment of the management cluster finishes successfully, Tanzu deletes the temporary `kind` cluster.

## Tanzu CLI
Tanzu CLI provides commands that facilitate many of the operations that you can perform with your management cluster. However, for certain operations, you still need to use `kubectl`. 