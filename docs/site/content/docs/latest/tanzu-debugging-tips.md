# Cluster Debugging Tips

This section is a collection of common issues and how to debug them.

## Bootstrap cluster fails to successfully pivot to a management cluster

### Problem

The bootstrap cluster fails to successfully pivots to a management cluster by hanging or returning an error similar to the following:

```shell
Error: unable to set up management cluster: unable to wait for cluster and get the cluster kubeconfig: error waiting for cluster to be provisioned (this may take a few minutes): cluster creation failed, reason:'NatGatewaysReconciliationFailed', message:'3 of 8 completed'
```

### Solution

Access the pod log for your respective provider to identify the problem.

1. On the machine running the bootstrap cluster, run `docker ps` to list the control plane container(s):

```sh
docker ps
CONTAINER ID        IMAGE                                                             COMMAND                  CREATED             STATUS              PORTS                       NAMES
2f505b8b0c8a        projects-stg.registry.vmware.com/tkg/kind/node:v1.21.2_vmware.1   "/usr/local/bin/entr…"   17 minutes ago      Up 17 minutes       127.0.0.1:33876->6443/tcp   tkg-kind-c51pdgq1m0cj2rffu1d0-control-plane
3dd82b26ac04        projects-stg.registry.vmware.com/tkg/kind/node:v1.21.2_vmware.1   "/usr/local/bin/entr…"   About an hour ago   Up About an hour    127.0.0.1:46450->6443/tcp   tkg-kind-c51oina1m0co5l1khoa0-control-plane
```

The previous shows a list of active containers in docker (the output for your environment will vary). Make note of the container ID as it will be used in the following steps.

1. Next, identify the controller-manager pod for your provider (for instance, for AWS it is `capa-controller-manager-*` in the `capa-system` namespace)

```sh
docker exec <CONTAINER_ID> kubectl get po --namespace capa-system --kubeconfig /etc/kubernetes/admin.conf
NAMESPACE                           NAME                                                                  READY   STATUS    RESTARTS   AGE
capa-system                         capa-controller-manager-<UINQUE_ID>                              2/2     Running   0          19m
```

The previous command will vary depending on your management cluster provider, as shown below:

* AWS: `capa-controller-manager` / `capa-system`
* Azure: `capz-controller-manager` / `capz-system`
* docker: `capd-controller-manager` / `capd-system`
* vSphere: `capv-controller-manager` / `capv-system`

1. Retrieve the log for that controller to verify why it is unable to complete the management cluster provision:

```shell
docker exec <CONTAINDER_ID> kubectl logs capa-controller-manager-<UNIQUE_ID> -n capa-system -c manager --kubeconfig /etc/kubernetes/admin.conf
```

The log from the command above should provide hints for why the provider is unable to complete the management cluster provision.

## Cleanup after an unsuccessful management deployment

### Problem

When a management cluster fails to deploy successfully (or partially deploys), it may leave orphaned objects in your bootstrap environment.

### Solution

Clean the bootstrap environment prior to a subsequent attempt of redeploying the management cluster.

1. If the management cluster got partially created, attempt to delete the resources for the failed cluster:

```shell
tanzu management-cluster delete <YOUR-CLUSTER-NAME>
```

1. Next, if the bootstrap cluster still exists, delete it:

```shell
kind get clusters
tkg-kind-b4o9sn5948199qbgca8d
```

```sh
kind delete cluster --name tkg-kind-b4o9sn5948199qbgca8d
```

1. Use `docker` to stop and remove any running containers related to the bootstrap process

```shell
docker ps
CONTAINER ID        IMAGE                                                             COMMAND                  CREATED             STATUS              PORTS                       NAMES
2f505b8b0c8a        projects-stg.registry.vmware.com/tkg/kind/node:v1.21.2_vmware.1   "/usr/local/bin/entr…"   17 minutes ago      Up 17 minutes       127.0.0.1:33876->6443/tcp   tkg-kind-c51pdgq1m0cj2rffu1d0-control-plane
3dd82b26ac04        projects-stg.registry.vmware.com/tkg/kind/node:v1.21.2_vmware.1   "/usr/local/bin/entr…"   About an hour ago   Up About an hour    127.0.0.1:46450->6443/tcp   tkg-kind-c51oina1m0co5l1khoa0-control-plane
```

```shell
docker stop 2f505b8b0c8a && docker rm 2f505b8b0c8a
docker stop 3dd82b26ac04 && docker rm 3dd82b26ac04
```
