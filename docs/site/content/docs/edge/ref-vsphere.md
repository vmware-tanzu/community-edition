# Reference for vSphere account

If you encounter issues deploying a cluster to vSphere, review the following troubleshooting and reference content.

## Required Permissions for the vSphere Account {#permissions}

The vCenter Single Sign On account that you provide to Tanzu Comunity Edition when you deploy a management cluster must have the correct permissions in order to perform the required operations in vSphere.

It is not recommended to provide a vSphere administrator account to Tanzu Community Edition, because this provides Tanzu Community Edition with far greater permissions than it needs. The best way to assign permissions to Tanzu Community Edition is to create a custom role, then either use an existing user account or create a new user account, and then to grant that user account that role on vSphere objects.

All steps can be completed in the user interface of the vCenter using the vSphere Client or programmatically with command line interfaces like [govc](https://github.com/vmware/govmomi/tree/master/govc) or [PowerCLI](https://developer.vmware.com/docs/15315).

### vCenter Server custom role

You can use one of the available options to create a new vCenter Server custom role.

#### Option 1: Create the role in the vSphere client

In the vSphere Client, go to **Administration** > **Access Control** > **Roles**, and create a new role, for example `TCE`, with the following permissions:

   |vSphere Object|Required Permission|
   |--- |--- |
   |Cns|Searchable|
   |Datastore|Allocate space, Browse datastore,     Low level file operations|
   |Global (if using Velero for backup and restore)|Disable methods, Enable methods, Licenses|
   |Network|Assign network|
   |Profile-driven storage|Profile-driven storage view|
   |Resource|Assign virtual machine to resource pool|
   |Sessions|Message, Validate session|
   |vApp|Import|

   |vSphere Object|Required Permission|
   |:---------------------- |:------------------------------------- |
   |Virtual machine|Change Configuration > Add existing disk|
   ||Change Configuration > Add new disk|
   ||Change Configuration > Add or remove device|
   ||Change Configuration > Advanced configuration|
   ||Change Configuration > Change CPU count|
   ||Change Configuration > Change Memory|
   ||Change Configuration > Change Settings|
   ||Change Configuration > Configure Raw device|
   ||Change Configuration > Extend virtual disk|
   ||Change Configuration > Modify device settings|
   ||Change Configuration > Remove disk|
   ||Change Configuration > Toggle disk change tracking*|
   ||Edit Inventory > Create from existing|
   ||Edit Inventory > Remove|
   ||Interaction > Power On|
   ||Interaction > Power Off|
   ||Provisioning > Allow read-only disk access*|
   ||Provisioning > Allow virtual machine download*|
   ||Provisioning > Deploy template|
   ||Snapshot Management > Create snapshot*
   ||Snapshot Management > Remove snapshot*|
   ||*Required to activate the Velero plugin, as described in Back Up and Restore Clusters.  You can add these permissions when needed later.|

#### Option 2: Create the role using GOVC

Create a new role, for example `TCE`, using govc. The provided example includes the required permissions for the [Velero plugin](#velero).

```sh
govc role.create TCE \
Cns.Searchable \
Datastore.AllocateSpace \
Datastore.Browse \
Datastore.FileManagement \
Global.DisableMethods \
Global.EnableMethods \
Global.Licenses \
Network.Assign \
StorageProfile.View \
Resource.AssignVMToPool \
Sessions.GlobalMessage \
Sessions.ValidateSession \
VApp.Import \
VirtualMachine.Config.AddExistingDisk \
VirtualMachine.Config.AddNewDisk \
VirtualMachine.Config.AddRemoveDevice \
VirtualMachine.Config.AdvancedConfig \
VirtualMachine.Config.CPUCount \
VirtualMachine.Config.Memory \
VirtualMachine.Config.Settings \
VirtualMachine.Config.RawDevice \
VirtualMachine.Config.DiskExtend \
VirtualMachine.Config.EditDevice \
VirtualMachine.Config.RemoveDisk \
VirtualMachine.Config.ChangeTracking \
VirtualMachine.Inventory.CreateFromExisting \
VirtualMachine.Inventory.Delete \
VirtualMachine.Interact.PowerOff \
VirtualMachine.Interact.PowerOn \
VirtualMachine.Provisioning.GetVmFiles \
VirtualMachine.Provisioning.DeployTemplate \
VirtualMachine.Provisioning.DiskRandomRead \
VirtualMachine.State.CreateSnapshot \
VirtualMachine.State.RemoveSnapshot
```

#### Option 3: Create the role using PowerCLI

Create a new role, for example `TCE`, using PowerCLI. The provided example includes the required permissions for the [Velero plugin](#velero).

```powershell
$tceprivs=@(
'Global.Licenses',
'Global.DisableMethods',
'Global.EnableMethods',
'Datastore.Browse',
'Datastore.FileManagement',
'Datastore.AllocateSpace',
'Network.Assign',
'VirtualMachine.Inventory.CreateFromExisting',
'VirtualMachine.Inventory.Delete',
'VirtualMachine.Interact.PowerOn',
'VirtualMachine.Interact.PowerOff',
'VirtualMachine.Config.AddExistingDisk',
'VirtualMachine.Config.AddNewDisk',
'VirtualMachine.Config.RemoveDisk',
'VirtualMachine.Config.RawDevice',
'VirtualMachine.Config.CPUCount',
'VirtualMachine.Config.Memory',
'VirtualMachine.Config.AddRemoveDevice',
'VirtualMachine.Config.EditDevice',
'VirtualMachine.Config.Settings',
'VirtualMachine.Config.AdvancedConfig',
'VirtualMachine.Config.DiskExtend',
'VirtualMachine.Config.ChangeTracking',
'VirtualMachine.State.CreateSnapshot',
'VirtualMachine.State.RemoveSnapshot',
'VirtualMachine.Provisioning.DeployTemplate',
'VirtualMachine.Provisioning.DiskRandomRead',
'VirtualMachine.Provisioning.GetVmFiles',
'Resource.AssignVMToPool',
'Sessions.ValidateSession',
'Sessions.GlobalMessage',
'VApp.Import',
'Cns.Searchable',
'StorageProfile.View')

New-VIRole -name TCE -Privilege (Get-VIPrivilege -id $tceprivs)
```

### vSphere Single Sign On user

 Use one of the available options to create a new user in the vSphere Single Sign On domain. Skip this step if you plan to use an existing user account.

#### Option 1: Creating a new user in the vSphere client

In **Administration** > **Single Sign On** > **Users and Groups**, create a new user account in the appropriate domain, for example `tce-user`.

#### Option 2: Create a new user user using GOVC

Create a new user account in the vSphere Single Sign On domain, for example `tce-user`.

```sh
govc sso.user.create -p <password> tce-user
```

#### Option 3: Create a new user user using PowerCLI

Create a new user account in the vSphere Single Sign On domain, for example `tce-user`.  
You must install the [VMware.vSphere.SsoAdmin](https://github.com/vmware/PowerCLI-Example-Scripts/tree/master/Modules/VMware.vSphere.SsoAdmin) module to complete this step.

```powershell
New-SsoPersonUser -UserName tce-user
```

### vSphere permissions

After you create the role and user you must assign the user and their role to the relevant inventory objects.

#### Option 1: Assigning permissions to Inventory Objects using the vSphere client {#vsphere-permissions-vsphere_client}

In the **Hosts and Clusters**, **VMs and Templates**, **Storage**, and **Networking** views, right-click the objects that your Tanzu Community Edition deployment will use, select **Add Permission**, and assign the `tce-user`  with the `TCE` role to each object.

   * Hosts and Clusters
     * The root vCenter Server object
     * The Datacenter and all of the Host and Cluster folders, from the Datacenter object down to the cluster that manages the Tanzu Community Edition deployment
     * Target hosts and clusters
     * Target resource pools, with propagate to children selected
   * VMs and Templates
     * The deployed Tanzu Community Edition base image templates
     * Target VM and Template folders, with propagate to children selected
     * Storage
     * Datastores and all storage folders, from the Datacenter object down to the datastores that will be used for Tanzu Community Edition deployments
   * Networking
     * Networks or distributed port groups to which clusters will be assigned
     * Distributed switches

#### Option 2: Assigning permissions to Inventory Objects using GOVC

Assign the custom role `TCE` to the `tce-user` on all inventory objects outlined in [Assigning permissions using the vSphere client](#vsphere-permissions-vsphere_client). For information about the assigning permissions using GOVC, see the [permission.set](https://github.com/vmware/govmomi/blob/master/govc/USAGE.md#permissionsset) command in the VMware GOVC documentation.

```sh
govc permissions.set -principal tce-user@vsphere.local -role TCE <inventory object>
```

#### Option 3: Assigning permissions to Inventory Objects using PowerCLI

Assign the custom role `TCE` to the `tce-user` on all inventory objects outlined in [Assigning permissions using the vSphere client](#vsphere-permissions-vsphere_client). For information about the assigning permissions using PowerCLI, see the [New-VIPermission](https://developer.vmware.com/docs/powercli/latest/vmware.vimautomation.core/commands/new-vipermission/#Default)-cmdlet in the VMware PowerCLI documentation.

```powershell
New-VIPermission -Entity <inventory object> -Principal (Get-VIAccount -Domain <SSO domain> |Where-Object {$_.Name -like "*tce*"}) -Role (Get-VIRole -name TCE)
```

### Velero permissions {#velero}

If you intend to use Velero to back up and restore management or workload clusters, you must also set the permissions listed in [Credentials and Privileges for VMDK Access](https://code.vmware.com/docs/11750/virtual-disk-development-kit-programming-guide/GUID-8301C6CF-37C2-42CC-B4C5-BB1DD28F79C9.html) in the *Virtual Disk Development Kit Programming Guide*.

## Static VIPs and Load Balancers for vSphere {#static-ip}

Each management cluster and workload cluster that you deploy to vSphere requires one static virtual IP address for external requests to the cluster's API server. You must be able to assign this IP address, so it cannot be within your DHCP range, but it must be in the same subnet as the DHCP range.

The cluster control plane's [Kube-vip](https://kube-vip.io/) pod uses this static virtual IP address to serve API requests, and the API server certificate includes the address to activate secure TLS communication.  In workload clusters, Kube-vip runs in a basic, Layer-2 failover mode, assigning the virtual IP address to one control plane node at a time. In this mode, Kube-vip does not function as a true load balancer for control plane traffic.

Tanzu Community Edition also does not use Kube-vip as a load balancer for workloads in workload clusters.
Kube-vip is used solely by the cluster's API server.

To load-balance workloads on vSphere, use NSX Advanced Load Balancer, also known as Avi Load Balancer, Essentials Edition.
You must deploy the NSX Advanced Load Balancer in your vSphere instance before you deploy management clusters.

## Obtain vSphere Certificate Thumbprints {#certificates}

If your vSphere environment uses untrusted, self-signed certificates to authenticate connections, you must verify the thumbprint of the vCenter Server when you deploy a management cluster. If your vSphere environment uses trusted certificates that are signed by a known Certificate Authority (CA), you do not need to verify the thumbprint.

### vCenter Server Appliance with embedded Platform Services Controller

On vSphere 7, you can obtain a vCenter Server certificate thumbprint using SSH and OpenSSL or govc.

#### Option 1: Use SSH and OpenSSL

If SSH is enabled on the vCenter Server Appliance, you can use SSH and OpenSSL to obtain the certificate thumbprint the instance.

1. Use SSH to connect to the vCenter Server Appliance as `root` user.

   ```sh
   ssh root@vcsa_address
   ```

2. Use `openssl` to view the certificate thumbprint.

   ```sh
   openssl x509 -in /etc/vmware-vpx/ssl/rui.crt -fingerprint -sha1 -noout
   ```

3. Copy the certificate thumbprint so that you can verify it when you deploy a management cluster.

#### Option 2: Use govc

1. Use govc to extract the thumbprint information from the vCenter Server Appliance remotely.

```sh
# This command on any operating system with govc. The last line of the output contains the SHA1 thumbprint.
govc about.cert

# On MacOS/Linux you can directly extract the SHA1 thumbprint when jq is installed on your system
govc about.cert -k -json | jq -r .ThumbprintSHA1
```

2. Copy the certificate thumbprint so that you can verify it when you deploy a management cluster.


### Platform Services Controller

On vSphere 6.7u3, you can obtain a vCenter Server certificate thumbprint by logging into the Platform Services Controller for that vCenter Server instance. If you are deploying a management cluster to vSphere 7, there is no Platform Services Controller.

1. Log in to the Platform Services Controller interface.

   * Embedded Platform Services Controller: https://*vcenter_server_address*/psc
   * Standalone Platform Services Controller: https://*psc_address*/psc

1. Select **Certificate Management** and enter a vCenter Single Sign-On password.
1. Select **Machine Certificates**, select a certificate, and click **Show Details**.
1. Copy the certificate thumbprint so that you can verify it when you deploy a management cluster.

## Thick provisioning versus thin provisioning {#provisioning}

If you select thick provisioning as the disk format, when Tanzu Community Edition creates cluster node VMs from the template, the full size of each node's disk is reserved. This can rapidly consume storage if you deploy many clusters or clusters with many nodes. However, if you select thin provisioning, as you deploy clusters this can give a false impression of the amount of storage that is available. If you select thin provisioning, there might be enough storage available at the time that you deploy clusters, but storage might run out as the clusters run and accumulate data.