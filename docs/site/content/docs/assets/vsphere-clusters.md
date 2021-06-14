## Create vSphere Clusters

This section describes setting up management and workload/guest clusters for
vSphere. If your deployment target is AWS, skip this section and move on to the
next.

1. Download the machine image that matches the version of the Kubernetes you plan on deploying (1.20.1 is default).

    At this time, we cannot guarantee the plugin versions that will be
    used for cluster management.   
    While using the kickstart UI to bootstrap your cluster, you may be asked to add an `OVA` to your vSphere environment. The following links are to the most recent OVAs at the time of writing this Getting Started guide.   
    To access the OVAs, you must have a VMware Customer Connect account. Complete the following steps to register a new account and access the OVAs:  
        a. If you don't already have an account, register a new account on [VMware Customer Connect](http://my.vmware.com/).  
        b. Log in to VMware Customer Connect, click on Products and Accounts > All Products.   
        c. Search for "tanzu kubernetes grid", and from the search results, select "Tanzu Kubernetes Grid > Product Binaries > Tanzu Kubernetes Grid" to access the OVAs.  

    * [1.20.4
      OVA](http://build-squid.eng.vmware.com/build/mts/release/bora-17800251/publish/lin64/tkg_release/node/ova-photon-3-v1.20.4+vmware.1-tkg.0-2326554155028348692/photon-3-kube-v1.20.4+vmware.1-tkg.0-2326554155028348692.ova)
    * [1.19.8
      OVA](http://build-squid.eng.vmware.com/build/mts/release/bora-17759077/publish/lin64/tkg_release/node/ova-photon-3-v1.19.8+vmware.1-tkg.0-15338136437231643652/photon-3-kube-v1.19.8+vmware.1-tkg.0-15338136437231643652.ova)

    <!--If you're asked for another `OVA` version by the kickstart UI, you can
    download the OVA that corresponds to the rc version (e.g. 1,2,3,etc) at the [TKG
    daily builds confluence page](https://confluence.eng.vmware.com/pages/viewpage.action?spaceKey=TKG&title=TKG+Release+Daily+Build#TKGReleaseDailyBuild-TKG1.3.0RC.3(March/09/2021)).-->

    Please note, validation work so far has focused on the **Photon** based
    images.

1. In vCenter, right click on your datacenter and import OVF template.

1. After importing, right-click and covert to a template.

1. Initialize the Tanzu kickstart UI.

    ```sh
    tanzu management-cluster create --ui
    ```

1. Go through the installation process for vSphere. With the following
   considerations:

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

1. Validate the management cluster started successfully.

    ```sh
    tanzu management-cluster get
    ```

1. Create a cluster name that will be used throughout this Getting Started guide. This instance of `MGMT_CLUSTER_NAME` should be set to whatever value is returned by `tanzu management-cluster get` above.

    ```sh
    export MGMT_CLUSTER_NAME="<INSERT_MGMT_CLUSTER_NAME_HERE>"
    export GUEST_CLUSTER_NAME="<INSERT_GUEST_CLUSTER_NAME_HERE>"
    ```

1. Capture the management cluster's kubeconfig.

    ```sh
    tanzu management-cluster kubeconfig get ${MGMT_CLUSTER_NAME} --admin

    Credentials of workload cluster 'mtce' have been saved
    You can now access the cluster by running 'kubectl config use-context mtce-admin@mtce'
    ```

    > Note the context name `${MGMT_CLUSTER_NAME}-admin@${MGMT_CLUSTER_NAME}`, you'll use the above command in
    > future steps.

1. Set your kubectl context to the management cluster.

    ```sh
    kubectl config use-context ${MGMT_CLUSTER_NAME}-admin@${MGMT_CLUSTER_NAME}
    ```

1. Validate you can access the management cluster's API server.

    ```sh
    kubectl get nodes

    NAME         STATUS   ROLES                  AGE    VERSION
    10-0-1-133   Ready    <none>                 123m   v1.20.1+vmware.2
    10-0-1-76    Ready    control-plane,master   125m   v1.20.1+vmware.2
    ```

1. Setup a guest cluster config file.

    ```sh
    cp  ~/.tanzu/tkg/clusterconfigs/xw6nt8jduy.yaml ~/.tanzu/tkg/clusterconfigs/guest1.yaml
    ```

   > This takes the configuration used to create your management cluster and
   > duplicates for use in the guest cluster. You can edit values in this new
   > file `guest1` as you please.

   [](ignored)

   > Creation of guest clusters now require the use of workload cluster YAML
   > configuration files.  [Example configuration templates](https://gitlab.eng.vmware.com/TKG/tkg-cli-providers/-/tree/cluster-templates/docs/cluster-templates)
   > are available to help get you started. Review settings and populate fields
   > that are not set.

   [](ignored)

   > Validation is performed on the file prior to applying it, so the `tanzu`
   > command should give you any clues if something necessary is omitted.

1. Assign a name to your guest cluster in the `~/.tanzu/tkg/clusterconfigs/guest1.yaml` file. For example, 

   ```yaml
   CLUSTER_CIDR: 100.96.0.0/11
   CLUSTER_NAME: my-guest-cluster
   CLUSTER_PLAN: dev
   ```
   #### Note
   * If you did not specify a name for your management cluster, the installer generated a unique name, in this case, you must manually add the CLUSTER_NAME parameter and assign a guest cluster name. 
   * If you specified a name for your management cluster, the CLUSTER_NAME parameter is present and needs to be changed to the new guest cluster name.

1. Edit the guest cluster config file's
   (`~/.tanzu/tkg/clusterconfigs/guest1.yaml`) VSPHERE_CONTROL_PLANE_ENDPOINT to
   a viable IP.

   > This will be **the API Server IP** for you guest cluster. You must choose
   > an IP that is **1.) routable** and **2.) not used elsewhere in your network
   > (eg. out of your DHCP range)**.

   [](ignored)

   > For vSphere, the other settings are likely fine as-is. However, you can change
   > them as you'd like and/or reference the [Example configuration templates](https://gitlab.eng.vmware.com/TKG/tkg-cli-providers/-/tree/cluster-templates/docs/cluster-templates).

1. Create your guest cluster.

    ```sh
    tanzu cluster create ${GUEST_CLUSTER_NAME} --file ${HOME}/.tanzu/tkg/clusterconfigs/guest1.yaml
    ```

1. Validate the cluster starts successfully.

    ```sh
    tanzu cluster list
    ```

1. Capture the guest cluster's kubeconfig.

    ```sh
    tanzu cluster kubeconfig get ${GUEST_CLUSTER_NAME} --admin
    ```

1. Set your `kubectl` context accordingly.

    ```sh
    kubectl config use-context ${GUEST_CLUSTER_NAME}-admin@${GUEST_CLUSTER_NAME}
    ```

1. Verify you can see pods in the cluster.

    ```sh
    kubectl get pods --all-namespaces

    NAMESPACE     NAME                                                    READY   STATUS    RESTARTS   AGE
    kube-system   antrea-agent-9d4db                                      2/2     Running   0          3m42s
    kube-system   antrea-agent-vkgt4                                      2/2     Running   1          5m48s
    kube-system   antrea-controller-5d594c5cc7-vn5gt                      1/1     Running   0          5m49s
    kube-system   coredns-5d6f7c958-hs6vr                                 1/1     Running   0          5m49s
    kube-system   coredns-5d6f7c958-xf6cl                                 1/1     Running   0          5m49s
    kube-system   etcd-tce-guest-control-plane-b2wsf                      1/1     Running   0          5m56s
    kube-system   kube-apiserver-tce-guest-control-plane-b2wsf            1/1     Running   0          5m56s
    kube-system   kube-controller-manager-tce-guest-control-plane-b2wsf   1/1     Running   0          5m56s
    kube-system   kube-proxy-9825q                                        1/1     Running   0          5m48s
    kube-system   kube-proxy-wfktm                                        1/1     Running   0          3m42s
    kube-system   kube-scheduler-tce-guest-control-plane-b2wsf            1/1     Running   0          5m56s
    kube-system   kube-vip-tce-guest-control-plane-b2wsf                  1/1     Running   0          5m56s
    kube-system   vsphere-cloud-controller-manager-nwrg4                  1/1     Running   2          5m48s
    kube-system   vsphere-csi-controller-5b6f54ccc5-trgm4                 5/5     Running   0          5m49s
    kube-system   vsphere-csi-node-drnkz                                  3/3     Running   0          5m48s
    kube-system   vsphere-csi-node-flszf                                  3/3     Running   0          3m42s
    ```
