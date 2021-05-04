# Prepare to Deploy Management Clusters to vSphere

Before you can use the Tanzu CLI or installer interface to deploy a management cluster, you must prepare your vSphere environment. You must make sure that vSphere meets the general requirements, and import the base image templates from which Tanzu Kubernetes Grid creates cluster node VMs.
Each base image template contains a version of a machine OS and a version of Kubernetes.

## <a id="general-requirements"></a> General Requirements 

- A machine with the Tanzu CLI, Docker, and `kubectl` installed. See [Install the Tanzu CLI](../install-cli.md).
    - This is the _bootstrap machine_ from which you run `tanzu`, `kubectl` and other commands.
    - The bootstrap machine can be a local physical machine or a VM that you access via a console window or client shell.
- A vSphere 7, vSphere 6.7u3, VMware Cloud on AWS, or Azure VMware Solution account with:
    - **vSphere 6.7**: an Enterprise Plus license
    - **vSphere 7** see [vSphere with Tanzu Provides Management Cluster](#mc-vsphere7) below.
    - At least the permissions described in [Required Permissions for the vSphere Account](#vsphere-permissions).
- Your vSphere instance has the following objects in place:
    - Either a standalone host or a vSphere cluster with at least two hosts
        - If you are deploying to a vSphere cluster, ideally vSphere DRS is enabled.
    - Optionally, a resource pool in which to deploy the Tanzu Kubernetes Grid Instance
    - A VM folder in which to collect the Tanzu Kubernetes Grid VMs
    - A datastore with sufficient capacity for the control plane and worker node VM files
    - If you intend to deploy multiple Tanzu Kubernetes Grid instances to the same vSphere instance, create a dedicated resource pool, VM folder, and network for each instance that you deploy.
- You have done the following to prepare your vSphere environment:
    - Created a base image template that matches the management cluster's Kubernetes version.  See [Import the Base Image Template into vSphere](#import-base).
    - Created a vSphere account for Tanzu Kubernetes Grid, with a role and permissions that let it manipulate vSphere objects as needed.  See [Required Permissions for the vSphere Account](#vsphere-permissions).
- A network&#42; with:
    - A DHCP server to which to connect the cluster node VMs that Tanzu Kubernetes Grid deploys. The node VMs must be able to connect to vSphere.
    - A set of available static virtual IP (VIP) addresses for the clusters that you create, including management and Tanzu Kubernetes clusters. Each control plane and worker node requires a static IP address. Make sure that these IP addresses are not in the DHCP range, but are in the same subnet as the DHCP range. For more information, see [Static VIPs and Load Balancers for vSphere](#load-balancer).
    **NOTE:** To make DHCP-assigned IP addresses static, after you deploy the cluster, configure a DHCP reservation with enough IPs for each control plane and worker node in the cluster.
    - Traffic allowed out to vCenter Server from the network on which clusters will run.
    - Traffic allowed between your local bootstrap machine and port 6443 of all VMs in the clusters you create. Port 6443 is where the Kubernetes API is exposed.
    - Traffic allowed between port 443 of all VMs in the clusters you create and vCenter Server. Port 443 is where the vCenter Server API is exposed.
    - Traffic allowed between your local bootstrap machine out to the image repositories listed in the management cluster Bill of Materials (BoM) file, over port 443, for TCP. The BoM file is under `~/.tanzu/tkg/bom/` and its name includes the Tanzu Kubernetes Grid version, for example `bom-1.3.0+vmware.1.yaml` for v1.3.0.
    - The Network Time Protocol (NTP) service running on all hosts, and the hosts running on UTC. To check the time settings on hosts:
       1. Use SSH to log in to the ESXi host.
       1. Run the `date` command to see the timezone settings.
       1. If the timezone is incorrect, run `esxcli system time set`.
- If your vSphere environment runs NSX-T Data Center, you can use the NSX-T Data Center interfaces when you deploy management clusters. Make sure that your NSX-T Data Center setup includes a segment on which DHCP is enabled. Make sure that NTP is configured on all ESXi hosts, on vCenter Server, and on the bootstrap machine.

&#42;Or see [Deploying Tanzu Kubernetes Grid in an Internet-Restricted Environment](airgapped-environments.md) for installing without external network access.


## <a id="mc-vsphere7"></a> vSphere with Tanzu Provides Management Cluster

On vSphere 7 and later, the vSphere with Tanzu feature includes a Supervisor Cluster that you can configure as a management cluster for Tanzu Kubernetes Grid. This means that on vSphere 7, you do not need to use the `tanzu management-cluster create` to deploy a management cluster if vSphere with Tanzu is enabled. Deploying a Tanzu Kubernetes Grid management cluster to vSphere 7 when vSphere with Tanzu is not enabled is supported, but the preferred option is to enable vSphere with Tanzu and use the built-in Supervisor Cluster.

The Tanzu CLI works with both management clusters deployed through vSphere with Tanzu and management clusters deployed by Tanzu Kubernetes Grid on Azure, Amazon EC2, and vSphere 6.7, letting you deploy and manage workload clusters across multiple infrastructures using a single tool. For more information, see [Use the Tanzu CLI with a vSphere with Tanzu Supervisor Cluster](../tanzu-k8s-clusters/connect-vsphere7.md).

For information about the vSphere with Tanzu feature in vSphere 7, see [vSphere with Tanzu Configuration and Management](https://docs.vmware.com/en/VMware-vSphere/7.0/vmware-vsphere-with-kubernetes/GUID-152BE7D2-E227-4DAA-B527-557B564D9718.html) in the vSphere 7 documentation.

**NOTE**: On VMware Cloud on AWS and Azure VMware Solution, you cannot create a supervisor cluster, and need to deploy a management cluster to run `tanzu` commands.

## <a id="load-balancer"></a> Static VIPs and Load Balancers for vSphere

Each management cluster and Tanzu Kubernetes cluster that you deploy to vSphere requires one static virtual IP address for external requests to the cluster's API server. You must be able to assign this IP address, so it cannot be within your DHCP range, but it must be in the same subnet as the DHCP range.

The cluster control plane's [Kube-vip](https://kube-vip.io/) pod uses this static virtual IP address to serve API requests, and the API server certificate includes the address to enable secure TLS communication.  In Tanzu Kubernetes clusters, Kube-vip runs in a basic, Layer-2 failover mode, assigning the virtual IP address to one control plane node at a time. In this mode, Kube-vip does not function as a true load balancer for control plane traffic.

Tanzu Kubernetes Grid also does not use Kube-vip as a load balancer for workloads in workload clusters.
Kube-vip is used solely by the cluster's API server.

To load-balance workloads on vSphere, use NSX Advanced Load Balancer, also known as Avi Load Balancer, Essentials Edition.
You must deploy the NSX Advanced Load Balancer in your vSphere instance before you deploy management clusters.
See [Install VMware NSX Advanced Load Balancer on a vSphere Distributed Switch](install-nsx-adv-lb.md).

## <a id="import-base"></a> Import the Base Image Template into vSphere

Before you can deploy a management cluster to vSphere, you must import into vSphere a base image template for the OS and Kubernetes versions that the management cluster nodes run on.
VMware publishes base image templates in OVA format for each OS and Kubernetes version supported for management and workload clusters.
After importing the OVA, you must convert the resulting VM into a VM template.
For information about the versions of Kubernetes that each Tanzu Kubernetes Grid release supports, see the release notes for that release.

1. Go to [the Tanzu Kubernetes Grid downloads page](https://my.vmware.com/en/web/vmware/downloads/info/slug/infrastructure_operations_management/vmware_tanzu_kubernetes_grid/1_x), log in with your My VMware credentials, and click **Go to Downloads**.
1. Download a Tanzu Kubernetes Grid OVA for the management cluster nodes, which can be one of: 

   - Kubernetes v1.20.4: **Ubuntu v20.04 Kubernetes v1.20.4 OVA**
   - Kubernetes v1.20.4: **Photon v3 Kubernetes v1.20.4 OVA**
   
   You can also download base image templates for other OS and Kubernetes versions that you expect to create clusters from, or you can download them later.

1. In the vSphere Client, right-click an object in the vCenter Server inventory, select **Deploy OVF template**.
1. Select **Local file**, click the button to upload files, and navigate to the downloaded OVA file on your local machine.
1. Follow the installer prompts to deploy a VM from the OVA.

    - Accept or modify the appliance name
    - Select the destination datacenter or folder
    - Select the destination host, cluster, or resource pool
    - Accept the end user license agreements (EULA)
    - Select the disk format and destination datastore
    - Select the network for the VM to connect to
    
    **NOTE**: If you select thick provisioning as the disk format, when Tanzu Kubernetes Grid creates cluster node VMs from the template, the full size of each node's disk will be reserved. This can rapidly consume storage if you deploy many clusters or clusters with many nodes. However, if you select thin provisioning, as you deploy clusters this can give a false impression of the amount of storage that is available. If you select thin provisioning, there might be enough storage available at the time that you deploy clusters, but storage might run out as the clusters run and accumulate data.
1. Click **Finish** to deploy the VM.
1. When the OVA deployment finishes, right-click the VM and select **Template** > **Convert to Template**.

   **NOTE**: Do not power on the VM before you convert it to a template.
1. In the **VMs and Templates** view, right-click the new template, select **Add Permission**, and assign the `tkg-user` to the template with the `TKG` role.

   For information about how to create the user and role for Tanzu Kubernetes Grid, see [Required Permissions for the vSphere Account](#vsphere-permissions) above.
   
Repeat the procedure for each of the Kubernetes versions for which you downloaded the OVA file.


## <a id="vsphere-permissions"></a> Required Permissions for the vSphere Account

The vCenter Single Sign On account that you provide to Tanzu Kubernetes Grid when you deploy a management cluster must have at the correct permissions in order to perform the required operations in vSphere.  

It is not recommended to provide a vSphere administrator account to Tanzu Kubernetes Grid, because this provides Tanzu Kubernetes Grid with far greater permissions than it needs. The best way to assign permissions to Tanzu Kubernetes Grid is to create a role and a user account, and then to grant that user account that role on vSphere objects.

**NOTE**: If you are deploying Tanzu Kubernetes clusters to vSphere 7 and vSphere with Tanzu is enabled, you must set the **Global** > **Cloud Admin** permission in addition to the permissions listed below. If you intend to use Velero to back up and restore management or workload clusters, you must also set the permissions listed in [Credentials and Privileges for VMDK Access](https://code.vmware.com/docs/11750/virtual-disk-development-kit-programming-guide/GUID-8301C6CF-37C2-42CC-B4C5-BB1DD28F79C9.html) in the *Virtual Disk Development Kit Programming Guide*.

1. In the vSphere Client, go to **Administration** > **Access Control** > **Roles**, and create a new role, for example `TKG`, with the following permissions.
   
   <table width="100%" border="0">
   <tr>
    <th scope="col">vSphere Object </th>
    <th scope="col">Required Permission </th>
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

1. In **Administration** > **Single Sign On** > **Users and Groups**, create a new user account in the appropriate domain, for example `tkg-user`.
1.  In the **Hosts and Clusters**, **VMs and Templates**, **Storage**, and **Networking** views, right-click the objects that your Tanzu Kubernetes Grid deployment will use, select **Add Permission**, and assign the `tkg-user`  with the `TKG` role to each object.

   - Hosts and Clusters
      - The root vCenter Server object
      - The Datacenter and all of the Host and Cluster folders, from the Datacenter object down to the cluster that manages the Tanzu Kubernetes Grid deployment
      - Target hosts and clusters
      - Target resource pools, with propagate to children enabled
   - VMs and Templates
      - The deployed Tanzu Kubernetes Grid base image templates
      -  Target VM and Template folders, with propagate to children enabled
   - Storage
      - Datastores and all storage folders, from the Datacenter object down to the datastores that will be used for Tanzu Kubernetes Grid deployments 
   - Networking
      - Networks or distributed port groups to which clusters will be assigned
      - Distributed switches


## <a id="ssh-key"></a> Create an SSH Key Pair

In order for the Tanzu CLI to connect to vSphere from the machine on which you run it, you must provide the public key part of an SSH key pair to Tanzu Kubernetes Grid when you deploy the management cluster. If you do not already have one on the machine on which you run the CLI, you can use a tool such as `ssh-keygen` to generate a key pair.

1. On the machine on which you will run the Tanzu CLI, run the following `ssh-keygen` command.

   <pre>ssh-keygen -t rsa -b 4096 -C "<em>email@example.com</em>"</pre>
1. At the prompt `Enter file in which to save the key (/root/.ssh/id_rsa):` press Enter to accept the default.
1. Enter and repeat a password for the key pair.
1. Add the private key to the SSH agent running on your machine, and enter the password you created in the previous step.

   ```
   ssh-add ~/.ssh/id_rsa
   ```
1. Open the file `.ssh/id_rsa.pub` in a text editor so that you can easily copy and paste it when you deploy a management cluster.
 

## <a id="vc-thumbprint"></a> Obtain vSphere Certificate Thumbprints

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


## <a id="what-next"></a> What to Do Next

For production deployments, it is strongly recommended to enable identity management for your clusters. For information about the preparatory steps to perform before you deploy a management cluster, see [Enabling Identity Management in Tanzu Kubernetes Grid](enabling-id-mgmt.md).

If you are using Tanzu Kubernetes Grid in an environment with an external internet connection, once you have set up identity management, you are  ready to deploy management clusters to vSphere.

- [Deploy Management Clusters with the Installer Interface](deploy-ui.md). This is the preferred option for first deployments.
- [Deploy Management Clusters from a Configuration File](deploy-cli.md). This is the more complicated method,  that allows greater flexibility of configuration and automation.
- If you are using Tanzu Kubernetes Grid in an internet-restricted environment, see [Deploying Tanzu Kubernetes Grid in an Internet-Restricted Environment](airgapped-environments.md) for the additional steps to perform.
- If you want to deploy clusters to Amazon EC2 and Azure as well as to vSphere, see [Prepare to Deploy Management Clusters to Amazon EC2](aws.html) and [Prepare to Deploy Management Clusters to Microsoft Azure](azure.md) for the required setup for those platforms.
