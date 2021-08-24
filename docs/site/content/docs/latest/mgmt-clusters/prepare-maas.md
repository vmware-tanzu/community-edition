# Prepare a vSphere Management as a Service Infrastructure

Tanzu Kubernetes Grid runs on two Management as a Service (MaaS) products that provide a vSphere interface and environment to public cloud infrastructures: VMware Cloud on AWS and Azure VMware Solution.

This topic explains how to prepare these services and use them to create a bootstrap machine for deploying Tanzu Kubernetes Grid.
For both VMware Cloud on AWS and Azure VMware Solution, the bootstrap machine is not a local physical machine, but is instead a cloud VM jumpbox that connects vSphere with its underlying infrastructure.

## <a id="prep-vmc"></a> Preparing VMware Cloud on AWS

To run Tanzu Kubernetes Grid on VMware Cloud on AWS, set up a Software-Defined Datacenter (SDDC) and create a bootstrap VM as follows.
The bootstrap machine is a VM managed through vCenter:

1. Log into the VMC Console and create a new SDDC by following the procedure [Deploy an SDDC from the VMC Console](https://docs.vmware.com/en/VMware-Cloud-on-AWS/services/com.vmware.vmc-aws.getting-started/GUID-EF198D55-03E3-44D1-AC48-6E2ABA31FF02.html) in the VMware Cloud on AWS documentation.

   - After you click **Deploy SDDC**, the SDDC creation process typically takes 2-3 hours.

1. Once the SDDC is created, open its pane in the VMC Console and click **Networking & Security** > **Network** > **Segments**.

1. The **Segment List** shows `sddc-cgw-network-1` with a subnet CIDR of `192.168.1.1/24`, giving 256 addresses.
If you need more internal IP addresses, you can:

   - Open `sddc-cgw-network-1` and modify its subnet CIDR to something broader, like 192.168.1.1/20.
   - Click **Add Segment** and create another network segment with a different subnet.  Make sure the new subnet CIDR does not overlap with `sddc-cgw-network-1` or any other existing segments.

1. Open `sddc-cgw-network-1` and any other network segments.
For each segment, click **Edit DHCP Config**.  A **Set DHCP Config** pane appears.

1. In the **Set DHCP Config** pane:

   - Set **DHCP Config** to **Enabled**.
   - Set **DHCP Ranges** to an IP address range or CIDR within the segment's subnet, but that leaves a pool of addresses free to serve as static IP addresses for Tanzu Kubernetes clusters.
   Each management cluster and workload cluster that Tanzu Kubernetes Grid creates will require a unique static IP address from this pool.

1. To enable access to vCenter, add a firewall rule or set up a VPN, following the [Connect to vCenter Server](https://docs.vmware.com/en/VMware-Cloud-on-AWS/services/com.vmware.vmc-aws.getting-started/GUID-C057DF7D-8016-45C4-AE12-56490E013F95.html) instructions in the VMware Cloud on AWS documentation.

1. To confirm access to vCenter, click **OPEN VCENTER** at upper-right in the SDDC pane. The vCenter client should appear.

1. From the vCenter portal, deploy your bootstrap machine and enable access to it following [Deploy Workload VMs](https://docs.vmware.com/en/VMware-Cloud-on-AWS/services/com.vmware.vmc-aws.getting-started/GUID-80B5A391-83B6-4AB1-BA92-4111AB5A61F6.html) in the VMware Cloud on AWS documentation.

   - You can log into the bootstrap machine by clicking **Launch Web Console** on its vCenter summary pane.
   - (Optional) If you want to `ssh` into the bootstrap machine, in addition to using the web console within vCenter, see [Set Up a VMware Cloud Bootstrap Machine for `ssh`](#vmc-ssh), below.

1. When installing the Tanzu CLI, deploying management clusters, and performing other operations, follow the instructions for vSphere, not the instructions for Amazon EC2.

### <a id="vmc-ssh"></a> Set Up a VMware Cloud Bootstrap Machine for `ssh`

To set up your bootstrap machine for access via `ssh`, follow these procedures in the VMware Cloud for AWS documentation:

1. [Assign a Public IP Address to a VM](https://docs.vmware.com/en/VMware-Cloud-on-AWS/services/com.vmware.vmc-aws.getting-started/GUID-BFE71806-64FC-4CD3-BB21-F1FEFD1478E3.html) to request a public IP address for the bootstrap machine.

1. [Create or Modify NAT Rules](https://docs.vmware.com/en/VMware-Cloud-on-AWS/services/com.vmware.vmc-aws.networking-security/GUID-DD72B243-1D08-4B6F-8A73-A745E8B0DC81.html) to create a NAT rule for the bootstrap machine, configured with:

   - **Public IP**: The public IP address requested above.
   - **Internal IP**: The IP address of the bootstrap machine. Can be either a static or DHCP IP.

1. The [Procedure](https://docs.vmware.com/en/VMware-Cloud-on-AWS/services/com.vmware.vmc-aws.networking-security/GUID-A5114A98-C885-4244-809B-151068D6A7D7.html) in _Add or Modify Compute Gateway Firewall Rules_ to add a compute gateway rule allowing access to the VM.

## <a id="prep-avs"></a> Preparing Azure VMware Solution on Microsoft Azure

To run Tanzu Kubernetes Grid on [Azure VMware Solution](https://docs.microsoft.com/en-us/azure/azure-vmware/introduction) (AVS),
set up AVS and its Windows 10 jumphost as follows.
The jumphost serves as the bootstrap machine for Tanzu Kubernetes Grid:

1. Log into NSX-T Manager as `admin`.

1. Unless you are intentionally deploying to an airgapped environment, confirm that AVS is configured to allow internet connectivity for AVS-hosted VMs. This is not enabled by default. To configure this, you can either:

    - Route outbound internet traffic through your on-premises datacenter by configuring Express Route Global Reach.
    - Allow internet access via the AVS Express Route connection to the Azure network by logging into the Azure portal, navigating to the AVS Private Cloud object, selecting **Manage** > **Connectivity**, flipping the **Internet enabled** toggle to **Enabled**, and clicking **Save**.

     ![Configure AVS Private Cloud Connectivity](../images/avs-connect.png)

1. Under **Networking** > **Connectivity** > **Segments**, click **Add Segment**, and configure the new segment with:

    - **Segment Name**: An identifiable name, like `avs_tkg`
    - **Connected Gateway**: The Tier-1 gateway that was predefined as part of your AVS account
    - **Subnets**: A subnet such as `192.168.20.1/24`
    - **DHCP Config** > **DHCP Range**: An address range or CIDR within the subnet, for example `192.168.20.10-192.168.20.100`.
    This range must exclude a pool of subnet addresses that DHCP cannot assign, leaving them free to serve as static IP addresses for Tanzu Kubernetes clusters.<br />
    Each management cluster and workload cluster that Tanzu Kubernetes Grid creates will require a unique static IP address from the pool outside of this DHCP range.
    - **Transport Zone**: Select the Overlay transport zone that was predefined as part of your AVS account.

    **Note**: After you create the segment, it should be visible in vCenter.

1. From the **IP Management** > **DHCP** pane, click **Add Server**, and configure the new DHCP server with:

    - **Server Name**: An identifiable name, like `avs_tkg_dhcp`
    - **Server IP Address**: A range that does not overlap with the subnet of the segment created above, for example `192.168.30.1/24`.
    - **Lease Time**: 5400 seconds; shorter than the default interval, to release IP addresses sooner

1. Under **Networking** > **Connectivity** > **Tier-1 Gateways**, open the predefined gateway.

1. Click the Tier-1 gateway's **IP Address Management** setting and associate it with the DHCP server created above.

1. Configure a DNS forwarder in NSX-T Manager or the Azure portal:

  - **NSX-T Manager**:
      1. Under **Networking** > **IP Management** > **DNS**, click **DNS Zones**.
      1. Click **Add DNS Zone** > **Add Default Zone**, and provide the following:
          * **Zone Name**: An identifiable name like `avs_tkg_dns_zone`.
          * **DNS Servers**: Up to three comma-separated IP addresses representing valid DNS servers.
      1. Click **Save**, and then select the **DNS Services** tab
      1. Click **Add DNS Service**, and provide the following:
          * **Name**: An identifiable name, like `avs_tkg_dns_svc`.
          * **Tier0/Tier1 Gateway**: The Tier-1 gateway that was predefined as part of your AVS account.
          * **DNS Service IP**: An IP address that does not overlap with the any other subnets created, such as `192.168.40.1`.
          * **Default DNS Zone**: Select the Zone Name defined earlier.
      1. Click **Save**.

  - **Azure Portal**:
      1. Navigate to the AVS Private Cloud object and select **Workload Networking** > **DNS**.
      1. With the **DNS zones** tab selected, click **Add** and provide the following:
          * **Type**: Default DNS zone.
          * **DNS zone name**: An identifiable name like `avs_tkg_dns_zone`.
          * **DNS server IP**: Up to three DNS servers.
      1. Click **OK** and then click the **DNS service** tab.
      1. Click **Add** and provide the following:
          * **Name**: An identifiable name, like `avs_tkg_dns_svc`.
          * **DNS Service** IP: An IP address that does not overlap. with the any other subnets created, such as `192.168.40.1`
          * **Default DNS Zone**: Select the DNS zone name defined earlier.
      1. Click **OK**.

1. When installing the Tanzu CLI, deploying management clusters, and performing other operations, follow the instructions for vSphere, not the instructions for Azure.  Configure the management cluster with:

   - **Kubernetes Network Settings** > **Network Name**: The name of the new segment.
   - **Management Cluster Settings** > **Virtual IP Address** The IP address range of the new segment.

## <a id="what-next"></a> What to Do Next

Your infrastructure and bootstrap machine are ready for you to deploy the Tanzu CLI.
See [Install the Tanzu CLI and Other Tools](../install-cli.md) for instructions, and then proceed to [Deploy Management Clusters on vSphere](../mgmt-clusters/vsphere.md).
