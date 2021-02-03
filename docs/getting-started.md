# Getting Started with TCE

The initial release of TCE leverages both the `tkg` CLI and `tanzu` CLI.
Currently, TKG is working to move the cluster management functionality as a
plugin to `tanzu` CLI. When this happens, there will no longer be a need for
`tkg` CLI.

🚨🚨🚨

**Thank you for trying Tanzu Community Edition! Please be sure to [leave
feedback
here](https://github.com/vmware-tanzu/tce/issues/new?assignees=&labels=feedback&template=feedback-on-tanzu-community-edition-template.md&title=)
after trying this guide!**

🚨🚨🚨

## Installing Tanzu Command Line Interface

Please note, TCE currently works on **macOS** and **Linux**.

In order to TCE, you **must** have access to
https://github.com/vmware-tanzu/tce. If you cannot see this repository, ask to
be added in the (currently internal) #tce channel.

1. Download the release.

    Make sure you're logged into Github and then go to the [TCE Releases](https://github.com/vmware-tanzu/tce/releases/tag/v0.1.0) page and download the Tanzu CLI for either

    * [Linux](https://github.com/vmware-tanzu/tce/releases/download/v0.1.0/dist-linux-v0.1.0.tar.gz), or
    * [Mac](https://github.com/vmware-tanzu/tce/releases/download/v0.1.0/dist-mac-v0.1.0.tar.gz).

1. Unpack the release.

    **linux**

    ```sh
    tar xzvf ~/Downloads/dist-linux-v0.1.0.tar.gz
    ```

    **macOS**

    ```sh
    tar xzvf ~/Downloads/dist-mac-v0.1.0.tar.gz
    ```

1. Run the install script.

    **linux**

    ```sh
    cd dist-linux
    ./install.sh
    ```

    **macOS**

    ```sh
    cd dist-mac
    ./install.sh
    ```

    > This installs the `tanzu` CLI and puts all the plugins in their proper
    location.

## Creating a Kubernetes Cluster

1. Initialize the TKG kickstart UI.

    ```sh
    tkg init --ui
    ```

1. Go through the installation process for your target platform.

    > You can find the full TKG docs
      [here](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.2/vmware-tanzu-kubernetes-grid-12/GUID-mgmt-clusters-deploy-management-clusters.html).
      Once `tanzu` CLI contains the functionality for cluster bootrapping, we'll
      include docs on getting started here.

1. Create a guest cluster with TKG.

    ```sh
    export CLUSTERNAME=<My new cluster name>
    tkg create cluster ${CLUSTERNAME} --plan=dev
    ```

    > Default plans are `dev` and `prod`.

1. Once the cluster starts, get the credentials.

    ```sh
    tkg get credentials ${CLUSTERNAME}
    ```

1. Set your `kubectl` context accordingly.

    ```sh
    kubectl config use-context ${CLUSTERNAME}-admin@${CLUSTERNAME}
    ```

1. Verify you can see pods in the cluster.

    ```sh
    kubectl get pods --all-namespaces

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

> This manifest points to a custom kapp-controller build where we've introduced
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
    tanzu extension token <My GitHub Personal Access Token>
    ```

1. List the available extensions.

    ```sh
    tanzu extension list

    Extension: velero
    Extension: gatekeeper
    Extension: cert-manager
    Extension: contour
    ```


1. Install the extension to the cluster.

    ```sh
    tanzu extension install gatekeeper
    ```

1. Verify gatekeeper is installed in the cluster.

    ```sh
    kubectl -n gatekeeper-system get all

    NAME                                               READY   STATUS    RESTARTS   AGE
    pod/gatekeeper-audit-65584c8875-qwfz8              1/1     Running   0          109s
    pod/gatekeeper-controller-manager-f7556dc9-6mtpl   1/1     Running   0          109s

    NAME                                 TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
    service/gatekeeper-webhook-service   ClusterIP   100.66.61.43   <none>        443/TCP   109s

    NAME                                            READY   UP-TO-DATE   AVAILABLE   AGE
    deployment.apps/gatekeeper-audit                1/1     1            1           109s
    deployment.apps/gatekeeper-controller-manager   1/1     1            1           109s

    NAME                                                     DESIRED   CURRENT   READY   AGE
    replicaset.apps/gatekeeper-audit-65584c8875              1         1         1       109s
    replicaset.apps/gatekeeper-controller-manager-f7556dc9   1         1         1       109s
    ```

## How it works

The experience above was facilitated with a grouping of technologies including
`tanzu` CLI, [imgpkg](https://carvel.dev/imgpkg/), [kbld](https://carvel.dev/kbld/), and [kapp-controller](https://github.com/vmware-tanzu/carvel-kapp-controller).

![january-tce-flow.png](./images/january-tce-flow.png)

To see the capturing off the App CR, the following command may be run.

1. Download an extension using `tanzu` CLI.

    ```sh
    tanzu extension get gatekeeper
    ```

    > This puts the extension's App file in
    `$XDG_DATA_HOME/tanzu-repository/extensions/latest/gatekeeper`.
