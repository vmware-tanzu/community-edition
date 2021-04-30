# DRAFT DRAFT WIP WIP DRAFT DRAFT WIP WIP  

There are three main steps involved in deploying Tanzu Community Edition. The following section describes the main steps and how they are invoked: 

1. Install the Tanzu CLI.  
   You will download this from GitHub and install it on your desktop machine.
2. Create a cluster. There are two ways to approach this:  
    - Create a management cluster and then create a workload cluster:  
       Create the management cluster using the Tanzu Kubernetes Grid Installer. This installer is initiated from the Tanzu CLI.  
       Create a workload cluster using the Tanzu CLI.   
         
       or  
      
    - Create a stand-alone cluster using Tanzu Kubernetes Grid Installer.
4. Install and configure packages using the Tanzu CLI.




This section provides descriptions of the components you deploy.


## Management cluster description
The management cluster provides management and operations for your instance. It runs Cluster-API which is used to create workload clusters, as well as creating shared services for all the clusters within the instance.  The management cluster is not intended to be used for application workloads. A management cluster is deployed using the Tanzu Kubernetes Grid Installer.

When you create a management cluster, a bootstrap cluster is created on your local machine. This is a [Kind](https://kind.sigs.k8s.io/)  based cluster -  a cluster in a container.  This bootstrap cluster then creates a cluster on your specified provider. The Cluster APIs then pivots this cluster into a management cluster. 
At this point, the local bootstrap cluster is deleted.  The management cluster can now instantiate more workload clusters. 

## Workload cluster description

After you deploy the management cluster, you can deploy a workload cluster. The workload cluster is deployed by the management cluster. The workload cluster is used to run your application workloads. The workload clusters is deployed using the Tanzu CLI.

## Stand-alone cluster description
A stand-alone cluster is a faster way to get a functioning cluster with minimal resources. A stand-alone cluster functions as a workload cluster, it is capable of running application workloads. A stand-alone cluster does not contain any of the management components.  A stand-alone cluster is deployed using the Tanzu Kubernetes Grid Installer.

When you create a stand-alone cluster, a bootstrap cluster is created on your local machine. This is a [Kind](https://kind.sigs.k8s.io/)  based cluster -  a cluster in a container.  This bootstrap cluster then creates a cluster on your specified provider, but it does not pivot into a management cluster - it functions as a workload cluster.  A workload cluster can be pivoted back to be a management cluster at a later point.

