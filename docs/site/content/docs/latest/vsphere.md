# Prepare to Deploy a Management or Stand-alone Clusters to vSphere

This topic explains how to prepare your environment before you deploy a management or standalone cluster to vSphere. You must:

* Download and import the OVF image and covert it to a template.
* Create a ssh key pair and provide the public key part of the SSH key pair to Tanzu Installer.

## Before you begin

* [ ] Ensure the Tanzu CLI is installed locally on the bootstrap machine. See [Install the Tanzu CLI](installation-cli.md).
* [ ] Review the vSphere account reference information here: [Reference for vSphere account](ref-vsphere.md)
* [ ] Ensure that vSphere meets the following general requirements:
  * [ ] A vSphere 7, vSphere 6.7u3, VMware Cloud on AWS, or Azure VMware Solution account with:
    * [ ] **vSphere 6.7**: an Enterprise Plus license
    * [ ] **vSphere 7** see [vSphere with Tanzu Provides Management Cluster](#mc-vsphere7) below.
  * [ ] Your vSphere instance has the following objects in place:
    * [ ] Either a standalone host or a vSphere cluster
    * [ ] If you are deploying to a vSphere cluster, ideally vSphere DRS is enabled and two or more hosts compose the cluster.
    * [ ] Optionally, a resource pool in which to deploy the Tanzu Community Edition Instance
    * [ ] A VM folder in which to collect the Tanzu Community Edition VMs
    * [ ] A datastore with sufficient capacity for the control plane and worker node VM files
  * [ ] If you intend to deploy multiple Tanzu Community Edition instances to the same vSphere instance, create a dedicated resource pool, VM folder, and network for each instance that you deploy.
  * [ ] A network with:
    * [ ] A DHCP server to connect the cluster node VMs that Tanzu Community Edition deploys. The node VMs must be able to connect to vSphere.
    * [ ] A set of available static virtual IP (VIP) addresses for the clusters that you create,  one for each management and workload cluster. Each cluster requires a static IP address that is used to access the cluster's Kubernetes control plane. Make sure that these IP addresses are not in the DHCP range, but are in the same subnet as the DHCP range.
    > **NOTE:** To make DHCP-assigned IP addresses static, after you deploy the cluster, configure a DHCP reservation with enough IPs for each control plane and worker node in the cluster.
    * [ ] Traffic allowed out to vCenter Server from the network on which clusters will run.
    * [ ] Traffic allowed between your local bootstrap machine and port 6443 of all VMs in the clusters you create. Port 6443 is where the Kubernetes API is exposed.
    * [ ] Traffic allowed between port 443 of all VMs in the clusters you create and vCenter Server. Port 443 is where the vCenter Server API is exposed.
    * [ ] The Network Time Protocol (NTP) service running on all hosts, and the hosts running on UTC. For more information, see [Configuring Network Time Protocol (NTP) on an ESXi host using the vSphere Client](https://kb.vmware.com/s/article/57147).
  * [ ] If your vSphere environment runs NSX-T Data Center, you can use the NSX-T Data Center interfaces when you deploy management clusters. Make sure that your NSX-T Data Center setup includes a segment on which DHCP is enabled. Make sure that NTP is configured on all ESXi hosts, on vCenter Server, and on the bootstrap machine.
  * [ ] You will need a VMware Customer Connect to download the OVAs. Register [here](https://customerconnect.vmware.com/account-registration).

## Procedure

1. Download the OVA for the management cluster nodes, directly from [VMware Customer Connect](https://customerconnect.vmware.com/downloads/get-download?downloadGroup=TCE-090).  
Alternatively, you can open the [Tanzu Community Edition product page](https://customerconnect.vmware.com/downloads/info/slug/infrastructure_operations_management/vmware_tanzu_community_edition/0_9_0) in Customer Connect and select and download the OVA version that you require. You will need a VMware Customer Connect to download the OVAs. Register [here](https://customerconnect.vmware.com/account-registration).


2. Complete the following steps to deploy the OVF template:
   1. In the vSphere client, right-click an object in the vCenter Server inventory, and select **Deploy OVF template**.
   2. Select **Local file**, click the button to upload files, and navigate to the downloaded OVA file on your local machine.
   3. Follow the installer prompts to deploy a VM from the OVA:

      * Accept or modify the appliance name
      * Select the destination datacenter or folder
      * Select the destination host, cluster, or resource pool
      * Accept the end user license agreements (EULA)
      * Select the disk format and destination datastore
      * Select the network for the VM to connect to

   4. Click **Finish** to deploy the VM.
   5. When the OVA deployment finishes, right-click the VM and select **Template** > **Convert to Template**.

      > **NOTE**: Do not power on the VM before you convert it to a template.

3. Complete the following steps to create an SSH Key Pair:

   1. On the bootstrap machine on which you will run the Tanzu CLI, run the following `ssh-keygen` command:

      ```sh
      ssh-keygen -t rsa -b 4096 -C "email@example.com"
      ```

   2. At the prompt `Enter file in which to save the key (/root/.ssh/id_rsa):` press Enter to accept the default.
   3. Enter and repeat a password for the key pair.
   4. Add the private key to the SSH agent running on your machine, and enter the password you created in the previous step:

      ```sh
      ssh-add ~/.ssh/id_rsa
      ```

   5. Open the file `.ssh/id_rsa.pub` in a text editor so that you can easily copy and paste it when you deploy a management cluster.

      > For information about how to install OpenSSH on Windows, see [Install OpenSSH](https://docs.microsoft.com/en-us/windows-server/administration/openssh/openssh_install_firstuse).
