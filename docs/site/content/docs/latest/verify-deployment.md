# Examine the management cluster
Before you progress to installing packages, examine your cluster.

## Procedure

1. Run the following command to verify that management cluster started successfully. If you did not specify a name for the management cluster, it will be something similar to `tkg-mgmt-vsphere-20200323121503` or `tkg-mgmt-aws-20200323140554`.
<!--add content for docker here -what will docker file name be-->
```sh
tanzu management-cluster get
```

2. Examine the folder structure. When Tanzu creates a management cluster for the first time, it creates a folder `~/.config/tanzu/tkg/providers` that contains all of the files required by Cluster API to create the management cluster.
The Tanzu installer interface saves the settings for the management cluster that it creates into a cluster configuration file `~/.config/tanzu/tkg/clusterconfigs/UNIQUE-ID.yaml`, where `UNIQUE-ID` is a generated filename.

3. To view the management cluster objects in vSphere, or Amazon EC2, do the following:
   * If you deployed the management cluster to vSphere, go to the resource pool that you designated when you deployed the management cluster. You should see:

      * One or three control plane VMs, for development or production control plane, respectively, with names similar to `CLUSTER-NAME-control-plane-sx5rp`
      * A worker node VM with a name similar to `CLUSTER-NAME-md-0-6b8db6b59d-kbnk4`
   * If you deployed the management cluster to Amazon EC2, go to the **Instances** view of your EC2 dashboard. You should see the following VMs or instances.
      * One or three control plane VM instances, for development or production control plane, respectively, with names similar to `CLUSTER-NAME-control-plane-bcpfp`
      * A worker node instance with a name similar to `CLUSTER-NAME-md-0-dwfnm`
      * An EC2 bastion host instance with the name `CLUSTER-NAME-bastion`




<!--## <a id="networking"></a> Management Cluster Networking

When you deploy a management cluster, pod-to-pod networking with [Antrea](https://antrea.io/) is automatically enabled in the management cluster.

## <a id="dhcp"></a> Configure DHCP Reservations for the Control Plane Nodes (vSphere Only)

After you deploy a cluster to vSphere, each control plane node requires a static IP address. This includes both management and Tanzu Kubernetes clusters. These static IP addresses are required in addition to the static IP address that you assigned to Kube-VIP when you deploy a managment cluster.

To make the IP addresses that your DHCP server assigned to the control plane nodes static, you can configure a DHCP reservation for each control plane node in the cluster. For instructions on how to configure DHCP reservations, see your DHCP server documentation.

## <a id="verify-deployment"></a>Verify the Deployment of the Management Cluster

After the deployment of the management cluster completes successfully, you can obtain information about your management cluster by:

* Locating the management cluster objects in vSphere, Amazon EC2, or Azure
* Using the Tanzu CLI and `kubectl`-->


