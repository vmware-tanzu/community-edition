# Register Your Management Cluster with Tanzu Mission Control

You can register management clusters with Tanzu Mission Control. Registering management clusters with Tanzu Mission Control allows you to provision and manage workload clusters in the Tanzu Mission Control dashboard interface.

## Prerequisites

To register your Tanzu Kubernetes Grid management cluster with Tanzu Mission Control, you must be a member of VMware Cloud Services organization that has access to Tanzu Mission Control.

For more information, see [Getting Started with VMware Tanzu Mission Control](https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/services/tanzumc-getstart/GUID-6BCCD353-CE6A-494B-A1E4-72304DC9FA7F.html) in the Tanzu Mission Control documentation.

Management clusters that you register in Tanzu Mission Control must be **production** clusters with multiple control plane nodes. This configuration allows Tanzu Mission Control to support complete lifecycle management for workload clusters that are managed by the management cluster.

For more information, see [Requirements for Registering a Tanzu Kubernetes Cluster with Tanzu Mission Control](https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/services/tanzumc-concepts/GUID-3AE5F733-7FA7-4B34-8935-C25D41D15EF9.html) in the Tanzu Mission Control documentation.

### Supported Providers

You can register management clusters that are deployed on Azure, Amazon EC2 or vSphere. Registering standalone clusters is not supported.

For a list of currently supported providers, see [Requirements for Registering a Tanzu Kubernetes Cluster with Tanzu Mission Control](https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/services/tanzumc-concepts/GUID-3AE5F733-7FA7-4B34-8935-C25D41D15EF9.html) in the Tanzu Mission Control documentation.

## Procedure

To register your management cluster, perform the following steps:

   1. Obtain a Tanzu Mission Control registration URL by following the steps in [Register a Management Cluster with Tanzu Mission Control](https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/services/tanzumc-using/GUID-EB507AAF-5F4F-400F-9623-BA611233E0BD.html).

   1. If you are deploying a new management cluster, you can add the registration URL in the Installer Interface or by using the Tanzu CLI.
    - **Installer Interface**: Copy and paste the URL you obtained from Tanzu Mission Control into the **Registration URL** field of the **Register with Tanzu Mission Control** configuration pane. For more information, see [Deploy Management Clusters with the Installer Interface](deploy-ui.md).
    - **Tanzu CLI**: Copy and paste the URL you obtained from Tanzu Mission Control into the <code>TMC-REGISTRATION-URL</code> configuration variable in your management cluster's configuration file. This configuration is applied when you run `tanzu management-cluster create`. For more information, see [Deploy Management Clusters from a Configuration File](deploy-cli.md).

   1. If you want to register an already deployed Tanzu Kubernetes Grid management cluster, you can use one of the following commands:

    ```
    tanzu management-cluster register --tmc-registration-url "TMC-REGISTRATION-URL"
    ```

    The registration URL must be contained within quotes. For example:

    ```
    tanzu management-cluster register --tmc-registration-url "https://tmc-org.cloud.vmware.com/installer?id=9448627322axe82e2fb042f84517710390d02c9e677f09199a36e2cff659859e&source=registration"
    ```

    If you want to skip the interactive prompt that asks for registration confirmation with Tanzu Mission Control, specify the `-y` or `--yes` flag.

    ```
    tanzu management-cluster register --tmc-registration-url "TMC-REGISTRATION-URL" --yes
    ```

    Alternately, you can also use `kubectl`.

    ```
    kubectl apply -f "TMC-REGISTRATION-URL"
    ```

    For example:

    ```
    kubectl apply -f "https://tmc-org.cloud.vmware.com/installer?id=9448627322axe82e2fb042f84517710390d02c9e677f09199a36e2cff659859e&source=registration"
    ```

    The commands create a namespace called `vmware-system-tmc` and install the Tanzu Mission Control cluster agent on the management cluster. The installation process may take a few minutes.

   1. (Optional) After you successfully register a management cluster, you can add any existing Tanzu Kubernetes clusters that are currently managed by the management cluster to Tanzu Mission Control. To manage these clusters in Tanzu Mission Control, see [Add a Workload Cluster into Tanzu Mission Control Management](https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/services/tanzumc-using/GUID-78908829-CB4E-459F-AA81-BEA415EC9A11.html) in the Tanzu Mission Control documentation.

## What's Next

After you register your management cluster, you can use Tanzu Mission Control to deploy, manage, and monitor your Tanzu Kubernetes clusters.

For more information about how Tanzu Mission Control allows you to manage Tanzu Kubernetes clusters, see [Managing the Lifecycle of Tanzu Kubernetes Clusters](
https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/services/tanzumc-using/GUID-1F847180-1F98-4F8F-9062-46DE9AD8F79D.html) in the Tanzu Mission Control documentation.
