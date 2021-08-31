## Create Standalone Amazon EC2 Clusters

This section covers setting up a standalone cluster in Amazon EC2. A standalone cluster provides a workload cluster that is **not** managed by a centralized management cluster.

1. Initialize the Tanzu Community Edition Installer Interface.

    ```sh
    tanzu standalone-cluster create --ui
    ```

1. Complete the configuration steps in the installer interface and create the standalone cluster. The following configuration settings are recommended:


   * Check the "Automate creation of AWS CloudFormation Stack" box if you do not have an existing TKG CloudFormation stack. This stack is used to created IAM resources that Tanzu Community Edition clusters use in Amazon EC2.
     You only need 1 TKG CloudFormation stack per AWS account. CloudFormation is global and not locked to a region.

   * Set the instance type size to m5.xlarger or larger for the control plane node.

   * Disable **Enable Identity Management Settings**. You can disable identity management for proof-of-concept/development deployments, but it is strongly recommended to implement identity management in production deployments.

1. Set your kubectl context to the cluster.

    ```sh
    kubectl config use-context <STANDALONE-CLUSTER-NAME>-admin@<STANDALONE-CLUSTER-NAME>
    ```
    Where `<STANDALONE-CLUSTER-NAME>` is the name of the standalone cluster that you specified or if you didn't specify a name, it's the randomly generated name.
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
