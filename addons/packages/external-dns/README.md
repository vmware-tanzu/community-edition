# ExternalDNS

[ExternalDNS](https://github.com/kubernetes-sigs/external-dns) synchronizes exposed Kubernetes Services and Ingresses with DNS providers.

## Components

* ExternalDNS deployment

## Configuration

The following configuration values can be set to customize the external-dns installation.

### Global

| Value       | Required/Optional | Description                                    |
|-------------|-------------------|------------------------------------------------|
| `namespace` | Optional          | The namespace in which to deploy external-dns. |

### external-dns Configuration

| Value                        | Required/Optional | Description                                       |
|------------------------------|--------------------|--------------------------------------------------|
| `deployment.args`            | Required           | Args passed via command-line to external-dns     |
| `deployment.env`             | Optional           | Environment variables to pass to external-dns    |
| `deployment.securityContext` | Optional           | Security context of the external-dns container   |
| `deployment.volumeMounts`    | Optional           | Volume mounts of the external-dns container      |
| `deployment.volumes`         | Optional           | Volumes of the external-dns pod                  |

Follow [the external-dns docs](https://github.com/kubernetes-sigs/external-dns#running-externaldns)
for guidance on how to configure ExternalDNS for your DNS provider.

### Configuration sample

After installing this add-on with the name, for example, `external-dns`, the
following command will generate an empty configuration file in the current directory:

`tanzu package configure external-dns`

A sample of how to fill in that empty configuration file is given below, for a simple `bind`
(rfc2136) implementation. Note that comments which
begin with `#@` are important `ytt` directives and should remain unchanged in
your final configuration file.

```yaml
#@data/values
#@overlay/match-child-defaults missing_ok=True
---

#! The namespace in which to deploy ExternalDNS.
namespace: external-dns

#! Deployment related configuration
deployment:
  #@overlay/replace
  args:
  - --source=service
  - --source=contour-httpproxy
  - --txt-owner-id=k8s
  - --domain-filter=k8s.example.org
  - --namespace=tanzu-system-service-discovery
  - --provider=rfc2136
  - --rfc2136-host=100.69.97.77
  - --rfc2136-port=53
  - --rfc2136-zone=k8s.example.org
  - --rfc2136-tsig-secret=MTlQs3NNU=
  - --rfc2136-tsig-secret-alg=hmac-sha256
  - --rfc2136-tsig-keyname=externaldns-key
  - --rfc2136-tsig-axfr
  env: []
  securityContext: []
  volumeMounts: []
  volumes: []
```

### Configuring with Contour HTTPProxy

Follow [this tutorial](https://github.com/kubernetes-sigs/external-dns/blob/v0.7.6/docs/tutorials/contour.md)
for guidance on providing arguments to ExternalDNS to enable HTTPProxy support. The ExternalDNS package is
preconfigured with the correct RBAC permissions to watch for HTTPProxies, so this part of the tutorial
may be skipped.

## Usage Example

This walkthrough guides you through setting up a hostname for a Service. You must deploy the package before attempting this walkthrough.

⚠️ Note: For more advanced use cases and documentation, see the official ExternalDNS [documentation](https://github.com/kubernetes-sigs/external-dns).

Run an application and expose it via a Kubernetes Service:

```
kubectl run nginx --image=nginx --port=80
kubectl expose pod nginx --port=80 --target-port=80 --type=LoadBalancer
```

Annotate the Service with your desired external DNS name. Make sure to change example.org to your domain.

```
kubectl annotate service nginx "external-dns.alpha.kubernetes.io/hostname=nginx.example.org."
```

Check that ExternalDNS has created the desired DNS record for your Service and that it points to its load balancer's IP. Then try to resolve it:

```
dig +short nginx.example.org.
104.155.60.49
```
