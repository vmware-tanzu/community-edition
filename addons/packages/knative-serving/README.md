# Knative Serving

This package provides serverless functionality using [Knative](https://knative.dev/).

## Components

* Knative Serving

## Configuration

There are no configuration options in the first release of this package.

### Installation

The knative-serving package requires use of Contour for ingress. To successfully install and use the knative-serving package, you must first install Contour.

```shell
tanzu package install contour
```

After the Contour package has been installed, you can install knative-serving.

```shell
tanzu package install knative-serving
```

## Usage Example

This sample demonstrates the scale to zero feature of Knative. You can watch the pods for the application start and quit automatically based on usage. It follows the knative-serving [Hello World - Go](https://knative.dev/docs/serving/samples/hello-world/helloworld-go/index.html) instructions.

**Before beginning this guide, knative requires networking layer configuration
as a pre-req. Be sure to install the TCE Contour package before continuing
with this guide. Eventually, this will happen automatically, once we have a
solution for dependency resolution.**

1. Create a service YAML file for your application.

```shell
cat <<EOF >> helloworld-service.yaml
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

1. Apply the service manifest

```shell
kubectl apply --filename helloworld-service.yaml
```

1. In a separate terminal window, watch the pods.

```shell
watch kubectl get pods --namespace example
```

1. At this point you will have to configure DNS so that you are able to reach your service. The Knative serving documenation provides instructions for 3 different types of configurations. In this example, we will use Magic DNS (xip.io) as it is very easy to install and requires no configuration. Run the Magic DNS job provided by Knative.

```shell
kubectl apply --filename https://github.com/knative/serving/releases/download/v0.18.0/serving-default-domain.yaml
```

> xip only works under some conditions. For alternative DNS configurations for
knative see the [config DNS section
here](https://knative.dev/docs/install/any-kubernetes-cluster/#installing-the-serving-component).
As an alternative, you can use real DNS or temporary DNS.

1. Get the Knative services. This will show which applications are available and the URL to access them.

```shell
kubectl get ksvc --namespace example

NAMESPACE   NAME              URL                                                   LATESTCREATED           LATESTREADY             READY   REASON
example     helloworld-go     http://helloworld-go.default.18.204.46.247.xip.io     helloworld-go-00001     helloworld-go-00001     True
```

1. Make a request to the application. Make a request to the application via a cURL of the URL from the previous step.

```shell
curl http://helloworld-go.default.18.204.46.247.xip.io

Hello Go Sample v1!
```

1. If you still have a terminal window open a watch of the pods, keep an eye on the pod for helloworld-go. After about 2 minutes, it should go away. This is Knative _scaling to zero_, esssentially stopping pods that are not being used. Once the pod has terminated, make another cURL request to your application and you should see the pod re-created on demand.
