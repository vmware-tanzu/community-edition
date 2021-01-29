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
`
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
    kubectl get pods --all-namespaces=false

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

```sh
kubectl create namespace tanzu-extensions
kubectl create namespace kapp-controller
kubectl --namespace kapp-controller \
apply --file https://gist.githubusercontent.com/joshrosso/e6f73bee6ade35b1be5280be4b6cb1de/raw/b9f8570531857b75a90c1e961d0d134df13adcf1/kapp-controller-build.yaml
```

> This manifest points to a custom kapp-controller build where we've introducing
  imgpkg support.

## Installing extensions

In order to install extensions, you **must** have access to
https://github.com/vmware-tanzu/tce. If you cannot see this repository, ask to
be added in the (currently internal) #tce channel.

1. Get a [personal access
   token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token)
   from GitHub.

1. Register your token in `tanzu` CLI.

    ```sh
    tanzu extension token $PERSONAL_ACCESS_TOKEN
    ```

1. List the available extensions.

    ```sh
    tanzu extension list

    Extension: velero
    Extension: gatekeeper
    Extension: cert-manager
    Extension: contour
    ```

1. Download an extension using tanzu-cli.

    ```sh
    tanzu get gatekeeper
    ```

    > This puts the extension's App file in
    `$XDG_DATA_HOME/tanzu-repository/extensions/latest/gatekeeper`.

1. Install the extension to the cluster.

    ```sh
    tanzu install gatekeeper
    ```

1. Verify gatekeeper is installed in the cluster.

## How it works

The experience above was facilitated with a grouping of technologies including
tanzu CLI, imgpkg, kbld, and kapp-controller.

![january-tce-flow.png](./images/january-tce-flow.png)


