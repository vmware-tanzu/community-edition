## Create Local Docker Clusters

This section describes setting up a management cluster on your local workstation
using Docker.

⚠️: Tanzu Community Edition support for Docker is **experimental** and may require troubleshooting on your system.

### ⚠️  Warning on DockerHub Rate Limiting

When using the Docker (CAPD) provider, the load balancer image (HA Proxy) is
pulled from DockerHub. DockerHub limits pulls per user and this can especially
impact users who share a common IP, in the case of NAT or VPN. If DockerHub
rate-limiting is an issue in your environment, you can pre-pull the load
balancer image to your machine by running the following command.

```sh
docker pull kindest/haproxy:v20210715-a6da3463
```

This behavior will eventually be addressed in
[https://github.com/vmware-tanzu/community-edition/issues/897](https://github.com/vmware-tanzu/community-edition/issues/897).

### Local Docker Bootstrapping

1. Ensure your Docker engine has adequate resources. The  minimum requirements with no other containers running are: 6 GB of RAM and 4 CPUs.
   * **Linux**: Run ``docker system info``
   * **Mac**: Select Preferences > Resources > Advanced

1. Create the management cluster.

    ```sh
    CLUSTER_PLAN=dev tanzu management-cluster create -i docker <STANDALONE-CLUSTER-NAME>
    ```

    > For increased logs, you can append `-v 10`.
    <!-- -->
    > ⚠️ Capture the name of your cluster, it will be referenced as
    > ${CLUSTER_NAME} going forward.
    <!-- -->
    > ⚠️ The deployment will fail due to the CLI client being unable
    > to reach the API server running in the WSL VM. This is expected.

1. Let the deployment report failure.

1. Retrieve the address of the WSL VM.

    > ⚠️ Capture the VM IP of your cluster, it will be referenced as
    > ${WSL_VM_IP} going forward.

1. Query the docker daemon to get the forwarded port for HA Proxy. In the
   following example, the port is `44393`

    ```sh
    docker ps | grep -i ha
    44c0a71735ef   kindest/haproxy:v20210715-a6da3463

    "haproxy -sf 7 -W -d…"   2 days ago     Up 2 days     35093/tcp,
    0.0.0.0:44393->6443/tcp     muuhmuh-lb
    ```

    > ⚠️ Capture the port mentioned above, it will be referenced as
    > ${HA_PROXY_PORT} going forward.

1. Edit your `~/.kube/config` file.

1. Locate the YAML entry for your `${CLUSTER_NAME}`

1. In that YAML entry, replace `certificate-authority-data: < BASE64 DATA >`
   with `insecure-skip-tls-verify: true`.

1. In the YAML entry, replace `server: < api server value >` with
   `${WSL_VM_IP}:${HA_PROXY_PORT}`. Assuming the `${CLUSTER_NAME}` was
test, the entry would now look as follows.

    ```yaml
    - cluster:
        insecure-skip-tls-verify: true
        server: https://192.0.1.1:44393
      name: test
    ```

1. Save the file and exit.

    > `kubectl` and `tanzu` CLI should now be able to interact with your
    > cluster.

1. Validate the management cluster started:

    ```sh
    tanzu management-cluster get
    ```

    The output should look similar to the following:

    ```sh
    NAME                                               READY  SEVERITY  REASON  SINCE  MESSAGE
    /tkg-mgmt-docker-20210601125056                                                                 True                     28s
    ├─ClusterInfrastructure - DockerCluster/tkg-mgmt-docker-20210601125056                          True                     32s
    ├─ControlPlane - KubeadmControlPlane/tkg-mgmt-docker-20210601125056-control-plane               True                     28s
    │ └─Machine/tkg-mgmt-docker-20210601125056-control-plane-5pkcp                                  True                     24s
    │   └─MachineInfrastructure - DockerMachine/tkg-mgmt-docker-20210601125056-control-plane-9wlf2
      └─MachineDeployment/tkg-mgmt-docker-20210601125056-md-0
        └─Machine/tkg-mgmt-docker-20210601125056-md-0-5d895cbfd9-khj4s                              True                     24s
          └─MachineInfrastructure - DockerMachine/tkg-mgmt-docker-20210601125056-md-0-d544k


    Providers:

      NAMESPACE                          NAME                   TYPE                    PROVIDERNAME  VERSION  WATCHNAMESPACE
      capd-system                        infrastructure-docker  InfrastructureProvider  docker        v0.3.10
      capi-kubeadm-bootstrap-system      bootstrap-kubeadm      BootstrapProvider       kubeadm       v0.3.14
      capi-kubeadm-control-plane-system  control-plane-kubeadm  ControlPlaneProvider    kubeadm       v0.3.14
      capi-system                        cluster-api            CoreProvider            cluster-api   v0.3.14
    ```

1. Capture the management cluster's kubeconfig and take note of the command for accessing the cluster in the message, as you will use this for setting the context in the next step.

    ```sh
    tanzu management-cluster kubeconfig get <MGMT-CLUSTER-NAME> --admin
    ```

    * Where <``MGMT-CLUSTER-NAME>`` should be set to the name returned by `tanzu management-cluster get`.
    * For example, if your management cluster is called 'mtce', you will see a message similar to:

    ```sh
    Credentials of workload cluster 'mtce' have been saved.
    You can now access the cluster by running 'kubectl config use-context mtce-admin@mtce'
    ```

1. Set your kubectl context to the management cluster.

    ```sh
    kubectl config use-context <MGMT-CLUSTER-NAME>-admin@<MGMT-CLUSTER-NAME>
    ```

1. Validate you can access the management cluster's API server.

    ```sh
    kubectl get nodes
    ```

    You will see output similar to:

    ```sh
    NAME                         STATUS   ROLES                  AGE   VERSION
    guest-control-plane-tcjk2    Ready    control-plane,master   59m   v1.20.4+vmware.1
    guest-md-0-f68799ffd-lpqsh   Ready    <none>                 59m   v1.20.4+vmware.1
    ```

1. Create your workload cluster.

   ```shell
   tanzu cluster create <WORKLOAD-CLUSTER-NAME> --plan dev
   ```

1. Validate the cluster starts successfully.

    ```sh
    tanzu cluster list
    ```

1. Capture the workload cluster's kubeconfig.

    ```sh
    tanzu cluster kubeconfig get <WORKLOAD-CLUSTER-NAME> --admin
    ```

1. Set your `kubectl` context accordingly.

    ```sh
    kubectl config use-context <WORKLOAD-CLUSTER-NAME>-admin@<WORKLOAD-CLUSTER-NAME>
    ```

1. Verify you can see pods in the cluster.

    ```sh
    kubectl get pods --all-namespaces
    ```

    The output will look similar to the following:

    ```sh
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

You now have local clusters running on Docker. The nodes can be seen by running the  following command:

```shell
docker ps
```

The output will be similar to the following:

```sh
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
docker exec -it 4ae /bin/bash

root@guest-control-plane-tcjk2:/# kubectl --kubeconfig=/etc/kubernetes/admin.conf get nodes

NAME                         STATUS   ROLES                  AGE   VERSION
guest-control-plane-tcjk2    Ready    control-plane,master   67m   v1.20.4+vmware.1
guest-md-0-f68799ffd-lpqsh   Ready    <none>                 67m   v1.20.4+vmware.1
```

> In the above `4ae` is a control plane node.
