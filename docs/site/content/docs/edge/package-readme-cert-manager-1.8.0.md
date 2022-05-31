# cert-manager

This documentation provides information about the specific TCE package. Please visit the [TCE package management page](https://tanzucommunityedition.io/docs/v0.11/package-management/) for general information about installation, removal, troubleshooting, and other topics.

## Installation

### Installation of dependencies

The cert-manager package does not have any dependencies on other software or packages for installation.

### Installation of package

You can install the latest version of the cert-manager package with the Tanzu package CLI. More information for discovering packages and working with the Tanzu package CLI are available in the [Package Management](https://tanzucommunityedition.io/docs/v0.12/package-management/#discovering-available-packages) documentation.

```shell
tanzu package install cert-manager --package-name cert-manager.community.tanzu.vmware.com --version <<version>>
```

## Options

### Package configuration values

The cert-manager package does not require or have any command line options at this time.

### Application configuration values

You can set the following configuration values to customize the cert-manager installation:

| Config | Values | Default | Description |
|--------|--------| -----------|-------------|
| `namespace` | Any valid namespace | `cert-manager` | Optional. The namespace in which to deploy cert-manager. |

#### Multi-cloud configuration steps

There are currently no configuration steps necessary for installation of the cert-manager package to any provider.

## What This Package Does

[cert-manager](https://cert-manager.io/docs/) is a powerful and extensible X.509 certificate controller for Kubernetes and OpenShift workloads. It will obtain certificates from a variety of Issuers—both popular public Issuers as well as private Issuers—and ensure they are valid and up-to-date. It will also attempt to renew certificates at a configured time before expiry.

## Components

* cert-manager

### Supported Providers

The following table shows the providers this package can work with.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
| ✅  |  ✅  | ✅  |  ✅  |

## Files

Not relevant for this particular package.

## Package Limitations

This package provides cert-manager with basic, “off the shelf” functionality. There are no configurable parameters at this time other than the namespace that cert-manager can be deployed to.

## Usage Example

### ACME with Let's Encrypt and Contour

This example documents how to set up an Automated Certificate Management Environment (ACME) using cert-manager, Let's Encrypt, and Contour.

You’ll need to install the Contour package. You can do this with the Tanzu CLI.

```shell
tanzu package install contour --package-name contour.community.tanzu.vmware.com --namespace default --version 1.20.1
```

With Contour installed, you'll need to get the External IP of the Load Balancer for the Envoy service.
Depending on your cloud provider, this will either be an IP address or a DNS name. This example is from AWS.

```shell
kubectl get -n projectcontour service envoy
NAME    TYPE           CLUSTER-IP     EXTERNAL-IP                                                              PORT(S)                      AGE
envoy   LoadBalancer   100.69.86.69   ac98d7e4261e340668de971bbaab57b4-902577911.us-east-1.elb.amazonaws.com   80:31852/TCP,443:31911/TCP   3m22s
```

You'll have to use this External IP to configure the DNS for your domain. This of course will vary
based on your cloud provider and/or DNS service and is beyond the scope of this example.

Create a Cluster Issuer for the Let's Encrypt staging servers. We'll use staging because the production
servers are rate-limited. If you make certificate requests more than five times in a seven-day period you'll be limited.

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

> Be sure to use a valid email address, not the one with an `example.com` domain provided.

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

Modify the Ingress to add the TLS certificate.

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

Clean up your cluster.

```shell
kubectl delete ingress/nginx
kubectl delete service/nginx
kubectl delete deployment/nginx
kubectl delete clusterissuer/letsencrypt-staging
```

### Self-Signing Issuers

Create a self-signed Issuer to issue self-signed certificates. This can be used for creating your own PKI certificate within your
clusters or to generate a root certificate authority (CA).

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

Create a certificate using the self-signed Issuer. The `isCA` flag indicates that this certificate can be used for
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

Create a signed certificate.

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

## Troubleshooting

Here are some steps to troubleshoot an installation of cert-manager.

### Components

A successfully installed cert-manager package will have the following pods, services, deployments and replicasets which you can view by running a kubectl command.

```shell
kubectl get all -n cert-manager
NAME                                           READY   STATUS    RESTARTS   AGE
pod/cert-manager-565cf7b7f5-82ppk              1/1     Running   0          47s
pod/cert-manager-cainjector-5559d568d6-7qq5r   1/1     Running   0          47s
pod/cert-manager-webhook-7ccf4cc6b8-4dm76      1/1     Running   0          47s

NAME                           TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
service/cert-manager           ClusterIP   10.96.202.42   <none>        9402/TCP   47s
service/cert-manager-webhook   ClusterIP   10.96.76.165   <none>        443/TCP    47s

NAME                                      READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/cert-manager              1/1     1            1           47s
deployment.apps/cert-manager-cainjector   1/1     1            1           47s
deployment.apps/cert-manager-webhook      1/1     1            1           47s

NAME                                                 DESIRED   CURRENT   READY   AGE
replicaset.apps/cert-manager-565cf7b7f5              1         1         1       47s
replicaset.apps/cert-manager-cainjector-5559d568d6   1         1         1       47s
replicaset.apps/cert-manager-webhook-7ccf4cc6b8      1         1         1       47s
```

### cmctl

You can also use the `cmctl` process to validate that cert-manager is installed and configured correctly. `cmctl` is a command line application that can be used to control and check cert-manager functionality. You can get more information on `cmctl`, including installation instructions, [here](https://cert-manager.io/docs/usage/cmctl).

To use `cmctl` to verify a cert-manager installation, run the following command. If a component is not installed properly and ready, it will be reported here.

```shell
cmctl check api --wait=2m
The cert-manager API is ready
```

## Additional Documentation

See the [cert-manager documentation](https://cerr-manager.io/docs/) for more information.
