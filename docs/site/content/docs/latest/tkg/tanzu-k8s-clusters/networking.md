# Customize Tanzu Kubernetes Cluster Networking

This topic describes how to customize networking for Tanzu Kubernetes (workload) clusters,
including using a cluster network interface (CNI) other than the default [Antrea](https://antrea.io/),
and supporting publicly-routable, no-NAT IP addresses for workload clusters on vSphere with NSX-T networking.

## <a id="nondefault-cni"></a> Deploy a Cluster with a Non-Default CNI

When you use the Tanzu CLI to deploy a Tanzu Kubernetes cluster, an [Antrea](https://antrea.io/) cluster network interface (CNI) is automatically enabled in the cluster. Alternatively, you can enable a [Calico](https://www.projectcalico.org/) CNI or your own CNI provider.

Existing Tanzu Kubernetes clusters that you deployed with a version of Tanzu Kubernetes Grid earlier than 1.2.x and then upgrade to v1.3 continue to use Calico as the CNI provider. You cannot change the CNI provider for these clusters.

You can change the default CNI for a Tanzu Kubernetes cluster by specifying the `CNI` variable in the configuration file. The `CNI` variable supports the following options:

- (**Default**) `antrea`: Enables Antrea.
- `calico`: Enables Calico. See [Enable Calico](#calico) below.
- `none`: Allows you to enable a custom CNI provider. See [Enable a Custom CNI Provider](#custom-cni) below.

If you do not specify the `CNI` variable, Antrea is enabled by default.

```
CNI: antrea

#! ---------------------------------------------------------------------
#! Antrea CNI configuration
#! ---------------------------------------------------------------------
ANTREA_NO_SNAT: false
ANTREA_TRAFFIC_ENCAP_MODE: "encap"
ANTREA_PROXY: false
ANTREA_POLICY: true 
ANTREA_TRACEFLOW: false
```

### <a id="calico"></a> Enable Calico

To enable Calico in a Tanzu Kubernetes cluster, specify the following in the configuration file:

```
CNI: calico
```

After the cluster creation process completes, you can examine the cluster
as described in [Retrieve Tanzu Kubernetes Cluster `kubeconfig`](../cluster-lifecycle/connect.md#kubeconfig) and
[Examine the Deployed Cluster](../cluster-lifecycle/connect.md#examine).

### <a id="custom-cni"></a> Enable a Custom CNI Provider

To enable a custom CNI provider in a Tanzu Kubernetes cluster, follow the steps below:

1. Specify `CNI: none` in the configuration file when you create the cluster. For example:

   ```
   CNI: none
   ```

   The cluster creation process will not succeed until you apply a CNI to the cluster. You can monitor the cluster creation process in the Cluster API logs on the management cluster. For instructions on how to access the Cluster API logs, see [Monitor Workload Cluster Deployments in Cluster API Logs](../troubleshooting-tkg/tips.md#workload-logs).

1. After the cluster has been initialized, apply your CNI provider to the cluster:

   1. Get the `admin` credentials of the cluster. For example:

       ```
       tanzu cluster kubeconfig get my-cluster --admin
       ```

   1. Set the context of `kubectl` to the cluster. For example:

       ```
       kubectl config use-context my-cluster-admin@my-cluster
       ```

   1. Apply the CNI provider to the cluster:

       ```
       kubectl apply -f PATH-TO-YOUR-CNI-CONFIGURATION/example.yaml
       ```

1. Monitor the status of the cluster by using the `tanzu cluster list` command.
When the cluster creation completes, the cluster status changes from `creating`
to `running`. For more information about how to examine your cluster, see
[Connect to and Examine Tanzu Kubernetes Clusters](../cluster-lifecycle/connect.md).

## <a id="multiple-cnis"></a> Enable Multiple CNI Providers

To enable multiple CNI providers on a workload cluster, such as [macvlan](https://www.cni.dev/plugins/current/main/macvlan/), [ipvlan](https://www.cni.dev/plugins/current/main/ipvlan/), [SR-IOV](https://docs.vmware.com/en/VMware-vSphere/7.0/com.vmware.vsphere.networking.doc/GUID-CC021803-30EA-444D-BCBE-618E0D836B9F.html) or [DPDK](https://www.dpdk.org/), install the Multus package onto a cluster that is already running Antrea or Calico CNI, and create additional `NetworkAttachmentDefinition` resources for CNIs.
Then you can create new pods in the cluster that use different network interfaces for different address ranges.

For directions, see [Implementing Multiple CNIs with Multus](../extensions/cni-multus.md).

## <a id="routable-ip"></a> Deploy Pods with Routable, No-NAT IP Addresses (NSX-T)

On **vSphere** with **NSX-T** networking and the **Antrea** container network interface (CNI), you can configure a Kubernetes workload cluster with routable IP addresses for its worker pods, bypassing network address translation (NAT) for external requests from and to the pods.

Routable IP addresses on pods let you:

- Trace outgoing requests to common shared services, because their source IP address is the routable pod IP address, not a NAT address.
- Support authenticated incoming requests from the external internet directly to pods, bypassing NAT.

The following sections explain how to deploy Tanzu Kubernetes Grid workload clusters with routable-IP pods.
The range of routable IP addresses is set with the cluster's `CLUSTER_CIDR` configuration variable.

### <a id="routable-nsxt"></a> Configure NSX-T for Routable-IP Pods

1. Browse to your **NSX-T** server and open the **Networking** tab.

1. Under **Connectivity** > **Tier-1 Gateways**, click **Add Tier-1 Gateway** and configure a new Tier-1 gateway dedicated to routable-IP pods:

  - **Name**: Make up a name for your routable pods T1 gateway.
  - **Linked Tier-0 Gateway**: Select the Tier-0 gateway that your other Tier-1 gateways for Tanzu Kubernetes Grid use.
  - **Edge Cluster**: Select an existing edge cluster.
  - **Route Advertisement**: Enable **All Static Routes**, **All NAT IP's**, and **All Connected Segments & Service Ports**.
  
  Click **Save** to save the gateway.

1. Under **Connectivity** > **Segments**, click **Add Segment** and configure a new NSX-T segment, a logical switch, for the workload cluster nodes containing the routable pods:

  - **Name**: Make up a name for the network segment for the workload cluster nodes.
  - **Connectivity**: Select the Tier-1 gateway that you just created.
  - **Transport Zone**: Select an overlay transport zone, such as `tz-overlay`.
  - **Subnets**: Choose an IP address range for cluster nodes, such as `195.115.4.1/24`. This range should not overlap with DHCP profile **Server IP Address** values.
  - **Route Advertisement**: Enable **All Static Routes**, **All NAT IP's**, and **All Connected Segments & Service Ports**.

  Click **Save** to save the gateway.

### <a id="routable-config"></a> Deploy a Cluster with Routable-IP Pods

To deploy a workload cluster that has no-NAT, publicly-routable IP addresses for its worker  pods:

1. Create a workload cluster configuration file as described in [Create a Tanzu Kubernetes Cluster Configuration File](deploy.md#config) and as follows:

  - To set the block of routable IP addresses assigned to worker pods, you can either:
      - Set `CLUSTER_CIDR` in the workload cluster configuration file, or
      - Prepend your `tanzu cluster create` command with a `CLUSTER_CIDR=` setting, as shown in the following step.
  - Set `NSXT_POD_ROUTING_ENABLED` to `"true"`.
  - Set `NSXT_MANAGER_HOST` to your NSX-T manager IP address.
  - Set `NSXT_ROUTER_PATH` to the inventory path of the newly-added Tier-1 gateway for routable IPs. Obtain this from NSX-T manager > **Connectivity** > **Tier-1 Gateways** by clicking the menu icon (<img src="../images/ellipsis-vertical-line.png" alt="Clarity vertical ellipsis icon" width="18" height="18"/>) to the left of the gateway name and clicking **Copy Path to Clipboard**. The name starts with `"/infra/tier-1s/`
  - Set other `NSXT_` string variables for accessing NSX-T by following the [NSX-T Pod Routing](../tanzu-config-reference.md#nsxt-pod-routing) table in the _Tanzu CLI Configuration File Variable Reference_.
  Pods can authenticate with NSX-T in one of four ways, with the least secure listed last:
      - **Certificate**: Set `NSXT_CLIENT_CERT_KEY_DATA`, `NSXT_CLIENT_CERT_KEY_DATA`, and for a CA-issued certificate, `NSXT_ROOT_CA_DATA_B64`.
      - **VMware Identity Manager** token on VMware Cloud (VMC): Set `NSXT_VMC_AUTH_HOST` and `NSXT_VMC_ACCESS_TOKEN`.
      - **Username/password** stored in a Kubernetes secret: Set `NSXT_SECRET_NAMESPACE`, `NSXT_SECRET_NAME`, `NSXT_USERNAME`, and `NSXT_PASSWORD`.
      - **Username/password** as plaintext in configuration file: Set `NSXT_USERNAME` and `NSXT_PASSWORD`.

1. Run `tanzu cluster create` as described in [Deploy Tanzu Kubernetes Clusters](deploy.md).  For example:

  ```
  $ CLUSTER_CIDR=100.96.0.0/11 tanzu cluster create my-routable-work-cluster -f my-routable-work-cluster-config.yaml
  Validating configuration...
  Creating workload cluster 'my-routable-work-cluster'...
  Waiting for cluster to be initialized...
  Waiting for cluster nodes to be available...
  ```

### <a id="routable-test"></a> Validate Routable IPs

To test routable IP addresses for your workload pods:

1. Deploy a webserver to the routable workload cluster.

1. Run `kubectl get pods --o wide` to retrieve `NAME`, `INTERNAL-IP` and `EXTERNAL-IP` values for your routable pods, and verify that the IP addresses listed are identical and are within the routable `CLUSTER_CIDR` range.

1. Run `kubectl get nodes --o wide` to retrieve `NAME`, `INTERNAL-IP` and `EXTERNAL-IP` values for the workload cluster nodes, which contain the routable-IP pods.

1. Log in to a different workload cluster's control plane node:

  1. Run `kubectl config use-context CLUSTER-CONTEXT` to change context to the different cluster.
  1. Run `kubectl get nodes` to retrieve the IP address of the current cluster's control plane node.
  1. Run `ssh capv@CONTROLPLANE-IP` using the IP address you just retrieved.
  1. `ping` and send `curl` requests to the routable IP address where you deployed the webserver, and confirm its responses.
      - `ping` output should list the webserver's routable pod IP as the source address.

1. From a browser, log in to **NSX-T** and navigate to the Tier-1 gateway that you created for routable-IP pods.

1. Click **Static Routes** and confirm that the following routes were created within the routable `CLUSTER_CIDR` range:

  1. A route for pods in the workload cluster's control plane node, with **Next Hops** shown as the address of the control plane node itself.
  1. A route for pods in the workload cluster's worker nodes, with **Next Hops** shown as the addresses of the worker nodes themselves.

### <a id="routable-delete"></a> Delete Routable IPs

After you delete a workload cluster that contains routable-IP pods, you may need to free the routable IP addresses by deleting them from T1 router:

1. In the NSX-T manager > **Connectivity** > **Tier-1 Gateways** select your routable-IP gateway.

1. Under **Static Routes** click the number of routes to open the list.

1. Search for routes that include the deleted cluster name, and delete each one from the menu icon (<img src="../images/ellipsis-vertical-line.png" alt="Clarity vertical ellipsis icon" width="18" height="18"/>) to the left of the route name.

  1. If if a permissions error prevents you from deleting the route from the menu, which may happen if the route is created by a certificate, delete the route via the API:
      1. From the menu next to the route name, select **Copy Path to Clipboard**.
      1. Run `curl -i -k -u 'NSXT_USERNAME:NSXT_PASSWORD' -H 'Content-Type: application/json' -H 'X-Allow-Overwrite: true' -X DELETE https://NSXT_MANAGER_HOST/policy/api/v1/STATIC-ROUTE-PATH` where:
          - `NSXT_MANAGER_HOST`, `NSXT_USERNAME`, and `NSXT_PASSWORD` are your NSX-T manager IP address and credentials
          - `STATIC_ROUTE_PATH` is the path that you just copied to the clipboard. The name starts with `/infra/tier-1s/` and includes `/static-routes/`.

## <a id="node-ip"></a> Customize Cluster Node IP Addresses

You can configure cluster-specific IP address blocks for management or workload cluster nodes.
How you do this depends on the cloud infrastructure that the cluster runs on:

### <a id="node-ip-vsphere"></a> vSphere

On vSphere, the cluster configuration file's `VSPHERE_NETWORK` sets the VM network that Tanzu Kubernetes Grid uses for cluster nodes and other Kubernetes objects.
IP addresses are allocated to nodes by a DHCP server that runs in this VM network, deployed separately from Tanzu Kubernetes Grid.

If you are using NSX-T networking, you can configure DHCP bindings for your cluster nodes by following [Configure DHCP Static Bindings on a Segment](https://docs.vmware.com/en/VMware-NSX-T-Data-Center/3.1/administration/GUID-99ED7B1C-9F3C-4FCB-A088-394F7CBC7CFE.html) in the _VMware NSX-T Data Center_ documentation.

### <a id="node-ip-aws"></a> Amazon EC2

To configure cluster-specific IP address blocks on Amazon EC2, set the following variables in the cluster configuration file as described in the [Amazon EC2](../tanzu-config-reference.md#aws) table in the _Tanzu CLI Configuration File Variable Reference_.

- Set `AWS_PUBLIC_NODE_CIDR` to set an IP address range for public nodes.
    - Make additional ranges available by setting `AWS_PRIVATE_NODE_CIDR_1` or `AWS_PRIVATE_NODE_CIDR_2`
- Set `AWS_PRIVATE_NODE_CIDR` to set an IP address range for private nodes.
    - Make additional ranges available by setting `AWS_PRIVATE_NODE_CIDR_1` and `AWS_PRIVATE_NODE_CIDR_2`
- All node CIDR ranges must lie within the cluster's VPC range, which defaults to `10.0.0.0/16`.
    - Set this range with `AWS_VPC_CIDR` or assign nodes to an existing VPC and address range with `AWS_VPC_ID`.

### <a id="node-ip-azure"></a> Microsoft Azure

To configure cluster-specific IP address blocks on Azure, set the following variables in the cluster configuration file as described in the [Microsoft Azure](../tanzu-config-reference.md#azure) table in the _Tanzu CLI Configuration File Variable Reference_.

- Set `AZURE_NODE_SUBNET_CIDR` to create a new VNET with a CIDR block for worker node IP addresses.
- Set `AZURE_CONTROL_PLANE_SUBNET_CIDR` to create a new VNET with a CIDR block for control plane node IP addresses.
- Set `AZURE_NODE_SUBNET_NAME` to assign worker node IP addresses from the range of an existing VNET.
- Set `AZURE_CONTROL_PLANE_SUBNET_NAME` to assign control plane node IP addresses from the range of an existing VNET.
