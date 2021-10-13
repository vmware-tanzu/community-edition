## Create Standalone Docker Clusters

This section describes setting up a standalone cluster on your local workstation
using Docker. This provides you a workload cluster that is **not** managed by a centralized management cluster.

⚠️: Tanzu Community Edition support for Docker is **experimental** and may require troubleshooting on your system.

Note: Bootstrapping a cluster to Docker from a Windows bootstrap machine is currently experimental.

## Prerequisites

The following additional configuration is needed for the Docker engine on your local client machine (with no other containers running):

| |
|:------------------------|
|6 GB of RAM |
|15 GB of local machine disk storage for images |
|4 CPUs|

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

### Before You Begin

To optimise your Docker system and ensure a successful deployment, you may wish to complete the next two optional steps.

1. (Optional): Stop all existing containers.

   ```shell
   docker kill $(docker ps -q)
   ```

1. (Optional): Run the following command to prune all existing containers, volumes, and images.

   Warning: Read the prompt carefully before running the command, as it erases the majority of what is cached in your Docker environment.
While this ensures your environment is clean before starting, it also significantly increases bootstrapping time if you already had the Docker images downloaded.

   ```sh
    docker system prune -a --volumes
   ```

### Local Docker Bootstrapping

1. Create the standalone cluster.

    ```sh
    tanzu standalone-cluster create -i docker <STANDALONE-CLUSTER-NAME>
    ```

    >`<STANDALONE-CLUSTER-NAME>` must end with a letter, not a numeric character, and must be compliant with DNS hostname requirements [RFC 952](https://tools.ietf.org/html/rfc952) and [RFC 1123](https://tools.ietf.org/html/rfc1123).
    > For increased logs, you can append `-v 10`.

    If the deployment is successful, you should see the following output:

    ```txt
    Standalone cluster created!
    ```

1. Set your kubectl context to the cluster.

    ```sh
    kubectl config use-context <STANDALONE-CLUSTER-NAME>-admin@<STANDALONE-CLUSTER-NAME>
    ```

1. Validate you can access the cluster's API server.

    ```sh
    kubectl get pod -A
    ```

    The output should look similar to the following:

    ```sh
    NAMESPACE     NAME                                                     READY   STATUS    RESTARTS   AGE
    kube-system   antrea-agent-4wwq9                                       2/2     Running   0          3m28s
    kube-system   antrea-agent-s9gbb                                       2/2     Running   0          3m28s
    kube-system   antrea-controller-58cdb9dc6d-mdn56                       1/1     Running   0          3m28s
    kube-system   coredns-8dcb5c56b-7dltt                                  1/1     Running   0          4m43s
    kube-system   coredns-8dcb5c56b-cvkpx                                  1/1     Running   0          4m43s
    kube-system   etcd-testme-control-plane-2fcfs                          1/1     Running   0          4m44s
    kube-system   kube-apiserver-testme-control-plane-2fcfs                1/1     Running   0          4m44s
    kube-system   kube-controller-manager-testme-control-plane-2fcfs       1/1     Running   0          4m44s
    kube-system   kube-proxy-7wfs8                                         1/1     Running   0          4m8s
    kube-system   kube-proxy-bzr2d                                         1/1     Running   0          4m43s
    kube-system   kube-scheduler-testme-control-plane-2fcfs                1/1     Running   0          4m44s
    tkg-system    kapp-controller-764fc6c69f-lpvn6                         1/1     Running   0          3m49s
    tkg-system    tanzu-capabilities-controller-manager-69f58566d9-8ks8q   1/1     Running   0          4m28s
    tkr-system    tkr-controller-manager-cc88b6968-hv8zg                   1/1     Running   0          4m28s
    ```

⚠️: If the Docker host machine is rebooted, the cluster will need to be
re-created. Support for clusters surviving a host reboot is tracked in issue [#832](https://github.com/vmware-tanzu/community-edition/issues/832).
