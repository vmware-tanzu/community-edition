Create a cluster on your infrastructure provider. There are two ways to approach this:  
   * Create a management cluster and then create a workload cluster. First, create the management cluster using the Tanzu Kubernetes Grid Installer. This installer is initiated from the Tanzu CLI. Then, create a workload cluster using the Tanzu CLI.  [workload](installation-planning/#standalone-cluster-description)
)
         
   or  
      
   * Create a standalone cluster using Tanzu Kubernetes Grid Installer. This is a faster approach to achieve a working cluster for a development environment. 

   There are three infrastructure providers:   

    * vSphere
    * Amazon EC2
    * Docker

    For descriptions cluster types, see , workload and standalone clusters,