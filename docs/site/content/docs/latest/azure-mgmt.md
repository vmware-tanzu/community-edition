# Preparing to Deploy Clusters to Microsoft Azure

This topic explains how to prepare your environment before you deploy a management or standalone cluster on Microsoft Azure.

<!--If you are installing Tanzu Community Edition on Azure VMware Solution (AVS), you are installing to a vSphere environment.
See [Preparing Azure VMware Solution on Microsoft Azure](prepare-maas.md#prep-avs) in _Prepare a vSphere Management as a Service Infrastructure_ to prepare your environment
and [Prepare to Deploy Management Clusters to vSphere](vsphere.md) to deploy management clusters.

## <a id="process diagram"></a> Installation Process Overview

The following diagram shows the high-level steps for installing a Tanzu Community Edition management cluster on Azure, and the interfaces you use to perform them.

These steps include the preparations listed below plus the procedures described in either [Deploy Management Clusters with the Installer Interface](deploy-ui.md) or [Deploy Management Clusters from a Configuration File](deploy-cli.md).

![Process Diagram: Start, Install the Tanzu CLI, Register a TKG App on Azure, Accept the Base Image License. If first deploy and no advanced config options, deploy with installer interface. Else deploy with config file.](../images/azure-install-process.png)-->

## <a id="general-requirements"></a> General Requirements

- [ ] Ensure Tanzu Community Edition is installed locally on your bootstrap machine. See [Install Tanzu Community Edition](cli-installation).

- [ ] Ensure the Azure CLI is installed locally.  See [Install the Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) in the Microsoft Azure documentation.

- [ ] Your Microsoft Azure account should meet the permissions and requirements described in the [Microsoft Azure account](ref-azure/#microsoft-azure-account) topic.

<!--&#42;Or see [Deploying Tanzu Community Edition in an Internet-Restricted Environment](airgapped-environments.md) for installing without external network access.-->

- [ ]Register Tanzu Community Edition as an Azure Client App. The full procedure is provided below: [Register Tanzu Community Edition as a Microsoft Azure Client App](azure-mgmt/#a-idtkg-appa-register-tanzu-community-edition-as-a-microsoft-azure-client-app).

- [ ] Accept the Base Image License.  The full procedure is provided below: [Accept the Base Image License](azure-mgmt/#accept-the-base-image-license).

- [ ] If you plan to use an existing VNET, see the [Network Security Groups on Microsoft Azure](ref-azure/#a-idnsgsa-network-security-groups-on-azure) topic for guidelines.

- [ ] (Optional) Create an SSH keypair. The full procedure is described below:[Create an SSH Key Pair](azure-mgmt/#a-idssh-keya-create-an-ssh-key-pair-optional).

- [ ] (Optional) For information about the configurations of the different sizes of node instances for Microsoft Azure, for example, Standard_D2s_v3 or Standard_D4s_v3, see [Sizes for virtual machines in Azure](https://docs.microsoft.com/en-us/azure/virtual-machines/sizes) in the  Microsoft Azure documentation.

## <a id="tkg-app"></a> Register Tanzu Community Edition as a Microsoft Azure Client App

Tanzu Community Edition manages Microsoft Azure resources as a registered client application that accesses Azure through a service principal account.
The following steps register your Tanzu Community Edition application with Microsoft Azure Active Directory, create its account, create a client secret for authenticating communications, and record information needed later to deploy a management cluster.

1. Log in to the [Azure Portal](https://portal.azure.com).

1. Record your **Tenant ID** by hovering over your account name at upper-right, or else browse to **Azure Active Directory** > \<Your Azure Org\> > **Properties** > **Tenant ID**.  The value is a GUID, for example `b39138ca-3cee-4b4a-a4d6-cd83d9dd62f0`.

1. Browse to **Active Directory** > **App registrations** and click **+ New registration**.

1. Enter a display name for the app, such as `tce`, and select who else can use it.  You can leave the **Redirect URI (optional)** field blank.

1. Click **Register**.  This registers the application with an Microsoft Azure service principal account as described in [How to: Use the portal to create an Azure AD application and service principal that can access resources](https://docs.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal) in the Microsoft Azure documentation.

1. An overview pane for the app appears. Record its **Application (client) ID** value, which is a GUID.

2. From the Microsoft Azure Portal, browse to **Subscriptions**.  At the bottom of the pane, select one of the subscriptions you have access to, and record its **Subscription ID**.  Click the subscription listing to open its overview pane.

3. Select to **Access control (IAM)** and click **Add a role assignment**.

4. In the **Add role assignment** pane
    - Select the **Owner** role
    - Leave **Assign access to** selection as "Azure AD user, group, or service principal"
    - Under **Select** enter the name of your app, `tce`.  It appears underneath under **Selected Members**

5. Click **Save**. A popup appears confirming that your app was added as an owner for your subscription.

6. From the Microsoft Azure Portal > **Azure Active Directory** > **App Registrations**, select your `tce` app under **Owned applications**. The app overview pane opens.

7. From **Certificates & secrets** > **Client secrets** click **+ New client secret**.

8. In the **Add a client secret** popup, enter a **Description**, choose an expiration period, and click **Add**.

9. The new secret is listed with its generated value under **Client Secrets**.  Record the value.

## Accept the Base Image License

To run management cluster VMs on Microsoft Azure, accept the license for their base Kubernetes version and machine OS.

1. Sign in to the Azure CLI as your `tce` client application.

   ```bash
   az login --service-principal --username AZURE_CLIENT_ID --password AZURE_CLIENT_SECRET --tenant AZURE_TENANT_ID
   ```

   Where `AZURE_CLIENT_ID`, `AZURE_CLIENT_SECRET`, and `AZURE_TENANT_ID` are your `tce` app's client ID and secret and your tenant ID, as recorded in [Register Tanzu Community Edition as an Azure Client App](azure-mgmt/#a-idtkg-appa-register-tanzu-community-edition-as-an-azure-client-app).

1. Run the `az vm image terms accept` command, specifying the `--plan` and your Subscription ID.

   In Tanzu Community Edition v1.3.1, the default cluster image `--plan` value is `k8s-1dot20dot5-ubuntu-2004`, based on Kubernetes version 1.20.5 and the  machine OS, Ubuntu 20.04. Run the following command:

   ```sh
   az vm image terms accept --publisher vmware-inc --offer tkg-capi --plan k8s-1dot20dot5-ubuntu-2004 --subscription AZURE_SUBSCRIPTION_ID
   ```

   Where `AZURE_SUBSCRIPTION_ID` is your Azure subscription ID.

You must repeat this to accept the base image license for every version of Kubernetes or OS that you want to use when you deploy clusters, and every time that you upgrade to a new version of Tanzu Community Edition.

## <a id="ssh-key"></a> Create an SSH Key Pair (Optional)

You will need OpenSSL installed locally, to create a new keypair or validate the download package thumbprint.  See [OpenSSL](https://www.openssl.org).

You deploy management clusters from a machine referred to as the _bootstrap machine_, using the Tanzu CLI.
To connect to Microsoft Azure, the bootstrap machine must provide the public key part of an SSH key pair. If your bootstrap machine does not already have an SSH key pair, you can use a tool such as `ssh-keygen` to generate one.

1. On your bootstrap machine, run the following `ssh-keygen` command.

   <pre>ssh-keygen -t rsa -b 4096 -C "<em>email@example.com</em>"</pre>
1. At the prompt `Enter file in which to save the key (/root/.ssh/id_rsa):` press Enter to accept the default.
1. Enter and repeat a password for the key pair.
1. Add the private key to the SSH agent running on your machine, and enter the password you created in the previous step.

   ```sh
   ssh-add ~/.ssh/id_rsa
   ```

1. Open the file `.ssh/id_rsa.pub` in a text editor so that you can easily copy and paste it when you deploy a management cluster.



