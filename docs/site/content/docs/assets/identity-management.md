
<!-- ??I don't know if this is something we want to reference or if we need to supply our own??? I presume this full section needs to be reworked for TCE  For information about how Tanzu Kubernetes Grid implements identity management, see [Enabling Identity Management in Tanzu Kubernetes Grid](enabling-id-mgmt.md).--> 
**(++ENG TEAM - NEEDS REWORK/INPUT++)** 
1. In the **Identity Management** section, optionally disable **Enable Identity Management Settings**. You can disable identity management for proof-of-concept deployments, but it is strongly recommended to implement identity management in production deployments. If you disable identity management, you can reenable it later.
<!--  ![Configure external Identity Provider](../images/install-v-7id.png)--> 
2. If you enable identity management, select **OIDC** or **LDAPS**.

   **OIDC**:

   Provide details of your OIDC provider account, for example, Okta.

   - **Issuer URL**: The IP or DNS address of your OIDC server.
   - **Client ID**: The `client_id` value that you obtain from your OIDC provider. For example, if your provider is Okta, log in to Okta, create a Web application, and select the **Client Credentials** options in order to get a `client_id` and `secret`.
   - **Client Secret**: The `secret` value that you obtain from your OIDC provider.
   - **Scopes**: A comma separated list of additional scopes to request in the token response. For example, `openid,groups,email`.
   - **Username Claim**: The name of your username claim. This is used to set a user's username in the JSON Web Token (JWT) claim. Depending on your provider, enter claims such as `user_name`, `email`, or `code`.
   - **Groups Claim**: The name of your groups claim. This is used to set a user's group in the JWT claim. For example, `groups`.

   ![Configure external Identity Provider](../images/install-v-7id-oidc.png)

   **LDAPS**:

   Provide details of your company's LDAPS server. All settings except for **LDAPS Endpoint** are optional.

   - **LDAPS Endpoint**: The IP or DNS address of your LDAPS server. Provide the address and port of the LDAP server, in the form `host:port`.
   - **Bind DN**: The DN for an application service account. The connector uses these credentials to search for users and groups. Not required if the LDAP server provides access for anonymous authentication.
   - **Bind Password**: The password for an application service account, if **Bind DN** is set.

   Provide the user search attributes.

   - **Base DN**: The point from which to start the LDAP search. For example, `OU=Users,OU=domain,DC=io`.
   - **Filter**: An optional filter to be used by the LDAP search.
   - **Username**: The LDAP attribute that contains the user ID. For example, `uid, sAMAccountName`.

   Provide the group search attributes.

   - **Base DN**: The point from which to start the LDAP search. For example, `OU=Groups,OU=domain,DC=io`.
   - **Filter**: An optional filter to be used by the LDAP search.
   - **Name Attribute**: The LDAP attribute that holds the name of the group. For example, `cn`.
   - **User Attribute**: The attribute of the user record that is used as the value of the membership attribute of the group record. For example, `distinguishedName, dn`.
   - **Group Attribute**:  The attribute of the group record that holds the user/member information. For example, `member`.

   Paste the contents of the LDAPS server CA certificate into the **Root CA** text box.

   <!--![Configure external Identity Provider](../images/install-v-7id-ldap.png)-->



