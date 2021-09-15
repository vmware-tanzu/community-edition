## Create Standalone Docker Clusters

This section describes setting up a standalone cluster on your local workstation
using Docker. This provides you a workload cluster that is **not** managed by a centralized management cluster.

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

1. Create the standalone cluster.

    ```sh
    tanzu standalone-cluster create -i docker <STANDALONE-CLUSTER-NAME>
    ```
    >``<STANDALONE-CLUSTER-NAME>`` must end with a letter, not a numeric character, and must be compliant with DNS hostname requirements [RFC 952](https://tools.ietf.org/html/rfc952) and [RFC 1123](https://tools.ietf.org/html/rfc1123).
    > For increased logs, you can append `-v 10`.

    > ⚠️ Capture the name of your cluster, it will be referenced as
    > ${CLUSTER_NAME} going forward.

    > ⚠️ The deployment will fail due to the CLI client being unable
    > to reach the API server running in the WSL VM. This is expected.

1. Let the deployment report failure.

1. Retrieve the address of the WSL VM.

    > ⚠️ Capture the VM IP of your cluster, it will be referenced as
    > ${WSL_VM_IP} going forward.

1. Query the docker daemon to get the forwarded port for HA Proxy. In the
   following example, the port is `44393`

    ```sh
    $ docker ps | grep -i ha
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
