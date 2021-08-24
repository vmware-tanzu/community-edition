# Upgrade Management Clusters

To upgrade your Tanzu Kubernetes Grid instance, you must first upgrade the management cluster. You cannot upgrade Tanzu Kubernetes clusters until you have upgraded the management cluster that manages them.

**IMPORTANT**: Management clusters and Tanzu Kubernetes clusters use client certificates to authenticate clients. These certificates are valid for one year. To renew them, upgrade your clusters at least once a year.

## <a id="prereqs"></a> Prerequisites

- You performed the steps in [Upgrading Tanzu Kubernetes Grid](index.md) that occur before the step for upgrading management clusters.
- If you deployed the previous version of Tanzu Kubernetes Grid in an Internet-restricted environment, you have performed the steps in [Deploying Tanzu Kubernetes Grid in an Internet-Restricted Environment](../mgmt-clusters/airgapped-environments.md#procedure) to recreate and run the `gen-publish-images.sh` and `publish-images.sh` scripts with the new component image versions.

## <a id="procedure"></a> Procedure

1. Run the `tanzu login` command to see an interactive list of management clusters available for upgrade.

   ```
   tanzu login
   ```

1. Select the management cluster that you want to upgrade. See [List Management Clusters and Change Context](../cluster-lifecycle/multiple-management-clusters.md#login) for more information.

1. Run the `tanzu cluster list` command with the `--include-management-cluster` option.
This command shows the versions of Kubernetes running on the management cluster and all of the clusters that it manages:

   ```
   $ tanzu cluster list --include-management-cluster
     NAME                 NAMESPACE   STATUS    CONTROLPLANE  WORKERS  KUBERNETES         ROLES       PLAN
     k8s-1-17-13-cluster  default     running   1/1           1/1      v1.17.13+vmware.1  <none>      dev
     k8s-1-18-10-cluster  default     running   1/1           1/1      v1.18.10+vmware.1  <none>      dev
     k8s-1-19-3-cluster   default     running   1/1           1/1      v1.19.3+vmware.1   <none>      dev
     mgmt-cluster         tkg-system  running   1/1           1/1      v1.20.4+vmware.1   management  dev
   ```

1. Before you run the upgrade command, remove all unmanaged `kapp-controller` deployment artifacts from the management cluster. An unmanaged `kapp-controller` deployment is a deployment that exists outside of the `vmware-system-tmc` namespace.

    1. Delete the `kapp-controller` deployment.

       ```
       kubectl delete deployment kapp-controller -n kapp-controller
       ```

       **Note:** If you receive a `NotFound` error message, ignore the error. You should continue with the following deletion steps in case you have any orphaned objects related to a pre-existing `kapp-controller` deployment.

       ```
       Error from server (NotFound): deployments.apps "kapp-controller" not found
       ```

    1. Delete all `kapp-controller` objects.

       ```
       kubectl delete clusterrole kapp-controller-cluster-role
       kubectl delete clusterrolebinding kapp-controller-cluster-role-binding
       kubectl delete serviceaccount kapp-controller-sa -n kapp-controller
       ```

1. If you set up Harbor access by installing a connectivity API on your v1.2 management cluster, follow the [Replace Connectivity API with a Load Balancer](#connectivity-api) procedure below.
If your workload clusters access Harbor via a load balancer, proceed to the next step.

1. Run the `tanzu management-cluster upgrade` command and enter `y` to confirm.

   The following command upgrades the current management cluster.

   ```
   tanzu management-cluster upgrade
   ```

   If multiple base VM images in your IaaS account have the same version of Kubernetes that you are upgrading to, use the `--os-name` option to specify the OS you want.
  See [Selecting an OS During Cluster Upgrade](cluster-os-upgrade.md) for more information.

   For example, on vSphere if you have uploaded both Photon and Ubuntu OVA templates with Kubernetes v1.20.5, specify `--os-name ubuntu` to upgrade your management cluster to run on an Ubuntu VM.

   ```
   tanzu management-cluster upgrade --os-name ubuntu
   ```

   To skip the confirmation step when you upgrade a cluster, specify the `--yes` option.

   ```
   tanzu management-cluster upgrade --yes
   ```

   The upgrade process first upgrades the Cluster API providers for vSphere, Amazon EC2, or Azure that are running in the management cluster. Then, it upgrades the version of Kubernetes in all of the control plane and worker nodes of the management cluster.

   If the upgrade times out before it completes, run `tanzu management-cluster upgrade` again and specify the `--timeout` option with a value greater than the default of 30 minutes.

   ```
   tanzu management-cluster upgrade --timeout 45m0s
   ```

1. When the upgrade finishes, run the `tanzu cluster list` command with the `--include-management-cluster` option again to check that the management cluster has been upgraded.

   ```
   tanzu cluster list --include-management-cluster
   ```

   You see that the management cluster is now running the new version of Kubernetes, but that the Tanzu Kubernetes clusters are still running previous versions of Kubernetes.

   ```
     NAME                 NAMESPACE   STATUS    CONTROLPLANE  WORKERS  KUBERNETES         ROLES       PLAN
     k8s-1-17-13-cluster  default     running   1/1           1/1      v1.17.13+vmware.1  <none>      dev
     k8s-1-18-10-cluster  default     running   1/1           1/1      v1.18.10+vmware.1  <none>      dev
     k8s-1-19-3-cluster   default     running   1/1           1/1      v1.19.3+vmware.1   <none>      dev
     mgmt-cluster         tkg-system  running   1/1           1/1      v1.20.5+vmware.2   management  dev
   ```

## <a id="connectivity-api"></a> Replace Connectivity API with a Load Balancer

In Tanzu Kubernetes Grid v1.2 workload clusters access the Harbor service via a load balancer or a connectivity API installed in the management cluster.
Tanzu Kubernetes Grid v1.3 only supports a load balancer for Harbor access.
If your v1.2 installation uses the connectivity API, you need to remove it and set up a load balancer for the Harbor domain name before you upgrade:

1. Set up a load balancer for your workload clusters.
For vSphere, see [Install VMware NSX Advanced Load Balancer on a vSphere Distributed Switch](../mgmt-clusters/install-nsx-adv-lb.md).

1. Remove the `tkg-connectivity` operator and `tanzu-registry` webhook from the management cluster:

   1. Set the context of `kubectl` to the context of your management cluster:

      ```
      kubectl config use-context MGMT-CLUSTER-admin@MGMT-CLUSTER
      ```

      Where `MGMT-CLUSTER` is the name of your management cluster.

   1. Run the following commands to remove the resources and related objects:

      ```
      kubectl delete mutatingwebhookconfiguration tanzu-registry-webhook
      kubectl delete namespace tanzu-system-connectivity
      kubectl delete namespace tanzu-system-registry
      kubectl delete clusterrolebinding tanzu-registry-webhook
      kubectl delete clusterrole tanzu-registry-webhook
      kubectl delete clusterrolebinding tkg-connectivity-operator
      kubectl delete clusterrole tkg-connectivity-operator
      ```

1. Undo the effects of the `tanzu-registry` webhook:

   1. List all cluster control plane resources in the management cluster:

      ```
      kubectl get kubeadmcontrolplane -A | grep -v tkg-system
      ```

      These resources correspond to the workload clusters and shared service cluster.

   1. For each workload cluster control plane listed, run `kubectl edit` to edit its resource manifest as follows.
   For example, for a workload cluster `my_cluster_1`, run:

      ```
      kubectl -n NAMESPACE edit kubeadmcontrolplane my_cluster_1-control-plane
      ```

      Where `NAMESPACE` is the namespace of the management cluster.
      If the namespace is `default`, you can omit this option.

   1. In the manifest's `files:` section, delete `/opt/tkg/tanzu-registry-proxy.sh` from from the `files:` array.

   1. In the manifest's `preKubeadmCommands:` section, delete the two lines that start with the following commands:

      - `echo` - appends an IP address for Harbor into the `/etc/hosts` file.
      - `/opt/tkg/tanzu-registry-proxy.sh` - executes the `tanzu-registry-proxy.sh` script.

## <a id="update-callbackurl"></a> Update the Callback URL for Management Clusters with OIDC Authentication

In Tanzu Kubernetes Grid v1.3.0, Pinniped used Dex as the endpoint for both OIDC and LDAP providers. In Tanzu Kubernetes Grid v1.3.1 and later, Pinniped no longer requires Dex and uses the Pinniped endpoint for OIDC providers. In Tanzu Kubernetes Grid v1.3.1 and later, Dex is only used if you use an LDAP provider. If you used Tanzu Kubernetes Grid v1.3.0 to deploy management clusters that implement OIDC authentication, when you upgrade those management clusters to v1.3.1, the `dexsvc` service running in the management cluster is removed and replaced by the `pinniped-supervisor` service. Consequently, you must update the callback URLs that you specified in your OIDC provider after you deployed the management clusters with Tanzu Kubernetes Grid v1.3.0, so that it connects to the `pinniped-supervisor` service rather than to the `dexsvc` service.

### Obtain the Address of the Pinniped Service

Before you can update the callback URL, you must obtain the address of the Pinniped service that is running in the upgraded cluster.

1. Get the `admin` context of the management cluster.

   ```
   tanzu management-cluster kubeconfig get --admin
   ```

   If your management cluster is named `id-mgmt-test`, you should see the confirmation `Credentials of workload cluster 'id-mgmt-test' have been saved. You can now access the cluster by running 'kubectl config use-context id-mgmt-test-admin@id-mgmt-test'`. The `admin` context of a cluster gives you full access to the cluster without requiring authentication with your IDP.

1. Set `kubectl` to the `admin` context of the management cluster.

   ```
   kubectl config use-context id-mgmt-test-admin@id-mgmt-test
   ```

1. Get information about the services that are running in the management cluster.

    In Tanzu Kubernetes Grid v1.3.1 and later, the identity management service runs in the `pinniped-supervisor` namespace:

    ```
    kubectl get all -n pinniped-supervisor
    ```

    You see the following entry in the output:

    vSphere:

    ```
    NAME                          TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
    service/pinniped-supervisor   NodePort   100.70.70.12   <none>        5556:31234/TCP   84m
    ```

    Amazon EC2:

    ```
    NAME                          TYPE           CLUSTER-IP     EXTERNAL-IP                              PORT(S)         AGE
    service/pinniped-supervisor   LoadBalancer   100.69.13.66   ab1[...]71.eu-west-1.elb.amazonaws.com   443:30865/TCP   56m
    ```

    Azure:

    ```
    NAME                          TYPE           CLUSTER-IP       EXTERNAL-IP      PORT(S)         AGE
    service/pinniped-supervisor   LoadBalancer   100.69.169.220   20.54.226.44     443:30451/TCP   84m
    ```

1. Note the following information:

   - For management clusters that are running on vSphere, note the port on which the `pinniped-supervisor` service is running. In the example above, the port listed under `EXTERNAL-IP` is `31234`.
   - For clusters that you deploy to Amazon EC2 and Azure, note the external  address of the `LoadBalancer` node of the `pinniped-supervisor` service is running, that is listed under `EXTERNAL-IP`.

### Update the Callback URL

Once you have obtained information about the address at which `pinniped-supervisor` is running, you must update the callback URL for your OIDC provider. For example, if your IDP is Okta, perform the following steps:

1. Log in to your Okta account.
1. In the main menu, go to **Applications**.
1. Select the application that you created for Tanzu Kubernetes Grid.
1. In the General Settings panel, click **Edit**.
1. Under Login, update **Login redirect URIs** to include the address of the node in which the `pinniped-supervisor` is running.

    - On vSphere, update the `pinniped-supervisor` port number that you noted in the previous procedure.

     ```
     https://<API-ENDPOINT-IP>:31234/callback
     ```

    - On Amazon EC2 and Azure, update the external address of the `LoadBalancer` node on which the `pinniped-supervisor` is running, that you noted in the previous procedure.

     ```
     https://<EXTERNAL-IP>/callback
     ```

     Specify `https`, not `http`.
1. Click **Save**.

## <a id="what-next"></a> What to Do Next

You can now [upgrade the Tanzu Kubernetes clusters](workload-clusters.md) that this management cluster manages and [deploy new Tanzu Kubernetes clusters](../tanzu-k8s-clusters/deploy.md). By default, any new clusters that you deploy with this management cluster will run the new default version of Kubernetes.

However, if required, you can use the `tanzu cluster create` command with the `--tkr` option to deploy new clusters that run different versions of Kubernetes. For more information, see [Deploy Tanzu Kubernetes Clusters with Different Kubernetes Versions](../tanzu-k8s-clusters/k8s-versions.md).
