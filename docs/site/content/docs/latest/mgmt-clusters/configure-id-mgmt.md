# Configure Identity Management After Management Cluster Deployment

If you enabled identity management when you deployed a management cluster, you must perform additional post-deployment steps on the management cluster so that authenticated users can access it.

To configure identity management on a management cluster, you must perform the following steps:

- Make sure that the authentication service is running correctly.
- For OIDC deployments, provide the callback URL for the management cluster to your OIDC identity provider.
- Generate a `kubeconfig` file to share with users by running `tanzu management-cluster kubeconfig get` with the `--export-file` option.
  - You can generate an administrator `kubeconfig` that contains embedded credentials, or a regular `kubeconfig` that prompts users to authenticate with an external identity provider.
  - See [Retrieve Management Cluster `kubeconfig`](./verify-deployment.md#kubeconfig) for more information.
- Set up role-based access control (RBAC) by creating a role binding on the management cluster, that assigns role-based permissions to individual authenticated users or user groups.

## Prerequisites

- You have deployed a management cluster with either OIDC or LDAPS identity management configured.
- If you configured an OIDC server as the identity provider, you have followed the procedures in [Enabling Identity Management in Tanzu Kubernetes Grid](enabling-id-mgmt.md) to add users in the OIDC server.

## Connect `kubectl` to the Management Cluster

To configure identity management, you must obtain and use the `admin` context of the management cluster.

1. Get the `admin` context of the management cluster.

   The procedures in this topic use a management cluster named `id-mgmt-test`.

   ```
   tanzu management-cluster kubeconfig get id-mgmt-test --admin
   ```

   If your management cluster is named `id-mgmt-test`, you should see the confirmation `Credentials of workload cluster 'id-mgmt-test' have been saved. You can now access the cluster by running 'kubectl config use-context id-mgmt-test-admin@id-mgmt-test'`. The `admin` context of a cluster gives you full access to the cluster without requiring authentication with your IDP.

1. Set `kubectl` to the `admin` context of the management cluster.

   ```
   kubectl config use-context id-mgmt-test-admin@id-mgmt-test
   ```

The next steps depend on whether you are using an OIDC or LDAP identity management service.

- [Check the Status of an OIDC Identity Management Service](#oidc)
- [Check the Status of an LDAP Identity Management Service](#ldap)

## <a id="oidc"></a> Check the Status of an OIDC Identity Management Service

In Tanzu Kubernetes Grid v1.3.0, Pinniped used Dex as the endpoint for both OIDC and LDAP providers. In Tanzu Kubernetes Grid v1.3.1 and later, Pinniped with OIDC no longer requires Dex. In Tanzu Kubernetes Grid v1.3.1 and later, Dex is only used if you use an LDAP provider. Because of this change, the way in which you check the status of an OIDC identity management service is different in Tanzu Kubernetes Grid v1.3.1 and later compared to Tanzu Kubernetes Grid v1.3.0.

For new management cluster deployments with OIDC authentication, it is **strongly recommended** to use Tanzu Kubernetes Grid v1.3.1 or later.

When you check the status of the service, you must note the address at which the service is exposed to your OIDC identity provider.

1. Get information about the services that are running in the management cluster.

    **Tanzu Kubernetes Grid 1.3.1 or later**:

    In Tanzu Kubernetes Grid v1.3.1 and later, the identity management service runs in the `pinniped-supervisor` namespace:

    ```
    kubectl get all -n pinniped-supervisor
    ```

    You see the following entry in the output:

    vSphere:

    ```
    NAME             TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
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

    **Tanzu Kubernetes Grid 1.3.0**:

    In Tanzu Kubernetes Grid v1.3.0, the identity management service runs in the `tanzu-system-auth` namespace:

    ```
    kubectl get all -n tanzu-system-auth
    ```

    You see the following entry in the output:

    vSphere:

    ```
    NAME             TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
    service/dexsvc   NodePort   100.70.70.12   <none>        5556:30167/TCP   84m
    ```

    Amazon EC2:

    ```
    NAME             TYPE           CLUSTER-IP       EXTERNAL-IP                                         PORT(S)         AGE
    service/dexsvc   LoadBalancer   100.65.184.107   a6e[...]cc6-921316974.eu-west-1.elb.amazonaws.com   443:32547/TCP   84m
    ```

    Azure:

    ```
    NAME             TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)         AGE
    service/dexsvc   LoadBalancer   100.69.169.220   20.54.226.44  443:30451/TCP   84m
    ```

1. Note the following information:

   - For management clusters that are running on vSphere, note the port on which the `pinniped-supervisor` or `dexsvc` service is running. In the example above, the port listed under `EXTERNAL-IP` is `31234` for the `pinniped-supervisor` service in Tanzu Kubernetes Grid v1.3.1 and later, or `30167` for the `dexsvc` service in v1.3.0.
   - For clusters that you deploy to Amazon EC2 and Azure, note the external  address of the `LoadBalancer` node of the `pinniped-supervisor` or `dexsvc` service is running, that is listed under `EXTERNAL-IP`.

1. Check that all services in the management cluster are running.

    ```
    kubectl get pods -A
    ```

    It can take several minutes for the Pinniped service to be up and running. For example, on Amazon EC2 and Azure deployments the service must wait for the `LoadBalancer` IP addresses to be ready. Wait until you see that `pinniped-post-deploy-job` is completed before you proceed to the next steps.

    ```
    NAMESPACE             NAME                                   READY  STATUS      RESTARTS  AGE
    [...]
    pinniped-supervisor   pinniped-post-deploy-job-hq8fc         0/1    Completed   0         85m
    ```

**NOTE**: You are able to run `kubectl get pods` because you are using the `admin` context for the management cluster. Users who attempt to connect to the management cluster with the regular context will not be able to access its resources, because they are not yet authorized to do so.

## <a id="ldap"></a> Check the Status of an LDAP Identity Management Service

If you use an LDAP identity management service, Pinniped uses Dex as the endpoint to expose to your provider. In Tanzu Kubernetes Grid v1.3.0, Pinniped uses Dex as the endpoint for both OIDC and LDAP providers. In Tanzu Kubernetes Grid v1.3.1 and later, Dex is only used if you use an LDAP provider. This procedure applies to LDAP identity management for all v1.3.x versions, and to OIDC identity management for Tanzu Kubernetes Grid v1.3.0. If you are using OIDC identity management with Tanzu Kubernetes Grid v1.3.1 or later, see [Check the Status of an OIDC Identity Management Service (v1.3.1 and later)](#oidc) above.

1. Get information about the services that are running in the management cluster in the `tanzu-system-auth` namespace.

    ```
    kubectl get all -n tanzu-system-auth
    ```

    You see the following entry in the output:

    vSphere:

    ```
    NAME             TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
    service/dexsvc   NodePort   100.70.70.12   <none>        5556:30167/TCP   84m
    ```

    Amazon EC2:

    ```
    NAME             TYPE           CLUSTER-IP       EXTERNAL-IP                              PORT(S)         AGE
    service/dexsvc   LoadBalancer   100.65.184.107   a6e[...]74.eu-west-1.elb.amazonaws.com   443:32547/TCP   84m
    ```

    Azure:

    ```
    NAME             TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)         AGE
    service/dexsvc   LoadBalancer   100.69.169.220   20.54.226.44  443:30451/TCP   84m
    ```

1. Check that all services in the management cluster are running.

    ```
    kubectl get pods -A
    ```

    It can take several minutes for the Pinniped service to be up and running. For example, on Amazon EC2 and Azure deployments the service must wait for the `LoadBalancer` IP addresses to be ready. Wait until you see that `pinniped-post-deploy-job` is completed before you proceed to the next steps.

    ```
    NAMESPACE             NAME                                   READY  STATUS      RESTARTS  AGE
    [...]
    pinniped-supervisor   pinniped-post-deploy-job-hq8fc         0/1    Completed   0         85m
    ```

**NOTE**: You are able to run `kubectl get pods` because you are using the `admin` context for the management cluster. Users who attempt to connect to the management cluster with the regular context will not be able to access its resources, because they are not yet authorized to do so.

## <a id="callback"></a> Provide the Callback URI to the OIDC Provider

If you configured an LDAP server as your identity provider, you do not need to configure a callback URI. For the next steps, go to [Generate a `kubeconfig` to Allow Authenticated Users to Connect to the Management Cluster](#gen-kubeconfig).

If you configured the management cluster to use OIDC authentication, you must provide the callback URI for that management cluster to your OIDC identity provider.

For example, if you are using OIDC and your IDP is Okta, perform the following steps:

1. Log in to your Okta account.
1. In the main menu, go to **Applications**.
1. Select the application that you created for Tanzu Kubernetes Grid.
1. In the General Settings panel, click **Edit**.
1. Under Login, update **Login redirect URIs** to include the address of the node in which the `pinniped-supervisor` is running:

  **NOTE**: In Tanzu Kubernetes Grid v1.3.0, you must provide the address of the `dexsvc` node. The port number of the API endpoint is also different for the `pinniped-supervisor` and `dexsvc` services.

  - On vSphere, add the IP address that you set as the API endpoint and the `pinniped-supervisor` or `dexsvc` port number that you noted in the previous procedure.

     - Tanzu Kubernetes Grid v1.3.1 and later:

         ```
         https://<API-ENDPOINT-IP>:31234/callback
         ```

     - Tanzu Kubernetes Grid v1.3.0:

         ```
         https://<API-ENDPOINT-IP>:30167/callback
         ```

  - On Amazon EC2 and Azure, add the external IP address of the `LoadBalancer` node on which the `pinniped-supervisor` or `dexsvc` is running, that you noted in the previous procedure.

     NOTE:

     ```
     https://<EXTERNAL-IP>/callback
     ```

     In all cases, you must specify `https`, not `http`.
1. Click **Save**.

## <a id="gen-kubeconfig"></a> Generate a `kubeconfig` to Allow Authenticated Users to Connect to the Management Cluster

To allow users to access the management cluster, you export the management cluster's `kubeconfig` to a file that you can share with those users.

- If you export the `admin` version of the `kubeconfig`, any users with whom you share it will have full access to the management cluster and IDP authentication is bypassed.
- If you export the regular version of the `kubeconfig`, it is populated with the necessary authentication information, so that the user's identity will be verified with your IDP before they can access the cluster's resources.

This procedure allows you to test the login step of the authentication process if a browser is present on the machine on which you are running `tanzu` and `kubectl` commands. If the machine does not have a browser, see [Authenticate Users on a Machine Without a Browser](#no-browser) below.

1. Export the regular `kubeconfig` for the management cluster to the local file `/tmp/id_mgmt_test_kubeconfig`.

   Note that the command does not include the `--admin` option, so the `kubeconfig` that is exported is the regular `kubeconfig`, not the `admin` version.

   ```
   tanzu management-cluster kubeconfig get --export-file /tmp/id_mgmt_test_kubeconfig
   ```

   You should see confirmation that `You can now access the cluster by running 'kubectl config use-context tanzu-cli-id-mgmt-test@id-mgmt-test' under path '/tmp/id_mgmt_test_kubeconfig'`.

1. Connect to the management cluster by using the newly-created `kubeconfig` file.

   ```
   kubectl get pods -A --kubeconfig /tmp/id_mgmt_test_kubeconfig
   ```

   The authentication process requires a browser to be present on the machine from which users connect to clusters, because running `kubectl` commands automatically opens the IDP login page so that users can log in to the cluster.

   Your browser should open and display the login page for your OIDC provider or an LDAPS login page.

   LDAPS:

   ![LDAPS login page](../images/id-mgmt-ldap-login.png)

   OIDC:

   ![OIDC login page](../images/id-mgmt-oidc-login.png)

   Enter the credentials of a user account that exists in your OIDC or LDAP server.

   After a successful login, the browser should display the following message:

   ```
   you have been logged in and may now close this tab
   ```

1. Go back to the terminal in which you run `tanzu` and `kubectl` commands.

   If you already configured a role binding on the cluster for the authenticated user, the output of `kubectl get pods -A`  appears, displaying the pod information.

   If you have not configured a role binding on the cluster, you see a message denying the user account access to the pods: `Error from server (Forbidden): pods is forbidden: User "user@example.com" cannot list resource "pods" in API group "" at the cluster scope`. This happens because the user has been successfully authenticated, but they are not yet authorized to access any resources on the cluster. To authorize the user to access the cluster resources, you must [Create a Role Binding on the Management Cluster](#create-rolebinding).

### <a id="no-browser"></a> Authenticate Users on a Machine Without a Browser

If the machine on which you are running `tanzu` and `kubectl` commands does not have a browser, you can skip the automatic opening of a browser during the authentication process.  

1. Set the `TANZU_CLI_PINNIPED_AUTH_LOGIN_SKIP_BROWSER=true` environment variable.

   This adds the `--skip-browser` option to the `kubeconfig` for the cluster.

   ```
   export TANZU_CLI_PINNIPED_AUTH_LOGIN_SKIP_BROWSER=true
   ```

   On Windows systems, use the `SET` command instead of `export`.
1. Export the regular `kubeconfig` for the management cluster to the local file `/tmp/id_mgmt_test_kubeconfig`.

   Note that the command does not include the `--admin` option, so the `kubeconfig` that is exported is the regular `kubeconfig`, not the `admin` version.

   ```
   tanzu management-cluster kubeconfig get --export-file /tmp/id_mgmt_test_kubeconfig
   ```

   You should see confirmation that `You can now access the cluster by running 'kubectl config use-context tanzu-cli-id-mgmt-test@id-mgmt-test' under path '/tmp/id_mgmt_test_kubeconfig'`.
1. Connect to the management cluster by using the newly-created `kubeconfig` file.

   ```
   kubectl get pods -A --kubeconfig /tmp/id_mgmt_test_kubeconfig
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

1. Connect to the management cluster again by using the same `kubeconfig` file as you used previously.

   ```
   kubectl get pods -A --kubeconfig /tmp/id_mgmt_test_kubeconfig
   ```

   If you already configured a role binding on the cluster for the authenticated user, the output shows the pod information.

   If you have not configured a role binding on the cluster, you will see a message denying the user account access to the pods: `Error from server (Forbidden): pods is forbidden: User "user@example.com" cannot list resource "pods" in API group "" at the cluster scope`. This happens because the user has been successfully authenticated, but they are not yet authorized to access any resources on the cluster. To authorize the user to access the cluster resources, you must configure Role-Based Access Control (RBAC) on the cluster by creating a cluster role binding.

## <a id="create-rolebinding"></a> Create a Role Binding on the Management Cluster

To complete the identity management configuration of the management cluster, you must create cluster role bindings for the users who use the `kubeconfig` that you generated in the preceding step. There are many roles with which you can associate users, but the most useful roles are the following:

- `cluster-admin`: Can perform any operation on the cluster.
- `admin`: Permission to view most resources but can only modify resources like roles and bindings. Cannot modify pods or deployments.
- `edit`: The opposite of `admin`. Can create, update, and delete resources like deployments, services, and pods. Cannot change roles or permissions.
- `view`: Read-only.

You can assign any of these roles to users. For more information about RBAC and cluster role bindings, see [Using RBAC Authorization](https://kubernetes.io/docs/reference/access-authn-authz/rbac/) in the Kubernetes documentation.

1. Make sure that you are using the `admin` context of the management cluster.

   ```
   kubectl config current-context
   ```

   If the context is not the management cluster `admin` context, set `kubectl` to use that context. For example:

   ```
   kubectl config use-context id-mgmt-test-admin@id-mgmt-test
   ```

1. To see the full list of roles that are available on a cluster, run the following command:

   ```
   kubectl get clusterroles
   ```

1. Create a cluster role binding to associate a given user with a role.

   The following command creates a role binding named `id-mgmt-test-rb` that binds the role `cluster-admin` for this cluster to the user `user@example.com`. For OIDC the username is usually the email address of the user. For LDAPS it is the LDAP username, not the email address.

   **OIDC**:

   ```
   kubectl create clusterrolebinding id-mgmt-test-rb --clusterrole cluster-admin --user user@example.com
   ```

   **LDAP**:  

   ```
   kubectl create clusterrolebinding id-mgmt-test-rb --clusterrole cluster-admin --user <username>
   ```

1. Attempt to connect to the management cluster again by using the `kubeconfig` file that you created in the previous procedure.

   ```
   kubectl get pods -A --kubeconfig /tmp/id_mgmt_test_kubeconfig
   ```

   This time, because the user is bound to the `cluster-admin` role on this management cluster, the list of pods should be displayed.

## <a id="lb-mgmt-cluster"></a> Add a Load Balancer Service to a Management Cluster on vSphere

When deploying a management cluster on vSphere, you may want to use a load balancer service with the identity management services provided by Pinniped for OIDC or by Pinniped and Dex for LDAP. Setting up a load balancer service on the management cluster for Pinniped and Dex identity management services can simplify your DNS and firewall configuration requirements.

This procedure modifies Pinniped by updating the app secret that contains the deployment configuration. This update ensures that any configuration changes made to Pinniped components are preserved during future upgrades of the management cluster.

### Prerequisites

Before you begin this procedure, you must have the following:

* An external load balancer service configured and available for use as a provider by the management cluster. For example, you may have set up MetalLB.
* Successfully installed and configured Pinniped for OIDC or Pinniped and Dex for LDAP on the management cluster.

### Procedure

1. Make sure that you are using the `admin` context of the management cluster.

   ```
   kubectl config current-context
   ```

   If the context is not the management cluster `admin` context, set `kubectl` to use that context. For example:

   ```
   kubectl config use-context id-mgmt-test-admin@id-mgmt-test
   ```

1. Obtain the name of the secret containing the Pinniped configuration.

   ```
   kubectl get secret -n tkg-system -l tkg.tanzu.vmware.com/addon-name=pinniped
   ```

1. Save the existing values from the secret.

   ```
   kubectl get secret tkg-mgmt-vc-pinniped-addon -n tkg-system -o jsonpath={.data.values\\.yaml} | base64 -D > values.yaml
   ```

1. Update the configuration for the Dex service located in `values.yaml`. This step is not required if you are using OIDC because Dex is not enabled for OIDC.

   **LDAP only**
   1. Locate and copy the entire section for Dex that is labeled as follows:  

      ```
      dex:
         app: dex
         ....
         ....
      ```

   1. In the Dex configuration section, add or update a `service:` section to include `name: dexsvc` and `type: LoadBalancer`. This configuration update should resemble the following:

      ```
      dex:
         app: dex
         ....
         ....
         service:
           name: dexsvc
           type: LoadBalancer
      ```

1. In a text editor, prepare the `stringData` section for the secret update. Make sure indentation matches the examples below.

   **OIDC:**

   ```
   stringData:
     overlays.yaml: |
       #@ load("@ytt:overlay", "overlay")
       #@overlay/match by=overlay.subset({"kind": "Service", "metadata": {"name": "pinniped-supervisor", "namespace": "pinniped-supervisor"}})
       ---
       #@overlay/replace
       spec:
         type: LoadBalancer
         selector:
           app: pinniped-supervisor
         ports:
           - name: https
             protocol: TCP
             port: 443
             targetPort: 8443
      values.yaml: |
        #@data/values
        #@overlay/match-child-defaults missing_ok=True
        ---
        infrastructure_provider: vsphere
        tkg_cluster_role: management
    ```

   **LDAP:**

   In the `values.yml:` section under `stringData:`, include the Dex configuration that you prepared in the previous step.

   ```
   stringData:
     overlays.yaml: |
       #@ load("@ytt:overlay", "overlay")
       #@overlay/match by=overlay.subset({"kind": "Service", "metadata": {"name": "pinniped-supervisor", "namespace": "pinniped-supervisor"}})
       ---
       #@overlay/replace
       spec:
         type: LoadBalancer
         selector:
           app: pinniped-supervisor
         ports:
           - name: https
             protocol: TCP
             port: 443
             targetPort: 8443
      values.yaml: |
        #@data/values
        #@overlay/match-child-defaults missing_ok=True
        ---
        infrastructure_provider: vsphere
        tkg_cluster_role: management
        dex:
          app: dex
          ....
          ....
          service:
            name: dexsvc
            type: LoadBalancer
   ```

1. Run `kubectl edit` to edit the secret.

   ```
   kubectl edit secret tkg-mgmt-vc-pinniped-addon -n tkg-system
   ```

1. When editing the secret text, add the `stringData` configuration YAML that you prepared in the previous step. Leave the existing encoded Base64 content for `values.yaml`.

   ```
   apiVersion: v1
      data:
     values.yaml: LEAVE EXISTING BASE64 ENCODED DATA
     stringData: ADD THE YAML PREPARED DURING PREVIOUS STEP

   kind: Secret
   ```

1. After you save the secret, check the status of Pinniped app. The app should eventually show `Reconcile succeeded`.

   ```
   kubectl get app pinniped -n tkg-system

   NAME      DESCRIPTION          SINCE-DEPLOY   AGE
   pinniped  Reconcile succeeded   3m23s          7h50m
   ```

1. If the returned status is `Reconcile failed`, run the following command to get details on the failure.

   ```
   kubectl get app pinniped -n tkg-system -o yaml
   ```

   If the failed status is due to a ytt template format error, edit the secret to correct the format in `values.yaml` or `stringData` and save the secret again.

1. Check the configuration map for pinniped-info.

   ```
   kubectl get cm pinniped-info -n kube-public -o yaml | grep issuer
   ```

   This command returns the IP address or DNS name for the load balancer service that you configured for the Pinniped. For example:

   ```

  issuer: <https://10.186.131.117>

   ```


## What to Do Next

Share the generated `kubeconfig` file with other users, to allow them to access the management cluster. You can also start creating workload clusters, assign users to roles on those clusters, and share their `kubeconfig` files with those users.

- For information about creating workload clusters, see [Deploying Tanzu Kubernetes Clusters](../tanzu-k8s-clusters/deploy.md).
- For information about how to grant users access to workload clusters on which you have implemented identity management, see [Authenticate Connections to a Workload Cluster](../cluster-lifecycle/connect.md#id-mgmt).
