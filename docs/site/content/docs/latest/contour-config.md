# Configuring the Contour Package

This package provides an ingress controller using [Contour](https://projectcontour.io/). The contour-operator is used.

## Components

* Contour operator
* Envoy reverse proxy and load balancer

## Usage Example

This example provides steps for setting up a very basic ingress route.

### Before you begin
Ensure the Contour package is installed, for more information about installing packages, see [Packages Introduction](packages-intro.md).

### Procedure
1. Create a namespace for the example

    ```shell
    kubectl create namespace ingress
    ```

1. Create an example deployment, in this case, Nginx

    ```shell
    kubectl create deployment ingress --image nginx --namespace ingress
    ```

1. Create a service for Nginx. This will map port 80 of the service to port 80 of the Nginx app.

    ```shell
    kubectl create service clusterip ingress --tcp=80:80 --namespace=ingress
   ```

1. Create a YAML file for the ingress configuration.

    ```shell
    cat <<EOF >> ingress.yaml
    apiVersion: networking.k8s.io/v1beta1
    kind: Ingress
    metadata:
        name: ingress
        namespace: ingress
        labels:
            app: ingress
    spec:
        backend:
            serviceName: ingress
            servicePort: 80
    EOF
    ```

1. Apply the ingress YAML

    ```shell
    kubectl apply --file ingress.yaml
    ```

1. Get the external address for the ingress.

    ```shell
    k get ingress -A                                                                                                                                                               ─╯
    NAMESPACE   NAME      CLASS    HOSTS   ADDRESS                                                                PORTS   AGE
    ingress     ingress   <none>   *       abc123abc123abc123abc123abc123ab-1234567.us-east-1.elb.amazonaws.com   80      52s
    ```

1. Enjoy the new ingress!

    ```shell
    curl abc123abc123abc123abc123abc123ab-1234567.us-east-1.elb.amazonaws.com 2>/dev/null | grep title                                                                             ─╯
    <title>Welcome to nginx!</title>
    ```
For more more information on advanced use cases and documentation, see the official Contour [documentation](https://projectcontour.io/docs/).