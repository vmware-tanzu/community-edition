# Deploying Tanzu Community Edition on vSphere with LDAP Identity Management and NSX-ALB

The purpose of this document is to guide the reader through the configuration of LDAP Identity Management during the deployment of a Tanzu Community Edition management and workload clusters. Since this is a Tanzu Community Edition deployment on vSphere, an NSX Advanced Load Balancer (NSX ALB) provides load-balancing services to the cluster. As we shall see, a number of additional configuration steps are needed in this release to enable LDAP Identity Management on vSphere when an NSX ALB provides load-balancing services to the cluster.

Identity Management in Tanzu Community Edition is provided via two packages, [Pinniped](https://pinniped.dev) and [Dex](https://dexidp.io/). Pinniped provides the authentication service, which uses Dex to connect to identity providers such as Active Directory. Tanzu Community Edition automatically deploys these components once Identity Manager is enable and configured during the Tanzu Community Edition deployment.

In this example, the Tanzu Community Edition cluster will be integrated with Microsoft Active Directory. You will need to retrieve a Certificate Authority (CA) from the Active Directory Certificate Service (CA). This CA must be encoded in Base64, which should be an available download option from the Active Directory Certificate Service. This Active Directory / LDAP  deployment does not use Anonymous Authentication, so an actual username will be required to search/browse the directory. This means that many of the configuration fields in the UI for Identity Management will need to be populated.

## Populating the Identity Management Settings for LDAP

The following fields should be populated in the Identity Management settings.

- LDAPS Endpoint IP or FQDN
- LDAPS Endpoint Port
- BIND DN (Distinguished Name)
- BIND Password
- User Search BASE DN (Distinguished Name)
- User Search Username
- Group Search BASE DN (Distinguished Name)
- Group Search Filter
- ROOT CA (Certificate Authority)

### LDAPS Endpoint IP or FQDN

This is the URL of your secure LDAP server. In my example, this service is provided by host `dc01.rainpole.com`.

### LDAPS Endpoint Port

LDAP secure communication uses port `636`.

### BIND DN (Distinguished Name)

In Active Directory, a name that includes an objects entire path to the root of the LDAP namespace is called its `distinguished name`, or `DN` for short. The BIND DN is the distinguished name of the credentials that will be used to search for users and groups. The BIND DN used here (its entire path) is: `cn=Administrator,cn=Users,dc=rainpole,dc=com`.  In Active directory, cn is common name, and dc is domain controller. The default container for User objects is `cn=Users`. An Active Directory domain with the DNS name `rainpole.com` would have the designator `dc=rainpole,dc=com`.

### BIND Password

This is the password for the user specified in the BIND DN section.

### User Search BASE DN (Distinguished Name)

This is the location in the LDAP tree where the user search begins, and this typically matches your domain name. In this example, the following user search attributes were added: `cn=Users,dc=rainpole,dc=com`.

### User Search Username

*Caution*: While the "Validate LDAP Configuration" will succeed without this setting, Dex will fail to deploy if a valid `User Search Username` is not provided. The Dex pod will throw the following error on deployment:

```sh
failed to initialize server: server: Failed to open connector ldap: failed to open connector: \
failed to create connector ldap: ldap: missing required field "userSearch.username"
```

In this configuration, the `User Search Username` is set to `userPrincipalName`. Other options are to leave it at the default value of `uid`.

### Group Search BASE DN (Distinguished Name)

This is similar to the `User Search BASE DN` but this is where group searching begins. Again, this typically matches your domain name. In this example, the following group search attributes were added: `dc=rainpole,dc=com`.

### Group Search Filter

While this should work with the default group filter setting, the "Validate LDAP Configuration" in the UI seems to look for the Group Search `Filter` field to be populated. Manually add a Group Search Filter of `(objectClass=group)` including the rounded brackets.

### ROOT CA (Certificate Authority)

This is the Root CA, which was retrieved via the Active Directory Certificate Service earlier.

## Completed UI View

With all of these settings in place, the Identity Management Settings for LDAPS integration should look similar to the following:

![LDAPS Identity Management Configuration](/docs/img/ldaps-im.png?raw=true)

## Verify LDAP Configuration

If everything has been configured correctly, and your LDAP service is working, then the `Verify LDAP Configuration` utility will run a check to ensure that it can correctly connect to your Active Directory/LDAPS service, do a BIND operation with the BIND user and credentials, and verify that it can search the user and group directories. If successful, it should report something simialr to the following. Note that an LDAP user has been added to the *Test User Name* field:

![LDAPS Configuration Verification](/docs/img/ldap-verify.png?raw=true)

If there are issues, then determine which part of the verification is failing. If it is an X509 error during the Connect phase, check that the certificate is correct, and in base64 format. If it is a BIND issue, check the BIND DN and Password. If it is in the User Search or Group Search, check the Distinguished Names, Filters and Username/Name Attributues accordingly.

### Management Cluster Manifest

These are the resulting settings in the management cluster manifest file. These might be useful if you want to create the management cluster from the CLI rather than the UI. The Certificate Authority and other pertinent information have been obfuscated.

```sh
IDENTITY_MANAGEMENT_TYPE: ldap
LDAP_BIND_DN: cn=Administrator,cn=Users,dc=rainpole,dc=com
LDAP_BIND_PASSWORD: <encoded:Vn[...]Iz>
LDAP_GROUP_SEARCH_BASE_DN: dc=rainpole,dc=com
LDAP_GROUP_SEARCH_FILTER: (objectClass=group)
LDAP_GROUP_SEARCH_GROUP_ATTRIBUTE: ""
LDAP_GROUP_SEARCH_NAME_ATTRIBUTE: ""
LDAP_GROUP_SEARCH_USER_ATTRIBUTE: ""
LDAP_HOST: dc01.rainpole.com:636
LDAP_ROOT_CA_DATA_B64: LS0t[...]tLQ==
LDAP_USER_SEARCH_BASE_DN: cn=Users,dc=rainpole,dc=com
LDAP_USER_SEARCH_FILTER: ""
LDAP_USER_SEARCH_NAME_ATTRIBUTE: userPrincipalName
LDAP_USER_SEARCH_USERNAME: userPrincipalName
OIDC_IDENTITY_PROVIDER_CLIENT_ID: ""
OIDC_IDENTITY_PROVIDER_CLIENT_SECRET: ""
OIDC_IDENTITY_PROVIDER_GROUPS_CLAIM: ""
OIDC_IDENTITY_PROVIDER_ISSUER_URL: ""
OIDC_IDENTITY_PROVIDER_NAME: ""
OIDC_IDENTITY_PROVIDER_SCOPES: ""
OIDC_IDENTITY_PROVIDER_USERNAME_CLAIM: ""
```

## Configuration Steps on the Management Cluster for NSX-ALB

Whilst the initial deployment of the Tanzu Community Edition Management Cluster to vSphere should now proceed successfully, there are some additional steps that need to be provided to enable "non-admin" users to access the management cluster when the NSX Advanced Load Balancer (ALB) is used to provide load-balancing services. In a nutshell, we need to change both `Pinniped` and `Dex` services to type Load Balancer rather than their default NodePort service type.

The assumption at this point is that the Tanzu Community Edition management cluster has deployed. Change context the Tanzu Community Edition management cluster, as admin. The name of this management cluster is *mgmt*.

```sh
% kubectl config use-context mgmt-admin@mgmt
Switched to context "mgmt-admin@mgmt".
```

### Create overlay to change Pinniped Supervisor and Dex services to type Load Balancer

This is the overlay manifest where both the Pinniped Supervisor and Dex services are changed to use Load Balancer services. This is placed in a file called `pinniped-supervisor-svc-overlay.yaml`.

```yaml
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


#@ load("@ytt:overlay", "overlay")
#@overlay/match by=overlay.subset({"kind": "Service", "metadata": {"name": "dexsvc", "namespace": "tanzu-system-auth"}}), missing_ok=True
---
#@overlay/replace
spec:
  type: LoadBalancer
  selector:
    app: dex
  ports:
    - name: dex
      protocol: TCP
      port: 443
      targetPort: https
```

### Convert overlay manifest to base64

The newly created `pinniped-supervisor-svc-overlay.yaml` manifest contents must be converted to base64, as it needs to be included as part of pinniped-addon secret. Run the following command to convert it to base64 (or whatever the appropriate conversion command is for your operating system)

```sh
% cat pinniped-supervisor-svc-overlay.yaml | base64
```

### Patch the Pinipped Addon Secret

Use the base64 output from the previous command to patch the `<mgmt-cluster-name>-pinniped-addon` secret that is already present in the management cluster. The base64 content is represented by *"I0Agb[...]wo="* in the command below, but it has been obfuscated. Typically, you will expect to see the base64 output from the previous command will be much larger.

```sh
% kubectl patch secret mgmt-pinniped-addon -n tkg-system -p '{"data": {"overlays.yaml": "I0Agb[...]wo="}}'
secret/mgmt-pinniped-addon patched
```

### Monitor the Pinniped and Dex Services

 It should now be possible to observe the Pinniped and Dex services changing from NodePort to Load Balancer services using `kubectl get svc -A`:

From:

```sh
pinniped-supervisor pinniped-supervisor   NodePort   100.66.101.230   <none>    443:31234/TCP 7m
tanzu-system-auth   dexsvc                NodePort   100.69.254.209   <none>   5556:30167/TCP 7m
```

To:

```sh
pinniped-supervisor   pinniped-supervisor  LoadBalancer  100.66.101.230   xx.yy.62.20   443:31916/TCP  7m
tanzu-system-auth     dexsvc               LoadBalancer  100.69.254.209   xx.yy.62.21   443:32349/TCP  8m
```

### Relaunch the pinniped-post-deploy-job

The final step is to relaunch the `pinniped-post-deploy-job`. This is done by simply deleting the original job on the management cluster and waiting for the new one to automatically start. This can take up to a few minutes to complete so be patient. Use the watch option (-w) to `kubectl` to monitor the jobs. Once the job restarts, that completes the extra configuration steps required to integrate Pinniped and Dex services with the NSX ALB on vSphere.

```sh
% kubectl get jobs -A
NAMESPACE NAME COMPLETIONS DURATION AGE
pinniped-supervisor pinniped-post-deploy-job 1/1 6m9s 55m


% kubectl get job pinniped-post-deploy-job -n pinniped-supervisor
NAME COMPLETIONS DURATION AGE
pinniped-post-deploy-job 1/1 6m9s 55m


% kubectl delete jobs pinniped-post-deploy-job -n pinniped-supervisor
job.batch "pinniped-post-deploy-job" deleted


% kubectl get jobs -A -w
NAMESPACE NAME COMPLETIONS DURATION AGE
pinniped-supervisor pinniped-post-deploy-job 0/1 0s
pinniped-supervisor pinniped-post-deploy-job 0/1 0s
pinniped-supervisor pinniped-post-deploy-job 0/1 0s 0s
pinniped-supervisor pinniped-post-deploy-job 1/1 11s 11s
^C%
```

At this point, if possible, logon to the NSX ALB management interface, and two additional Virtual Services for Pinniped and Dex should be visible.

## Setting up the management cluster for non-admin users

The next step is to try to interact with the management cluster as a non-admin user, using an LDAP user. There are a number of additional steps required around authentication to allow this.

### Create a ClusterRoleBinding on the management cluster

First, a `ClusterRoleBinding` must be created and applied to the management cluster. An example of a ClusterRoleBinding is below. This is the same user *cormac* that was used to "Verify LDAP Configuration" in the UI earlier but we are using the users login name which is *chogan@rainpole.com*. This is also the user that will authenticate in the DEX portal shortly.

```yaml
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: chogan
subjects:
  - kind: User
    name: chogan@rainpole.com
    apiGroup:
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: rbac.authorization.k8s.io
```

```sh
% kubectl apply -f chogan-crb.yaml
clusterrolebinding.rbac.authorization.k8s.io/chogan created
```

### Switch contexts to non-admin on management cluster

Get the kubeconfig of the management cluster as a "non-admin", and switch to that "non-admin" context as shown below.

```sh
% tanzu management-cluster kubeconfig get
You can now access the cluster by running 'kubectl config use-context tanzu-cli-mgmt@mgmt'

% kubectl config use-context tanzu-cli-mgmt@mgmt
Switched to context "tanzu-cli-mgmt@mgmt".
```

### Add LDAP credentials to Dex Portal for management cluster

As soon as an attempt is made to query the cluster in the "non-admin" context, for example `kubectl get nodes`, the Dex portal launches in a browser window. Add the user credentials from the `ClusterRoleBinding` which was created in first step of this section, such as *chogan@rainpole.com*.

![DEX LDAP Portal](/docs/img/dex-portal.png?raw=true)

The `kubectl` command run previously to query the cluster should now complete once the authentication step with Dex has completed successfully.

```sh
% kubectl get nodes
NAME                        STATUS ROLES                AGE VERSION
mgmt-control-plane-p6pdl    Ready  control-plane,master 68m v1.21.2+vmware.1
mgmt-control-plane-wz9hc    Ready  control-plane,master 74m v1.21.2+vmware.1
mgmt-control-plane-xq477    Ready  control-plane,master 70m v1.21.2+vmware.1
mgmt-md-0-77589686bc-m6lq8  Ready  <none>               72m v1.21.2+vmware.1
```

Everything should now work, and the non-admin/LDAP user should be able to communicate with the cluster successfully. Revisit the steps in this section if you get errors such as permissions errors, for example:

```sh
% kubectl get nodes
Error from server (Forbidden): nodes is forbidden: User "chogan@rainpole.com" cannot list resource "nodes" in API group "" at the cluster scope
```

If problems persist, verify the overlay manifest and patching steps. Another option is to retry the deletion of the pinniped job once again and allowing it to restart once more to see if resolves any issues that might be experienced with authentication.

## Configuration Steps on the Workload Cluster

Note that the above steps are for the management cluster only. Some additional configuration steps are required for workload clusters.

### IDENTITY_MANAGEMENT_TYPE variable in manifest

 If you want to use a "non-admin" LDAP user on a workload clusters, a manifest for the workload cluster which contains an environment variable IDENTITY_MANAGEMENT_TYPE needs to be created. Since this configuration is using LDAP, it should be set to `ldap`.

Here is an example of such a manifest, with certain items obfuscated:

```yaml
#! -- See https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.3/vmware-tanzu-kubernetes-grid-13/GUID-tanzu-k8s-clusters-vsphere.html
##
##! ---------------------------------------------------------------------
##! Basic cluster creation configuration
##! ---------------------------------------------------------------------
##
CLUSTER_NAME: workload
CLUSTER_PLAN: dev
CNI: antrea
#
##! ---------------------------------------------------------------------
##! Identity Management (set if management cluster is configured)
##! ---------------------------------------------------------------------
#
IDENTITY_MANAGEMENT_TYPE: ldap
#
##! ---------------------------------------------------------------------
##! Node configuration
##! ---------------------------------------------------------------------
#
CONTROL_PLANE_MACHINE_COUNT: 1
WORKER_MACHINE_COUNT: 1
VSPHERE_CONTROL_PLANE_NUM_CPUS: 2
VSPHERE_CONTROL_PLANE_DISK_GIB: 40
VSPHERE_CONTROL_PLANE_MEM_MIB: 8192
VSPHERE_WORKER_NUM_CPUS: 2
VSPHERE_WORKER_DISK_GIB: 40
VSPHERE_WORKER_MEM_MIB: 4096
#
##! ---------------------------------------------------------------------
##! vSphere configuration
##! ---------------------------------------------------------------------
#
VSPHERE_DATACENTER: /OCTO-Datacenter
VSPHERE_DATASTORE: /OCTO-Datacenter/datastore/vsan-OCTO-Cluster-A
VSPHERE_FOLDER: /OCTO-Datacenter/vm/TKG
VSPHERE_NETWORK: "VM Network"
VSPHERE_PASSWORD: <encoded:Vk13YXJlMTIzIQ==>
VSPHERE_RESOURCE_POOL: /OCTO-Datacenter/host/OCTO-Cluster-A/Resources
VSPHERE_SERVER: vcsa-06.rainpole.com
VSPHERE_SSH_AUTHORIZED_KEY: ssh-rsa AAAA[...]iqJlH chogan@rainpole.com
VSPHERE_TLS_THUMBPRINT: AA:BB:CC:DD:EE:FF:GG:HH:II:JJ:KK
VSPHERE_USERNAME: administrator@vsphere.local
#
##! ---------------------------------------------------------------------
##! Common configuration
##! ---------------------------------------------------------------------
#
ENABLE_DEFAULT_STORAGE_CLASS: true
CLUSTER_CIDR: 100.96.102.0/11
SERVICE_CIDR: 100.64.102.0/13
#
##! ---------------------------------------------------------------------
##! AVI (NSX ALB) configuration
##! ---------------------------------------------------------------------
#
AVI_CA_DATA_B64: LS0t[...]Cg==
AVI_CLOUD_NAME: Default-Cloud
AVI_CONTROL_PLANE_HA_PROVIDER: "true"
AVI_CONTROLLER: XX.YY.51.163
AVI_DATA_NETWORK: VM-62-DPG
AVI_DATA_NETWORK_CIDR: XX.YY.62.0/26
AVI_ENABLE: "true"
AVI_LABELS: ""
AVI_MANAGEMENT_CLUSTER_VIP_NETWORK_CIDR: XX.YY.62.0/26
AVI_MANAGEMENT_CLUSTER_VIP_NETWORK_NAME: VM-62-DPG
AVI_PASSWORD: <encoded:Vk[...]==>
AVI_SERVICE_ENGINE_GROUP: Default-Group
AVI_USERNAME: admin
```

### Create the workload cluster

Use the `tanzu` command to create the workload cluster, using the above manifest.

```sh
% tanzu cluster create -f workload.yaml
Validating configuration...
Creating workload cluster 'workload'...
Waiting for cluster to be initialized...
Waiting for cluster nodes to be available...
Waiting for addons installation...
Waiting for packages to be up and running...

Workload cluster 'workload' created
```

### Switch to workload cluster admin context

Once the workload cluster has been successfully created, change to the workload cluster admin context.

```sh
% tanzu cluster kubeconfig get workload --admin
Credentials of cluster 'workload' have been saved
You can now access the cluster by running 'kubectl config use-context workload-admin@workload'

% kubectl config use-context workload-admin@workload
Switched to context "workload-admin@workload".
```

### Create a ClusterRoleBinding on the workload cluster

As was done on the management cluster, a `ClusterRoleBinding` must also be created and applied to the workload cluster. Once more, this is the same user *cormac* that was used to "Verify LDAP Configuration" in the UI earlier. This is also the user that will be used to authenticate with in the DEX portal later on, so that this LDAP user can access the workload cluster.

```yaml
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: chogan
subjects:
  - kind: User
    name: chogan@rainpole.com
    apiGroup:
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: rbac.authorization.k8s.io
```

```sh
% kubectl apply -f chogan-crb.yaml
clusterrolebinding.rbac.authorization.k8s.io/chogan created
```

### Switch to non-admin context on workload cluster

Retrieve and switch to non-admin user context, then try to access cluster. Just like we saw on the management cluster, you will not be able to access the workload cluster until you have authenticated through Dex.

```sh
% tanzu cluster kubeconfig get workload
ℹ  You can now access the cluster by running 'kubectl config use-context tanzu-cli-workload@workload'

% kubectl config use-context tanzu-cli-workload@workload
Switched to context "tanzu-cli-workload@workload".
```

### Add LDAP credentials to Dex Portal for workload cluster

As seen previously, as soon as you attempt to query the workload cluster in the "non-admin" context, e.g. by running `kubectl get nodes`, the Dex portal launches and prompts for AD authentication in a browser window to determine who is allowed to access the cluster. Once again, the user credentials from the ClusterRoleBinding which was created earlier in this section should be added, *chogan@rainpole.com*.

This user, possibly a developer persona, should now be able to access the workload cluster.

```sh
% kubectl get nodes
NAME                             STATUS   ROLES                  AGE   VERSION
workload-control-plane-bcv2k    Ready    control-plane,master   29m   v1.21.2+vmware.1
workload-md-0-646d6d9f5-p2cww   Ready    <none>                 28m   v1.21.2+vmware.1
```

That completes the setup. LDAP users can now use their non-admin privileges to access the workload clusters. The `ClusterRolebinding`, managed by the actual cluster admin, controls the privileges that each user has on the cluster, rather than allowing the cluster users to have full admin privileges.
