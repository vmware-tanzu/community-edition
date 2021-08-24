# Connect to and Examine Tanzu Kubernetes Clusters

After you have deployed Tanzu Kubernetes clusters, you use the `tanzu cluster list` and `tanzu cluster kubeconfig get` commands to obtain the list of running clusters and their credentials. Then, you can connect to the clusters by using `kubectl` and start working with your clusters.

## <a id="get"></a> Obtain Lists of Deployed Tanzu Kubernetes Clusters

To see lists of Tanzu Kubernetes clusters and the management clusters that manage them, use the `tanzu cluster list` command.

* To list all of the Tanzu Kubernetes clusters that are running in the `default` namespace of this management cluster, run the `tanzu cluster list` command.

   ```
   tanzu cluster list
   ```

   The output lists all of the Tanzu Kubernetes clusters that are managed by the management cluster. The output lists the cluster names, the namespace in which they are running, their current status, the numbers of actual and requested control plane and worker nodes, and the Kubernetes version that the cluster is running.

    ```
    NAME              NAMESPACE  STATUS   CONTROLPLANE  WORKERS  KUBERNETES        ROLES
    vsphere-cluster   default    running  1/1           1/1      v1.20.5+vmware.1  <none>
    vsphere-cluster2  default    running  1/1           1/1      v1.20.5+vmware.1  <none>
    my-vsphere-tkc    default    running  1/1           1/1      v1.20.5+vmware.1  <none>
    ```

   Clusters can be in the following states:

   - `creating`: The control plane is being created
   - `createStalled`: The process of creating control plane has stalled
   - `deleting`: The cluster is in the process of being deleted
   - `failed`: The creation of the control plane has failed
   - `running`: The control plane has initialized fully
   - `updating`: The cluster is in the process of rolling out an update or is scaling nodes
   - `updateFailed`: The cluster update process failed
   - `updateStalled`: The cluster update process has stalled
   - No status: The creation of the cluster has not started yet

   If a cluster is in a stalled state, check that there is network connectivity to the external registry, make sure that there are sufficient resources on the target platform for the operation to complete, and ensure that DHCP is issuing IPv4 addresses correctly.

* To list only those clusters that are running in a given namespace, specify the `--namespace` option.

   ```
   tanzu cluster list --namespace=my-namespace
   ```

* To include the current management cluster in the output of `tanzu cluster list`, specify the `--include-management-cluster` option.

   ```
   tanzu cluster list --include-management-cluster
   ```

   You can see that the management cluster is running in the `tkg-system` namespace and has the `management` role.

   ```
   NAME                  NAMESPACE   STATUS   CONTROLPLANE  WORKERS  KUBERNETES          ROLES       
   vsphere-cluster       default     running  1/1           1/1      v1.19.1+vmware.2    <none>
   vsphere-cluster2      default     running  3/3           3/3      v1.19.1+vmware.2    <none> 
   vsphere-mgmt-cluster  tkg-system  running  1/1           1/1      v1.19.1+vmware.2    management 
   ```

