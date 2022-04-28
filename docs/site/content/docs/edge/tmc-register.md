# Register a Management Cluster in Tanzu Mission Control

VMware Tanzu Mission Control is a licensed VMware product that you can integrate with.

Use Tanzu Mission Control to manage your entire Kubernetes footprint, regardless of where your clusters reside.
Bring all your Kubernetes clusters together to monitor and manage them from a single console backed by an API-driven service.

## Before You Begin

Note:

- You can register a management cluster in Tanzu Mission Control that was deployed on one of the following Cloud Infrastructure Providers: AWS, Azure, vSphere.

- For clusters deployed on vSphere, Tanzu Mission Control can only support clusters that use the Photon OS image. Tanzu Mission Control will not work properly with vSphere clusters using the Ubuntu image.

Complete the following steps to register a Tanzu Community Edition management cluster in Tanzu Mission Control.

1. Log into Tanzu Mission Control, for more information, see [Log In to the Tanzu Mission Control Console](https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/services/tanzumc-using/GUID-855A8998-E19A-46AC-A833-12C347486EF7.html) in the Tanzu Mission Control documentation.

## Procedure

1. On the **Management clusters** tab of the **Administration** page, click **REGISTER MANAGEMENT CLUSTERS** and select **Tanzu Kubernetes Grid**. This launches a wizard.

1. Name and assign: Enter a name and group for your management cluster and click **Next**.

1. Proxy Configuration: **Ensure Set proxy for management cluster** is set to No and click **Next**.

1. Register: Copy the registration URL.

1. Create a new file called `tmc-registration.yaml` and paste the contents of the Tanzu Mission Control registration URL from the previous step into this new file.

    ```sh
    vim tmc-registration.yaml
    ```

1. Before you proceed, ensure your context is set to the management cluster. Run:

    ```sh
    tanzu mc kubeconfig get --admin
    ```

    Copy the run the output, for example

    ```sh
    kubectl config  use-context mgmttest-admin@mgmttest
    ```

1. Now that you are interacting with the management cluster, apply the `tmc-register.yaml` file:

    ```sh
    apply -f tmc-register.yaml
    ```

1. (Optional)  To watch the Tanzu Mission Control pods get created, run:

    ```sh
    kubectl get po -A
    ```

    From the output copy the `vmware-system-tmc` namespace and run

    ```sh
    kubectl -n vmware-system-tmc get po
    ```

1. Once all the pods associated with Tanzu Mission Control are running, in the Tanzu Mission Control console, you should be able to click **View the Management Cluster** and click **Verify Connection**.

For more information on how to work with Tanzu Mission Control, see the [VMware Tanzu Mission Control Documentation](https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/services/tanzumc-concepts/GUID-9E6DEA00-C368-4B06-B93E-BA1916EB2929.html).
