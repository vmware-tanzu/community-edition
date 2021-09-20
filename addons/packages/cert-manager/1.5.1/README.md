# cert-manager

This package provides certificate management functionality using [cert-manager](https://cert-manager.io/docs/).

## Supported Providers

The following tables shows the providers this package can work with.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
| ✅  |  ✅  | ✅  | ❌  |

## Components

* cert-manager version: `1.5.1`

## Configuration

The following configuration values can be set to customize the cert-manager installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Optional | The namespace in which to deploy cert-manager. |

## Usage Examples

### ACME with Let's Encrypt and Contour

This example will document how to setup an Automated Certificate Management Environment (ACME) using
cert-manager, Let's Encryot and Contour.

First step is to install cert-manager and Contour. You can do this with the Tanzu CLI.

```shell
tanzu package install cert-manager --package-name cert-manager.community.tanzu.vmware.com --namespace default --version 1.3.1
tanzu package install contour --package-name contour.community.tanzu.vmware.com --namespace default --version 1.17.1
```

With Contour installed, you'll need to get the External IP of the Load Balancer for the envoy service.
Depending on your cloud provider, this will either be an IP address or a DNS name. This example is from AWS.

```shell
kubectl get -n projectcontour service envoy
NAME    TYPE           CLUSTER-IP     EXTERNAL-IP                                                              PORT(S)                      AGE
envoy   LoadBalancer   100.69.86.69   ac98d7e4261e340668de971bbaab57b4-902577911.us-east-1.elb.amazonaws.com   80:31852/TCP,443:31911/TCP   3m22s
```

You'll have to use this External IP to configure your DNS for your domain. This of course will vary
based on your cloud provider and/or DNS service and is beyond the scope of this example.

Create a Cluster Issuer for the Let's Encrypt staging servers. We'll use staging because the production
servers are rate limited. Make certificate requests more than 5 times in a 7 day period and you'll be limited.

```shell
cat <<EOF | kubectl apply --filename -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-staging
  namespace: cert-manager
spec:
  acme:
    email: j.engineer@example.com
    privateKeySecretRef:
      name: letsencrypt-staging
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    solvers:
    - http01:
        ingress:
          class: contour
EOF
```

> Be sure to use a valid email address other than one with an `example.com` domain

For simplicity, we'll deploy the typical nginx server and service.

```shell
kubectl create deployment nginx --image nginx
```

```shell
cat <<EOF | kubectl apply --filename -
apiVersion: v1
kind: Service
metadata:
  name: nginx
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: nginx
EOF
```

Next we'll create an Ingress to route traffic to nginx.

```shell
cat <<EOF | kubectl apply --filename -
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: nginx
spec:
  rules:
  - host: nginx.example.com
    http:
      paths:
      - backend:
          service:
            name: nginx
            port:
              number: 80
        pathType: ImplementationSpecific
EOF
```

Test that you can reach the nginx server.

```shell
curl nginx.example.com
```

Modify the ingress to add the TLS certificate.

```shell
cat <<EOF | kubectl apply --filename -
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-staging
    ingress.kubernetes.io/force-ssl-redirect: "true"
    kubernetes.io/ingress.class: contour
    kubernetes.io/tls-acme: "true"
spec:
  tls:
  - secretName: nginx-tls
    hosts:
    - nginx.example.com
  rules:
  - host: nginx.example.com
    http:
      paths:
      - pathType: Prefix
        path: "/"
        backend:
          service:
            name: nginx
            port:
              number: 80
EOF
```

Try the curl command again, this time with `https`. We use the `-k` option because certificates from
the Let's Encrypt staging servers are not trusted.

```shell
curl -k https://nginx.example.com
```

Cleanup your cluster.

```shell
kubectl delete ingress/nginx
kubectl delete service/nginx
kubectl delete deployment/nginx
kubectl delete clusterissuer/letsencrypt-staging
```

### Self Signing Issuers

Create a self-signed Issuer to issue self-signed certificates. This can be used for creating your own PKI within your
clusters or to generate a root certificate authority (CA)

Create a self-signed issuer.

```shell
cat <<EOF | kubectl apply --filename -
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: self-signed
spec:
  selfSigned: {}
EOF
```

Create an Issuer for your private certificate authority.

```shell
cat <<EOF | kubectl apply --filename -
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: private-ca
spec:
  ca:
    secretName: private-ca
EOF
```

Create a certificate using the self-signed issuer. The `isCA` flag indicates that this certificate can be used for
signing other certificates.

```shell
cat <<EOF | kubectl apply --filename -
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: private-ca
spec:
  isCA: true
  duration: 2160h
  secretName: private-ca
  commonName: private-ca
  subject:
    organizations:
      - cert-manager
  issuerRef:
    name: self-signed
    kind: Issuer
    group: cert-manager.io
EOF
```

Create a signed certificate

```shell
cat <<EOF | kubectl apply --filename -
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: example-com
  namespace: default
spec:
  secretName: example-com-tls
  issuerRef:
    name: private-ca
    kind: Issuer
  commonName: example.com
  dnsNames:
    - example.com
    - www.example.com
EOF
```

Review your new certificate.

```shell
kubectl get certificate/example-com

NAME          READY   SECRET            AGE
example-com   True    example-com-tls   5s
```
