## Create vSphere Clusters

This section describes setting up management and workload clusters for
vSphere.

1. Download the machine image that matches the version of the Kubernetes you plan on deploying.

    At this time, we cannot guarantee the plugin versions that will be used for cluster management.  While running the installer interface to bootstrap a cluster, you are required to add an `OVA` to your vSphere environment.

    The official OVA publishing location is still to be determined. In the meantime, to get access to the necessary OVAs for the current build, please ask on the `#tanzu-community-edition` Slack channel.

    Please note, validation work so far has focused on the **Photon** based
    images.

1. In vCenter, right-click on your datacenter and deploy the OVF template.

1. After importing, right-click and covert to a template.

1. Initialize the Tanzu Community Edition installer interface.

    ```sh
    tanzu management-cluster create --ui
    ```

1. Complete the configuration steps in the installer interface for vSphere and create the management cluster. The following configuration settings are recommended:

   * If you do not specify a name, the installer generates a unique name. If you do specify a name, the name must end with a letter, not a numeric character, and must be compliant with DNS hostname requirements described here: [RFC 1123](https://tools.ietf.org/html/rfc1123).
   * Set all instance profiles to large or larger. In our testing, we found resource constraints caused bootstrapping issues. Choosing a large profile or more will give a better chance for
     successful bootstrapping.
   * Set your control plane IP. The control plane IP is a virtual IP that fronts the Kubernetes API
     server. You **must** set an IP that is routable and won't be taken by
     another system (e.g. DHCP).
   * Disable **Enable Identity Management Settings**. You can disable identity management for proof-of-concept/development deployments, but it is strongly recommended to implement identity management in production deployments. For more information about enabling Identity Management, see [Identity Management ](vsphere-install-mgmt/#step-7-identity-management).

1. Validate the management cluster started successfully.

    ```sh
    tanzu management-cluster get
    ```

1. Capture the management cluster's kubeconfig.

    ```sh
    tanzu management-cluster kubeconfig get <MGMT-CLUSTER-NAME> --admin

    Where <``<MGMT-CLUSTER-NAME>`` should be set to the name returned by `tanzu management-cluster get` above.  <br><br>
    For example, if your management cluster is called 'mtce', you will see a message similar to:
    ```sh
    Credentials of workload cluster 'mtce' have been saved.
    You can now access the cluster by running 'kubectl config use-context mtce-admin@mtce'

1. Set your kubectl context to the management cluster.

    ```sh
    kubectl config use-context <MGMT-CLUSTER-NAME>-admin@<MGMT-CLUSTER-NAME>
    ```
    Where <``MGMT-CLUSTER-NAME>`` should be set to the name returned by `tanzu management-cluster get`.
1. Validate you can access the management cluster's API server.

    ```sh
    kubectl get nodes

    NAME         STATUS   ROLES                  AGE    VERSION
    10-0-1-133   Ready    <none>                 123m   v1.20.1+vmware.2
    10-0-1-76    Ready    control-plane,master   125m   v1.20.1+vmware.2
    ```

1. Next you will create a workload cluster. First, create a workload cluster configuration file by taking a copy of the management cluster YAML configuration file that was created when you deployed your management cluster. This example names the workload cluster configuration file ``workload1.yaml``.

    ```sh
    cp  ~/.config/tanzu/tkg/clusterconfigs/<MGMT-CONFIG-FILE> ~/.config/tanzu/tkg/clusterconfigs/workload1.yaml
    ```

    * Where ``<MGMT-CONFIG-FILE>`` is the name of the management cluster YAML config file. The management cluster YAML configuration file will either have the name you assigned to the management cluster, or if no name was assigned, it will be a randomly generated name.

    * The duplicated file (``workload1.yaml``) will be used as the configuration file for your workload cluster. You can edit the parameters in this new  file as required. For an example of a workload cluster template, see  [vSphere Workload Cluster Template](../vsphere-wl-template).

      * In the next two steps you will edit the parameters in this new file (`workload1.yaml`) and then use the file to deploy a workload cluster.

1. In the new workload cluster file (`~/.config/tanzu/tkg/clusterconfigs/workload.yaml`), edit the CLUSTER_NAME parameter to assign a name to your workload cluster. For example,


   ```yaml
   CLUSTER_CIDR: 100.96.0.0/11
   CLUSTER_NAME: my-workload-cluster
   CLUSTER_PLAN: dev
   ```
   * If you did not specify a name for your management cluster, the installer generated a random unique name. In this case, you must manually add the CLUSTER_NAME parameter and assign a workload cluster name. The workload cluster names must be must be 42 characters or less and must comply with DNS hostname requirements as described here: [RFC 1123](https://tools.ietf.org/html/rfc1123)
   * If you specified a name for your management cluster, the CLUSTER_NAME parameter is present and needs to be changed to the new workload cluster name.

1. In the workload cluster file (`~/.config/tanzu/tkg/clusterconfigs/workload1.yaml`), edit the VSPHERE_CONTROL_PLANE_ENDPOINT parameter to apply a viable IP.


   * This will be the API Server IP for your workload cluster. You must choose an IP that is routable and not used elsewhere in your network, e.g., out of your DHCP range.


   * The other parameters in ``workload1.yaml`` are likely fine as-is. Validation is performed on the file prior to applying it, so the `tanzu` command will return a message if something necessary is omitted. However, you can parameters as required. Reference an example configuration template here:  [vSphere Workload Cluster Template](../vsphere-wl-template).


4. Create your workload cluster.

    ```sh
    tanzu cluster create <WORKLOAD-CLUSTER-NAME> --file ~/.config/tanzu/tkg/clusterconfigs/workload1.yaml
    ```

5. Validate the cluster starts successfully.

    ```sh
    tanzu cluster list
    ```

6. Capture the workload cluster's kubeconfig.

    ```sh
    tanzu cluster kubeconfig get <WORKLOAD-CLUSTER-NAME> --admin
    ```

7. Set your `kubectl` context accordingly.

    ```sh
    kubectl config use-context <WORKLOAD-CLUSTER-NAME>-admin@<WORKLOAD-CLUSTER-NAME>
    ```

8. Verify you can see pods in the cluster.

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
