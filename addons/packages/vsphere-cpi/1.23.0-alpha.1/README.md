# vSphere CPI

> This package provides cloud provider interface using vsphere-cpi.

For more information, see the [GitHub page](https://github.com/kubernetes-sigs/cluster-api-provider-vsphere) of vSphere CPI.

## Configuration

The following configuration values can be set to customize the vsphere CPI installation.

### Global

None

### vSphere CPI Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `vsphereCPI.mode`    | Optional | The vSphere mode. Either `vsphereCPI` or `vsphereParavirtualCPI`. Default value is `vsphereCPI` |
| `vsphereCPI.server` | Required | The IP address or FQDN of the vSphere endpoint. Default value is `""`. |
| `vsphereCPI.datacenter` | Required | The datacenter in which VMs are created/located. Default value is `""`. |
| `vsphereCPI.username` | Required | Username used to access a vSphere endpoint. Default value is `""`. |
| `vsphereCPI.password` | Required | Password used to access a vSphere endpoint. Default value is `""`. |
| `vsphereCPI.tlsThumbprint` | Optional | The cryptographic thumbprint of the vSphere endpoint's certificate. Default value is `""`. |
| `vsphereCPI.insecureFlag` | Optional | The flag that disables TLS peer verification. Default value is `False`. |
| `vsphereCPI.ipFamily` | Optional | IP family. Default value is `""`. |
| `vsphereCPI.vmInternalNetwork` | Optional | Internal VM network name. Default value is `""`. |
| `vsphereCPI.vmExternalNetwork` | Optional | External VM network name. Default value is `""`. |
| `vsphereCPI.vmExcludeInternalNetworkSubnetCidr` | Optional | External VM network name. Default value is `""`. |
| `vsphereCPI.vmExcludeExternalNetworkSubnetCidr` | Optional | Internal VM network name. Default value is `""`. |
| `vsphereCPI.cloudProviderExtraArgs.tls-cipher-suites` | Optional | External arguments for cloud provider. Default: `TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384` |
| `vsphereCPI.nsxt.podRoutingEnabled` | Optional | A flag that enables pod routing. Default: `false`. |
| `vsphereCPI.nsxt.routes.routerPath` | Optional | NSX-T T0/T1 logical router path. Default: `""`. |
| `vsphereCPI.nsxt.routes.clusterCidr` | Optional | Cluster CIDR. Default: `""` . |
| `vsphereCPI.nsxt.username` | Optional | The username used to access NSX-T. Default: `""`. |
| `vsphereCPI.nsxt.password` | Optional | The password used to access NSX-T. Default: `""`. |
| `vsphereCPI.nsxt.host`| Optional | The NSX-T server. Default: `""`. |
| `vsphereCPI.nsxt.insecureFlag` | Optional | InsecureFlag is to be set to true if NSX-T uses self-signed cert. Default: `false`. |
| `vsphereCPI.nsxt.remoteAuth` | Optional | RemoteAuth is to be set to true if NSX-T uses remote authentication (authentication done through the vIDM). Default: `false`. |
| `vsphereCPI.nsxt.vmcAccessToken`| Optional | VMCAccessToken is VMC access token for token based authentification. Default: `""`. |
| `vsphereCPI.nsxt.vmcAuthHost` | Optional | VMCAuthHost is VMC verification host for token based authentification. Default: `""`. |
| `vsphereCPI.nsxt.clientCertKeyData` | Optional | Client certificate key. Default: `""`. |
| `vsphereCPI.nsxt.clientCertData`| Optional | Client certificate data. Default: `""`. |
|`vsphereCPI.nsxt.rootCAData` | Optional | The certificate authority for the server certificate for locally signed certificates. Default: `""`. |
| `vsphereCPI.nsxt.secretName` | Optional | The name of secret that stores CPI configuration. Default: `cloud-provider-vsphere-nsxt-credentials`. |
| `vsphereCPI.nsxt.secretNamespace`| Optional | The namespace of secret that stores CPI configuration. Default: `True`. |
| `vsphereCPI.clusterAPIVersion`| Optional for `vsphereParavirtual` | Used in `vsphereParavirtual` mode, defines the Cluster API versions. Default: `cluster.x-k8s.io/v1beta1` |
| `vsphereCPI.clusterKind`| Optional for `vsphereParavirtual` | Used in `vsphereParavirtual` mode, defines the Cluster kind. Default: `Cluster` |
| `vsphereCPI.clusterName`| Required for `vsphereParavirtual` | Used in `vsphereParavirtual` mode, defines the Cluster name. Default: `test-cluster` |
| `vsphereCPI.clusterUID`| Required for `vsphereParavirtual` | Used in `vsphereParavirtual` mode, defines the Cluster UID. Default: `""` |
| `vsphereCPI.supervisorMasterEndpointIP`| Required for `vsphereParavirtual` | Used in `vsphereParavirtual` mode, the endpoint IP of supervisor cluster's API server. Default: `""` |
| `vsphereCPI.supervisorMasterPort`| Required for `vsphereParavirtual` | Used in `vsphereParavirtual` mode, the endpoint port of supervisor cluster's API server port. Default: `""` |

## Usage Example

To learn more about how to use vSphere CPI refer to [vSphere CPI website](https://github.com/kubernetes-sigs/cluster-api-provider-vsphere)
