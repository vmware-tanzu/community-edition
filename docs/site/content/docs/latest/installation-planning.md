# DRAFT DRAFT WIP WIP DRAFT DRAFT WIP WIP  

There are three main steps involved in deploying Tanzu Community Edition. The following section describes the main steps and how they are invoked: 

1. Install the Tanzu CLI.  
   You will download this from GitHub and install it on your desktop machine.
2. Create a cluster on your infrastructure. There are two ways to approach this:  
    - Create a management cluster and then create a workload cluster. First, create the management cluster using the Tanzu Kubernetes Grid Installer. This installer is initiated from the Tanzu CLI. Then, create a workload cluster using the Tanzu CLI.   
         
       or  
      
    - Create a stand-alone cluster using Tanzu Kubernetes Grid Installer.

    There are three infrastructure providers:   

    - vSphere
    - Amazon EC2
    - Docker as Cluster API provider (CAPD) 
4. Install and configure packages using the Tanzu CLI.


This section provides descriptions of the components you deploy, and the elements required in the deployment.


## Management cluster
The management cluster provides management and operations for your instance. It runs [Cluster-API](https://cluster-api.sigs.k8s.io/) which is used to create workload clusters, as well as creating shared services for all the clusters within the instance.  The management cluster is not intended to be used for application workloads. A management cluster is deployed using the Tanzu Kubernetes Grid Installer.

When you create a management cluster, a bootstrap cluster is created on your local machine. This is a [Kind](https://kind.sigs.k8s.io/)  based cluster -  a cluster in a container.  This bootstrap cluster then creates a cluster on your specified provider. The Cluster APIs then pivots this cluster into a management cluster. 
At this point, the local bootstrap cluster is deleted.  The management cluster can now instantiate more workload clusters. 

## Workload cluster description

After you deploy the management cluster, you can deploy a workload cluster. The workload cluster is deployed by the management cluster. The workload cluster is used to run your application workloads. The workload clusters is deployed using the Tanzu CLI.

## Stand-alone cluster description
A stand-alone cluster is a faster way to get a functioning cluster with minimal resources. A stand-alone cluster functions as a workload cluster, it can run application workloads. A stand-alone cluster does not contain any of the management components.  A stand-alone cluster is deployed using the Tanzu Kubernetes Grid Installer.

When you create a stand-alone cluster, a bootstrap cluster is created on your local machine. This is a [Kind](https://kind.sigs.k8s.io/)  based cluster - a cluster in a container.  This bootstrap cluster then creates a cluster on your specified provider, but it does not pivot into a management cluster - it functions as a workload cluster.  A workload cluster can be pivoted back to be a management cluster at a later point.


## Bootstrap Machine
The bootstrap machine is the laptop, host, or server on which you download and run the Tanzu CLI. This is where the initial bootstrapping of a management or stand-alone cluster occurs before it is pushed to the platform where it will run. You run tanzu, kubectl and other commands on the bootstrap machine. The bootstrap machine can be a local physical machine or a VM that you access via a console window or client shell.


## Tanzu Kubernetes Grid Installer
The Tanzu Kubernetes Grid installer is a graphical wizard that you start up by running the ``tanzu management-cluster create --ui`` command. The installer wizard runs locally in a browser on the bootstrap machine and provides a user interface to guide you through the process of deploying a management or stand-alone cluster. 

The installer interface launches in a browser and takes you through steps to configure the management or standalone cluster.



## Package
