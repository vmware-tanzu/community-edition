## Create Standalone AWS Clusters

This section covers setting up a standalone cluster in AWS. This provides you
a workload cluster that is **not** managed by a centralized management cluster.

1. Store the name of your cluster (set in the configuration file) to a
   `GUEST_CLUSTER_NAME` environment variable.

    ```sh
    export GUEST_CLUSTER_NAME="<INSERT_GUEST_CLUSTER_NAME_HERE>"
    ```

1. Initialize the Tanzu kickstart UI.

    ```sh
    tanzu standalone-cluster create ${GUEST_CLUSTER_NAME} --ui
    ```

1. Go through the configuration steps, considering the following.

   * Check the "Automate creation of AWS CloudFormation Stack" box if you do not have an existing TKG CloudFormation stack. This stack is used to created IAM resources that TCE clusters use in Amazon EC2.
      * You only need 1 TKG CloudFormation stack per AWS account. CloudFormation is global and not locked to a region.
   * Set all instance sizes to m5.xlarge or larger.
   * Disable OIDC configuration.

    > Until we have more TCE documentation, you can find the full TKG docs
    > [here](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.2/vmware-tanzu-kubernetes-grid-12/GUID-mgmt-clusters-deploy-management-clusters.html).
    > We will have more complete `tanzu` cluster bootstrapping documentation available here in the near future.

1. At the end of the UI, deploy the cluster.

1. Set your kubectl context to the cluster.

    ```sh
    kubectl config use-context ${GUEST_CLUSTER_NAME}-admin@${GUEST_CLUSTER_NAME}
    ```

1. Validate you can access the cluster's API server.

    ```sh
    kubectl get nodes

    NAME                                       STATUS   ROLES                  AGE    VERSION
    ip-10-0-1-133.us-west-2.compute.internal   Ready    <none>                 123m   v1.20.1+vmware.2
    ip-10-0-1-76.us-west-2.compute.internal    Ready    control-plane,master   125m   v1.20.1+vmware.2
    ```
