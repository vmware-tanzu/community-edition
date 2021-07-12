# Prepare to Deploy Management Clusters to Microsoft Azure

This topic explains how to prepare your environment before you deploy a management or standalone cluster on Azure.

<!--If you are installing Tanzu Community Edition on Azure VMware Solution (AVS), you are installing to a vSphere environment.
See [Preparing Azure VMware Solution on Microsoft Azure](prepare-maas.md#prep-avs) in _Prepare a vSphere Management as a Service Infrastructure_ to prepare your environment
and [Prepare to Deploy Management Clusters to vSphere](vsphere.md) to deploy management clusters.

## <a id="process diagram"></a> Installation Process Overview

The following diagram shows the high-level steps for installing a Tanzu Community Edition management cluster on Azure, and the interfaces you use to perform them.

These steps include the preparations listed below plus the procedures described in either [Deploy Management Clusters with the Installer Interface](deploy-ui.md) or [Deploy Management Clusters from a Configuration File](deploy-cli.md).

![Process Diagram: Start, Install the Tanzu CLI, Register a TKG App on Azure, Accept the Base Image License. If first deploy and no advanced config options, deploy with installer interface. Else deploy with config file.](../images/azure-install-process.png)-->

## <a id="general-requirements"></a> General Requirements

- Ensure the Tanzu CLI is installed locally on the bootstrap machine. See [Install the Tanzu CLI](installation-cli.md).
- A Microsoft Azure account with:
   - Permissions required to register an app. See [Permissions required for registering an app](https://docs.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal#permissions-required-for-registering-an-app) in the Azure documentation.
   - Sufficient VM core (vCPU) quotas for your clusters. A standard Azure account has a quota of 10 vCPU per region. Tanzu Community Edition clusters require 2 vCPU per node, which translates to:
     - Management cluster:
         - `dev` plan: 4 vCPU (1 main, 1 worker)
         - `prod` plan: 8 vCPU (3 main , 1 worker)
     - Each workload cluster:
         - `dev` plan: 4 vCPU (1 main, 1 worker)
         - `prod` plan: 12 vCPU (3 main , 3 worker)
     - For example, assuming a **single management cluster** and all clusters with the same plan:
   <table width="100%" border="0">
   <tr>
     <th width="17%">Plan</th>
     <th width="22%">Workload Clusters</th>
     <th width="22%">vCPU for Workload</th>
     <th width="22%">vCPU for Management</th>
     <th width="17%">Total vCPU</th>
   </tr>
   <tr>
     <td rowspan="2">Dev</td>
     <td>1</td>
     <td>4</td>
     <td rowspan="2">4</td>
     <td>8</td>
   </tr>
   <tr>
     <td>5</td>
     <td>20</td>
     <td>24</td>
   </tr>
   <tr>
     <td rowspan="2">Prod</td>
     <td>1</td>
     <td>12</td>
     <td rowspan="2">8</td>
     <td>20</td>
   </tr>
   <tr>
     <td>5</td>
     <td>60</td>
     <td>68</td>
   </tr>
   </table>
   - Sufficient public IP address quotas for your clusters, including the quota for Public IP Addresses - Standard, Public IP Addresses - Basic, and Static Public IP Addresses. A standard Azure account has a quota of 10 public IP addresses per region. Every Tanzu Community Edition cluster requires 2 Public IP addresses regardless of how many control plane nodes and worker nodes it has. For each Kubernetes Service object with type `LoadBalancer`, 1 Public IP address is required.
   - Run a DNS lookup on all `imageRepository` values to find their CNAMEs.
- (Optional) OpenSSL installed locally, to create a new keypair or validate the download package thumbprint.  See [OpenSSL](https://www.openssl.org).
- (Optional) A VNET with:
   - A subnet for the management cluster control plane node
   - A Network Security Group on the control plane subnet with the following inbound security rules, to enable SSH and Kubernetes API server connections:
      - Allow TCP over port 22 for any source and destination
      - Allow TCP over port 6443 for any source and destination.
      Port 6443 is where the Kubernetes API is exposed on VMs in the clusters you create.
   - A subnet and Network Security Group for the management cluster worker nodes.

   If you do not use an existing VNET, the installation process creates a new one.

- The Azure CLI installed locally. See [Install the Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) in the Microsoft Azure documentation.

<!--&#42;Or see [Deploying Tanzu Community Edition in an Internet-Restricted Environment](airgapped-environments.md) for installing without external network access.-->

## <a id="nsgs"></a> Network Security Groups on Azure

Tanzu Community Edition management and workload clusters on Azure require the following Network Security Groups (NSGs) to be defined on their VNET:

- One control plane NSG shared by the control plane nodes of all clusters, including the management cluster and the workload clusters that it manages.
- One worker NSG for each cluster, for the cluster's worker nodes.

If you do not specify a VNET when deploying a management cluster, the deployment process creates a new VNET along with the NSGs required for the management cluster.
If you optionally create a VNET for Tanzu Community Edition before deploying a management cluster, you must also create these NSGs as described in the [General Requirements](#requirements) above.

For each workload cluster that you deploy later, you need to create a worker NSG named `CLUSTER-NAME-node-nsg`, where `CLUSTER-NAME` is the name of the workload cluster.
This worker NSG must have the same VNET and region as its management cluster.

## <a id="tkg-app"></a> Register Tanzu Community Edition as an Azure Client App

Tanzu Community Edition manages Azure resources as a registered client application that accesses Azure through a service principal account.
The following steps register your Tanzu Community Edition application with Azure Active Directory, create its account, create a client secret for authenticating communications, and record information needed later to deploy a management cluster.

1. Log in to the [Azure Portal](https://portal.azure.com).

1. Record your **Tenant ID** by hovering over your account name at upper-right, or else browse to **Azure Active Directory** > \<Your Azure Org\> > **Properties** > **Tenant ID**.  The value is a GUID, for example `b39138ca-3cee-4b4a-a4d6-cd83d9dd62f0`.

1. Browse to **Active Directory** > **App registrations** and click **+ New registration**.

1. Enter a display name for the app, such as `tkg`, and select who else can use it.  You can leave the **Redirect URI (optional)** field blank.

1. Click **Register**.  This registers the application with an Azure service principal account as described in [How to: Use the portal to create an Azure AD application and service principal that can access resources](https://docs.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal) in the Azure documentation.

1. An overview pane for the app appears. Record its **Application (client) ID** value, which is a GUID.

1. From the Azure Portal top level, browse to **Subscriptions**.  At the bottom of the pane, select one of the subscriptions you have access to, and record its **Subscription ID**.  Click the subscription listing to open its overview pane.

1. Select to **Access control (IAM)** and click **Add a role assignment**.

1. In the **Add role assignment** pane
    - Select the **Owner** role
    - Leave **Assign access to** selection as "Azure AD user, group, or service principal"
    - Under **Select** enter the name of your app, `tkg`.  It appears underneath under **Selected Members**

1. Click **Save**. A popup appears confirming that your app was added as an owner for your subscription.

1. From the Azure Portal > **Azure Active Directory** > **App Registrations**, select your `tkg` app under **Owned applications**. The app overview pane opens.

1. From **Certificates & secrets** > **Client secrets** click **+ New client secret**.

1. In the **Add a client secret** popup, enter a **Description**, choose an expiration period, and click **Add**.

1. Azure lists the new secret with its generated value under **Client Secrets**.  Record the value.

## <a id="license"></a> Accept the Base Image License

To run management cluster VMs on Azure, accept the license for their base Kubernetes version and machine OS.

1. Sign in to the Azure CLI as your `tkg` client application.

   ```bash
   az login --service-principal --username AZURE_CLIENT_ID --password AZURE_CLIENT_SECRET --tenant AZURE_TENANT_ID
   ```

   Where `AZURE_CLIENT_ID`, `AZURE_CLIENT_SECRET`, and `AZURE_TENANT_ID` are your `tkg` app's client ID and secret and your tenant ID, as recorded in [Register Tanzu Community Edition as an Azure Client App](#tkg-app).

1. Run the `az vm image terms accept` command, specifying the `--plan` and your Subscription ID.

   In Tanzu Community Edition v1.3.1, the default cluster image `--plan` value is `k8s-1dot20dot5-ubuntu-2004`, based on Kubernetes version 1.20.5 and the  machine OS, Ubuntu 20.04. Run the following command:

   ```
   az vm image terms accept --publisher vmware-inc --offer tkg-capi --plan k8s-1dot20dot5-ubuntu-2004 --subscription AZURE_SUBSCRIPTION_ID
   ```

   Where `AZURE_SUBSCRIPTION_ID` is your Azure subscription ID.

You must repeat this to accept the base image license for every version of Kubernetes or OS that you want to use when you deploy clusters, and every time that you upgrade to a new version of Tanzu Community Edition.

## <a id="ssh-key"></a> Create an SSH Key Pair (Optional)

You deploy management clusters from a machine referred to as the _bootstrap machine_, using the Tanzu CLI.
To connect to Azure, the bootstrap machine must provide the public key part of an SSH key pair. If your bootstrap machine does not already have an SSH key pair, you can use a tool such as `ssh-keygen` to generate one.

1. On your bootstrap machine, run the following `ssh-keygen` command.

   <pre>ssh-keygen -t rsa -b 4096 -C "<em>email@example.com</em>"</pre>
1. At the prompt `Enter file in which to save the key (/root/.ssh/id_rsa):` press Enter to accept the default.
1. Enter and repeat a password for the key pair.
1. Add the private key to the SSH agent running on your machine, and enter the password you created in the previous step.

   ```sh
   ssh-add ~/.ssh/id_rsa
   ```

1. Open the file `.ssh/id_rsa.pub` in a text editor so that you can easily copy and paste it when you deploy a management cluster.

## <a id="checklist"></a> Preparation Checklist

Use this checklist to make sure you are prepared to deploy a Tanzu Community Edition management cluster to Azure:

- Tanzu CLI installed

   - Run `tanzu version`. The output should list `version: v1.3.1`.

- Azure account

   - Log in to the Azure web portal at `https://portal.azure.com`.

- Azure CLI installed

   - Run `az version`. The output should list the current version of the Azure CLI as listed in [Install the Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli), in the Microsoft Azure documentation.

- Registered `tkg` app

   - In the Azure portal, select **Active Directory** > **App Registrations** > **Owned applications** and confirm that your `tkg` app is listed as configured in [Register Tanzu Community Edition as an Azure Client App](#tkg-app) above, and with a current certificate.

- Base VM image license accepted

   - Run `az vm image terms show --publisher vmware-inc --offer tkg-capi --plan k8s-1dot20dot5-ubuntu-2004`. The output should contain `"accepted": true`.

