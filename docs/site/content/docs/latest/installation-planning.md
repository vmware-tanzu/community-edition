# DRAFT DRAFT WIP WIP DRAFT DRAFT WIP WIP  

There are three main steps involved in deploying Tanzu Community Edition. The following section describes the main steps and how they are invoked: 

1. Install the Tanzu CLI. You will download this from GitHub.
2. Create a cluster for your workloads. There are two ways to approach this:  
    a. Create a management cluster and then create a workload cluster:  
        - Create the management cluster using the Tanzu Kubernetes Grid Installer. This installer is initiated from the Tanzu CLI using the ``tanzu management-cluster create --ui`` command.  
        - Create a workload cluster using the Tanzu CLI.   
    b. Create a stand-alone cluster for your workloads.  A stand-alone cluster can be a quicker method to get a cluster up and running.    
4. Install and configure packages on your cluster using the Tanzu CLI.



This section provides descriptions of the components you deploy.


## Management cluster description
After you deploy the Tanzu ClI, this is the first element you deploy. The management cluster provides management and operations for your instance. It runs Cluster-API which is used to create workload clusters, as well as creating shared services for all the clusters within the instance.  The management cluster is not intended to be used for application workloads. A management cluster is deployed using the Tanzu Kubernetes Grid Installer.

## Workload cluster description

After you deploy the management cluster, you can deploy a workload cluster. The workload cluster is deployed by the management cluster. The workload cluster is used to run your application workloads. The workload clusters is deployed using the Tanzu CLI.

## Stand-alone cluster description
A stand-alone cluster is a faster way to get a functioning cluster with minimal resources. A stand-alone cluster does not contain any of the management components.  A stand-alone cluster is deployed using the Tanzu Kubernetes Grid Installer.

