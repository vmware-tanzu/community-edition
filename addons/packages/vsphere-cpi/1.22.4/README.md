# vSphere CPI

> This package provides cloud provider interface using vsphere-cpi.

For more information, see the [Kubernetes vSphere Cloud Provider](https://github.com/kubernetes/cloud-provider-vsphere).

## Configuration

The following configuration values can be set to customize the vsphere CPI installation.

### Global

None

### vSphere CPI Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `vsphereCPI.server` | Required | The IP address or FQDN of the vSphere endpoint. Default value is `""`. |
| `vsphereCPI.datacenter` | Required | The datacenter in which VMs are created/located. Default value is `""`. |
| `vsphereCPI.username` | Required | Username used to access a vSphere endpoint. Default value is `""`. |
| `vsphereCPI.password` | Required | Password used to access a vSphere endpoint. Default value is `""`. |
| `vsphereCPI.tlsThumbprint` | Optional | The cryptographic thumbprint of the vSphere endpoint's certificate. Default value is `""`. |
| `vsphereCPI.insecureFlag` | Optional | The flag that disables TLS peer verification. Default value is `False`. |
| `vsphereCPI.ipFamily` | Optional | The IP family configuration. Default value is `null`. |
| `vsphereCPI.vmInternalNetwork` | Optional | Internal VM network name. Default value is `null`. |
| `vsphereCPI.vmExternalNetwork` | Optional | External VM network name. Default value is `null`. |
| `vsphereCPI.cloudProviderExtraArgs.tls-cipher-suites` | Optional | External arguments for cloud provider. Default: `TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384` |
| `vsphereCPI.nsxt.podRoutingEnabled` | Optional | A flag that enables pod routing. Default: `false`. |
| `vsphereCPI.nsxt.routes.routerPath` | Optional | NSX-T T0/T1 logical router path. Default: `""`. |
| `vsphereCPI.nsxt.routes.clusterCidr` | Optional | Cluster CIDR. Default: `""` . |
| `vsphereCPI.nsxt.username` | Optional | The username used to access NSX-T. Default: `""`. |
| `vsphereCPI.nsxt.password` | Optional | The password used to access NSX-T. Default: `""`. |
| `vsphereCPI.nsxt.host`| Optional | The NSX-T server. Default: `null`. |
| `vsphereCPI.nsxt.insecureFlag` | Optional | InsecureFlag is to be set to true if NSX-T uses self-signed cert. Default: `false`. |
| `vsphereCPI.nsxt.remoteAuth` | Optional | RemoteAuth is to be set to true if NSX-T uses remote authentication (authentication done through the vIDM). Default: `false`. |
| `vsphereCPI.nsxt.vmcAccessToken`| Optional | VMCAccessToken is VMC access token for token based authentification. Default: `""`. |
| `vsphereCPI.nsxt.vmcAuthHost` | Optional | VMCAuthHost is VMC verification host for token based authentification. Default: `""`. |
| `vsphereCPI.nsxt.clientCertKeyData` | Optional | Client certificate key. Default: `""`. |
| `vsphereCPI.nsxt.clientCertData`| Optional | Client certificate data. Default: `""`. |
|`vsphereCPI.nsxt.rootCAData` | Optional | The certificate authority for the server certificate for locally signed certificates. Default: `""`. |
| `vsphereCPI.nsxt.secretName` | Optional | The name of secret that stores CPI configuration. Default: `cloud-provider-vsphere-nsxt-credentials`. |
| `vsphereCPI.nsxt.secretNamespace`| Optional | The namespace of secret that stores CPI configuration. Default: `True`. |

## Usage Example

To learn more about how to use vSphere CPI, refer to [Kubernetes vSphere Cloud Provider](https://github.com/kubernetes/cloud-provider-vsphere).
