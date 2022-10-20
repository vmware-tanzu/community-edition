# Knative-Serving

This package provides serverless functionality using Knative.

This documentation provides information about the specific TCE package. Please visit the [TCE package management page](https://tanzucommunityedition.io/docs/v0.11/package-management/) for general information about installation, removal, troubleshooting, and other topics.

## Installation

### Installation of dependencies

To successfully install and use the Knative Serving package, you must first install [Contour](https://tanzucommunityedition.io/docs/v0.12/package-readme-contour-1.20.1/) for ingress.

### Installation of package

After you’ve installed the Contour package, you can install Knative-Serving. Knative-Serving requires configuration of the type of DNS resolution that you will be using. This will require you to create a data values file. You can learn more about configuration in the [Options](#options) section.

Create a data values file.

```shell
cat <<EOF >knative-values.yaml
domain:
  type: real
  name: example.com
EOF
```

Install Knative-Serving

```shell
tanzu package install knative-serving \
–package-name knative-serving.community.tanzu.vmware.com \
–version 1.0.0 \
–values-file knative-values.yaml
```

## Options

### Package configuration values

| Config | Values | Default | Description |
|--------|--------|--------|-------------|
| namespace | any namespace | `knative-serving`| Namespace where you want to install knative. |
| domain.type | real, sslip.io, nip.io | nip.io | Type of DNS resolution to use for your Knative services. If you use real DNS, you need to provide a domain.name or else use sslip.io or nip.io. |
| domain.name | any domain name | empty | If you have a valid domain, make sure that it's properly configured to your ingress controller. |
| ingress.external.namespace | any namespace | `projectcontour` | Namespace of the ingress controller for external services. |
| ingress.internal.namespace | any namespace | `projectcontour` | Namespace of the ingress controller for internal services. If you don't want to have internal services separated from external, use the same namespace for both. |
| tls.certmanager.clusterissuer | `ClusterIssuerName` | empty | Name of a cert-manager ClusterIssuer to provide wildcard certificates for your cluster. |

### Application configuration values

There are no additional application configuration options that can be passed to the Knative-Serving application at this time.

## What This Package Does

Knative Serving provides components that enable rapid deployment of serverless containers and autoscaling of pods.

## Components

* Knative Serving

### Supported Providers

The following table shows the providers this package can work with.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
| ✅  |  ✅  | ✅  | ✅ |

## Files

The Knative-Serving package accepts a data values file to customize the package's behavior. See Package configuration values for details of the configurable settings. An example of the data values file looks like:

```yaml
    domain:
      type: real
      name: example.com
```

## Package Limitations

The package only supports Knative-Serving. Knative-Eventing is not currently supported.

## Usage Example

### Scale to Zero

This sample demonstrates Knative’s scale to zero feature, which  enables you to watch the pods for your application start and quit automatically based on usage. It follows the Knative-Serving [Hello World - Go](https://knative.dev/docs/serving/samples/hello-world/helloworld-go/index.html) instructions.

#### Networking Configuration

##### Amazon Web Services Configuration

If you have a domain name and are using Route53, it is simple to configure networking.

1. Start by identifying the external address of the Envoy service:

   ```shell
   > kubectl get services/envoy --namespace projectcontour \
     -o custom-columns=EXTERNALIP:".status.loadBalancer.ingress[*].hostname"
   EXTERNALIP
   a096a354b2818471e91c583ff2fb01f5-975360063.us-east-1.elb.amazonaws.com
   ```

2. Go to your hosted domain on [Route53](https://console.aws.amazon.com/route53/v2/home#Dashboard) and create/edit the `A` record for the domain.

3. In the `Edit Record` pane, in the `Route traffic to` section, select the load balancer/address that Envoy is using.

4. You should now have DNS for your domain pointed to the load balancer for the Contour/Envoy service.

#### Deploying Knative Serving

1. Optional: If you are using AWS and have a domain name, create a values.yaml file to configure your domain. Replace `example.com` with your domain.

    ```shell
    > cat > knative-values.yaml <<EOF
    domain:
      type: real
      name: example.com
    EOF
    ```

2. Install the knative-serving package. Refer to the [installation](#installation) section for more details.

3. Create a service YAML file for your application.

    ```shell
    cat <<EOF > helloworld-service.yaml
    ---
    apiVersion: v1
    kind: Namespace
    metadata:
      name: example
    ---
    apiVersion: serving.knative.dev/v1
    kind: Service
    metadata:
      name: helloworld-go
      namespace: example
    spec:
      template:
        spec:
          containers:
            - image: gcr.io/knative-samples/helloworld-go
              env:
                - name: TARGET
                  value: "Go Sample v1"
    EOF
    ```

4. Apply the service manifest.

    ```shell
    kubectl apply --filename helloworld-service.yaml
    ```

5. In a separate terminal window, watch the pods.

    ```shell
    watch kubectl get pods --namespace example
    ```

6. At this point you will have a [Magic DNS](https://nip.io) configured to provide you with access to your service. If you want to use real DNS or any other Magic DNS you'll need to provide the appropriate configuration to the knative-serving package when installing it.

   > __NOTE__: nip only works if your Contour has a load balancer with a reachable IP address. If you don't have a load balancer, or have one that can only be reached by CNAME, you can use real DNS; see the configuration options for this package.

7. Get the Knative services. This will show which applications are available and the URL to access them.

    ```shell
    kubectl get ksvc --namespace example

    NAMESPACE   NAME              URL                                                   LATESTCREATED           LATESTREADY             READY   REASON
    example     helloworld-go     http://helloworld-go.default.18.204.46.247.nip.io     helloworld-go-00001     helloworld-go-00001     True
    ```

8. Make a request to the application via a cURL of the URL from the previous step.

    ```shell
    > curl http://helloworld-go.example.18.204.46.247.nip.io
    Hello Go Sample v1!
    ```

9. Another means of testing is to cURL from a temporary pod within your cluster.

    ```shell
    kubectl run tmp --restart=Never --image=nginx -i --rm -- curl http://helloworld-go.example.svc.cluster.local
    ```

   > If you installed Knative-Serving to AWS with your own domain name, you can `curl http://helloworld-go.example.com`.

10. If you still have a terminal window open with a watch of the pods, keep an eye on it for the helloworld-go pod. After about two minutes, it should go away. This is Knative _scaling to zero_—essentially stopping pods that are not being used. Once the pod has terminated, make another cURL request to your application and you should see the pod recreated on demand.

## Troubleshooting

See Knative’s [Debugging application issues page](https://knative.dev/docs/serving/troubleshooting/debugging-application-issues/) for information on troubleshooting Knative-Serving.

## Additional Documentation

You can find official “getting started” documentation on the Knative [site](https://knative.dev/docs/getting-started/). Provided are examples for:

* [Deploying a Knative service](https://knative.dev/docs/getting-started/first-service/)
* [Scaling to Zero](https://knative.dev/docs/getting-started/first-autoscale/)
* [Traffic Splitting](https://knative.dev/docs/getting-started/first-traffic-split/)
