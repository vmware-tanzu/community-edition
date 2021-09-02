## Create Standalone Docker Clusters

This section describes setting up a standalone cluster on your local workstation
using Docker. This provides you a workload cluster that is **not** managed by a centralized management cluster.

⚠️: Tanzu Community Edition support for Docker is **experimental** and may require troubleshooting on your system.

<<<<<<< HEAD

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
## Deployment Procedure

=======
>>>>>>> 914070b68dfd5e17fc8557afc0bf32bf1e41f0ed
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

<<<<<<< HEAD
=======
1. Ensure your Docker engine has adequate resources. The  minimum requirements with no other containers running are: 6 GB of RAM and 4 CPUs.
    * **Linux**: Run ``docker system info``
    * **Mac**: Select Preferences > Resources > Advanced

>>>>>>> 914070b68dfd5e17fc8557afc0bf32bf1e41f0ed
1. Create the standalone cluster.

    ```sh
    tanzu standalone-cluster create -i docker <STANDALONE-CLUSTER-NAME>
    ```
    >``<STANDALONE-CLUSTER-NAME>`` must end with a letter, not a numeric character, and must be compliant with DNS hostname requirements [RFC 952](https://tools.ietf.org/html/rfc952) and [RFC 1123](https://tools.ietf.org/html/rfc1123).
    > For increased logs, you can append `-v 10`.

   If the deployment is successful, you should see the following output:

    ```txt
    Standalone cluster created!
    ```

2. Set your kubectl context to the cluster.

    ```sh
    kubectl config use-context <STANDALONE-CLUSTER-NAME>-admin@<STANDALONE-CLUSTER-NAME>
    ```

3. Validate you can access the cluster's API server.

    ```sh
    kubectl get pod -A
    ```
    The output should look similar to the following:

    ```sh
    NAMESPACE         NAME                                                                         READY   STATUS    RESTARTS   AGE
    kapp-controller   kapp-controller-5c66dcc7cf-62jl2                                             1/1     Running   0          3m52s
    kube-system       antrea-agent-7vs9l                                                           2/2     Running   0          3m52s
    kube-system       antrea-agent-zkgv7                                                           2/2     Running   0          3m28s
    kube-system       antrea-controller-785dbc59b8-6vj86                                           1/1     Running   0          3m52s
    kube-system       coredns-68d49685bd-sjp7t                                                     1/1     Running   0          3m52s
    kube-system       coredns-68d49685bd-xr5b2                                                     1/1     Running   0          3m52s
    kube-system       etcd-tkg-mgmt-docker-20210429071830-control-plane-vd8nl                      1/1     Running   0          4m12s
    kube-system       kube-apiserver-tkg-mgmt-docker-20210429071830-control-plane-vd8nl            1/1     Running   0          4m12s
    kube-system       kube-controller-manager-tkg-mgmt-docker-20210429071830-control-plane-vd8nl   1/1     Running   0          4m12s
    kube-system       kube-proxy-7r54w                                                             1/1     Running   0          3m28s
    kube-system       kube-proxy-m6l64                                                             1/1     Running   0          3m52s
    kube-system       kube-scheduler-tkg-mgmt-docker-20210429071830-control-plane-vd8nl            1/1     Running   0          4m12s
    tkr-system        tkr-controller-manager-96445c85d-8qh44                                       1/1     Running   0          3m52s
    ```

⚠️: If the Docker host machine is rebooted, the cluster will need to be
re-created. Support for clusters surviving a host reboot is track in issue
[#832](https://github.com/vmware-tanzu/community-edition/issues/832).