* To see all of the management clusters and change the context of the Tanzu CLI to a different management cluster, run the `tanzu login` command. See [List Management Clusters and Change Context](../cluster-lifecycle/multiple-management-clusters.md#login) for more information.

## <a id="output"></a> Export Tanzu Kubernetes Cluster Details to a File

You can export the details of the clusters that are managed by a management cluster in either JSON or YAML format. You can save the JSON or YAML to a file so that you can use it in scripts to run bulk operations on clusters.

1. To export cluster details as JSON, run `tanzu cluster list` with the `--output` option, specifying `json`.

   ```
   tanzu cluster list --output json
   ```

   The output shows the cluster information as JSON:

   ```
   [
     {
       "name": "vsphere-cluster",
       "namespace": "default",
       "status": "running",
       "plan": "",
       "controlplane": "1/1",
       "workers": "1/1",
       "kubernetes": "v1.19.1+vmware.2",
       "roles": []
     },
     {
       "name": "vsphere-cluster2",
       "namespace": "default",
       "status": "running",
       "plan": "",
       "controlplane": "3/3",
       "workers": "3/3",
       "kubernetes": "v1.19.1+vmware.2",
       "roles": []
     }
   ]
   ```

1. To export cluster details as YAML, run `tanzu cluster list` with the `--output` option, specifying `yaml`.

   ```
   tanzu cluster list --output yaml
   ```

   The output shows the cluster information as YAML:

   ```
   - name: vsphere-cluster
     namespace: default
     status: running
     plan: ""
     controlplane: 1/1
     workers: 1/1
     kubernetes: v1.19.1+vmware.2
     roles: []
   - name: vsphere-cluster2
     namespace: default
     status: running
     plan: ""
     controlplane: 3/3
     workers: 3/3
     kubernetes: v1.19.1+vmware.2
     roles: []
   ```

1. Save the output as a file.

   ```
   tanzu cluster list --output json > clusters.json
   ```

   ```
   tanzu cluster list --output yaml > clusters.yaml
   ```

For how to save the details of multiple management clusters, including their context and `kubeconfig` files, see [Save Management Cluster Details to a File](../cluster-lifecycle/multiple-management-clusters.md#output).

## <a id="kubeconfig"></a> Retrieve Tanzu Kubernetes Cluster `kubeconfig`

After you create a Tanzu Kubernetes cluster, you can obtain its cluster, context, and user `kubeconfig` settings by running the `tanzu cluster kubeconfig get` command, specifying the name of the cluster.

By default, the command adds the cluster's `kubeconfig` settings to your current `kubeconfig` file.

To generate a standalone _administrator_ `kubeconfig` file with embedded credentials, add the `--admin` option. This `kubeconfig` file grants its user full access to the cluster's resources and lets them access the cluster without logging in to an identity provider.

**IMPORTANT**: If identity management is not configured on the cluster, you must specify the `--admin` option.

```
tanzu cluster kubeconfig get my-cluster --admin
```

You should see the following output:

```
You can now access the cluster by running 'kubectl config use-context my-cluster-admin@my-cluster'

```

If identity management is enabled on a cluster, you can generate a regular `kubeconfig` that requires the user to authenticate with your external identity provider, and grants them access to cluster resources based on their assigned roles. In this case, run `tanzu cluster kubeconfig get` without the `--admin` option.

```
tanzu cluster kubeconfig get my-cluster
```

You should see the following output:

```
You can now access the cluster by running 'kubectl config use-context tanzu-cli-my-cluster@my-cluster'
```

If the cluster is running in a namespace other than the `default` namespace, you must specify the `--namespace` option to get the credentials of that cluster.

```
tanzu cluster kubeconfig get my-cluster --namespace=my-namespace
```

To save the configuration information in a standalone `kubeconfig` file, for example to distribute them to developers, specify the `--export-file` option.
This `kubeconfig` file requires the user to authenticate with an external identity provider, and grants access to cluster resources based on their assigned roles.

```
tanzu cluster kubeconfig get my-cluster --export-file my-cluster-credentials
```

**IMPORTANT**: By default, unless you specify the `--export-file` option to save the `kubeconfig` for a cluster to a specific file, the credentials for all clusters that you deploy from the Tanzu CLI are added to a shared `kubeconfig` file. If you delete the shared `kubeconfig` file, all clusters become unusable.

To retrieve a `kubeconfig` for a management cluster, run `tanzu management-cluster kubeconfig get` as described in [Retrieve Tanzu Kubernetes Cluster `kubeconfig`](../mgmt-clusters/verify-deployment.md#kubeconfig).

## <a id="id-mgmt"></a> Authenticate Connections to a Workload Cluster

If you deployed the management cluster with identity management enabled or enabled identity management on the management cluster as a post-deployment step, any workload clusters that you create from your management cluster are automatically configured to use the same identity management service. When you provide users with the `admin` `kubeconfig` for a management cluster or workload cluster, they have full access to the cluster and do not need to be authenticated. However, if you provide users with the regular `kubeconfig`, they must have a user account in your OIDC or LDAP identity provider and you must configure Role-Based Access Control (RBAC) on the cluster to grant access permissions to the designated user.

The authentication process requires a browser to be present on the machine from which users connect to clusters, because running `kubectl` commands automatically opens the IDP login page so that users can log in to the cluster. If the machine on which you are running `tanzu` and `kubectl` commands does not have a browser, see [Authenticate Users on a Machine Without a Browser](#no-browser) below.

To authenticate users on a workload cluster on which identity management is enabled, perform the following steps.

1. Obtain the regular `kubeconfig` for the workload cluster and export it to a file.

    This example exports the `kubeconfig` for the cluster `my-cluster` to the file `my-cluster-credentials`.

    ```
    tanzu cluster kubeconfig get my-cluster --export-file my-cluster-credentials
    ```

1. Use the generated file to attempt to run an operation on the cluster.

   For example, run:

   ```
   kubectl get pods -A --kubeconfig my-cluster-credentials
   ```

    You should be redirected to the log in page for your identity provider.

    After successfully logging in with a user account from your identity provider, if you already configured a role binding on the cluster for the authenticated user, the output shows the pod information.

   If you have not configured a role binding on the cluster, you see the message `Error from server (Forbidden): pods is forbidden: User "<user>" cannot list resource "pods" in API group "" at the cluster scope`. This happens because this user does not have any permissions on the cluster yet. To authorize the user to access the cluster resources, you must [Configure a Role Binding](#rbac) on the cluster.

### <a id="no-browser"></a> Authenticate Users on a Machine Without a Browser

If the machine on which you are running `tanzu` and `kubectl` commands does not have a browser, you can skip the automatic opening of a browser during the authentication process.  

1. If it is not set already, set the `TANZU_CLI_PINNIPED_AUTH_LOGIN_SKIP_BROWSER=true` environment variable.

   This adds the `--skip-browser` option to the `kubeconfig` for the cluster.

   ```
   export TANZU_CLI_PINNIPED_AUTH_LOGIN_SKIP_BROWSER=true
   ```

   On Windows systems, use the `SET` command instead of `export`.
1. Export the regular `kubeconfig` for the cluster to the local file `my-cluster-credentials`.

   Note that the command does not include the `--admin` option, so the `kubeconfig` that is exported is the regular `kubeconfig`, not the `admin` version.

   ```
   tanzu cluster kubeconfig get my-cluster --export-file my-cluster-credentials
   ```

1. Connect to the cluster by using the newly-created `kubeconfig` file.

   ```
   kubectl get pods -A --kubeconfig my-cluster-credentials
   ```

   The login URL is displayed in the terminal. For example:

   ```
   Please log in: https://ab9d82be7cc2443ec938e35b69862c9c-10577430.eu-west-1.elb.amazonaws.com/oauth2/authorize?access_type=offline&client_id=pinniped-cli&code_challenge=vPtDqg2zUyLFcksb6PrmE8bI9qF8it22KQMy52hB6DE&code_challenge_method=S256&nonce=2a66031e3075c65ea0361b3ba30bf174&redirect_uri=http%3A%2F%2F127.0.0.1%3A57856%2Fcallback&response_type=code&scope=offline_access+openid+pinniped%3Arequest-audience&state=01064593f32051fee7eff9333389d503
   ```  

1. Copy the login URL and paste it into a browser on a machine that does have one.
1. In the browser, log in to your identity provider.

  You will see a message that the identity provider could not send the authentication code because there is no localhost listener on your workstation.

1. Copy the URL of the authenticated session from the URL field of the browser.
1. On the machine that does not have a browser, use the URL that you copied in the preceding step to get the authentication code from the identity provider.

   ```
   curl -L '<copied_URL>'
   ```

   Wrap the URL in quotes, to escape any special characters. For example, the command will resemble the following:

   ```
   curl - L 'http://127.0.0.1:37949/callback?code=FdBkopsZwYX7w5zMFnJqYoOlJ50agmMWHcGBWD-DTbM.8smzyMuyEBlPEU2ZxWcetqkStyVPjdjRgJNgF1-vODs&scope=openid+offline_access+pinniped%3Arequest-audience&state=a292c262a69e71e06781d5e405d42c03'
   ```

   After running `curl -L '<copied_URL>'`, you should see the following message:

   ```
   you have been logged in and may now close this tab
   ```

1. Connect to the cluster again by using the same `kubeconfig` file as you used previously.

   ```
   kubectl get pods -A --kubeconfig my-cluster-credentials
   ```

   If you already configured a role binding on the cluster for the authenticated user, the output shows the pod information.

   If you have not configured a role binding on the cluster, you will see a message denying the user account access to the pods: `Error from server (Forbidden): pods is forbidden: User "user@example.com" cannot list resource "pods" in API group "" at the cluster scope`. This happens because the user has been successfully authenticated, but they are not yet authorized to access any resources on the cluster. To authorize the user to access the cluster resources, you must configure Role-Based Access Control (RBAC) on the cluster by creating a cluster role binding.

## <a id="rbac"></a> Configure a Role Binding on a Workload Cluster

To complete the identity management configuration of the workload cluster, you must create cluster role bindings for the users who use the `kubeconfig` that you generated in the preceding step. There are many roles with which you can associate users, but the most useful roles are the following:

- `cluster-admin`: Can perform any operation on the cluster.
- `admin`: Permission to view most resources but can only modify resources like roles and bindings. Cannot modify pods or deployments.
- `edit`: The opposite of `admin`. Can create, update, and delete resources like deployments, services, and pods. Cannot change roles or permissions.
- `view`: Read-only.

You can assign any of these roles to users. For more information about RBAC and cluster role bindings, see [Using RBAC Authorization](https://kubernetes.io/docs/reference/access-authn-authz/rbac/) in the Kubernetes documentation.

1. Set the `kubectl` context to the workload cluster's `admin` `kubeconfig`.

   You need to switch to the workload cluster's `admin` context so that you can create a role binding. For example, run the following two commands to change to the `admin` context:

   Get the `kubeconfig`:

   ```
   tanzu cluster kubeconfig get my-cluster --admin
   ```

   Switch context:

   ```
   kubectl config use-context my-cluster-admin@my-cluster
   ```

1. To see the full list of roles that are available on the cluster, run the following command:

   ```
   kubectl get clusterroles
   ```

1. Create a cluster role binding to associate a given user with a role on the cluster.

   The following command creates a role binding named `workload-test-rb` that binds the role `cluster-admin` for this cluster to the user `user@example.com`. For OIDC the username is usually the email address of the user. For LDAPS it is the LDAP username, not the email address.

   **OIDC**:

   ```
   kubectl create clusterrolebinding workload-test-rb --clusterrole cluster-admin --user user@example.com
   ```

   **LDAP**:  

   ```
   kubectl create clusterrolebinding workload-test-rb --clusterrole cluster-admin --user <username>
   ```

1. Use the regular `kubeconfig` file that you generated above to attempt to run an operation on the cluster again.

   For example, run:

   ```
   kubectl get pods -A --kubeconfig my-cluster-credentials
   ```

   This time, you should see the list of pods that are running in the workload cluster. This is because the user of the `my-cluster-credentials` `kubeconfig` file has both been authenticated by your identity provider, and has the necessary permissions on the cluster. You can share the `my-cluster-credentials` `kubeconfig` file with any users for whom you configure role bindings on the cluster.

For information about how to configure RBAC on management clusters, see [Configure Identity Management After Management Cluster Deployment](../mgmt-clusters/configure-id-mgmt.md).

## <a id="examine"></a> Examine the Deployed Cluster

1. After you have added the credentials to your `kubeconfig`, you can connect to the cluster by using `kubectl`.

    ```
    kubectl config use-context my-cluster-admin@my-cluster
    ```

1. Use `kubectl` to see the status of the nodes in the cluster.

   ```
   kubectl get nodes
   ```

   For example, if you deployed the `my-prod-cluster` in [Deploy a Cluster with a Highly Available Control Plane](../tanzu-k8s-clusters/deploy.md#ha-control-plane) with the `prod` plan and the default 3 control plane nodes and worker nodes, you see the following output.

    ```
    NAME                                    STATUS   ROLES    AGE     VERSION
    my-prod-cluster-control-plane-gp4rl     Ready    master   8m51s   v1.19.1+vmware.2
    my-prod-cluster-control-plane-n8bh7     Ready    master   5m58s   v1.19.1+vmware.2
    my-prod-cluster-control-plane-xflrg     Ready    master   3m39s   v1.19.1+vmware.2
    my-prod-cluster-md-0-6946bcb48b-dk7m6   Ready    <none>   6m45s   v1.19.1+vmware.2
    my-prod-cluster-md-0-6946bcb48b-dq8s9   Ready    <none>   7m23s   v1.19.1+vmware.2
    my-prod-cluster-md-0-6946bcb48b-nrdlp   Ready    <none>   7m8s    v1.19.1+vmware.2
    
    ```

    Because networking with Antrea is enabled by default in Tanzu Kubernetes clusters, all clusters are in the `Ready` state without requiring any additional configuration.
1. Use `kubectl` to see the status of the pods running in the cluster.

   ```
   kubectl get pods -A
   ```

   The example below shows the pods running in the `kube-system` namespace in the `my-prod-cluster` cluster on vSphere.

    ```
    NAMESPACE     NAME                                                             READY   STATUS    RESTARTS   AGE
    kube-system   antrea-agent-2mw42                                               2/2     Running   0          4h41m
    kube-system   antrea-agent-4874z                                               2/2     Running   1          4h45m
    kube-system   antrea-agent-9qfr6                                               2/2     Running   0          4h48m
    kube-system   antrea-agent-cf7cf                                               2/2     Running   0          4h46m
    kube-system   antrea-agent-j84mz                                               2/2     Running   0          4h46m
    kube-system   antrea-agent-rklbg                                               2/2     Running   0          4h46m
    kube-system   antrea-controller-5d594c5cc7-5pttm                               1/1     Running   0          4h48m
    kube-system   coredns-5bcf65484d-7dp8d                                         1/1     Running   0          4h48m
    kube-system   coredns-5bcf65484d-pzw8p                                         1/1     Running   0          4h48m
    kube-system   etcd-my-prod-cluster-control-plane-frsgd                         1/1     Running   0          4h48m
    kube-system   etcd-my-prod-cluster-control-plane-khld4                         1/1     Running   0          4h44m
    kube-system   etcd-my-prod-cluster-control-plane-sjvx7                         1/1     Running   0          4h41m
    kube-system   kube-apiserver-my-prod-cluster-control-plane-frsgd               1/1     Running   0          4h48m
    kube-system   kube-apiserver-my-prod-cluster-control-plane-khld4               1/1     Running   1          4h45m
    kube-system   kube-apiserver-my-prod-cluster-control-plane-sjvx7               1/1     Running   0          4h41m
    kube-system   kube-controller-manager-my-prod-cluster-control-plane-frsgd      1/1     Running   1          4h48m
    kube-system   kube-controller-manager-my-prod-cluster-control-plane-khld4      1/1     Running   0          4h45m
    kube-system   kube-controller-manager-my-prod-cluster-control-plane-sjvx7      1/1     Running   0          4h41m
    kube-system   kube-proxy-hzqlt                                                 1/1     Running   0          4h48m
    kube-system   kube-proxy-jr4w6                                                 1/1     Running   0          4h45m
    kube-system   kube-proxy-lx8bp                                                 1/1     Running   0          4h46m
    kube-system   kube-proxy-rzbgh                                                 1/1     Running   0          4h46m
    kube-system   kube-proxy-s684n                                                 1/1     Running   0          4h41m
    kube-system   kube-proxy-z9v9t                                                 1/1     Running   0          4h46m
    kube-system   kube-scheduler-my-prod-cluster-control-plane-frsgd               1/1     Running   1          4h48m
    kube-system   kube-scheduler-my-prod-cluster-control-plane-khld4               1/1     Running   0          4h45m
    kube-system   kube-scheduler-my-prod-cluster-control-plane-sjvx7               1/1     Running   0          4h41m
    kube-system   kube-vip-my-prod-cluster-control-plane-frsgd                     1/1     Running   1          4h48m
    kube-system   kube-vip-my-prod-cluster-control-plane-khld4                     1/1     Running   0          4h45m
    kube-system   kube-vip-my-prod-cluster-control-plane-sjvx7                     1/1     Running   0          4h41m
    kube-system   vsphere-cloud-controller-manager-4nlsw                           1/1     Running   0          4h41m
    kube-system   vsphere-cloud-controller-manager-gw7ww                           1/1     Running   2          4h48m
    kube-system   vsphere-cloud-controller-manager-vp968                           1/1     Running   0          4h44m
    kube-system   vsphere-csi-controller-555595b64c-l82kb                          5/5     Running   3          4h48m
    kube-system   vsphere-csi-node-5zq47                                           3/3     Running   0          4h41m
    kube-system   vsphere-csi-node-8fzrg                                           3/3     Running   0          4h46m
    kube-system   vsphere-csi-node-8zs5l                                           3/3     Running   0          4h45m
    kube-system   vsphere-csi-node-f2v55                                           3/3     Running   0          4h46m
    kube-system   vsphere-csi-node-khtwv                                           3/3     Running   0          4h48m
    kube-system   vsphere-csi-node-shtqj                                           3/3     Running   0          4h46m
    ```

    You can see from the example above that the following services are running in the `my-prod-cluster` cluster:

    - [Antrea](https://antrea.io/), the container networking interface
    - [`coredns`](https://coredns.io/), for DNS
    - [`etcd`](https://etcd.io/), for key-value storage
    - [`kube-apiserver`](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/), the Kubernetes API server
    - [`kube-proxy`](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-proxy/), the Kubernetes network proxy
    - [`kube-scheduler`](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-scheduler/), for scheduling and availability
    - [`vsphere-cloud-controller-manager`](https://cloud-provider-vsphere.sigs.k8s.io/), the Kubernetes cloud provider for vSphere
    - [`kube-vip`](https://kube-vip.io/), load balancing services for the Cluster API server
    - `vsphere-csi-controller` and `vsphere-csi-node`, the [container storage interface for vSphere](https://cloud-provider-vsphere.sigs.k8s.io/container_storage_interface.html)
