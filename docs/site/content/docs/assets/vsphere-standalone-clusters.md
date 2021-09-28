## Create vSphere Clusters

This section describes setting up standalone clusters on vSphere.

1. Download the OVA for the management cluster nodes directly from [VMware Customer Connect](https://customerconnect.vmware.com/downloads/
get-download?downloadGroup=TCE-090).  
Alternatively, you can open the [Tanzu Community Edition product page](https://customerconnect.vmware.com/downloads/info/slug/
infrastructure_operations_management/
vmware_tanzu_community_edition/0_9_0) in Customer Connect and select and download the OVA version that you require. You will need a VMware Customer
Connect account to download the OVA, register [here](https://customerconnect.vmware.com/account-registration).
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
   tanzu standalone-cluster create --ui
   ```

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

    ![kickstart vsphere management cluster settings](/docs/img/kickstart-vsphere-sa-cluster.png)

    * `A`: Choose between Development profile, with 1 control plane node or
      Production, which features a highly-available three node control plane.
      Additionally, choose the instance type you'd like to use for control plane and
    workload nodes.
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
    * `F`: Choose whether you'd like to enable [Kubernetes API server
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
    * `D`: Setup a proxy that cluster traffic should egress through to access
      extrenal network(s).

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

1. Click the Review Configuration button.

    > For your record, the configuration settings have been saved to
    > `${HOME}/.config/tanzu/tkg/clusterconfigs`

1. Deploy the cluster.

    > If you experience issues deploying your cluster, visit the [Troubleshooting
    > documentation](../tsg-bootstrap).

1. Once complete, Set your kubectl context to the cluster.

   ```sh
   kubectl config use-context <STANDALONE-CLUSTER-NAME>-admin@<STANDALONE-CLUSTER-NAME>
   ```

1. Validate you can access the cluster's API server.

   ```sh
   kubectl get nodes

   NAME         STATUS   ROLES                  AGE    VERSION
   10-0-1-133   Ready    <none>                 123m   v1.20.1+vmware.2
   10-0-1-76    Ready    control-plane,master   125m   v1.20.1+vmware.2
   ```
