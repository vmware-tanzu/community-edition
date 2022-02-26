# pinniped Package

This package provides user authentication using [pinniped](https://pinniped.dev).

## Components

## Configuration

The following configuration values can be set to customize the Pinniped installation.
See `bundle/config/values.yaml` for descriptions of the configuration values.

| Value | Default |
|-------| ------- |
| `custom_cluster_issuer` |  |
| `custom_tls_secret` |  |
| `dex.app` |  `dex`  |
| `dex.certificate.duration` |  `2160h`  |
| `dex.certificate.renewBefore` |  `360h`  |
| `dex.commonname` |  `tkg-dex`  |
| `dex.config.connector` |  `<nil>`  |
| `dex.config.enablePasswordDB` |  `false`  |
| `dex.config.expiry.authRequests` |  `90m`  |
| `dex.config.expiry.deviceRequests` |  `5m`  |
| `dex.config.expiry.idTokens` |  `5m`  |
| `dex.config.expiry.signingKeys` |  `90m`  |
| `dex.config.frontend.theme` |  `tkg`  |
| `dex.config.issuerPort` |  `30167`  |
| `dex.config.ldap.BIND_PW_ENV_VAR` |  `<nil>`  |
| `dex.config.ldap.bindDN` |  `<nil>`  |
| `dex.config.ldap.bindPW` |  `<nil>`  |
| `dex.config.ldap.groupSearch.baseDN` |  `<nil>`  |
| `dex.config.ldap.groupSearch.filter` |  `(objectClass=posixGroup)`  |
| `dex.config.ldap.groupSearch.nameAttr` |  `cn`  |
| `dex.config.ldap.groupSearch.scope` |  `sub`  |
| `dex.config.ldap.groupSearch.userMatchers` |  `[]`  |
| `dex.config.ldap.host` |  `<nil>`  |
| `dex.config.ldap.insecureNoSSL` |  `false`  |
| `dex.config.ldap.insecureSkipVerify` |  `false`  |
| `dex.config.ldap.rootCA` |  `<nil>`  |
| `dex.config.ldap.rootCAData` |  `<nil>`  |
| `dex.config.ldap.startTLS` |  `<nil>`  |
| `dex.config.ldap.userSearch.baseDN` |  `<nil>`  |
| `dex.config.ldap.userSearch.emailAttr` |  `mail`  |
| `dex.config.ldap.userSearch.filter` |  `(objectClass=posixAccount)`  |
| `dex.config.ldap.userSearch.idAttr` |  `uid`  |
| `dex.config.ldap.userSearch.nameAttr` |  `givenName`  |
| `dex.config.ldap.userSearch.scope` |  `sub`  |
| `dex.config.ldap.userSearch.username` |  `uid`  |
| `dex.config.ldap.usernamePrompt` |  `LDAP Username`  |
| `dex.config.logger.format` |  `json`  |
| `dex.config.logger.level` |  `info`  |
| `dex.config.oauth2.responseTypes` |  `[]`  |
| `dex.config.oauth2.skipApprovalScreen` |  `true`  |
| `dex.config.oidc.CLIENT_ID` |  `<nil>`  |
| `dex.config.oidc.CLIENT_SECRET` |  `<nil>`  |
| `dex.config.oidc.basicAuthUnsupported` |  `<nil>`  |
| `dex.config.oidc.claimMapping.email` |  `email`  |
| `dex.config.oidc.claimMapping.email_verified` |  `email_verified`  |
| `dex.config.oidc.claimMapping.groups` |  `DEPRECATED`  |
| `dex.config.oidc.clientID` |  `$OIDC_CLIENT_ID`  |
| `dex.config.oidc.clientSecret` |  `$OIDC_CLIENT_SECRET`  |
| `dex.config.oidc.getUserInfo` |  `<nil>`  |
| `dex.config.oidc.hostedDomains` |  `[DEPRECATED]`  |
| `dex.config.oidc.insecureEnableGroups` |  `true`  |
| `dex.config.oidc.insecureSkipEmailVerified` |  `false`  |
| `dex.config.oidc.issuer` |  `<nil>`  |
| `dex.config.oidc.scopes` |  `[DEPRECATED]`  |
| `dex.config.oidc.userIDKey` |  `<nil>`  |
| `dex.config.oidc.userNameKey` |  `<nil>`  |
| `dex.config.staticClients` |  `[]`  |
| `dex.config.storage.config.inCluster` |  `true`  |
| `dex.config.storage.type` |  `kubernetes`  |
| `dex.config.web.https` |  `0.0.0.0:5556`  |
| `dex.config.web.tlsCert` |  `/etc/dex/tls/tls.crt`  |
| `dex.config.web.tlsKey` |  `/etc/dex/tls/tls.key`  |
| `dex.create_namespace` |  `true`  |
| `dex.deployment.replicas` |  `1`  |
| `dex.dns.aws.DEX_SVC_LB_HOSTNAME` |  `<nil>`  |
| `dex.dns.aws.dnsNames` |  `[]`  |
| `dex.dns.azure.DEX_SVC_LB_HOSTNAME` |  `<nil>`  |
| `dex.dns.azure.dnsNames` |  `[]`  |
| `dex.dns.vsphere.DEX_SVC_LB_HOSTNAME` |  `<nil>`  |
| `dex.dns.vsphere.dnsNames` |  `[]`  |
| `dex.dns.vsphere.ipAddresses` |  `[]`  |
| `dex.image.name` |  `DEPRECATED`  |
| `dex.image.pullPolicy` |  `DEPRECATED`  |
| `dex.image.repository` |  `DEPRECATED`  |
| `dex.image.tag` |  `DEPRECATED`  |
| `dex.namespace` |  `tanzu-system-auth`  |
| `dex.organization` |  `vmware`  |
| `dex.service.name` |  `dexsvc`  |
| `dex.service.type` |  `<nil>`  |
| `http_proxy` |  |
| `https_proxy` |  |
| `identity_management_type` |  `<nil>`  |
| `imageInfo.imagePullPolicy` |  `IfNotPresent`  |
| `imageInfo.imageRepository` |  `projects-stg.registry.vmware.com/tkg`  |
| `imageInfo.images.dexImage.imagePath` |  `dex`  |
| `imageInfo.images.dexImage.tag` |  `v2.27.0_vmware.1`  |
| `imageInfo.images.pinnipedImage.imagePath` |  `pinniped`  |
| `imageInfo.images.pinnipedImage.tag` |  `v0.4.1_vmware.1`  |
| `imageInfo.images.tkgPinnipedPostDeployImage.imagePath` |  `tkg-pinniped-post-deploy`  |
| `imageInfo.images.tkgPinnipedPostDeployImage.tag` |  `v0.4.1_vmware.1`  |
| `infrastructure_provider` |  `<nil>`  |
| `no_proxy` |  |
| `pinniped.cert_duration` |  `2160h`  |
| `pinniped.cert_renew_before` |  `360h`  |
| `pinniped.image.name` |  `DEPRECATED`  |
| `pinniped.image.pull_policy` |  `DEPRECATED`  |
| `pinniped.image.repository` |  `DEPRECATED`  |
| `pinniped.image.tag` |  `DEPRECATED`  |
| `pinniped.post_deploy_job_image.name` |  `DEPRECATED`  |
| `pinniped.post_deploy_job_image.pull_policy` |  `DEPRECATED`  |
| `pinniped.post_deploy_job_image.repository` |  `DEPRECATED`  |
| `pinniped.post_deploy_job_image.tag` |  `DEPRECATED`  |
| `pinniped.supervisor_ca_bundle_data` |  `ca_bundle_data_of_pinniped_supervisor_svc`  |
| `pinniped.supervisor_svc_endpoint` |  `https://0.0.0.0:31234`  |
| `pinniped.supervisor_svc_external_dns` |  `<nil>`  |
| `pinniped.supervisor_svc_external_ip` |  `0.0.0.0`  |
| `pinniped.upstream_oidc_additional_scopes` |  `[]`  |
| `pinniped.upstream_oidc_claims.groups` |  |
| `pinniped.upstream_oidc_claims.username` |  |
| `pinniped.upstream_oidc_client_id` |  |
| `pinniped.upstream_oidc_client_secret` |  |
| `pinniped.upstream_oidc_issuer_url` |  `https://0.0.0.0:30167`  |
| `pinniped.upstream_oidc_provider_name` |  `DEPRECATED`  |
| `pinniped.upstream_oidc_tls_ca_data` |  `ca_bundle_data_of_dex_svc`  |
| `tkg_cluster_role` |  `<nil>`  |

## Usage Example

See bundle/examples directory for example configurations of the Pinniped package.

## Building the templates

Build the templates using `oidc` or `ldap` overlay:

```bash
cd $THIS_DIRECTORY && ytt -f bundle/config -f ../examples/mc-oidc.yaml
```

## Generate image package

`kbld` will generate the `bundle/.imgpkg/images.yml` file via the following:

```bash
  cd $THIS_DIRECTORY && ytt -f bundle/config/ -f bundle/examples/mc-ldap.yaml | kbld -f bundle/kbld-config.yaml -f - --imgpkg-lock-output bundle/.imgpkg/images.yml
```

---

This file was generated by regen-readme.sh on Thu, 17 Feb 2022 16:23:06 EST.
