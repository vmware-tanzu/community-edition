# Getting Started with TCE

This guide walks you through standing up a management and guest cluster using
Tanzu CLI. It then demonstrates how you can deploy add-ons into the cluster.
Currently we have getting started guides for [vSphere](#vsphere) and
[AWS](#aws). For detailed documentation on tanzu-cli and deployment of clusters,
see the [TKG docs](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/index.html).

🚨🚨🚨

**Thank you for trying Tanzu Community Edition! Please be sure to [leave
feedback
here](https://github.com/vmware-tanzu/tce/issues/new?assignees=&labels=feedback&template=feedback-on-tanzu-community-edition-template.md&title=)
after trying this guide!**

🚨🚨🚨

## Prerequistites

Please note, TCE currently works on **macOS** and **Linux** AMD64 (also known as
x64) environments. Windows and other architectures may be supported in the
future.

The Docker runtime is required on the deployment machine, regardless of your
final deployment environment. Before proceeding, please ensure [Docker has
been installed](https://docs.docker.com/engine/install/).

## CLI Installation

1. Download the release.

    Make sure you're logged into GitHub and then go to the [TCE Releases](https://github.com/vmware-tanzu/tce/releases/tag/v0.2.0) page and download the Tanzu CLI for either

    * [Linux](https://github.com/vmware-tanzu/tce/releases/download/v0.2.0/tce-linux-amd64-v0.2.0.tar.gz)
    * [Mac](https://github.com/vmware-tanzu/tce/releases/download/v0.2.0/tce-darwin-amd64-v0.2.0.tar.gz)

1. Unpack the release.

    **linux**

    ```sh
    tar xzvf ~/Downloads/tce-linux-amd64-v0.2.0.tar.gz
    ```

    **macOS**

    ```sh
    tar xzvf ~/Downloads/tce-darwin-amd64-v0.2.0.tar.gz
    ```

1. Run the install script (make sure to use the appropriate directory for your platform).

    **linux**

    ```sh
    cd tce-linux-amd64-v0.2.0
    ./install.sh
    ```

    **macOS**

    ```sh
    cd tce-darwin-amd64-v0.2.0
    ./install.sh
    ```

    > This installs the `tanzu` CLI and puts all the plugins in their proper location.
    > The first time you run the `tanzu` command the installed plugins and plugin repositories will be initialized. This action might take a minute.

1. If you wish to run commands against any of the Kubernetes clusters that are created, you will need to download and install `kubectl`.

    **linux**

    ```sh
    curl -LO https://dl.k8s.io/release/v1.20.1/bin/linux/amd64/kubectl
    sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
    ```

    **macOS**

    ```sh
    curl -LO https://dl.k8s.io/release/v1.20.1/bin/darwin/amd64/kubectl
    sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
    ```

## Creating vSphere Clusters

This section describes setting up management and workload/guest clusters for
vSphere. If your deployment target is AWS, skip this section and move on to the
next.

1. Download the machine image that matches the version of the Kubernetes you plan on deploying (1.20.1 is default).

    At this time, we cannot guarantee the exact plugin versions that will be
    used for cluster management. While using the kickstart UI to bootstrap your
    cluster, you may be asked add an `ova` to your vSphere environment. The
    following links are points to the most recent ovas at the time of writing
    this getting started guide.

    * [1.20.4
      OVA](http://build-squid.eng.vmware.com/build/mts/release/bora-17759077/publish/lin64/tkg_release/node/ova-photon-3-v1.20.4+vmware.1-tkg.0-2326554155028348692/photon-3-kube-v1.20.4+vmware.1-tkg.0-2326554155028348692.ova)
    * [1.19.8
      OVA](http://build-squid.eng.vmware.com/build/mts/release/bora-17759077/publish/lin64/tkg_release/node/ova-photon-3-v1.19.8+vmware.1-tkg.0-15338136437231643652/photon-3-kube-v1.19.8+vmware.1-tkg.0-15338136437231643652.ova)

    If you're asked for another `ova` version by the kickstart UI, you can
    download the ova that corresponds to the rc version (e.g. 1,2,3,etc) at the [TKG
    daily builds confluence
    page](https://confluence.eng.vmware.com/pages/viewpage.action?spaceKey=TKG&title=TKG+Release+Daily+Build#TKGReleaseDailyBuild-TKG1.3.0RC.3(March/09/2021)).

    Please note, validation work so far has focused on the **Photon** based
    images.

1. In vCenter, right click on your datacenter and import OVF template.

1. After importing, right-click and covert to a template.

1. Initialize the Tanzu kickstart UI.

    ```sh
    tanzu management-cluster create --ui
    ```

1. Go through the installation process for vSphere. With the following
   considerations:

   * Set all instance profile to large or larger.
     * In our testing, we found resource constraints caused bootstrapping
     issues. Choosing a large profile or more will give a better chance for
     successful bootstrapping.
   * Set your control plane IP
     * The control plane IP is a virtual IP that fronts the Kubernetes API
     server. You **must** set an IP that is routable and won't be taken by
     another system (e.g. DHCP).
   * Disable OIDC configuration.

    > Until we have more TCE documentation, you can find the full TKG docs
    > [here](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.2/vmware-tanzu-kubernetes-grid-12/GUID-mgmt-clusters-deploy-management-clusters.html).
    > We will have more complete `tanzu` cluster bootstrapping documentation available here in the near future.

1. Validate the management cluster started successfully.

    ```sh
    tanzu management-cluster get
    ```

1. Create a cluster names that will be used throughout this getting-started.md guide. This instance of `MGMT_CLUSTER_NAME` should be set to whatever value is returned by `tanzu management-cluster get` above.

    ```sh
    export MGMT_CLUSTER_NAME="<INSERT_MGMT_CLUSTER_NAME_HERE>"
    export GUEST_CLUSTER_NAME="<INSERT_GUEST_CLUSTER_NAME_HERE>"
    ```

1. Capture the management cluster's kubeconfig.

    ```sh
    tanzu management-cluster kubeconfig get ${MGMT_CLUSTER_NAME} --admin

    Credentials of workload cluster 'mtce' have been saved
    You can now access the cluster by running 'kubectl config use-context mtce-admin@mtce'
    ```

    > Note the context name `${MGMT_CLUSTER_NAME}-admin@mtce`, you'll use the above command in
    > future steps. Your management cluster name may be different than
    > `${MGMT_CLUSTER_NAME}`.

1. Set your kubectl context to the management cluster.

    ```sh
    kubectl config use-context ${MGMT_CLUSTER_NAME}-admin@${MGMT_CLUSTER_NAME}
    ```

1. Validate you can access the management cluster's API server.

    ```sh
    kubectl get no

    NAME         STATUS   ROLES                  AGE    VERSION
    10-0-1-133   Ready    <none>                 123m   v1.20.1+vmware.2
    10-0-1-76    Ready    control-plane,master   125m   v1.20.1+vmware.2
    ```

1. Setup a guest cluster config file.

    ```sh
    cp  ~/.tanzu/tkg/clusterconfigs/xw6nt8jduy.yaml ~/.tanzu/tkg/clusterconfigs/guest1.yaml
    ```

   > This takes the configuration used to create your management cluster and
   > duplicates for use in the guest cluster. You can edit values in this new
   > file `guest1` as you please.

   [](ignored)

   > Creation of guest clusters now require the use of workload cluster YAML
   > configuration files.  [Example configuration templates](https://gitlab.eng.vmware.com/TKG/tkg-cli-providers/-/tree/cluster-templates/docs/cluster-templates)
   > are available to help get you started. Review settings and populate fields
   > that are not set.

   [](ignored)

   > Validation is performed on the file prior to applying it, so the `tanzu`
   > command should give you any clues if something necessary is omitted.

1. Edit the guest cluster config file's
   (`~/.tanzu/tkg/clusterconfigs/guest1.yaml`) CLUSTER_NAME.

   ```yaml
   CLUSTER_CIDR: 100.96.0.0/11
   CLUSTER_NAME: my-guest-cluster
   CLUSTER_PLAN: dev
   ```

1. Edit the guest cluster config file's
   (`~/.tanzu/tkg/clusterconfigs/guest1.yaml`) VSPHERE_CONTROL_PLANE_ENDPOINT to
   a viable IP.

   > This will be **the API Server IP** for you guest cluster. You must choose
   > an IP that is **1.) routable** and **2.) not used elsewhere in your network
   > (eg. out of your DHCP range)**.

   [](ignored)

   > For vSphere, the other settings are likely fine as-is. However, you can change
   > them as you'd like and/or reference the [Example configuration templates](https://gitlab.eng.vmware.com/TKG/tkg-cli-providers/-/tree/cluster-templates/docs/cluster-templates).

1. Create your guest cluster.

    ```sh
    tanzu cluster create ${GUEST_CLUSTER_NAME} --file ${HOME}/.tanzu/tkg/clusterconfigs/guest1.yaml
    ```

1. Validate the cluster starts successfully.

    ```sh
    tanzu cluster list
    ```

1. Capture the guest cluster's kubeconfig.

    ```sh
    tanzu cluster kubeconfig get ${GUEST_CLUSTER_NAME} --admin
    ```

1. Set your `kubectl` context accordingly.

    ```sh
    kubectl config use-context ${GUEST_CLUSTER_NAME}-admin@${GUEST_CLUSTER_NAME}
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

## Create AWS Clusters

This section describes setting up management and workload/guest clusters for
AWS. If your deployment target is vSphere, skip this section.

1. Initialize the Tanzu kickstart UI.

    ```sh
    tanzu management-cluster create --ui
    ```

1. Go through the installation process for AWS. With the following
   considerations:

   * Set all instance sizes to m5.xlarge or larger.
   * Disable OIDC configuration.

    > Until we have more TCE documentation, you can find the full TKG docs
    > [here](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.2/vmware-tanzu-kubernetes-grid-12/GUID-mgmt-clusters-deploy-management-clusters.html).
    > We will have more complete `tanzu` cluster bootstrapping documentation available here in the near future.

1. Validate the management cluster started successfully.

    ```sh
    tanzu management-cluster get

    NAME  NAMESPACE   STATUS   CONTROLPLANE  WORKERS  KUBERNETES        ROLES
    mtce  tkg-system  running  1/1           1/1      v1.20.1+vmware.2  management

    Details:

    NAME                                                     READY  SEVERITY  REASON  SINCE  MESSAGE
    /mtce                                                    True                     113m
    ├─ClusterInfrastructure - AWSCluster/mtce                True                     113m
    ├─ControlPlane - KubeadmControlPlane/mtce-control-plane  True                     113m
    │ └─Machine/mtce-control-plane-r7k52                     True                     113m
    └─Workers
      └─MachineDeployment/mtce-md-0
        └─Machine/mtce-md-0-fdfc9f766-6n6lc                  True                     113m

    Providers:

    NAMESPACE                          NAME                   TYPE                    PROVIDERNAME  VERSION  WATCHNAMESPACE
    capa-system                        infrastructure-aws     InfrastructureProvider  aws           v0.6.4
    capi-kubeadm-bootstrap-system      bootstrap-kubeadm      BootstrapProvider       kubeadm       v0.3.14
    capi-kubeadm-control-plane-system  control-plane-kubeadm  ControlPlaneProvider    kubeadm       v0.3.14
    capi-system                        cluster-api            CoreProvider            cluster-api   v0.3.14
    ```

1. Create a cluster names that will be used throughout this getting-started.md guide. This instance of `MGMT_CLUSTER_NAME` should be set to whatever value is returned by `tanzu management-cluster get` above.

    ```sh
    export MGMT_CLUSTER_NAME="<INSERT_MGMT_CLUSTER_NAME_HERE>"
    export GUEST_CLUSTER_NAME="<INSERT_GUEST_CLUSTER_NAME_HERE>"
    ```

1. Capture the management cluster's kubeconfig.

    ```sh
    tanzu management-cluster kubeconfig get ${MGMT_CLUSTER_NAME} --admin

    Credentials of workload cluster 'mtce' have been saved
    You can now access the cluster by running 'kubectl config use-context mtce-admin@mtce'
    ```

    > Note the context name `${MGMT_CLUSTER_NAME}-admin@mtce`, you'll use the above command in
    > future steps. Your management cluster name may be different than
    > `${MGMT_CLUSTER_NAME}`.

1. Set your kubectl context to the management cluster.

    ```sh
    kubectl config use-context ${MGMT_CLUSTER_NAME}-admin@${MGMT_CLUSTER_NAME}
    ```

1. Validate you can access the management cluster's API server.

    ```sh
    kubectl get no

    NAME                                       STATUS   ROLES                  AGE    VERSION
    ip-10-0-1-133.us-west-2.compute.internal   Ready    <none>                 123m   v1.20.1+vmware.2
    ip-10-0-1-76.us-west-2.compute.internal    Ready    control-plane,master   125m   v1.20.1+vmware.2
    ```

1. Setup a guest cluster config file.

    ```sh
    cp  ~/.tanzu/tkg/clusterconfigs/xw6nt8jduy.yaml ~/.tanzu/tkg/clusterconfigs/guest1.yaml
    ```

   > This takes the configuration used to create your management cluster and
   > duplicates for use in the guest cluster. You can edit values in this new
   > file `guest1` as you please.

   [](ignored)

   > Creation of guest clusters now require the use of workload cluster YAML
   > configuration files.  [Example configuration templates](https://gitlab.eng.vmware.com/TKG/tkg-cli-providers/-/tree/cluster-templates/docs/cluster-templates)
   > are available to help get you started. Review settings and populate fields
   > that are not set.

   [](ignored)

   > Validation is performed on the file prior to applying it, so the `tanzu`
   > command should give you any clues if something necessary is omitted.

1. Edit the guest cluster config file's
   (`~/.tanzu/tkg/clusterconfigs/guest1.yaml`) CLUSTER_NAME.

   ```yaml
   CLUSTER_CIDR: 100.96.0.0/11
   CLUSTER_NAME: my-guest-cluster
   CLUSTER_PLAN: dev
   ```

   > For AWS, the other settings are likely fine as-is. However, you can change
   > them as you'd like and/or reference the [Example configuration
   > templates](https://gitlab.eng.vmware.com/TKG/tkg-cli-providers/-/tree/cluster-templates/docs/cluster-templates).

1. Create your guest cluster.

    ```sh
    tanzu cluster create ${GUEST_CLUSTER_NAME} --file ${HOME}/.tanzu/tkg/clusterconfigs/guest1.yaml
    ```

1. Validate the cluster starts successfully.

    ```sh
    tanzu cluster list
    ```

1. Capture the guest cluster's kubeconfig.

    ```sh
    tanzu cluster kubeconfig get ${GUEST_CLUSTER_NAME} --admin
    ```

1. Set your `kubectl` context accordingly.

    ```sh
    kubectl config use-context ${GUEST_CLUSTER_NAME}-admin@${GUEST_CLUSTER_NAME}
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
    ```

## Configure kapp-controller

To use the Packaging APIs, TCE requires an alpha build of kapp-controller In
order to make this work, you need to **stop** the management cluster from
managing the guest clusters's kapp-controller. This enables you to mutate the
version of kapp-controller on the guest cluster.

1. Set your kube context to the **management cluster**.

    ```sh
    kubectl config use-context ${MGMT_CLUSTER_NAME}-admin@${MGMT_CLUSTER_NAME}
    ```

1. Set the `kapp-controller` App CR to pause reconciliation.

    ```sh
    kubectl patch app/${GUEST_CLUSTER_NAME}-kapp-controller --patch '{"spec":{"paused":true}}' --type=merge
    ```

1. Validate `kapp-controller` is not actively managed.

    ```sh
    $ kubectl get app -A
    NAMESPACE    NAME                        DESCRIPTION           SINCE-DEPLOY   AGE
    default      tce-guest-kapp-controller   Canceled/paused       128m           135m
    tkg-system   antrea                      Reconcile succeeded   2m40s          152m
    tkg-system   metrics-server              Reconcile succeeded   2m49s          149m
    tkg-system   tanzu-addons-manager        Reconcile succeeded   2m53s          153m
    ```

1. Set your kube context to the **workload/guest** cluster.

    ```sh
    kubectl config use-context ${MGMT_CLUSTER_NAME}-admin@${MGMT_CLUSTER_NAME}
    ```

1. Delete the existing `kapp-controller`.

   ```sh
   kubectl delete deploy -n tkg-system kapp-controller
   ```

1. Apply the alpha `kapp-controller` into the cluster.

   ```sh
   kubectl apply -f https://raw.githubusercontent.com/vmware-tanzu/carvel-kapp-controller/dev-packaging/alpha-releases/v0.17.0-alpha.1.yml
   ```

## Installing packages

1. Make sure your `kubectl` context is set to the workload cluster.

    ```sh
    kubectl config use-context ${GUEST_CLUSTER_NAME}-admin@${GUEST_CLUSTER_NAME}
    ```

1. Add TCE packages.

    ```sh
    tanzu package repository install --default
    ```

1. List the available packages.

    ```sh
    tanzu pacakge list

    NAME                 VERSION          DESCRIPTION
    cert-manager         1.1.0-vmware0
    contour-operator     1.11.0-vmware0
    fluent-bit           1.7.2-vmware0
    gatekeeper           3.2.3-vmware0
    grafana              7.4.3-vmware0
    knative-serving      0.21.0-vmware0
    prometheus           2.25.0-vmware0
    velero               1.5.2-vmware0
    ```

1. Install the package to the cluster.

    ```sh
    tanzu package install gatekeeper
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

## Cleaning up

After going through this guide, the following enables you to clean-up resources.

1. Delete any deployed workload clusters.

    ```sh
    tanzu cluster delete ${GUEST_CLUSTER_NAME}
    ```

1. Once all workload clusters have been deleted, the management cluster can
   then be removed as well.

    ```sh
    tanzu management-cluster get

    NAME                         NAMESPACE   STATUS   CONTROLPLANE  WORKERS  KUBERNETES        ROLES
    tkg-mgmt-aws-20210226062452  tkg-system  running  1/1           1/1      v1.20.1+vmware.2  management
    ```

    ```sh
    tanzu management-cluster delete tkg-mgmt-aws-20210226062452
    ```
