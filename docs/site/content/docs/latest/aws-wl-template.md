# Amazon EC2 Workload Cluster Template

When you deploy workload clusters to Amazon EC2, you must specify options in the cluster configuration file to connect to your AWS account and identify the resources that the cluster will use. You can also specify the sizes for the control plane and worker node VMs, distribute nodes across availability zones, and share VPCs between clusters.

The template below includes all of the options that are relevant to deploying workload clusters on Amazon EC2. You can copy this template and update it to deploy workload clusters to Amazon EC2.

Mandatory options are uncommented. Optional settings are commented out.  Default values are included where applicable.

```
#! ---------------------------------------------------------------------
#! Cluster creation basic configuration
#! ---------------------------------------------------------------------

#! CLUSTER_NAME:
CLUSTER_PLAN: dev
NAMESPACE: default
CNI: antrea

#! ---------------------------------------------------------------------
#! Node configuration
#! AWS-only MACHINE_TYPE settings override cloud-agnostic SIZE settings.
#! ---------------------------------------------------------------------

# SIZE:
# CONTROLPLANE_SIZE:
# WORKER_SIZE:
CONTROL_PLANE_MACHINE_TYPE: t3.small
NODE_MACHINE_TYPE: m5.large
# CONTROL_PLANE_MACHINE_COUNT: 1
# WORKER_MACHINE_COUNT: 1
# WORKER_MACHINE_COUNT_0:
# WORKER_MACHINE_COUNT_1:
# WORKER_MACHINE_COUNT_2:

#! ---------------------------------------------------------------------
#! AWS Configuration
#! ---------------------------------------------------------------------

AWS_REGION:
AWS_NODE_AZ: ""
# AWS_NODE_AZ_1: ""
# AWS_NODE_AZ_2: ""
# AWS_VPC_ID: ""
# AWS_PRIVATE_SUBNET_ID: ""
# AWS_PUBLIC_SUBNET_ID: ""
# AWS_PUBLIC_SUBNET_ID_1: ""
# AWS_PRIVATE_SUBNET_ID_1: ""
# AWS_PUBLIC_SUBNET_ID_2: ""
# AWS_PRIVATE_SUBNET_ID_2: ""
# AWS_VPC_CIDR: 10.0.0.0/16
# AWS_PRIVATE_NODE_CIDR: 10.0.0.0/24
# AWS_PUBLIC_NODE_CIDR: 10.0.1.0/24
# AWS_PRIVATE_NODE_CIDR_1: 10.0.2.0/24
# AWS_PUBLIC_NODE_CIDR_1: 10.0.3.0/24
# AWS_PRIVATE_NODE_CIDR_2: 10.0.4.0/24
# AWS_PUBLIC_NODE_CIDR_2: 10.0.5.0/24
AWS_SSH_KEY_NAME:
BASTION_HOST_ENABLED: true

#! ---------------------------------------------------------------------
#! Machine Health Check configuration
#! ---------------------------------------------------------------------

ENABLE_MHC: true
MHC_UNKNOWN_STATUS_TIMEOUT: 5m
MHC_FALSE_STATUS_TIMEOUT: 12m

#! ---------------------------------------------------------------------
#! Common configuration
#! ---------------------------------------------------------------------

# TKG_CUSTOM_IMAGE_REPOSITORY: ""
# TKG_CUSTOM_IMAGE_REPOSITORY_CA_CERTIFICATE: ""

# TKG_HTTP_PROXY: ""
# TKG_HTTPS_PROXY: ""
# TKG_NO_PROXY: ""

ENABLE_AUDIT_LOGGING: true
ENABLE_DEFAULT_STORAGE_CLASS: true

CLUSTER_CIDR: 100.96.0.0/11
SERVICE_CIDR: 100.64.0.0/13

# OS_NAME: ""
# OS_VERSION: ""
# OS_ARCH: ""

#! ---------------------------------------------------------------------
#! Autoscaler configuration
#! ---------------------------------------------------------------------

ENABLE_AUTOSCALER: false
# AUTOSCALER_MAX_NODES_TOTAL: "0"
# AUTOSCALER_SCALE_DOWN_DELAY_AFTER_ADD: "10m"
# AUTOSCALER_SCALE_DOWN_DELAY_AFTER_DELETE: "10s"
# AUTOSCALER_SCALE_DOWN_DELAY_AFTER_FAILURE: "3m"
# AUTOSCALER_SCALE_DOWN_UNNEEDED_TIME: "10m"
# AUTOSCALER_MAX_NODE_PROVISION_TIME: "15m"
# AUTOSCALER_MIN_SIZE_0:
# AUTOSCALER_MAX_SIZE_0:
# AUTOSCALER_MIN_SIZE_1:
# AUTOSCALER_MAX_SIZE_1:
# AUTOSCALER_MIN_SIZE_2:
# AUTOSCALER_MAX_SIZE_2:

#! ---------------------------------------------------------------------
#! Antrea CNI configuration
#! ---------------------------------------------------------------------

# ANTREA_NO_SNAT: false
# ANTREA_TRAFFIC_ENCAP_MODE: "encap"
# ANTREA_PROXY: false
# ANTREA_POLICY: true
# ANTREA_TRACEFLOW: false
```
<!--
## <a id="plans-azs"></a>Tanzu Kubernetes Cluster Plans and Node Distribution across AZs

