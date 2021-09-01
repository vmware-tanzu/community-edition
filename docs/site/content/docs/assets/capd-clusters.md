## Create Local Docker Clusters

This section describes setting up a management cluster on your local workstation
using Docker.

⚠️: Tanzu Community Edition support for Docker is **experimental** and may require troubleshooting on your system.

**Note: You cannot bootstrap a cluster to Docker from a Windows bootstrap machine, only Linux and Mac are supported at this time for Docker cluster deployments.**

## Prerequisites
The following additional configuration is needed for the Docker engine on your local client machine (with no other containers running):
| |
|:------------------------|
|6 GB of RAM |
|15 GB of local machine disk storage for images |
|4 CPUs|

Check your Docker configuration as follows:
- Linux: Run ``docker system info``
-  Mac: Select Preferences > Resources > Advanced
## Before You Begin
To optimise your Docker system and ensure a successful deployment, you may wish to complete the next two optional steps.

1. (Optional): Stop all existing containers.

   ```shell
   docker kill $(docker ps -q)
   ```
1. (Optional): Run the following command to prune all existing containers, volumes, and images.

    Warning: Read the prompt carefully before running the command, as it erases the majority of what is cached in your Docker environment. While this ensures your environment is clean before starting, it also significantly increases bootstrapping time if you already had the Docker images downloaded.

   ```sh
    docker system prune -a --volumes
   ```
## Deployment Procedure
1. Initialize the Tanzu Community Edition installer interface.

   ```sh
   tanzu management-cluster create --ui
   ```
1. Complete the configuration steps in the installer interface for Docker and create the management cluster. The following configuration settings are recommended:

   * The Kubernetes Network Settings are auto-filled with a default CNI Provider and Cluster Service CIDR.
   * Docker Proxy settings are experimental and are to be used at your own risk.
   * We will have more complete `tanzu` cluster bootstrapping documentation available here in the near future.
   * If you ran the `prune` command in the previous step, expect this to take some time, as it'll download an image that is over 1GB.

1. (Alternative method) It is also possible to use the command line to create a Docker based management cluster:
    ```sh
    tanzu management-cluster create -i docker --name <MGMT-CLUSTER-NAME> -v 10 --plan dev --ceip-participation=false
    ```
    -  ``<MGMT-CLUSTER-NAME>`` must end with a letter, not a numeric character, and must be compliant with DNS hostname requirements described here: [RFC 1123](https://tools.ietf.org/html/rfc1123).

2. Validate the management cluster started:
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

3. Capture the management cluster's kubeconfig and take note of the command for accessing the cluster in the message, as you will use this for setting the context in the next step.

    ```sh
    tanzu management-cluster kubeconfig get <MGMT-CLUSTER-NAME> --admin
    ```
    - Where <``MGMT-CLUSTER-NAME>`` should be set to the name returned by `tanzu management-cluster get`.
    - For example, if your management cluster is called 'mtce', you will see a message similar to:
    ```sh
    Credentials of workload cluster 'mtce' have been saved.
    You can now access the cluster by running 'kubectl config use-context mtce-admin@mtce'
    ```

4. Set your kubectl context to the management cluster.

    ```sh
    kubectl config use-context <MGMT-CLUSTER-NAME>-admin@<MGMT-CLUSTER-NAME>
    ```

5. Validate you can access the management cluster's API server.

    ```sh
    kubectl get nodes
    ```
    You will see output similar to:
    ```sh
    NAME                         STATUS   ROLES                  AGE   VERSION
    guest-control-plane-tcjk2    Ready    control-plane,master   59m   v1.20.4+vmware.1
    guest-md-0-f68799ffd-lpqsh   Ready    <none>                 59m   v1.20.4+vmware.1
    ```

6. Create your workload cluster.

   ```shell
   tanzu cluster create <WORKLOAD-CLUSTER-NAME> --plan dev
   ```

7.  Validate the cluster starts successfully.

    ```sh
    tanzu cluster list
    ```

8.  Capture the workload cluster's kubeconfig.

    ```sh
    tanzu cluster kubeconfig get <WORKLOAD-CLUSTER-NAME> --admin
    ```

9.  Set your `kubectl` context accordingly.

    ```sh
    kubectl config use-context <WORKLOAD-CLUSTER-NAME>-admin@<WORKLOAD-CLUSTER-NAME>
    ```

10. Verify you can see pods in the cluster.

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
$ docker ps
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
$ docker exec -it 4ae /bin/bash

root@guest-control-plane-tcjk2:/# kubectl --kubeconfig=/etc/kubernetes/admin.conf get nodes

NAME                         STATUS   ROLES                  AGE   VERSION
guest-control-plane-tcjk2    Ready    control-plane,master   67m   v1.20.4+vmware.1
guest-md-0-f68799ffd-lpqsh   Ready    <none>                 67m   v1.20.4+vmware.1
```

> In the above `4ae` is a control plane node.
