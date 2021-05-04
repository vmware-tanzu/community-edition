# Prepare to Deploy Management Clusters

Before you can use the Tanzu CLI or installer interface to deploy a management cluster, you must make sure that your infrastructure provider is correctly set up.

- For information about how to set up a vSphere infrastructure provider, see [Prepare to Deploy Management Clusters to vSphere](vsphere.md).
- For information about how to set up an Amazon EC2 infrastructure provider, see [Prepare to Deploy Management Clusters to Amazon EC2](aws.md).
- For information about how to set up a Microsoft Azure infrastructure provider, see [Prepare to Deploy Management Clusters to Microsoft Azure](azure.md).

For production deployments, it is strongly recommended to enable identity management for your clusters. For information about the preparatory steps to perform before you deploy a management cluster, see [Enabling Identity Management in Tanzu Kubernetes Grid](enabling-id-mgmt.md).

If you need to deploy Tanzu Kubernetes Grid in an environment with no external Internet access, see [Deploying Tanzu Kubernetes Grid in an Internet-Restricted Environment](airgapped-environments.md). 

To deploy Tanzu Kubernetes Grid to VMware Cloud on AWS or to Azure VMware Solution, see [Prepare a vSphere Management as a Service Infrastructure](prepare-maas.md).
