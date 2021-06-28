## Create Standalone Docker Clusters

This section describes setting up a standalone cluster on your local workstation
using Docker. This provides you a workload cluster that is **not** managed by a centralized management cluster.

⚠️: Tanzu Community Edition support for Docker is **experimental** and may require troubleshooting on your system.

1. Ensure your Docker engine has adequate resources. The  minimum requirements with no other containers running are: 6 GB of RAM and 4 CPUs.
    * **Linux**: Run ``docker system info``
    * **Mac**: Select Settings > Resources > Advanced

1. Store a name for your standalone cluster.

    ```sh
    export GUEST_CLUSTER_NAME="<INSERT_GUEST_CLUSTER_NAME_HERE>"
    ```

1. Create the standalone cluster.

    ```sh
    tanzu standalone-cluster create -i docker ${GUEST_CLUSTER_NAME}
    ```

    > For increased logs, you can append `-v 10`.

1. Validate the cluster started successfully.

    ```txt
    Standalone cluster created!
    ```

1. Set your kubectl context to the cluster.

    ```sh
    kubectl config use-context ${GUEST_CLUSTER_NAME}-admin@${GUEST_CLUSTER_NAME}
    ```

1. Validate you can access the cluster's API server.

    ```sh
    kubectl get pod -A

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
