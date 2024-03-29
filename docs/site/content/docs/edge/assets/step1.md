This topic describes how to use the Tanzu installer interface to deploy a management cluster. The installer interface launches in a browser and takes you through the steps to configure the management. The input values are saved in: `~/.config/tanzu/tkg/clusterconfigs/cluster-config.yaml`.

## Before you begin

* [ ] Make sure that you have installed Tanzu Community Edition. See [Plan Your Install](../installation-planning)
* [ ] Make sure that you have completed steps to prepare to deploy a cluster. See [Plan Your Install](../installation-planning).
* [ ] Make sure you have met the following installer prerequisites:
  * [ ] NTP is running on the bootstrap machine on which you are running `tanzu management-cluster create` and on the hypervisor.
  * [ ] A DHCP server is available.
  * [ ] The host where the CLI is being run has unrestricted Internet access in order to pull down container images.
  * [ ] Docker is running.
* [ ] By default Tanzu saves the `kubeconfig` for all management clusters in the `~/.kube-tkg/config` file. If you want to save the `kubeconfig` file to a different location, set the `KUBECONFIG` environment variable before running the installer, for example:

   ```sh
   KUBECONFIG=/path/to/mc-kubeconfig.yaml
   ```
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

The Tanzu Installer opens, click the **Deploy** button for **VMware vSphere**, **AWS**, **Azure**, or **Docker**.

**Note**: If you are bootstrapping from a Windows machine and you encounter the following error, see this [troubleshooting entry](../faq-cluster-bootstrapping/#x509-certificate-signed-by-unknown-authority-when-deploying-management-cluster-from-windows) for a workaround.

```sh
Error: unable to ensure prerequisites: unable to ensure tkg BOM file: failed to download TKG compatibility file from the registry: failed to list TKG compatibility image tags: Get "https://projects.registry.vmware.com/v2/": x509: certificate signed by unknown authority
```

Complete the Installer steps as follows:
