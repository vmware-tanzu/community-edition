# Deploy Tanzu Kubernetes Clusters with Different Kubernetes Versions

Tanzu Kubernetes Grid can create Tanzu Kubernetes clusters that run on:

- A Kubernetes version that Tanzu Kubernetes Grid ships with, including the default version that the management cluster runs, or
- A Kubernetes version that comes out after the current version of Tanzu Kubernetes Grid, and that VMware publishes a Bill of Materials (BoM) for in a public registry.

## <a id="k8s-vers-list"></a> List Available Versions

To list all available Kubernetes releases with their current compatibility and upgrade status, run `tanzu kubernetes-release get` with an optional version match argument, for example:

  - `tanzu kubernetes-release get`: list all releases
  - `tanzu kubernetes-release get v1.19`: list all releases matching `v1.19`
  - `tanzu kubernetes-release get v1.19.1+vmware.1`: list the `v1.19.1+vmware.1` release

Sample output:

  ```
  $ tanzu kubernetes-release get
   NAME                                       VERSION                                COMPATIBLE  UPGRADEAVAILABLE
   v1.17.16---vmware.2                        v1.17.16+vmware.2                      True        True
   v1.18.14---vmware.2                        v1.18.14+vmware.2                      True        True
   v1.18.16---vmware.2                        v1.18.16+vmware.2                      True        False
   v1.19.6---vmware.2                         v1.19.6+vmware.2                       True        True
   v1.19.8---vmware.2                         v1.19.8+vmware.2                       True        False
   v1.20.1---vmware.2                         v1.20.1+vmware.2                       True        True
   v1.20.4---vmware.2                         v1.20.4+vmware.2                       False       False
  ```

## <a id="k8s-upgrades-list"></a> List Available Upgrades

To list the available upgrades for a Kubernetes release, run `tanzu kubernetes-release available-upgrades get` with the full name of the version, for example:

  ```
  tanzu kubernetes-release available-upgrades get v1.19.6---vmware.2
  NAME                                     VERSION
  v1.19.8---vmware.2                       v1.19.8+vmware.2
  v1.19.9---vmware.2                       v1.19.9+vmware.2
  v1.20.1---vmware.2                       v1.20.1+vmware.2 
  v1.20.4---vmware.2                       v1.20.4+vmware.2
  v1.20.5---vmware.2                       v1.20.5+vmware.1
  ```

## <a id="k8s-vers-process"></a> How Tanzu Kubernetes Grid Updates Kubernetes Versions

Tanzu Kubernetes Grid manages Kubernetes versions as custom resources definition (CRD) objects called Tanzu Kubernetes releases (TKr).

A TKr controller periodically polls a public registry for new Kubernetes version BoM files.
When it detects a new version, it downloads the BoM and creates a corresponding TKr.
The controller then saves the new BoM and TKr in the management cluster, as a ConfigMap and custom resource, respectively.

The `tanzu` CLI queries the management cluster to list available Kubernetes versions.
When the CLI needs to create a new cluster, it downloads the TKr and BoM, uses the TKr to create the cluster, and saves the BoM to the local `~/.tanzu/tkg/bom` directory.

## <a id="non-default"></a> Deploy a Cluster with a Non-Default Kubernetes Version

Each release of Tanzu Kubernetes Grid provides a default version of Kubernetes. The default version for Tanzu Kubernetes Grid v1.3.1 is Kubernetes v1.20.5.

As upstream Kubernetes releases patches or new versions, VMware publishes them in a public registry and the Tanzu Kubernetes release controller imports them into the management cluster.
This lets the `tanzu` CLI create clusters based on the new versions.

To list available Kubernetes versions, see [Available Kubernetes Versions](#k8s-version) above.

To deploy clusters that run a non-default version of Kubernetes different from the default version, follow the steps below.

### <a id="k8s-version-publish"></a> Publish the Kubernetes Version to your Infrastructure

On vSphere and Azure, you need to take an additional step before you can deploy clusters that run non-default versions of Kubernetes:

* **vSphere**: Import the appropriate base image template OVA file into vSphere and convert it to a VM template. For information about importing base OVA files into vSphere, see [Import the Base Image Template into vSphere](../mgmt-clusters/vsphere.md#import-base-ova).

* **Azure**: Run the [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) command to accept the license for the base OS version. Once you have accepted a license, you can skip this step in the future:
   1. Convert your target Kubernetes version listed in the output of the `tanzu kubernetes-release get` command into its Azure image SKU as follows:
      * Change leading `v` to `k8s-`.
      * Change `.` to `dot` in the version number.
      * Change trailing `+vmware.*` to `-ubuntu-2004`, to designate Ubuntu v20.04, the default OS version for all Tanzu Kubernetes Grid VMs on Azure.
      * Examples: `k8s-1dot19dot8-ubuntu-2004`, `k8s-1dot20dot5-ubuntu-2004`.
   1. Run `az vm image terms accept`. For example:

      ```
      az vm image terms accept --publisher vmware-inc --offer tkg-capi --plan k8s-1dot20dot5-ubuntu-2004
      ```

* **Amazon EC2**: No action required.  The Amazon Linux 2 Amazon Machine Images (AMI) that includes the supported Kubernetes versions is publicly available to all Amazon EC2 users, in all supported AWS regions. Tanzu Kubernetes Grid automatically uses the appropriate AMI for the Kubernetes version that you specify.

### <a id="k8s-version-deploy"></a> Deploy the Kubernetes Cluster

To deploy a Tanzu Kubernetes cluster with a version of Kubernetes that is not the default for your Tanzu Kubernetes Grid release, specify the version in the `--tkr` option.

- Deploy a Kubernetes v1.19.1 cluster to vSphere:

    ```
    tanzu cluster create my-1-19-1-cluster --tkr v1.19.1---vmware.1-tkg.1-60d2ffd
    ```

- Deploy a Kubernetes v1.19.1 cluster to Amazon EC2 or Azure:

    ```
    tanzu cluster create my-1-19-1-cluster --tkr v1.19.1---vmware.1-tkg.1-60d2ffd
    ```

For more details on how to create a Tanzu Kubernetes cluster, see [Deploy Tanzu Kubernetes Clusters](deploy.md).

## <a id="alt-os-custom"></a> Deploy a Cluster with an Alternate OS or Custom Machine Image

With out-of-the-box Tanzu Kubernetes Grid, the `--tkr` option to `tanzu cluster create` supports common Kubernetes versions running on common base machine OSes.
But you can build custom machine images and TKr to create new clusters with.

Reasons to do this include:

- To create clusters on a base OS that VMware supports but does not distribute, such as Red Hat Enterprise Linux (RHEL) v7.
- To install additional packages into the base machine image, or otherwise customize it as described in [Customization](https://image-builder.sigs.k8s.io/capi/capi.html#customization) in the Image Builder documentation.

To deploy a cluster with an alternate OS or custom machine image, you build a custom image, create a TKr for it, and deploy clusters with it as described in [Building Machine Images](../build-images/index.md).
