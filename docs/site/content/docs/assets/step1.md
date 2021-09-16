This topic describes how to use the Tanzu installer interface to deploy a management cluster. The installer interface launches in a browser and takes you through the steps to configure the management. The input values are saved in: `~/.config/tanzu/tkg/clusterconfigs/cluster-config.yaml`.

## Before you begin

* [ ] Make sure that you have installed Tanzu Community Edition. See [Plan Your Install](installation-planning)
* [ ] Make sure that you have completed steps to prepare to deploy a cluster. See [Plan Your Install](installation-planning).
* [ ] Make sure you have met the following installer prerequisites:
  * [ ] NTP is running on the bootstrap machine on which you are running `tanzu management-cluster create` and on the hypervisor.
  * [ ] A DHCP server is available.
  * [ ] The host where the CLI is being run has unrestricted Internet access in order to pull down container images.
  * [ ] Docker is running.
* [ ] By default Tanzu saves the `kubeconfig` for all management clusters in the `~/.kube-tkg/config` file. If you want to save the `kubeconfig` file to a different location, set the `KUBECONFIG` environment variable before running the installer, for example:

   ```sh
   KUBECONFIG=/path/to/mc-kubeconfig.yaml
   ```

<!--- For production deployments, it is strongly recommended to enable identity management for your clusters. For information about the preparatory steps to perform before you deploy a management cluster, see [Enabling Identity Management in Tanzu Kubernetes Grid](enabling-id-mgmt.md).
* If you want to register your management cluster with Tanzu Mission Control, follow the procedure in [Register Your Management Cluster with Tanzu Mission Control](register_tmc.md).
* If you are deploying clusters in an internet-restricted environment to either vSphere or Amazon EC2, you must also perform the steps in [Deploying Tanzu Kubernetes Grid in an Internet-Restricted Environment](airgapped-environments.md).-->

<!--- **NOTE**: On vSphere with Tanzu, you do not need to deploy a management cluster. See [Use the Tanzu CLI with a vSphere with Tanzu Supervisor Cluster](../tanzu-k8s-clusters/connect-vsphere7.md).-->

## Procedure

### Start the Installer in your Browser

On the machine on which you downloaded and installed the Tanzu CLI, run the `tanzu management-cluster create` command with the `--ui` option:

```sh
tanzu management-cluster create --ui
```

If the prerequisites are met, the installer interface opens locally, at [http://127.0.0.1:8080](http://127.0.0.1:8080) in your default browser. To change where the installer interface runs, including running it on a different machine from the Tanzu CLI, use the following parameters:

* `--browser` specifies the local browser to open the interface in. Supported values are `chrome`, `firefox`, `safari`, `ie`, `edge`, or `none`. You can use `none` with `--bind` to run the interface on a different machine.
* `--bind` specifies the IP address and port to serve the interface from. For example, if another process is already using [http://127.0.0.1:8080](http://127.0.0.1:8080), use `--bind` to serve the interface from a different local port.

Example:

```sh
tanzu management-cluster create --ui --bind 192.168.1.87:5555 --browser none
```

The Tanzu Installer opens, click the **Deploy** button for **VMware vSphere**, **Amazon EC2**, **Azure**, or **Docker**.

<!--  ![Tanzu Kubernetes Grid installer interface welcome page with Deploy to vSphere button](../images/deploy-management-cluster.png)-->
Complete the Installer steps as follows:
