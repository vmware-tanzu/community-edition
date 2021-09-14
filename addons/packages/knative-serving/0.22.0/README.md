# Knative Serving

This package provides serverless functionality using [Knative](https://knative.dev/).

## Components

* Knative Serving version: 0.22.0

## Configuration

| Config | Values | Default | Description |
|--------|--------|--------|-------------|
| namespace | any namespace | `knative-serving`| Namespace where you want to install knative |
| domain.type | real, sslip.io, nip.io | nip.io | Type of DNS resolution to use for your knative services. We can use real dns, in which case, you need to provide a domain.name or else use sslip.io or nip.io |
| domain.name | any domain name | empty | If you have a valid domain, make sure that it's properly configure to your ingress controller |
| ingress.external.namespace | any namespace | `projectcontour` | Namespace of the ingress controller for external services |
| ingress.internal.namespace | any namespace | `projectcontour` | Namespace of the ingress controller for internal services. If you don't want to have internal services separated from external, use the same namespace for both. |
| tls.certmanager.clusterissuer | `ClusterIssuerName` | empty | Name of a cert-manager ClusterIssuer to provide wildcard certificates for your cluster |

### Installation

The Knative Serving package requires use of Contour for ingress. To successfully install and use the Knative Serving package, you must first install Contour.

        tanzu package install contour --package-name contour.community.tanzu.vmware.com --namespace contour-external --version 1.17.1

After the Contour package has been installed, you can install Knative-Serving.

        tanzu package install knative-serving --package-name knative-serving.community.tanzu.vmware.com --namespace knative-serving --version 0.22.0

### Removal

To remove the Knative-Serving and Contour packages, issue package delete commands.

        tanzu package installed delete knative-serving
        tanzu package installed delete contour

## Usage Example

Official getting started documentation can be found on the Knative [site](https://knative.dev/docs/getting-started/).
Provided there are examples for:

* [Deploying a Knative service](https://knative.dev/docs/getting-started/first-service/)
* [Scaling to Zero](https://knative.dev/docs/getting-started/first-autoscale/)
* [Traffic Splitting](https://knative.dev/docs/getting-started/first-traffic-split/)

### Scale to Zero

This sample demonstrates the scale to zero feature of Knative. You can watch the pods for the application start and quit
automatically based on usage. It follows the Knative-Serving
[Hello World - Go](https://knative.dev/docs/serving/samples/hello-world/helloworld-go/index.html) instructions.

**Before beginning this guide, knative requires networking layer configuration
as a pre-requisite. Be sure to install the Contour package before continuing
with this guide.**

1. Create a service YAML file for your application.

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

1. Apply the service manifest

        kubectl apply --filename helloworld-service.yaml

1. In a separate terminal window, watch the pods.

        watch kubectl get pods --namespace example

1. At this point you will have a Magic DNS (nip.io) configured to provide you access to your service. If you want to use real DNS or any other Magic DNS you'll need to provide the appropriate configuration to the knative-serving package when installing it.

    > __NOTE__: nip only works if your Contour has a loadbalancer with a reachable IP address. If you don't have a loadbalancer or have a loadbalancer that only be reached by CNAME, you can use real DNS, see the configuration options for this package.

1. Get the Knative services. This will show which applications are available and the URL to access them.

        kubectl get ksvc --namespace example
        
        NAMESPACE   NAME              URL                                                   LATESTCREATED           LATESTREADY             READY   REASON
        example     helloworld-go     http://helloworld-go.default.18.204.46.247.nip.io     helloworld-go-00001     helloworld-go-00001     True

1. Make a request to the application via a cURL of the URL from the previous step.

        curl http://helloworld-go.default.18.204.46.247.nip.io
        
        Hello Go Sample v1!

1. If you still have a terminal window open a watch of the pods, keep an eye on the pod for helloworld-go. After about 2 minutes, it should go away. This is Knative _scaling to zero_, esssentially stopping pods that are not being used. Once the pod has terminated, make another cURL request to your application and you should see the pod re-created on demand.
