## Create Standalone Azure Clusters

This section covers setting up a standalone cluster in Azure. A standalone cluster provides a workload cluster that is **not** managed by a centralized management cluster.

There are some prerequisites this process will assume. Refer to the
[Prepare to Deploy a Cluster to Azure](../azure-mgmt) docs for instructions on
accepting image licenses and preparing your Azure account.

1. Initialize the Tanzu Community Edition installer interface.

    ```sh
    tanzu standalone-cluster create --ui
    ```

1. Choose Azure from the provider tiles.

    ![kickstart azure tile](/docs/img/kickstart-azure-tile.png)

1. Fill out the IaaS Provider section.

    ![kickstart azure iaas](/docs/img/kickstart-azure-iaas.png)

    * `A`: Your account's Tenant ID.
    * `B`: Your Client ID.
    * `C`: Your Client secret.
    * `D`: Your Subscription ID.
    * `E`: The Azure Cloud in which to deploy. For example, "Public Cloud", "US
      Government Cloud", etc.
    * `F`: [The region of
      Azure](https://azure.microsoft.com/en-us/global-infrastructure/geographies/#geographies)
      you'd like all networking, compute, etc to be created within.
    * `G`: The public key you'd like to use for your VM instances. This is how
      you'll SSH into control plane and worker nodes.
    * `H`: Whether to use an existing
      [resource group](https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/manage-resource-groups-portal#what-is-a-resource-group)
      or create a new one.
    * `I`: The existing resource group, or the name to provide the new resource group.

1. Fill out the VNET settings.

    <!-- TODO: Swap image with resolution of https://github.com/vmware-tanzu/community-edition/issues/1886 -->
    ![kickstart azure vnet](/docs/img/kickstart-azure-vnet.png)

    * `A`: Whether to create a new
      [Virtual Network in Azure](https://docs.microsoft.com/en-us/azure/virtual-network/virtual-networks-overview)
      or use an existing one.
      If using an existing one, you must provide its VNET name. For initial
      deployments, it is recomended to create a new Virtual Network. This will
      ensure the installer takes care of all networking creation and configuration.
    * `B`: The Resource Group under which to create the VNET.
    * `C`: The name to use when creating a new VNET.
    * `D`: The CIDR block to use for this VNET.
    * `E`: The name for the control plane subnet.
    * `F`: The CIDR block to use for the control plane subnet.
    * `G`: Whether to deploy without a publicly accessible IP address.
      Access to the cluster will be limited to your Azure private network only.
      Various ways for connecting to your private cluster
      [can be found in the Azure private cluster
      documentation](https://docs.microsoft.com/en-us/azure/aks/private-clusters#options-for-connecting-to-the-private-cluster).

1. Fill out the Standalone Cluster Settings.

    ![kickstart azure standalone cluster settings](/docs/img/kickstart-azure-sa-cluster.png)

    * `A`: Choose between Development profile with one control plane node, or
      Production, which features a highly-available three node control plane.
      Additionally, choose the instance type you'd like to use for control plane nodes.
    * `B`: Name the cluster. This is a friendly name that will be used to
      reference your cluster in the Tanzu CLI and `kubectl`.
    * `C`: Whether to enable [Cluster API's machine health
      checks](https://cluster-api.sigs.k8s.io/tasks/healthcheck.html).
    * `D`: Choose whether you'd like to enable [Kubernetes API server
      auditing](https://kubernetes.io/docs/tasks/debug-application-cluster/audit/).

1. If you would like additional metadata to be tagged in your soon-to-be-created
   Azure infrastructure, fill out the Metadata section.

1. Fill out the Kubernetes Network section.

    ![kickstart kubernetes networking](/docs/img/kickstart-azure-network.png)

    * `A`: Set the CIDR for Kubernetes [Services (Cluster
      IPs)](https://kubernetes.io/docs/concepts/services-networking/service/).
These are internal IPs that, by default, are only exposed and routable within
Kubernetes.
    * `B`: Set the CIDR range for Kubernetes Pods. These are internal IPs that, by
      default, are only exposed and routable within Kubernetes.
    * `C`: Set a network proxy that internal traffic should egress through to
      access external network(s).

1. Fill out the Identity Management section.

    ![kickstart identity management](/docs/img/kickstart-identity.png)

    * `A`: Select whether you want to enable identity management. If this is
      off, certificates (via kubeconfig) are used to authenticate users. For
      most development scenarios, it is preferred to keep this off.
    * `B`: If identity management is on, choose whether to authenticate using
      [OIDC](https://openid.net/connect/) or [LDAPS](https://ldap.com).
    * `C`: Fill out connection details for identity management.

1. Fill out the OS Image section.

    ![kickstart azure os](/docs/img/kickstart-azure-os.png)

    * `A`: The Azure image to use for Kubernetes host VMs. This list should
      populate based on known images uploaded by VMware.
      These images are publicly accessible for your use. Choose
      based on your preferred Linux distribution.

1. Skip the TMC Registration section.

1. Click the Review Configuration button.

    > For your record, the configuration settings have been saved to
    > `${HOME}/.config/tanzu/tkg/clusterconfigs`.

1. Deploy the cluster.

    > If you experience issues deploying your cluster, visit the [Troubleshooting
    > documentation](../tsg-bootstrap).

1. Set your kubectl context to the cluster.

    ```sh
    kubectl config use-context <STANDALONE-CLUSTER-NAME>-admin@<STANDALONE-CLUSTER-NAME>
    ```

1. Validate you can access the cluster's API server.

    ```sh
    kubectl get nodes
    ```

    The output will look similar to the following:

    ```sh
    NAME                                       STATUS   ROLES                  AGE    VERSION
    ip-10-0-1-133.us-west-2.compute.internal   Ready    <none>                 123m   v1.20.1+vmware.2
    ip-10-0-1-76.us-west-2.compute.internal    Ready    control-plane,master   125m   v1.20.1+vmware.2
    ```
