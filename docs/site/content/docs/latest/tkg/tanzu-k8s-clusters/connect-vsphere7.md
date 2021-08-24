# Use the Tanzu CLI with a vSphere with Tanzu Supervisor Cluster

You can use Tanzu CLI with a vSphere with Tanzu Supervisor Cluster that is running in a vSphere 7 instance. In this way, you can deploy Tanzu Kubernetes clusters to vSphere with Tanzu and manage their lifecycle directly from the Tanzu CLI.

vSphere with Tanzu provides a vSphere Plugin for `kubectl`. The vSphere Plugin for `kubectl` extends the standard `kubectl` commands so that you can connect to the Supervisor Cluster from `kubectl` by using vCenter Single Sign-On credentials. Once you have installed the vSphere Plugin for `kubectl`, you can connect the Tanzu CLI to the Supervisor Cluster. Then, you can use the Tanzu CLI to deploy and manage Tanzu Kubernetes clusters running in vSphere.

**NOTE**: On VMware Cloud on AWS and Azure VMware Solution, you cannot create a Supervisor Cluster, and need to deploy a management cluster to run `tanzu` commands.

## <a id="prereqs"></a> Prerequisites

- Perform the steps described in [Install the Tanzu CLI and Other Tools](../install-cli.md).
- Make sure that you have a vSphere account that has the correct permissions for deployment of clusters to vSphere 7. For information about how to create a user account, see [Required Permissions for the vSphere Account](../mgmt-clusters/vsphere.html#vsphere-permissions). Alternatively, you can use a vSphere with Tanzu DevOps account. For information about the DevOps user role, see [vSphere with Tanzu User Roles and Workflows](https://docs.vmware.com/en/VMware-vSphere/7.0/vmware-vsphere-with-tanzu/GUID-A6EFA324-7C58-4DB4-9166-CE18462A8CCF.html).
- You have access to a vSphere 7 instance on which the vSphere with Tanzu feature is enabled.
- Download and install the `kubectl vsphere` CLI utility on the bootstrap machine on which you run Tanzu CLI commands.

   For information about how to obtain and install the vSphere Plugin for `kubectl`, see <a href="https://docs.vmware.com/en/VMware-vSphere/7.0/vmware-vsphere-with-kubernetes/GUID-0F6E45C4-3CB1-4562-9370-686668519FCA.html" target="_blank">Download and Install the Kubernetes CLI Tools for vSphere</a> in the vSphere with Tanzu documentation.

## <a id="add"></a> Step 1: Add the Supervisor Cluster

Connect to the supervisor cluster and add it as a management cluster to the `tanzu` CLI:

1. From vCenter **Hosts and Clusters** view, in the left column, expand the nested Datacenter, the vCenter cluster that hosts the supervisor cluster, and its **Namespaces** object.

1. Under **Namespaces**, select the namespace <img src="../images/namespace-line.png" alt="Clarity namespace icon" width="18" height="18"/> containing or adjacent to the three **SupervisorControlPlaneVM** instances.  In the main pane, select the **Summary** tab.

1. Under **Summary** > **Status** > **Link to CLI Tools** click **Copy link** and record the URL, for example `https://192.168.123.3`. Remove the `https://` to obtain the supervisor cluster API endpoint, `SUPERVISOR_IP` below, which serves as the download page for the Kubernetes CLI tools.

1. On the bootstrap machine, run the `kubectl vsphere login` command to log in to vSphere 7 with your vCenter Single Sign-On user account.

   Specify a vCenter Single Sign-On user account with the required privileges for Tanzu Kubernetes Grid operation, and the virtual IP (VIP) address for the control plane of the supervisor cluster. For example:

   ```
   kubectl vsphere login --vsphere-username administrator@vsphere.local --server=SUPERVISOR_IP --insecure-skip-tls-verify=true
   ```

1. Enter the password you use to log in to your vCenter Single Sign-On user account.

   When you have successfully logged in, `kubectl vsphere` displays all of the contexts to which you have access. The list of contexts should include the IP address of the supervisor cluster.

1. Set the context of `kubectl` to the supervisor cluster.

   ```
   kubectl config use-context SUPERVISOR_IP
   ```

1. Collect information to run the `tanzu login` command, which adds the supervisor cluster to your Tanzu Kubernetes Grid instance:

  - Decide on a name for the `tanzu` CLI to use for the supervisor cluster, serving as a Tanzu Kubernetes Grid management cluster.
  - The path to the local management cluster `kubeconfig` file, which defaults to `~/.kube/config` and is set by the `KUBECONFIG` environment variable.
  - The context of the supervisor cluster, which is the same as `SUPERVISOR_IP`.

1. Run the `tanzu login` command, passing in the values above.  

   In the example below, the `KUBECONFIG_PATH` defaults to `~/.kube/config` if the `KUBECONFIG` env variable is not set.

   ```
   $ tanzu login --name my-super --kubeconfig <KUBECONFIG_PATH> --context 10.161.90.119
   ✔  successfully logged in to management cluster using the kubeconfig my-super
   ```
  
1. Check that the supervisor cluster was added by running `tanzu login` again.

   The supervisor cluster should be listed by the name that you provided in the preceding step:

   ```
   tanzu login
   ? Select a server  [Use arrows to move, type to filter]
   > my-vsphere-mgmt-cluster  ()
     my-aws-mgmt-cluster      ()
     SUPERVISOR_IP            ()
     + new server
   ```

## <a id="config"></a> Step 2: Configure Cluster Parameters

Configure the Tanzu Kubernetes clusters that the `tanzu` CLI calls the supervisor cluster to create:

1. Obtain information about the storage classes that are defined in the supervisor cluster.

   ```
   kubectl get storageclasses
   ```

1. Set variables to define the storage classes, VM classes, service domain, namespace, and other required values with which to create your cluster. For information about all of the configuration parameters that you can set when deploying Tanzu Kubernetes clusters to vSphere with Tanzu, see <a href="https://docs.vmware.com/en/VMware-vSphere/7.0/vmware-vsphere-with-kubernetes/GUID-4E68C7F2-C948-489A-A909-C7A1F3DC545F.html" target="_blank">Configuration Parameters for Provisioning Tanzu Kubernetes Clusters</a> in the vSphere with Tanzu documentation.

  The following table lists the required variables:

<table width="100%" border="0">
  <tbody>
    <tr>
      <th scope="col" width="30%">Option</th>
      <th scope="col" width="30%">Value Type or Example</th>
      <th scope="col" width="40%">Description</th>
    </tr>
    <tr>
      <td><code>CONTROL_PLANE_STORAGE_CLASS</code></td>
      <td rowspan=2>Value returned from CLI: <code>kubectl get storageclasses</code></td>
      <td>Default storage class for control plane nodes</td>
    </tr>
    <tr>
      <td><code>WORKER_STORAGE_CLASS</code></td>
      <td>Default storage class for worker nodes</td>
    </tr>
    <tr>
      <td><code>DEFAULT_STORAGE_CLASS</code></td>
      <td>Empty string <code>""</code> for no default, or value from CLI, as above.</td>
      <td>Default storage class for control plane or workers</td>
    </tr>
    <tr>
      <td><code>STORAGE_CLASSES</code></td>
      <td>Empty string <code>""</code> lets clusters use any storage classes in the namespace, or comma-separated list string of values from CLI, <code>"SC-1,SC-2,SC-3"</code></td>
      <td>Storage classes available for node customization</td>
    </tr>
    <tr>
      <td><code>CONTROL_PLANE_VM_CLASS</code></td>
      <td rowspan=2>A standard VM class for vSphere with Tanzu, for example <code>guaranteed-large</code>.<br />
      See <a href="https://docs.vmware.com/en/VMware-vSphere/7.0/vmware-vsphere-with-kubernetes/GUID-7351EEFF-4EF0-468F-A19B-6CEA40983D3D.html">Virtual Machine Class Types for Tanzu Kubernetes Clusters</a> in the vSphere with Tanzu documentation.</td>
      <td>VM class for control plane nodes</td>
    </tr>
    <tr>
      <td><code>WORKER_VM_CLASS</code></td>
      <td>VM class for worker nodes</td>
    </tr>
    <tr>
      <td><code>SERVICE_CIDR</code></td>
      <td>CIDR range</td>
      <td>The CIDR range to use for the Kubernetes services. The recommended range is <code>100.64.0.0/13</code>. Change this value only if the recommended range is unavailable.</td>
    </tr>
    <tr>
      <td><code>CLUSTER_CIDR</code></td>
      <td>CIDR range</td>
      <td>The CIDR range to use for pods. The recommended range is <code>100.96.0.0/11</code>. Change this value only if the recommended range is unavailable.</td>
    </tr>
    <tr>
      <td><code>SERVICE_DOMAIN</code></td>
      <td>Domain</td>
      <td>e.g. <code>my.example.com</code>, or <code>cluster.local</code> if no DNS.  If you are going to assign FQDNs with the nodes, DNS lookup is required.</td>
    </tr>
        <tr>
      <td><code>NAMESPACE</code></td>
      <td>Namespace</td>
      <td>The namespace in which to deploy the cluster. </td>
    </tr>
    <tr>
      <td><code>CLUSTER_PLAN</code></td>
      <td><code>dev</code>, <code>prod</code>, or a custom plan</td>
      <td rowspan=2>See <a href="../tanzu-config-reference.md#all_iaases">Tanzu CLI Configuration File Variable Reference</a> for variables required for all Tanzu Kubernetes cluster configuration files.</td>
    </tr>
    <tr>
      <td><code>INFRASTRUCTURE_PROVIDER</code></td>
      <td><code>tkg-service-vsphere</code></td>
    </tr>
  </table>

   You can set the variables above by doing either of the following:

   - Include them in the cluster configuration file passed to the `tanzu` CLI `--file` option.  For example:

      ```
      CONTROL_PLANE_VM_CLASS: guaranteed-large
      ```

   - From command line, set them as local environment variables by running `export` (on Linux and macOS) or `SET` (on Windows) on the command line.  For example:

      ```
      export CONTROL_PLANE_VM_CLASS=guaranteed-large
      ```

      **Note:** If you want to configure unique proxy settings for a Tanzu Kubernetes cluster, you can set `TKG_HTTP_PROXY`, `TKG_HTTPS_PROXY`, and `NO_PROXY` as environment variables and then use the Tanzu CLI to create the cluster. These variables take precedence over your existing proxy configuration in vSphere with Tanzu.

## <a id="create"></a> Step 3: Create a Cluster

Run `tanzu cluster create` to create a Tanzu Kubernetes cluster.

1. Determine the versioned Tanzu Kubernetes release (TKr) for the cluster:

  1. Obtain the list of TKr that are available in the supervisor cluster.

      ```
      tanzu kubernetes-release get 
      ```

  1. From the command output, record the desired value listed under `NAME`, for example `v1.18.9---vmware.1-tkg.1.a87f261`.  The `tkr` `NAME` is the same as its `VERSION` but with `+` changed to `---`.

1. Determine the namespace for the cluster.

  1. Obtain the list of namespaces.

      ```
      kubectl get namespaces
      ```

  1. From the command output, record the namespace that includes the Supervisor cluster, for example `test-gc-e2e-demo-ns`.

1. Decide on the cluster plan: `dev`, `prod`, or a custom plan.

  - You can customize or create cluster plans with files in the `~/.tanzu/tkg/providers/infrastructure-tkg-service-vsphere` directory.
  See [Configure Tanzu Kubernetes Plans and Clusters](config-plans.md) for details.

1. Run `tanzu cluster create` with the namespace and `tkr` `NAME` values above to create a Tanzu Kubernetes cluster:

   ```
   tanzu cluster create my-vsphere7-cluster --tkr=TKR-NAME
   ```

```
#! ---------------------------------------------------------------------
#! Settings for creating clusters on vSphere with Tanzu
#! ---------------------------------------------------------------------
#! Identifies the storage class to be used for storage of the disks that store the root file systems of the worker nodes.
CONTROL_PLANE_STORAGE_CLASS:
#! Specifies the name of the VirtualMachineClass that describes the virtual
#! hardware settings to be used each control plane node in the pool.
CONTROL_PLANE_VM_CLASS:
#! Specifies a named storage class to be annotated as the default in the
#! cluster. If you do not specify it, there is no default.
DEFAULT_STORAGE_CLASS:
#! Specifies the service domain for the cluster
SERVICE_DOMAIN:
#! Specifies named persistent volume (PV) storage classes for container
#! workloads. Storage classes associated with the Supervisor Namespace are
#! replicated in the cluster. In other words, each storage class listed must be
#! available on the Supervisor Namespace for this to be a valid value
STORAGE_CLASSES:
#! Identifies the storage class to be used for storage of the disks that store the root file systems of the worker nodes.
WORKER_STORAGE_CLASS:
#! Specifies the name of the VirtualMachineClass that describes the virtual
#! hardware settings to be used each worker node in the pool
WORKER_VM_CLASS:
NAMESPACE: 
```

## <a id="what-next"></a> What to Do Next

You can now use the Tanzu CLI to deploy more Tanzu Kubernetes clusters to the vSphere with Tanzu Supervisor Cluster. You can also use the Tanzu CLI to manage the lifecycles of clusters that are already running there. For information about how to manage the lifecycle of clusters, see the other topics in [Deploying Tanzu Kubernetes Clusters](index.md).
