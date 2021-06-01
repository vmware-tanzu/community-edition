## Create vSphere Clusters

This section describes setting up standalone clusters for
vSphere. These clusters are not managed by a management cluster.

1. Download the machine image that matches the version of the Kubernetes you plan on deploying (1.20.1 is default).

    At this time, we cannot guarantee the exact plugin versions that will be
    used for cluster management. While using the kickstart UI to bootstrap your
    cluster, you may be asked add an `ova` to your vSphere environment. The
    following links are points to the most recent ovas at the time of writing
    this getting started guide.

    * [1.20.4
      OVA](http://build-squid.eng.vmware.com/build/mts/release/bora-17759077/publish/lin64/tkg_release/node/ova-photon-3-v1.20.4+vmware.1-tkg.0-2326554155028348692/photon-3-kube-v1.20.4+vmware.1-tkg.0-2326554155028348692.ova)
    * [1.19.8
      OVA](http://build-squid.eng.vmware.com/build/mts/release/bora-17759077/publish/lin64/tkg_release/node/ova-photon-3-v1.19.8+vmware.1-tkg.0-15338136437231643652/photon-3-kube-v1.19.8+vmware.1-tkg.0-15338136437231643652.ova)

    If you're asked for another `ova` version by the kickstart UI, you can
    download the ova that corresponds to the rc version (e.g. 1,2,3,etc) at the [TKG
    daily builds confluence
    page](https://confluence.eng.vmware.com/pages/viewpage.action?spaceKey=TKG&title=TKG+Release+Daily+Build#TKGReleaseDailyBuild-TKG1.3.0RC.3(March/09/2021)).

    Please note, validation work so far has focused on the **Photon** based
    images.

1. In vCenter, right click on your datacenter and import OVF template.

1. After importing, right-click and covert to a template.

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

   * Set all instance profile to large or larger.
     * In our testing, we found resource constraints caused bootstrapping
     issues. Choosing a large profile or more will give a better chance for
     successful bootstrapping.
   * Set your control plane IP
     * The control plane IP is a virtual IP that fronts the Kubernetes API
     server. You **must** set an IP that is routable and won't be taken by
     another system (e.g. DHCP).
   * Disable OIDC configuration.

    > Until we have more TCE documentation, you can find the full TKG docs
    > [here](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.2/vmware-tanzu-kubernetes-grid-12/GUID-mgmt-clusters-deploy-management-clusters.html).
    > We will have more complete `tanzu` cluster bootstrapping documentation available here in the near future.

1. At the end of the UI, create the standalone cluster.

1. Set your kubectl context to the cluster.

    ```sh
    kubectl config use-context ${GUEST_CLUSTER_NAME}-admin@${GUEST_CLUSTER_NAME}
    ```

1. Validate you can access the cluster's API server.

    ```sh
    kubectl get nodes

    NAME         STATUS   ROLES                  AGE    VERSION
    10-0-1-133   Ready    <none>                 123m   v1.20.1+vmware.2
    10-0-1-76    Ready    control-plane,master   125m   v1.20.1+vmware.2
    ```