When you create a `prod` Tanzu Kubernetes cluster on Amazon EC2, Tanzu Kubernetes Grid evenly distributes its control plane and worker nodes across the three Availability Zones (AZs) that you specified in your management cluster configuration. This includes Tanzu Kubernetes clusters that are configured with any of the following:

* The default number of control plane nodes
* The `CONTROL_PLANE_MACHINE_COUNT` setting that is greater than the default number of control plane nodes
* The default number of worker nodes
* The `WORKER_MACHINE_COUNT` setting that is greater than the default number of worker nodes

For example, if you specify `WORKER_MACHINE_COUNT: 5`, Tanzu Kubernetes Grid
deploys two worker nodes in the first AZ, two worker nodes in the second AZ, and one worker node in the third AZ. You can optionally customize this default AZ placement mechanism for worker nodes by following the
instructions in [Configure AZ Placement Settings for Worker Nodes ](#az-placement) below. You cannot customize the default AZ placement mechanism for control plane nodes.

### <a id="az-placement"></a> Configure AZ Placement Settings for Worker Nodes

When creating a `prod` Tanzu Kubernetes cluster on Amazon EC2, you can optionally specify how many worker nodes the `tanzu cluster create` command deploys in each of the three AZs you selected in the Tanzu Kubernetes Grid installer interface or configured in the cluster configuration file.

To do this:

1. Set the following variables in the cluster configuration file:

   * `WORKER_MACHINE_COUNT_0`: Sets the number of worker nodes in the first AZ, `AWS_NODE_AZ`.
   * `WORKER_MACHINE_COUNT_1`: Sets the number of worker nodes in the second AZ, `AWS_NODE_AZ_1`.
   * `WORKER_MACHINE_COUNT_2`: Sets the number of worker nodes in the third AZ, `AWS_NODE_AZ_2`.

1. Create the cluster. For example:

   ```
   tanzu cluster create my-prod-cluster
   ```

## <a id="aws-vpc"></a> Deploy a Cluster that Shares a VPC and NAT Gateway(s) with the Management Cluster

By default, Amazon EC2 imposes a limit of 5 NAT gateways per availability zone. For more information about this limit, see [Resource Usage in Your Amazon Web Services Account](../mgmt-clusters/aws.md#aws-resources). If you used the option to create a new VPC when you deployed the management cluster, by default, all Tanzu Kubernetes clusters that you deploy from this management cluster will also create a new VPC and one or three NAT gateways: one NAT gateway for development clusters and three NAT gateways, one in each of your availability zones, for production clusters. So as not to hit the limit of 5 NAT gateways per availability zone, you can modify the configuration with which you deploy Tanzu Kubernetes clusters so that they reuse the VPC and NAT gateway(s) that were created when the management cluster was deployed.

Configuring Tanzu Kubernetes clusters to share a VPC and NAT gateway(s) with their management cluster depends on how the management cluster was deployed:

- It was deployed with the option to create a new VPC, either by selecting the option in the UI or by specifying `AWS_VPC_CIDR` in the cluster configuration file.
- Ideally, `tanzu cluster create` was used with the `--file` option to save the cluster configuration to a different location than the default `.tanzu/tkg/cluster-config.yaml` file.

To deploy Tanzu Kubernetes clusters that reuse the same VPC as the management cluster, you must modify the configuration file from which you deploy Tanzu Kubernetes clusters.

If you deployed the management cluster with the option to reuse an existing VPC, all Tanzu Kubernetes clusters will share that VPC and its NAT gateway(s), and no action is required.

1. Open the cluster configuration file for the management cluster in a text editor.
1. Update the setting for `AWS_VPC_ID` with the ID the VPC that was created when the management cluster was deployed.

   You can obtain this ID from your Amazon EC2 dashboard. Alternatively, you can obtain it by running `tanzu management-cluster create --ui`, selecting **Deploy to AWS EC2** and consulting the value that is provided if you select **Select an existing VPC** in the VPC for AWS section of the installer interface. Cancel the deployment when you have copied the VPC ID.

   ![Configure the connection to AWS](../images/aws-existing-vpc.png)
1. Update the settings for the `AWS_PUBLIC_SUBNET_ID` and `AWS_PRIVATE_SUBNET_ID` variables.
If you are deploying a `prod` Tanzu Kubernetes cluster, update `AWS_PUBLIC_SUBNET_ID`, `AWS_PUBLIC_SUBNET_ID_1`, and
`AWS_PUBLIC_SUBNET_ID_2` and `AWS_PRIVATE_SUBNET_ID`, `AWS_PRIVATE_SUBNET_ID_1`, and `AWS_PRIVATE_SUBNET_ID_2`.

   You can obtain the network information from the VPC dashboard.
1. Save the cluster configuration file.
1. Run the `tanzu cluster create` command with the `--file` option, specifying the modified cluster configuration file.

   ```
   tanzu cluster create my-cluster --file my-cluster-config.yaml
   ```

## <a id="own-vpc"></a> Deploy a Cluster to an Existing VPC and Add Subnet Tags

If both of the following are true, you must add the `kubernetes.io/cluster/YOUR-CLUSTER-NAME=shared` tag to the public subnet or subnets that you intend to use for your Tanzu Kubernetes cluster:

* You want to deploy the cluster to an existing VPC that was not created by Tanzu Kubernetes Grid.
* You want to create services of type `LoadBalancer` in the cluster.

Adding the `kubernetes.io/cluster/YOUR-CLUSTER-NAME=shared` tag to the public subnet or subnets enables you to create services of type `LoadBalancer` after you deploy the cluster.
To add this tag and then deploy the cluster, follow the steps below:

1. Gather the ID or IDs of the public subnet or subnets within your existing VPC that you want to use for the cluster. To deploy a `prod` Tanzu Kubernetes cluster, you must provide three subnets.

1. Create the required tag by running the following command:

    ```
    aws ec2 create-tags --resources YOUR-PUBLIC-SUBNET-ID-OR-IDS --tags Key=kubernetes.io/cluster/YOUR-CLUSTER-NAME,Value=shared
    ```

    Where:

    * `YOUR-PUBLIC-SUBNET-ID-OR-IDS` is the ID or IDs of the public subnet or subnets that you gathered in the previous step.
    * `YOUR-CLUSTER-NAME` is the name of the Tanzu Kubernetes cluster that you want to create.

    For example:

    ```
    aws ec2 create-tags --resources subnet-00bd5d8c88a5305c6 subnet-0b93f0fdbae3436e8 subnet-06b29d20291797698 --tags Key=kubernetes.io/cluster/my-cluster,Value=shared
    ```

1. Create the cluster. For example:

    ```
    tanzu cluster create my-cluster
    ```

## <a id="cluster-aws"></a> Deploy a Prod Cluster from a Dev Management Cluster

When you create a `prod` Tanzu Kubernetes cluster from a
`dev` management cluster that is running on Amazon EC2, you must define a subset of additional variables in the cluster configuration file, which defaults to `.tanzu/tkg/cluster-config.yaml`, before running the `tanzu cluster create` command. This enables Tanzu Kubernetes Grid
to create the cluster and spread its control plane and worker nodes across AZs.

To create a `prod` Tanzu Kubernetes cluster from a `dev` management cluster on Amazon EC2,
perform the steps below:

1. Set the following variables in the cluster configuration file:

   * Set `PLAN` to `prod`.
   * `AWS_NODE_AZ` variables: `AWS_NODE_AZ` was set when you deployed your `dev` management cluster.
   For the `prod` Tanzu Kubernetes cluster, add `AWS_NODE_AZ_1` and `AWS_NODE_AZ_2`.
   * `AWS_PUBLIC_NODE_CIDR` (new VPC) or `AWS_PUBLIC_SUBNET_ID` (existing VPC) variables:
   `AWS_PUBLIC_NODE_CIDR` or `AWS_PUBLIC_SUBNET_ID` was set when you deployed your `dev` management cluster.
   For the `prod` Tanzu Kubernetes cluster, add one of the following:
     * `AWS_PUBLIC_NODE_CIDR_1` and `AWS_PUBLIC_NODE_CIDR_2`
     * `AWS_PUBLIC_SUBNET_ID_1` and `AWS_PUBLIC_SUBNET_ID_2`
   * `AWS_PRIVATE_NODE_CIDR` (new VPC) or `AWS_PRIVATE_SUBNET_ID` (existing VPC) variables:
   `AWS_PRIVATE_NODE_CIDR` or `AWS_PRIVATE_SUBNET_ID` was set when you deployed your `dev` management cluster.
   For the `prod` Tanzu Kubernetes cluster, add one of the following:
     * `AWS_PRIVATE_NODE_CIDR_1` and `AWS_PRIVATE_NODE_CIDR_2`
     * `AWS_PRIVATE_SUBNET_ID_1` and `AWS_PRIVATE_SUBNET_ID_2`

1. (Optional) Customize the default AZ placement mechanism for the worker nodes that you intend to deploy by following the instructions
in [Configure AZ Placement Settings for Worker Nodes](#az-placement).
By default, Tanzu Kubernetes Grid distributes `prod` worker nodes evenly across the AZs.

1. Deploy the cluster by running the `tanzu cluster create` command. For example:

   ```
   tanzu cluster create my-cluster
   ```

## What to Do Next

Advanced options that are applicable to all infrastructure providers are described in the following topics:

- [Deploy Tanzu Kubernetes Clusters with Different Kubernetes Versions](k8s-versions.md)
- [Customize Tanzu Kubernetes Cluster Networking](networking.md)
- [Create Persistent Volumes with Storage Classes](storage.md)
- [Configure Tanzu Kubernetes Plans and Clusters](config-plans.md)

After you have deployed your cluster, see [Managing Cluster Lifecycles](../cluster-lifecycle/index.md).-->
