# Configuring the Knative Serving Package

This package provides serverless functionality using [Knative](https://knative.dev/).

## Components

* Knative Serving

## Configuration

There are no configuration options in the first release of this package.

### Installation

The knative-serving package requires Contour as an ingress controller. To successfully install and use the Knative Serving package, you must first install Contour.

```shell
tanzu package install contour-operator.tce.vmware.com
```

After the Contour package has been installed, you can install the Knative Serving package.

```shell
tanzu package install knative-serving.tce.vmware.com
```

## Usage Example

This example demonstrates the scale to zero feature of Knative Serving. You can watch the pods for the application start and quit automatically based on usage. It's based on the  ``helloworld-go`` sample app described in the [Knative documentation](https://knative.dev/docs/serving/samples/hello-world/helloworld-go/index.html).

### Before You Begin
Networking layer configuration is a Knative Serving prerequisite. You must install the Contour package before continuing with this example. Eventually, this will happen automatically, once we have a solution for dependency resolution.

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
1. At this point, you must configure DNS so that you are able to reach your service. Knative provides three possible types of DNS configurations; Magic DNS (xip.io), Real DNS, and Temporary DNS. In this example, we will use Magic DNS (xip.io) as it is easy to install and requires no configuration. Run the Magic DNS job provided by Knative.

    ```shell
    kubectl apply --filename https://github.com/knative/serving/releases/download/v0.18.0/serving-default-domain.yaml
    ```
    Magic DNS (xip.io) only works under certain conditions. For alternative DNS configurations for knative see the Configure DNS topic in the [Knative documentation](https://knative.dev/v0.22-docs/install/install-serving-with-yaml/).


1. Get the Knative services. This will show which applications are available and the URL to access them.

    ```shell
    kubectl get ksvc --namespace example


    NAMESPACE  NAME           URL
    example    helloworld-go  http://helloworld-go.default.18.204.46.247.xip.io
    ```

1. Make a request to the application. Make a request to the application via a cURL of the URL from the previous step.

    ```shell
    curl http://helloworld-go.default.18.204.46.247.xip.io

    Hello Go Sample v1!
    ```

1. If you still have a terminal window open a watch of the pods, keep an eye on the pod for helloworld-go. After about 2 minutes, it should go away. This is Knative _scaling to zero_, essentially stopping pods that are not being used. Once the pod has terminated, make another cURL request to your application and you should see the pod re-created on demand.
