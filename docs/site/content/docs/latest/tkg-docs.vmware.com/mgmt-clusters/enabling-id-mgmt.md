# Enabling Identity Management in Tanzu Kubernetes Grid

Tanzu Kubernetes Grid implements user authentication with [Pinniped](https://pinniped.dev/). Pinniped allows you to plug external OpenID Connect (OIDC) or LDAP identity providers (IDP) into Tanzu Kubernetes clusters, so that you can control user access to those clusters. Pinniped is an open-source authentication service for Kubernetes clusters. If you use LDAP authentication, Pinniped uses [Dex](https://github.com/dexidp/dex/blob/master/README.md) as the endpoint to connect to your upstream LDAP identity provider. If you use OIDC, Pinniped provides its own endpoint, so Dex is not required. Pinniped and Dex run automatically as in-cluster services in your management clusters if you enable identity management during management cluster deployment.

**IMPORTANT**:

- In Tanzu Kubernetes Grid v1.3.0, Pinniped used Dex as the endpoint for both OIDC and LDAP providers. In Tanzu Kubernetes Grid v1.3.1 and later, Pinniped no longer requires Dex and uses the Pinniped endpoint for OIDC providers. In Tanzu Kubernetes Grid v1.3.1 and later, Dex is only used if you use an LDAP provider. Consequently, it is **strongly recommended** to use Tanzu Kubernetes Grid v1.3.1 if you want to implement identity management on new management clusters. If you have already used Tanzu Kubernetes Grid v1.3.0 to deploy management clusters that implement OIDC authentication, when you upgrade those management clusters to v1.3.1, you must perform additional steps to update the Pinniped configuration. For information about the additional steps to perform, see [Update the Callback URL for Management Clusters with OIDC Authentication](../upgrade-tkg/management-cluster.md#update-callbackurl) in *Upgrade Management Clusters*.
- Previous versions of Tanzu Kubernetes Grid included optional Dex and Gangway extensions to provide identity management. These manually deployed Dex and Gangway extensions from previous versions are deprecated in this release. If you manually deployed the Dex and Gangway extensions on clusters in a previous release and you upgrade the clusters to this version of Tanzu Kubernetes Grid, it is **strongly recommended** to migrate your identity management implementation from Dex and Gangway to Pinniped and Dex. If you did not implement Dex and Gangway on clusters from a previous version of Tanzu Kubernetes Grid and you upgrade them to this version, it is also strongly recommended to implement Pinniped and Dex on those clusters.

## About Tanzu Kubernetes Grid Identity Management

The process for implementing identity management is as follows:

- The Tanzu Kubernetes Grid administrator creates a management cluster, specifying an external OIDC or LDAP IDP.
- Authentication service components are deployed into the management cluster, using the OIDC or LDAP IDP specified during deployment.
- The administrator creates a Tanzu Kubernetes (workload) cluster. The workload cluster inherits the authentication configuration from the management cluster.
- The administrator creates a role binding to associate a given user with a given role on the workload cluster.
- The administrator provides the `kubeconfig` for the workload cluster to the user.
- A user uses the `kubeconfig` to connect to the workload cluster, for example, by running `kubectl get pods --kubeconfig <kubeconfig-file>`.
- The management cluster authenticates the user with the IDP.
- The workload cluster either allows or denies the `kubectl get pods` request, depending on the permissions of the user's role.

In the image below, the blue arrows represent the authentication flow between the workload cluster, the management cluster and the external IDP. The green arrows represent Tanzu CLI and `kubectl` traffic between the workload cluster, the management cluster and the external IDP.

![Identity Management in Tanzu Kubernetes Grid](../images/tkg-id-mgmt.png)

## What Happens When You Enable Identity Management

The diagram below shows the identity management components that Tanzu Kubernetes Grid deploys in the management cluster and in Tanzu Kubernetes (workload) clusters when you enable identity management.

![Identity Management architecture in Tanzu Kubernetes Grid](../images/tkg-id-mgmt-architecture.png)

Understanding the diagram:

* The purple-bordered rectangles show the identity management components, which include Pinniped, Dex, and a post-deployment job in the management cluster and Pinniped and a post-deployment job in the workload cluster. In Tanzu Kubernetes Grid v1.3.0, Pinniped uses Dex as the endpoint for both OIDC and LDAP providers. In v1.3.1 and later, Dex is deployed only for LDAP providers.

* The gray-bordered rectangles show the components that Tanzu Kubernetes Grid uses to control the lifecycle of the identity management components, which include the Tanzu CLI, `tanzu-addons-manager`, and `kapp-controller`.

* The green-bordered rectangle shows the Pinniped add-on secret created for the management cluster.

* The orange-bordered rectangle in the management cluster shows the Pinniped add-on secret created for the workload cluster. The secret is mirrored to the workload cluster.

Internally, Tanzu Kubernetes Grid deploys the identity management components as a core add-on, `pinniped`. When you deploy a management cluster with identity management enabled, the Tanzu CLI creates a Kubernetes secret for the `pinniped` add-on in the management cluster. `tanzu-addons-manager` reads the secret, which contains your IDP configuration information, and instructs `kapp-controller` to configure the `pinniped` add-on using the configuration information from the secret.

The Tanzu CLI creates a separate `pinniped` add-on secret for each workload cluster that you deploy from the management cluster. All secrets are stored in the management cluster.

## <a id="idp"></a> Obtain Your Identity Provider Details

Before you can deploy a management cluster with identity management enabled, you must have an identity provider. Tanzu Kubernetes Grid supports LDAPS and OIDC identity providers.

To use your company's internal LDAPS server as the identity provider, obtain LDAPS information from your LDAP administrator.

To use OIDC as the identity provider, you must have an account with an IDP that supports the OpenID Connect standard, for example [Okta](https://www.okta.com/).

### Example: Register a Tanzu Kubernetes Grid Application in Okta

To use Okta as your OIDC provider, you must create an account with Okta and register an application for Tanzu Kubernetes Grid with your account.

1. If you do not have one, create an [Okta](https://www.okta.com/) account.
1. Go to the Admin portal by clicking the **Admin** button.
1. Go to Applications, and click **Add Application**.
1. Click **Create New App**.
1. For **Platform**, select **Web** and for **Sign on method**, select **OpenID Connect**, then click **Create**.
1. Give your application a name.
1. Enter a placeholder **Login redirect URI**.

   For example, enter `http://localhost:8080/callback`. You will update this with the real URL after you deploy the management cluster.
1. Click **Save**.
1. In the **General** tab for your application, copy and save the **Client ID** and **Client secret**.

   You will need these credentials when you deploy the management cluster.
1. In the **Assignments** tab, assign people and groups to the application.

   The people and groups that you assign to the application will be the users  who can access the management cluster and the Tanzu Kubernetes clusters that you use it to deploy.

## What to Do Next

You can now deploy management clusters that implement identity management, to restrict access to clusters to authorized users.

- [Deploy Management Clusters with the Installer Interface](deploy-ui.md)
- [Deploy Management Clusters from a Configuration File](deploy-cli.md)

If you implement identity management, after you deploy the management cluster, there are post-deployment steps to perform, that are described in [Configure Identity Management After Management Cluster Deployment](configure-id-mgmt.md).
