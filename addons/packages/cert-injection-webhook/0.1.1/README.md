# Cert Injection Webhook

The Cert Injection Webhook for Kubernetes extends kubernetes with a webhook that injects
CA certificates and proxy environment variables into pods. The webhook uses certificates and
environment variables defined in configmaps and injects them into pods with the desired labels or annotations.

## Components

* cert-injection-webhook

## Supported Providers

The following table shows the providers this package can work with.

| AWS | Azure | vSphere | Docker |
|-----|-------|---------|--------|
| ✅   | ✅     | ✅       | ✅      |

## Configuration

The following configuration values can be set to customize the cert-injection-webhook
installation.

| Value          | Required/Optional                        | Description                                                                                                   |
|----------------|------------------------------------------|---------------------------------------------------------------------------------------------------------------|
| `ca_cert_data` | Optional                                 | CA cert data to inject into pod trust store                                                                   |
| `labels`       | Required if annotations are not provided | Array of labels that will be used to match on pods that will have certs and proxy environment injected        |
| `annotations`  | Required if labels are not provided      | Array of annotations that will be used to match on pods that will have certs and proxy environment injected   |
| `http_proxy`   | Optional                                 | The HTTP proxy to inject into pod environment                                                                 |
| `https_proxy`  | Optional                                 | The HTTPS proxy to inject into pod environment                                                                |
| `no_proxy`     | Optional                                 | A comma-separated list of hostnames, IP addresses, or IP ranges in CIDR format to inject into pod environment |

## Installation

### Package Installation steps

1. Create a `cert-injection-webhook-config-values.yaml` with the labels or annotations (or both) that you would like to use.
   Any pod that matches one of these labels or annotations will have the provided cert injected. For example:

   ```yaml
   ---
   labels:
   - kpack.io/build
   annotations:
   - some-annotation
   ca_cert_data: |
     -----BEGIN CERTIFICATE-----
     MIICrDCCAZQCCQDcakcvwbW4UTANBgkqhkiG9w0BAQsFADAYMRYwFAYDVQQDDA1t
     eXdlYnNpdGUuY29tMB4XDTIyMDIxNDE2MjM1OVoXDTMyMDIxMjE2MjM1OVowGDEW
     MBQGA1UEAwwNbXl3ZWJzaXRlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
     AQoCggEBAMgWkhYr7OPSTuDwGSM5jMQtO5vnqfESPPh829IMTBNXkS0KV6Hi90ka
     T/gIbq0H+QO5Abzh8QDIOWqaTLLp5FedsU1xsGTiKQ+YVKfoQ7T7R/K+adWuJL6H
     i8kgb4ErzhYhDQqsPU6ZglKkTZTL+7fhpsc7ZewASa7TRJiSo51Qye9K1qsjj3Wd
     MB+0qH1vxvN2zs/117qowW/2YH2H++lJSfnEMH4Z67RQ5o56DpeHvE7mLz0LNVu/
     gyM8JXClgsPdr11Iiv17TevWoXSeoWa0ts6MGd/r376dtEZ60wGG+geXcf9szAx1
     GZLEQamRHnVyrGvb7U/AvLaJMnNY8PcCAwEAATANBgkqhkiG9w0BAQsFAAOCAQEA
     bc4XeX7sKvtEHK5tYKJDarP6suArgs7/IpfT2DiRB8JSBYX7rHD6NIB3433JxQfc
     SHD9FBpH9E8aSMDsCWKcuRRI7GeRarqwfblAqflCv85NJaiC9zu+haue7aNMNnwA
     uB+q0urjiKlEOM2OsLqgjXXmx5+nSrdwUhFXmyMsJC2eP4Dm1gJp5tQG2hSONC7w
     dX2wAQp7PYaq+h1ASkDNaKy3ZoeD7yEp3Mhbnh+fu0O06NpnJhUZPhdTtMD3LYPJ
     +iwL43iSAQt05ZK39u23zsdMc+RLFbqQYsULYZS2g/SmcSnw8CC3aer8X6x4lEw7
     FpCpA2Wta8mXHGKqmq0+og==
     -----END CERTIFICATE-----
   ```

You can install the cert-injection-webhook package using the command below -

`tanzu package install cert-injection-webhook --package-name cert-injection-webhook.community.tanzu.vmware.com --version <package-version> -f cert-injection-webhook-config-values.yaml`

### Injecting certificates into kpack builds

When providing ca_cert_data directly to kpack, that CA Certificate be injected into builds themselves.
If you want kpack builds to have CA Certificates for communicating with a self-signed registry,
make sure the values yaml has a label with `kpack.io/build`. This will match on any build pod that kpack creates.
