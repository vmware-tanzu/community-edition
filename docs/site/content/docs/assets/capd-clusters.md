## Create Local Docker Clusters (CAPD)

⚠️: TCE's support for Cluster-API Docker (CAPD) is **experimental** and may require troubleshooting on your system.

In order to use CAPD, you should ensure your Docker engine has plenty of resources and storage available.
We do not have exact numbers at this time, but most modern laptops and desktops should be able to run this configuration.
Even with a modern system, there is always potential you're using up most or all of what Docker has allocated. **For your
best chance at success, you may wish to do the steps marked `(Optional)`.**

1. Ensure your Docker engine has >= 6GB of RAM and 4+ CPUs. This is especially critical for Mac users, below you can see a screenshot from Docker Desktop's settings.

   ![docker settings](/docs/img/docker-settings.png)

   > These are the minimal settings we validated against with not other containers running.

1. (Optional): Stop all existing containers.

   ```shell
   docker kill $(docker ps -q)
   ```

1. (Optional): Prune all existing containers, volumes, and images

   ```sh
    docker system prune -a --volumes
   ```

   > Read the prompt carefully, this will erase a majority of what is cached in your Docker environment.
   > While this will ensure your environment is clean before starting, it will also significantly increase bootstrapping
   > time if you already had the docker images downloaded.

1. Create a cluster names that will be used throughout this getting-started.md guide. This instance of `MGMT_CLUSTER_NAME` should be set to whatever value is returned by `tanzu management-cluster get` above.

    ```sh
    export MGMT_CLUSTER_NAME="<INSERT_MGMT_CLUSTER_NAME_HERE>"
    export GUEST_CLUSTER_NAME="<INSERT_GUEST_CLUSTER_NAME_HERE>"
    ```

1. Create a Docker-based management cluster

   ```sh
   tanzu management-cluster create -i docker --name ${MGMT_CLUSTER_NAME} -v 10 --plan dev --ceip-participation=false
   ```

   > If you ran the `prune` command in the previous step, expect this to take some time, as it'll download an image that is over 1GB.

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
    kubectl get nodes

    NAME                         STATUS   ROLES                  AGE   VERSION
    guest-control-plane-tcjk2    Ready    control-plane,master   59m   v1.20.4+vmware.1
    guest-md-0-f68799ffd-lpqsh   Ready    <none>                 59m   v1.20.4+vmware.1
    ```

1. Create your guest cluster.

   ```shell
   tanzu cluster create ${GUEST_CLUSTER_NAME} --plan dev
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

With the above done, you now have local clusters running atop Docker. The nodes can be seen by running a command as follows.

```shell
$ docker ps
CONTAINER ID   IMAGE                                                         COMMAND                  CREATED             STATUS             PORTS                                  NAMES
33e4e422e102   projects.registry.vmware.com/tkg/kind/node:v1.20.4_vmware.1   "/usr/local/bin/entr…"   About an hour ago   Up About an hour                                          guest-md-0-f68799ffd-lpqsh
4ae2829ab6e1   projects.registry.vmware.com/tkg/kind/node:v1.20.4_vmware.1   "/usr/local/bin/entr…"   About an hour ago   Up About an hour   41637/tcp, 127.0.0.1:41637->6443/tcp   guest-control-plane-tcjk2
c0947823840b   kindest/haproxy:2.1.1-alpine                                  "/docker-entrypoint.…"   About an hour ago   Up About an hour   42385/tcp, 0.0.0.0:42385->6443/tcp     guest-lb
a2f156fe933d   projects.registry.vmware.com/tkg/kind/node:v1.20.4_vmware.1   "/usr/local/bin/entr…"   About an hour ago   Up About an hour                                          mgmt-md-0-b8689788f-tlv68
128bf25b9ae9   projects.registry.vmware.com/tkg/kind/node:v1.20.4_vmware.1   "/usr/local/bin/entr…"   About an hour ago   Up About an hour   40753/tcp, 127.0.0.1:40753->6443/tcp   mgmt-control-plane-9rdcq
e59ca95c14d7   kindest/haproxy:2.1.1-alpine                                  "/docker-entrypoint.…"   About an hour ago   Up About an hour   35621/tcp, 0.0.0.0:35621->6443/tcp     mgmt-lb
```

The above reflects 1 management cluster and 1 workload cluster, both featuring 1 control plane node and 1 worker node.
Each cluster gets an `haproxy` container fronting the control plane node(s). This enables scaling the control plane into
an HA configuration.

🛠️: For troubleshooting failed bootstraps, you can exec into a container and use the kubeconfig at `/etc/kubernetes/admin.conf` to access
the API server directly. For example:

```shell
$ docker exec -it 4ae /bin/bash

root@guest-control-plane-tcjk2:/# kubectl --kubeconfig=/etc/kubernetes/admin.conf get nodes

NAME                         STATUS   ROLES                  AGE   VERSION
guest-control-plane-tcjk2    Ready    control-plane,master   67m   v1.20.4+vmware.1
guest-md-0-f68799ffd-lpqsh   Ready    <none>                 67m   v1.20.4+vmware.1
```

> In the above `4ae` is a control plane node.

## Configure kapp-controller

To use the Packaging APIs, TCE requires an alpha build of kapp-controller. In
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
    kubectl config use-context ${GUEST_CLUSTER_NAME}-admin@${GUEST_CLUSTER_NAME}
    ```

1. Delete the existing `kapp-controller`.

   ```sh
   kubectl delete deployment -n tkg-system kapp-controller
   ```

1. Apply the alpha `kapp-controller` into the cluster. You can check for the latest version in kapp-controller's [installation documentation](https://carvel.dev/kapp-controller/docs/latest/install-alpha/).

   ```sh
   kubectl apply -f https://raw.githubusercontent.com/vmware-tanzu/carvel-kapp-controller/dev-packaging/alpha-releases/v0.18.0-alpha.4.yml
   ```
