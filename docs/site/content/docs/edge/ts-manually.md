# Troubleshooting a Bootstrap Cluster Manually

If the [Tanzu diagnostics](tanzu-diagnostics) command did not yield the information you needed or you want complete control over your troubleshooting steps, complete the following steps to troubleshoot a bootstrap cluster manually:

1. Run docker ps on your local Docker system to get the name of the bootstrap cluster container:

   ```sh
   docker ps
   ```

   The bootstrap cluster container name will begin with `tkg-cluster` followed by a unique ID, for example, `tkg-cluster-example1234567abcdef`. Copy the CONTAINER ID of the bootstrap cluster container.

1. Open a bash shell in the bootstrap cluster container:

   ```sh
   docker exec -it <BOOTSTRAP-CLUSTER-ID> bash
   ```

   Where ``<BOOTSTRAP-CLUSTER-ID>`` is the value copied in the previous step.

1. Before you can proceed to run ``kubectl`` commands against the pods inside the bootstrap cluster container, copy the `admin.conf` file to the default kubeconfig location:

   ```sh
   cp -v /etc/kubernetes/admin.conf ~/.kube/config
   ```

1. Now you are inside the bootstrap cluster container that is going to bootstrap your cluster to the target provider, you can run ``kubectl`` commands against this container. By watching the status of the pods, you can understand what might go wrong in the bootstrap process. Run the following command to see the pods being created inside the container:

   ```sh
   kubectl get po -A
   ```

1. Copy the name of the controller manager, it is usually first in the list. It will be named similarly to the following depending on your target provider:

   * cap**a**-controller-manager-12a3456789-b1cde (AWS)
   * cap**d**-controller-manager-12a3456789-b1cde (Docker)
   * cap**v**-controller-manager-12a3456789-b1cde (vSphere)
   * cap**z**-controller-manager-12a3456789-b1cde (Azure)

1. Next, you can examine the logs of the controller manager that communicates with the target provider. This step is important, if you are having problems bootstrapping, the errors in the controller logs will provide the detail.  Examine the logs for the controller manager:

   ```sh
   kubectl logs -n <NAMESPACE> <CONTROLLER-MANAGER> -c manager –f
   ```

   Where

   * ``<CONTROLLER-MANAGER>`` is the value copied in the previous step.
   * ``<NAMESPACE>`` will vary based on your provider, use:
     * ``capa-system`` (AWS)
     * ``capd-system`` (Docker)
     * ``capv-system`` (vSphere)
     * ``capz-system`` (Azure)

1. [Optional] Events are also reported based on actions taken in the target
   provider. You can view the known events by running:

   ```sh
   kubectl get events -n tkg-system
   ```
