# ingress Extension

This extension provides disaster recovery capabilities using [ingress](https://ingress.io/). At the moment, it leverages [minio](https://github.com/minio/minio) for object storage.

## Components

* ingress Namespace
* Contour CRDs
* Contour Deployment
* Envoy Deployment

## Configuration

Change the following values to suit your environment:

```yaml
# the cloud provider you are installed into (currently aws or kind)
tanzu_provider: kind
# the fqdn of your ingress domain
tanzu_ingress.domain: 127.0.0.1.xip.io
# If you have a static eip you wish to use in AWS
tanzu_ingress.aws_eip_allocation_id: eipalloc-xxxxxx
```

## Usage Example

Create a kind cluster:

```bash
kind create cluster --config ./kind-config.yaml
```

Create the tanzu-extensions namespace:

```bash
kubectl create namespace tanzu-extensions
```

Process cert-manager:

```bash
ytt --ignore-unknown-comments -f ./dependencies/cert-manager > deploy/cert-manager.yaml
```

Process ingress:

```bash
ytt --ignore-unknown-comments -f ./config -f ./values > deploy/ingress.yaml
```

Deploy:

```bash
kapp deploy -n tanzu-extensions -a tanzu-ingress -p ./deploy
```

Verify with Demo:

```bash
kubectl apply -f demo/deployment.yaml
```

Validate the Ingress and Certificates are created:

```bash
$ kubectl -n demo get ingress,certificate
NAME                       CLASS    HOSTS                   ADDRESS   PORTS     AGE
ingress.extensions/nginx   <none>   demo.127.0.0.1.xip.io             80, 443   62s

NAME                                         READY   SECRET           AGE
certificate.cert-manager.io/nginx-cert-tls   True    nginx-cert-tls   62s
```

Validate Ingress is working:

```bash
$ curl demo.127.0.0.1.xip.io
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```

Cleanup:

```bash
kind delete cluster
```