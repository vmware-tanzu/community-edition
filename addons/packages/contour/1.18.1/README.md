# Contour Package

This package provides an ingress controller using [Contour](https://projectcontour.io/).

## Supported Providers

The following tables shows the providers this package can work with.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
| ✅  |  ✅  | ✅  | ✅  |

## Components

* Contour controller
* Envoy reverse proxy and load balancer

## Installation

### Environments with support for LoadBalancer services (AWS, Azure, vSphere)

1. Install the Contour package using default configuration options:

    ```shell
    tanzu package install contour \
      --package-name contour.community.tanzu.vmware.com \
      --version 1.18.1
    ```

### Environments without support for LoadBalancer services (Docker)

1. Create a configuration file that sets the Envoy service type to `ClusterIP`:

    ```shell
    cat <<EOF >contour-values.yaml
    envoy:
      service:
        type: ClusterIP
    EOF
    ```

1. Install the Contour package using the configuration file:

    ```shell
    tanzu package install contour \
      --package-name contour.community.tanzu.vmware.com \
      --version 1.18.1 \
      --values-file contour-values.yaml
    ```

For additional configuration options, see the [Configuration Reference section](#configuration-reference).

## Usage Example

The follow is a basic guide for getting started with Contour. You must deploy the package before attempting this walkthrough.

⚠️ Note: For more advanced use cases and documentation, see the official Contour [documentation](https://projectcontour.io/docs/).

1. Create a namespace for the example workload:

    ```shell
    kubectl create namespace contour-example-workload
    ```

1. Create an example deployment, in this case, nginx, to route traffic to via Contour:

    ```shell
    kubectl create deployment nginx-example --image nginx --namespace contour-example-workload
    ```

1. Create a service for nginx. This will map port 80 of the service to port 80 of the nginx app.

    ```shell
    kubectl create service clusterip nginx-example --tcp 80:80 --namespace contour-example-workload
   ```

1. Create a Contour `HTTPProxy` that directs traffic to the nginx instance:

    ```shell
    kubectl apply -f - <<EOF
    apiVersion: projectcontour.io/v1
    kind: HTTPProxy
    metadata:
      name: nginx-example-proxy
      namespace: contour-example-workload
      labels:
        app: ingress
    spec:
      virtualhost:
        fqdn: nginx-example.projectcontour.io
      routes:
      - conditions:
        - prefix: /
        services:
        - name: nginx-example
          port: 80
    EOF
    ```

1. Get the address to send traffic to. If you are using a `LoadBalancer` Envoy service, get the service's external IP:

    ```shell
    kubectl --namespace projectcontour get service envoy
    NAME    TYPE           CLUSTER-IP    EXTERNAL-IP            PORT(S)                      AGE
    envoy   LoadBalancer   10.12.10.93   <ENVOY-EXTERNAL-IP>    80:31232/TCP,443:32459/TCP   1h
    ```

   If you are using a `NodePort` or `ClusterIP` Envoy service, in a separate terminal, set up a port forward from your local machine to the service:

    ```shell
    kubectl --namespace projectcontour port-forward svc/envoy <LOCAL-PORT>:80
    ```

1. Make a request. If you are using a `LoadBalancer` Envoy service:

   ```shell
   curl -s -H "Host: nginx-example.projectcontour.io" <ENVOY-EXTERNAL-IP> | grep title
   <title>Welcome to nginx!</title>
   ```

   If you are using a `NodePort` or `ClusterIP` Envoy service:

   ```shell
   curl -s -H "Host: nginx-example.projectcontour.io" http://localhost:<LOCAL-PORT> | grep title
   <title>Welcome to nginx!</title>
   ```

## Configuration Reference

You can configure the following:

| Config | Default | Description |
|--------|---------|-------------|
| `namespace` | `projectcontour` | The namespace in which to deploy Contour and Envoy. |
| `contour.configFileContents` | (none) | The YAML contents of the Contour config file. See [the Contour configuration documentation](https://projectcontour.io/docs/v1.18.1/configuration/#configuration-file) for more information. |
| `contour.replicas` | `2` | How many Contour pod replicas to have. |
| `contour.useProxyProtocol` | `false` | Whether to enable PROXY protocol for all Envoy listeners. |
| `contour.logLevel` | `info` | The Contour log level. Valid values are `info` and `debug`. |
| `envoy.service.type` | `LoadBalancer` | The type of Kubernetes service to provision for Envoy. Valid values are `LoadBalancer`, `NodePort`, and `ClusterIP`. |
| `envoy.service.externalTrafficPolicy` | `Local` | The external traffic policy for the Envoy service. Valid values are `Local` and `Cluster`.  If `envoy.service.type` is `ClusterIP`, this field is ignored. |
| `envoy.service.annotations` | (none) | Annotations to set on the Envoy service. |
| `envoy.service.nodePorts.http` | (none) | If `envoy.service.type` == `NodePort`, the node port number to expose Envoy's HTTP listener on. If not specified, a node port will be auto-assigned by Kubernetes. |
| `envoy.service.nodePorts.https` | (none) | If `envoy.service.type` == `NodePort`, the node port number to expose Envoy's HTTPS listener on. If not specified, a node port will be auto-assigned by Kubernetes. |
| `envoy.hostPorts.enable` | `false` | Whether to enable host ports for the Envoy pods. If false, `envoy.hostPorts.http` and `envoy.hostPorts.https` are ignored. |
| `envoy.hostPorts.http` | `80` | If `envoy.hostPorts.enable` == true, the host port number to expose Envoy's HTTP listener on. |
| `envoy.hostPorts.https` | `443` | If `envoy.hostPorts.enable` == true, the host port number to expose Envoy's HTTPS listener on. |
| `envoy.hostNetwork` | `false` | Whether to enable host networking for the Envoy pods. |
| `envoy.terminationGracePeriodSeconds` | `300` | The termination grace period, in seconds, for the Envoy pods. |
| `envoy.logLevel` | `info` | The Envoy log level. Valid values are `trace`, `debug`, `info`, `warn`, `error`, `critical`, and `off`. |
| `certificates.useCertManager` | `false` | Whether to use cert-manager to provision TLS certificates for securing communication between Contour and Envoy. If false, the upstream Contour certgen job will be used to provision certificates. If true, the `cert-manager` addon must be installed in the cluster. |
| `certificates.duration` | `8760h` |  If using cert-manager, how long the certificates should be valid for. If useCertManager is false, this field is ignored. |
| `certificates.renewBefore` | `360h` |  If using cert-manager, how long before expiration the certificates should be renewed. If useCertManager is false, this field is ignored. |
