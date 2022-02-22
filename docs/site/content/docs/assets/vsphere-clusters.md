## Create vSphere Clusters

This section describes setting up a management and workload cluster on
vSphere.

1. Open the [Tanzu Community Edition product page
   on](https://customerconnect.vmware.com/downloads/get-download?downloadGroup=TCE-090)
VMware Customer Connect.

    > If you do not have a Customer Connect account, [register
    > here](https://customerconnect.vmware.com/account-registration).

1. Ensure you have the version selected corresponding to your installation.

    ![customer connect download page](/docs/img/customer-connect-downloads.png)

1. Locate and download the machine image (OVA) for your desired operating system
   and Kubernetes version.

    ![customer connect ova downloads](/docs/img/customer-connect-ovas.png)

1. Log in to your vCenter instance.

1. In vCenter, right-click on your datacenter and choose Deploy OVF Template.

    ![vcenter deploy ovf](/docs/img/vcenter-deploy-ovf.png)

1. Follow the prompts, browsing to the local file that is the `.ova` downloaded
   in a previous step.

1. Allow the template deployment to complete.

    ![vcenter deploy ovf](/docs/img/vcenter-import-ovf.png)

1. Right-click on the newly imported OVF template and choose Template > Convert to Template.

    ![vcenter convert to template](/docs/img/vcenter-convert-to-template.png)

1. Verify the template is added by selecting the VMs and Templates icon and
   locating it within your datacenter.

    ![vcenter template import](/docs/img/vcenter-template-import.png)

1. Initialize the Tanzu Community Edition installer interface.

    ```sh
    tanzu management-cluster create --ui
    ```

    **Note**: If you are bootstrapping from a Windows machine and encounter an `unable to ensure prerequisites` error, see the following[troubleshooting topic](../faq-cluster-bootstrapping/#x509-certificate-signed-by-unknown-authority-when-deploying-management-cluster-from-windows).

1. Choose VMware vSphere from the provider tiles.

    ![kickstart vsphere tile](/docs/img/kickstart-vsphere-tile.png)

1. Fill out the IaaS Provider section.

    ![kickstart vsphere iaas](/docs/img/kickstart-vsphere-iaas.png)

    * `A`: The IP or DNS name pointing at your vCenter instance. This is the
      same instance you uploaded the OVA to in previous steps.
    * `B`: The username, with elevated privileges, that can be used to create
      infrastructure in vSphere.
    * `C`: The password corresponding to that username.
    * `D`: With the above filled out, connect to the instance to continue. You
      may be prompted to verify the SSL fingerprint for vCenter.
    * `E`: The datacenter you'll deploy Tanzu Community Edition into. This
      should be the same datacenter you uploaded the OVA to.
    * `F`: The public key you'd like to use for your VM instances. This is how
      you'll SSH into control plane and worker nodes.

1. Fill out the Management Cluster Settings.

    ![kickstart vsphere management cluster settings](/docs/img/kickstart-vsphere-mgmt-cluster.png)

    * `A`: Choose between Development profile, with 1 control plane node or
      Production, which features a highly-available three node control plane.
      Additionally, choose the instance type you'd like to use for control plane nodes.
    * `B`: Name the cluster. This is a friendly name that will be used to
      reference your cluster in the Tanzu CLI and `kubectl`.
    * `C`: Choose whether to enable [Cluster API's machine health
      checks](https://cluster-api.sigs.k8s.io/tasks/healthcheck.html).
    * `D`: Choose how to expose your control plane endpoint. If you have NSX
      available and would like to use it, choose NSX Advanced Load Balancer,
      otherwise choose [Kube-vip](https://kube-vip.io), which will expose a virtual IP in your network.
    * `E`: Set the IP address your Kubernetes API server should be accessible from. This
      should be an IP that **is routable** in your network but **excluded from
      your DHCP range**.
    * `F`: Set the instance type you'd like to use for workload nodes.
    * `G`: Choose whether you'd like to enable [Kubernetes API server
      auditing](https://kubernetes.io/docs/tasks/debug-application-cluster/audit/).

1. If you choose NSX as your Control Plane Endpoint Provider in the above step,
   fill out the VMware NSX Advanced Load Balancer section.

1. If you would like additional metadata to be tagged in your soon-to-be-created
   vSphere infrastructure, fill out the Metadata section.

1. Fill out the Resources section.

    ![kickstart vsphere resources](/docs/img/kickstart-vsphere-resources.png)

    * `A`: Set the [VM
      folder](https://docs.vmware.com/en/VMware-Workstation-Pro/16.0/com.vmware.ws.using.doc/GUID-016FF81D-4FE4-4D9E-92D6-A08E022AA6D4.html)
      you'd like new virtual machines to be created in. By default, this will be
      `${DATACENTER_NAME}/vm/`
    * `B`: Set the
      [Datastore](https://docs.vmware.com/en/VMware-vSphere/7.0/com.vmware.vsphere.storage.doc/GUID-D5AB2BAD-C69A-4B8D-B468-25D86B8D39CE.html) you'd like volumes to be created within.
    * `C`: Set the servers or resource pools within the data center you'd like
      VMs, networking, etc to be created in.

1. Fill out the Kubernetes Network section.

    ![kickstart kubernetes networking](/docs/img/kickstart-network.png)

    * `A`: Select the [vSphere
      network](https://docs.vmware.com/en/VMware-vSphere/7.0/com.vmware.vsphere.networking.doc/GUID-35B40B0B-0C13-43B2-BC85-18C9C91BE2D4.html)
where **host**/**virtual machine** networking will be setup in.
    * `B`: Set the CIDR for Kubernetes [Services (Cluster
      IPs)](https://kubernetes.io/docs/concepts/services-networking/service/).
These are internal IPs that, by default, are only exposed and routable within
Kubernetes.
    * `C`: Set the CIDR range for Kubernetes Pods. These are internal IPs that, by
      default, are only exposed and routable within Kubernetes.
    * `D`: Set a network proxy that internal traffic should egress through to
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

    ![kickstart vsphere os](/docs/img/kickstart-vsphere-os.png)

    * `A`: The OVA image to use for Kubernetes host VMs. This list should
      populate based on the OVA you uploaded in previous steps. If it's missing,
      you may have uploaded an incompatible OVA.

1. Skip the TMC Registration section.

1. Click the Review Configuration button.

    > For your record, the configuration settings have been saved to
    > `${HOME}/.config/tanzu/tkg/clusterconfigs`.

1. Deploy the cluster.

    > If you experience issues deploying your cluster, visit the [Troubleshooting
    > documentation](../tsg-bootstrap).

1. Validate the management cluster started successfully.

   ```sh
   tanzu management-cluster get
   ```

1. Capture the management cluster's kubeconfig.

   ```sh
   tanzu management-cluster kubeconfig get <MGMT-CLUSTER-NAME> --admin
   ```

   Where <`<MGMT-CLUSTER-NAME>` should be set to the name returned by `tanzu management-cluster get` above.  <br><br>
   For example, if your management cluster is called 'mtce', you will see a message similar to:

   ```sh
   Credentials of cluster 'mtce' have been saved.
   You can now access the cluster by running 'kubectl config use-context mtce-admin@mtce'
   ```

1. Set your kubectl context to the management cluster.

   ```sh
   kubectl config use-context <MGMT-CLUSTER-NAME>-admin@<MGMT-CLUSTER-NAME>
   ```

   Where <`MGMT-CLUSTER-NAME>` should be set to the name returned by `tanzu management-cluster get`.
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

1. Create your workload cluster.

    ```sh
    tanzu cluster create <WORKLOAD-CLUSTER-NAME> --file ~/.config/tanzu/tkg/clusterconfigs/workload1.yaml
    ```

1. Validate the cluster starts successfully.

    ```sh
    tanzu cluster list
    ```

1. Capture the workload cluster's kubeconfig.

    ```sh
    tanzu cluster kubeconfig get <WORKLOAD-CLUSTER-NAME> --admin
    ```

1. Set your `kubectl` context accordingly.

    ```sh
    kubectl config use-context <WORKLOAD-CLUSTER-NAME>-admin@<WORKLOAD-CLUSTER-NAME>
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
