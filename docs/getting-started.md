# Getting Started with TCE

The initial release of TCE leverages both the `tkg` CLI and `tanzu` CLI.
Currently, TKG is working to move the cluster management functionality as a
plugin to `tanzu` CLI. When this happens, there will no longer be a need for
`tkg` CLI.

## Installing Tanzu CLI

Please note, TCE currently work on **Mac** and **Linux**.

1. Download the release.

    ```sh
    wget TODO:githubURL
    ```

1. Unpack the release.

    ```sh
    tar xzvf tce-0.1.0.tar.gz
    ```

1. Run the install script.

    ```sh
    ./install.sh
    ```

    > This installs the tanzu-cli and puts all the plugins in their proper
    location.

## Creating a Kubernetes Cluster

1. Initialize the TKG kickstart UI.

    ```sh
    tkg init --ui
    ```

1. Go through the installation process for your target platform.

    > You can find the full TKG docs
      [here](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.2/vmware-tanzu-kubernetes-grid-12/GUID-mgmt-clusters-deploy-management-clusters.html).
      Once tanzu CLI contains the functionality for cluster boostrapping, we'll
      include docs on getting started here.

1. Create a guest cluster with TKG.

    ```
    tkg create cluster $CLUSTERNAME --plan=dev
    ```

    > Default plans are `dev` and `prod`.

1. Once the cluster starts, get the credentials.

    ```sh
    tkg get credentials $CLUSTERNAME
    ```

1. Set your `kubectl` context accordingly.

    ```sh
    kubectl config use-context ${CLUSTERNAME}-admin@${CLUSTERNAME}
    ```

1. Verify you can see pods in the cluster.

    ```sh
    kubectl get po -A

    NAMESPACE     NAME                                                    READY   STATUS    RESTARTS   AGE
    kube-system   antrea-agent-9d4db                                      2/2     Running   0          3m42s
    kube-system   antrea-agent-vkgt4                                      2/2     Running   1          5m48s
    kube-system   antrea-controller-5d594c5cc7-vn5gt                      1/1     Running   0          5m49s
    kube-system   coredns-5d6f7c958-hs6vr                                 1/1     Running   0          5m49s
    kube-system   coredns-5d6f7c958-xf6cl                                 1/1     Running   0          5m49s
    kube-system   etcd-tce-guest-control-plane-b2wsf                      1/1     Running   0          5m56s
    kube-system   kube-apiserver-tce-guest-control-plane-b2wsf            1/1     Running   0          5m56s
    kube-system   kube-controller-manager-tce-guest-control-plane-b2wsf   1/1     Running   0          5m56s
    kube-system   kube-proxy-9825q                                        1/1     Running   0          5m48s
    kube-system   kube-proxy-wfktm                                        1/1     Running   0          3m42s
    kube-system   kube-scheduler-tce-guest-control-plane-b2wsf            1/1     Running   0          5m56s
    kube-system   kube-vip-tce-guest-control-plane-b2wsf                  1/1     Running   0          5m56s
    kube-system   vsphere-cloud-controller-manager-nwrg4                  1/1     Running   2          5m48s
    kube-system   vsphere-csi-controller-5b6f54ccc5-trgm4                 5/5     Running   0          5m49s
    kube-system   vsphere-csi-node-drnkz                                  3/3     Running   0          5m48s
    kube-system   vsphere-csi-node-flszf                                  3/3     Running   0          3m42s
    ```

## Install kapp-controller

    ```
    kubectl create ns kapp-controller
    kubectl --namespace kapp-controller apply --file https://github.com/k14s/kapp-controller/releases/latest/download/release.yml
    ```

### Quick Start with Velero

The following steps guide you through installing a sample extension, Velero, into your TCE cluster.

1. Create the `tanzu-extensions` namespace:

    ```shell
    kubectl create ns tanzu-extensions
    ```

2. Deploy kapp controller into the target cluster:

    ```shell
    # TODO: deploy to tanzu-extensions namespace instead?
    kubectl create ns kapp-controller
    kubectl -n kapp-controller apply -f https://github.com/k14s/kapp-controller/releases/latest/download/release.yml
    ```

3. Validate that the kapp-controller started successfully:

    ```shell
    $ kubectl -n kapp-controller get deployment
    NAME              READY   UP-TO-DATE   AVAILABLE   AGE
    kapp-controller   1/1     1            1           8m34s
    ```

4. Deploy Velero extension:

    ```shell
    kubectl -n tanzu-extensions apply -f extensions/velero/extension.yaml
    ```

5. Validate the extension was deployed successfully by checking the App's description.

    ```shell
    $ kubectl -n tanzu-extensions get app velero
    NAME     DESCRIPTION           SINCE-DEPLOY   AGE
    velero   Reconcile succeeded   2m13s          2m47s
    ```

    If the description field shows an error, use the `kubectl describe` command to troubleshoot further:

    ```shell
    kubectl -n tanzu-extensions describe app velero
    ```

See the Velero extension [documentation](./extensions/velero) for more information, including a walkthrough that guides you through a usage example.

### Extension Structure

Each extension lives in its own directory within the `extensions/` directory.

An extension is composed of the following assets:

* `extension.yaml`: Contains the [kapp-controller](https://github.com/vmware-tanzu/carvel-kapp-controller) `App` resource. This file also includes a Service Account and a ClusterRoleBinding for kapp-controller to use to deploy the extension.
* `config/` directory: Contains the Kubernetes deployment manifests necessary to deploy the extension. Manifests can be templatized with [`ytt`](https://github.com/vmware-tanzu/carvel-ytt), when necessary.

### Extension Delivery

The deployment manifests of an extension are bundled into an OCI image using [imgpkg](https://github.com/vmware-tanzu/carvel-imgpkg). The OCI image is pushed to a container registry for consumption.

### Extension Deployment

Assuming you have a TCE cluster up and running, the extension deployment process works as follows:

1. Install the `extension.yaml` file of a given extension. This initiates the App resource reconciliation.
2. `kapp-controller` fetches the OCI image that contains the extension's deployment manifests
3. `kapp-controller` renders the templates using `ytt`
4. `kapp-controller` deploys the extension

### Creating a new Extension

To create a new extension, you must:

1. Create a new directory within the `extensions/` directory. The name of the extension must be the name of the directory.
2. Create an `extension.yaml` file which contains a kapp-controller App custom resource.
3. Create a `config` directory which contains the extension's deployment manifests.
4. Templatize the deployment manifests, when necessary.
5. Bundle the deployment manifests into an OCI image using [imgpkg](https://github.com/vmware-tanzu/carvel-imgpkg) and push to the OCI registry.
6. Test the extension.
