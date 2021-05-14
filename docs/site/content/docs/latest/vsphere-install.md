### <a id="iaas-vsphere"></a> Configure a vSphere Infrastructure Provider

1. In the **IaaS Provider** section, enter the IP address or fully qualified domain name (FQDN) for the vCenter Server instance on which to deploy the management cluster.

   Tanzu Kubernetes Grid does not support IPv6 addresses. This is because upstream Kubernetes only provides alpha support for IPv6. Always provide IPv4 addresses in the procedures in this topic.
1. Enter the vCenter Single Sign On username and password for a user account that has the required privileges for Tanzu Kubernetes Grid operation, and click **Connect**.

   ![Configure the connection to vSphere](../images/install-v-1iaas.png)

1. Verify the SSL thumbprint of the vCenter Server certificate and click **Continue** if it is valid.

   For information about how to obtain the vCenter Server certificate thumbprint, see [Obtain vSphere Certificate Thumbprints](vsphere.md#vc-thumbprint).

   ![Verify vCenter Server certificate thumbprint](../images/vsphere-thumprint.png)

1. If you are deploying a management cluster to a vSphere 7 instance, confirm whether or not you want to proceed with the deployment.   

   On vSphere 7, the vSphere with Tanzu option includes a built-in supervisor cluster that works as a management cluster and provides a better experience than a separate management cluster deployed by Tanzu Kubernetes Grid.  Deploying a Tanzu Kubernetes Grid management cluster to vSphere 7 when vSphere with Tanzu is not enabled is supported, but the preferred option is to enable vSphere with Tanzu and use the Supervisor Cluster. VMware Cloud on AWS and Azure VMware Solution do not support a supervisor cluster, so you need to deploy a management cluster.
   For information, see [Use the Tanzu CLI with a vSphere with Tanzu Supervisor Cluster](../tanzu-k8s-clusters/connect-vsphere7.md).

   To reflect the recommendation to use vSphere with Tanzu when deploying to vSphere 7, the Tanzu Kubernetes Grid installer behaves as follows:

      - **If vSphere with Tanzu is enabled**, the installer informs you that deploying a management cluster is not possible, and exits.
      - **If vSphere with Tanzu is not enabled**, the installer informs you that deploying a Tanzu Kubernetes Grid management cluster is possible but not recommended, and presents a choice:
          - **Configure vSphere with Tanzu** opens the vSphere Client so you can configure your Supervisor Cluster as described in [Configuring and Managing a Supervisor Cluster](https://docs.vmware.com/en/VMware-vSphere/7.0/vmware-vsphere-with-tanzu/GUID-21ABC792-0A23-40EF-8D37-0367B483585E.html) in the vSphere documentation.
          - **Deploy TKG Management Cluster** allows you to continue deploying a management cluster, against recommendation for vSphere 7, but as required for VMware Cloud on AWS and Azure VMware Solution. When using vSphere 7, the preferred option is to enable vSphere with Tanzu and use the built-in Supervisor Cluster instead of deploying a Tanzu Kubernetes Grid management cluster.

   ![Deploy management cluster to vSphere 7](../images/vsphere7-detected.png)

1. Select the datacenter in which to deploy the management cluster from the **Datacenter** drop-down menu.

1. Paste the contents of your SSH public key into the text box and click **Next**.

   ![Select datacenter and provide SSH public key](../images/dc-ssh-vsphere.png)

For the next steps, go to [Configure the Management Cluster Settings](#config-mgmt-cluster).


