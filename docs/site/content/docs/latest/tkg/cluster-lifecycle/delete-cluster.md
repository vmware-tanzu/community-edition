# Delete Tanzu Kubernetes Clusters

To delete a Tanzu Kubernetes cluster, run the `tanzu cluster delete` command.
Depending on the cluster contents and cloud infrastructure, you may need to delete in-cluster volumes and services before you delete the cluster itself.

## <a id="list"></a> Step One: List Clusters

To list all of the Tanzu Kubernetes clusters that a management cluster is managing, run the `tanzu cluster list` command.

   ```
   tanzu cluster list
   ```  

## <a id="vols-services"></a> Step Two: Delete Volumes and Services

If the cluster you want to delete contains persistent volumes or services such as load balancers and databases, you may need to manually delete them before you delete the cluster itself.
What you need to pre-delete depends on your cloud infrastructure:

* **vSphere**

    * **Load Balancer**: see [Delete Service type LoadBalancer](#servicetypelb) below.
    * **Persistent Volumes and Persistent Volume Claims**: see [Delete Persistent Volume Claims and Persistent Volumes](#pv), below.

* **Amazon EC2**

    * **Load Balancers**: Application or Network Load Balancers (ALBs or NLBs) in the cluster's VPC, but not Classic Load Balancers (ELB v1).
    * **Other Services**: Any subnet/EC2 backed service in cluster's VPC, such as an RDS.
    * **Persistent Volumes and Persistent Volume Claims**: see [Delete Persistent Volume Claims and Persistent Volumes](#pv), below.

* **Azure**

    * No action required.
    Deleting a cluster deletes everything that TKG created in the cluster's resource group.

### <a id="servicetypelb"></a> Delete Service type LoadBalancer

To delete Service type LoadBalancer (Service) in a cluster:

1. Set `kubectl` to the cluster's context.

   ```
   kubectl config set-context my-cluster@user
   ```

1. Retrieve the cluster's list of services.

   ```
   kubectl get service
   ```

1. Delete each Service type `LoadBalancer`.

    ```
    kubectl delete service <my-svc>
    ```

### <a id="pv"></a> Delete Persistent Volumes and Persistent Volume Claims

To delete Persistent Volume (PV) and Persistent Volume Claim (PVC) objects in a cluster:

1. Run `kubectl config set-context my-cluster@user` to set `kubectl` to the cluster's context.

1. Run `kubectl get pvc` to retrieve the cluster's Persistent Volume Claims (PVCs).

1. For each PVC:

    1. Run `kubectl describe pvc <my-pvc>` to identify the PV it is bound to.
    The PV is listed in the command output as **Volume**, after **Status: Bound**.

    1. Run `kubectl describe pv <my-pv>` to describe to determine if its bound PV `Reclaim Policy` is `Retain` or `Delete`.

    1. Run `kubectl delete pvc <my-pvc>` to delete the PVC.

    1. If the PV reclaim policy is `Retain`, run `kubectl delete pv <my-pvc>` and then log into your cloud portal and delete the PV object there.
    For example, delete a vSphere CNS volume from your datastore pane > **Monitor** > **Cloud Native Storage** > **Container Volumes**.
    For more information about vSphere CNS, see [Getting Started with VMware Cloud Native Storage](https://docs.vmware.com/en/VMware-vSphere/6.7/Cloud-Native-Storage/GUID-51D308C7-ECFE-4C04-AD56-64B6E00A6548.html).

## <a id="cluster"></a> Step Three: Delete Cluster

1. To delete a cluster, run `tanzu cluster delete`.

   ```
   tanzu cluster delete my-cluster
   ```

   If the cluster is running in a namespace other than the `default` namespace, you must specify the `--namespace` option to delete that cluster.

   ```
   tanzu cluster delete my-cluster --namespace=my-namespace
   ```

   To skip the `yes/no` verification step when you run `tanzu cluster delete`, specify the `--yes` option.

   ```
   tanzu cluster delete my-cluster --namespace=my-namespace --yes
   ```

**IMPORTANT**: Do not change context or edit the `.kube-tkg/config` file while Tanzu Kubernetes Grid operations are running.
