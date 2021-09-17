# Reference for vSphere account

If you encounter issues deploying a cluster to vSphere, review the following troubleshooting and reference content:

[Required Permissions for the vSphere Account](#permissions)
[Static VIPs and Load Balancers for vSphere](#static-ip)
[Obtain vSphere Certificate Thumbprints](#certificates)
[Thick provisioning versus thin provisioning](#provisioning)
## </a> Required Permissions for the vSphere Account {#permissions}

The vCenter Single Sign On account that you provide to Tanzu Comunity Edition when you deploy a management cluster must have the correct permissions in order to perform the required operations in vSphere.

It is not recommended to provide a vSphere administrator account to Tanzu Community Edition, because this provides Tanzu Community Edition with far greater permissions than it needs. The best way to assign permissions to Tanzu Community Edition is to create a role and a user account, and then to grant that user account that role on vSphere objects.

1. In the vSphere Client, go to **Administration** > **Access Control** > **Roles**, and create a new role, for example `TCE`, with the following permissions:

  <table width="100%" border="1" align="left">
   <tr>
    <th scope="col" align="left">vSphere Object </th>
    <th scope="col" align="left">Required Permission </th>
   </tr>
   <tr>
    <td>Cns</td>
    <td>Searchable</td>
   </tr>
   <tr>
    <td>Datastore</td>
    <td>Allocate space<br />
    Browse datastore<br />
	Low level file operations</td>
   </tr>
   <tr>
    <td>Global (if using Velero for backup and restore)</td>
    <td>Disable methods<br />
    Enable methods<br />
	Licenses</td>
   </tr>
   <tr>
    <td>Network</td>
    <td>Assign network</td>
   </tr>
   <tr>
    <td>Profile-driven storage</td>
    <td>Profile-driven storage view</td>
   </tr>
   <tr>
    <td>Resource</td>
    <td>Assign virtual machine to resource pool</td>
   </tr>
   <tr>
     <td>Sessions</td>
     <td>Message<br />
     Validate session</td>
   </tr>
   <tr>
    <td>Virtual machine</td>
    <td>
      Change Configuration &gt; Add existing disk<br />
      Change Configuration &gt; Add new disk<br />
      Change Configuration &gt; Add or remove device<br />
      Change Configuration &gt; Advanced configuration<br />
      Change Configuration &gt; Change CPU count<br />
      Change Configuration &gt; Change Memory<br />
      Change Configuration &gt; Change Settings<br />
      Change Configuration &gt; Configure Raw device<br />
      Change Configuration &gt; Extend virtual disk<br />
      Change Configuration &gt; Modify device settings<br />
      Change Configuration &gt; Remove disk<br />
      Change Configuration &gt; Toggle disk change tracking*<br />
      Edit Inventory &gt; Create from existing<br />
      Edit Inventory &gt; Remove<br />
      Interaction &gt; Power On<br />
      Interaction &gt; Power Off<br />
      Provisioning &gt; Allow read-only disk access*<br />
      Provisioning &gt; Allow virtual machine download*<br />
      Provisioning &gt; Deploy template<br />
      Snapshot Management &gt; Create snapshot*<br />
      Snapshot Management &gt; Remove snapshot*<br /><br />
      *Required to enable the Velero plugin, as described in <a href="../cluster-lifecycle/backup-restore-mgmt-cluster.md">Back Up and Restore Clusters</a>. You can add these permissions when needed later.
      </td>
   </tr>
   <tr>
    <td>vApp</td>
    <td>Import</td>
   </tr>
   </table>

2. In **Administration** > **Single Sign On** > **Users and Groups**, create a new user account in the appropriate domain, for example `tkg-user`.
3.  In the **Hosts and Clusters**, **VMs and Templates**, **Storage**, and **Networking** views, right-click the objects that your Tanzu Community Edition deployment will use, select **Add Permission**, and assign the `tce-user`  with the `TCE` role to each object.

- Hosts and Clusters
    - The root vCenter Server object
    - The Datacenter and all of the Host and Cluster folders, from the Datacenter object down to the cluster that manages the Tanzu Community Edition deployment
    - Target hosts and clusters
    - Target resource pools, with propagate to children enabled
- VMs and Templates
    - The deployed Tanzu Community Edition base image templates
    -  Target VM and Template folders, with propagate to children enabled
    - Storage
    - Datastores and all storage folders, from the Datacenter object down to the datastores that will be used for Tanzu Community Edition deployments
- Networking
    - Networks or distributed port groups to which clusters will be assigned
    - Distributed switches

If you intend to use Velero to back up and restore management or workload clusters, you must also set the permissions listed in [Credentials and Privileges for VMDK Access](https://code.vmware.com/docs/11750/virtual-disk-development-kit-programming-guide/GUID-8301C6CF-37C2-42CC-B4C5-BB1DD28F79C9.html) in the *Virtual Disk Development Kit Programming Guide*.
 

##  Static VIPs and Load Balancers for vSphere {#static-ip}

Each management cluster and Tanzu Kubernetes cluster that you deploy to vSphere requires one static virtual IP address for external requests to the cluster's API server. You must be able to assign this IP address, so it cannot be within your DHCP range, but it must be in the same subnet as the DHCP range.

The cluster control plane's [Kube-vip](https://kube-vip.io/) pod uses this static virtual IP address to serve API requests, and the API server certificate includes the address to enable secure TLS communication.  In Tanzu Kubernetes clusters, Kube-vip runs in a basic, Layer-2 failover mode, assigning the virtual IP address to one control plane node at a time. In this mode, Kube-vip does not function as a true load balancer for control plane traffic.

Tanzu Community Edition also does not use Kube-vip as a load balancer for workloads in workload clusters.
Kube-vip is used solely by the cluster's API server.

To load-balance workloads on vSphere, use NSX Advanced Load Balancer, also known as Avi Load Balancer, Essentials Edition.
You must deploy the NSX Advanced Load Balancer in your vSphere instance before you deploy management clusters.

## Obtain vSphere Certificate Thumbprints {#certificates}

If your vSphere environment uses untrusted, self-signed certificates to authenticate connections, you must verify the thumbprint of the vCenter Server when you deploy a management cluster. If your vSphere environment uses trusted certificates that are signed by a known Certificate Authority (CA), you do not need to verify the thumbprint.

You can use either SSH and OpenSSL or the Platform Services Controller to obtain certificate thumbprints.

### vCenter Server Appliance

You can use SSH and OpenSSL to obtain the certificate thumbprint for a vCenter Server Appliance instance.

1. Use SSH to connect to the vCenter Server Appliance as `root` user.<pre>$ ssh root@<i>vcsa_address</i></pre>
2. Use `openssl` to view the certificate thumbprint. <pre>openssl x509 -in /etc/vmware-vpx/ssl/rui.crt -fingerprint -sha1 -noout</pre>
3. Copy the certificate thumbprint so that you can verify it when you deploy a management cluster.

### Platform Services Controller

On vSphere 6.7u3, you can obtain a vCenter Server certificate thumbprint by logging into the Platform Services Controller for that vCenter Server instance. If you are deploying a management cluster to vSphere 7, there is no Platform Services Controller.

1. Log in to the Platform Services Controller interface.

    - Embedded Platform Services Controller: https://<i>vcenter_server_address</i>/psc
    - Standalone Platform Services Controller: https://<i>psc_address</i>/psc

1. Select **Certificate Management** and enter a vCenter Single Sign-On password.
1. Select **Machine Certificates**, select a certificate, and click **Show Details**.
1. Copy the certificate thumbprint so that you can verify it when you deploy a management cluster.

## Thick provisioning versus thin provisioning {#provisioning}

If you select thick provisioning as the disk format, when Tanzu Community Edition creates cluster node VMs from the template, the full size of each node's disk is reserved. This can rapidly consume storage if you deploy many clusters or clusters with many nodes. However, if you select thin provisioning, as you deploy clusters this can give a false impression of the amount of storage that is available. If you select thin provisioning, there might be enough storage available at the time that you deploy clusters, but storage might run out as the clusters run and accumulate data.