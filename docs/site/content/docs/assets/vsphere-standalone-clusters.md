## Create vSphere Clusters

This section describes setting up standalone clusters for
vSphere. These clusters are not managed by a management cluster.

1. Download the machine image that matches the version of the Kubernetes you plan on deploying (1.20.1 is default).

    At this time, we cannot guarantee the plugin versions that will be
    used for cluster management. While using the kickstart UI to bootstrap your
    cluster, you may be asked add an `ova` to your vSphere environment. The
    following links are points to the most recent ovas at the time of writing
    this Getting Started guide.

    The official OVA publishing location is still to be determined. In the meantime, to get access to the necessary OVAs
    for the current build, please ask on the `#tanzu-community-edition` Slack channel.

    Please note, validation work so far has focused on the **Photon** based
    images.

1. In vCenter, right-click on your datacenter and import the OVF template.

1. After importing, right-click and covert to a template.

1. Initialize the Tanzu kickstart UI.

    ```sh
    tanzu standalone-cluster create --ui
    ```

1. Go through the configuration steps, considering the following.

   * Set all instance profile to large or larger.
     * In our testing, we found resource constraints caused bootstrapping
     issues. Choosing a large profile or more will give a better chance for
     successful bootstrapping.
   * Set your control plane IP
     * The control plane IP is a virtual IP that fronts the Kubernetes API
     server. You **must** set an IP that is routable and won't be taken by
     another system (e.g. DHCP).
   * Disable **Enable Identity Management Settings**. You can disable identity management for proof-of-concept/development deployments, but it is strongly recommended to implement identity management in production deployments.

1. At the end of the UI, create the standalone cluster.

1. Store the name of your cluster (set during configuration or generated) to a
   `WORKLOAD_CLUSTER_NAME` environment variable.

    ```sh
    export WORKLOAD_CLUSTER_NAME="<INSERT_WORKLOAD_CLUSTER_NAME_HERE>"
    ```
1. Set your kubectl context to the cluster.

    ```sh
    kubectl config use-context ${WORKLOAD_CLUSTER_NAME}-admin@${WORKLOAD_CLUSTER_NAME}
    ```

1. Validate you can access the cluster's API server.

    ```sh
    kubectl get nodes

    NAME         STATUS   ROLES                  AGE    VERSION
    10-0-1-133   Ready    <none>                 123m   v1.20.1+vmware.2
    10-0-1-76    Ready    control-plane,master   125m   v1.20.1+vmware.2
    ```
