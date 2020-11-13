# Tanzu Community Edition

A Kubernetes distribution and modular application platform.

## Extensions

Extensions provide the additional functionality necessary to build an application platform atop Kubernetes. We follow a modular approach in which operators building a platform can deploy the extensions they need to fulfill their requirements.

| Name | Description | Documentation |
|------|-------------|---------------|
| Velero | Provides disaster recovery capabilities | [Velero extension docs](./extensions/velero) |

### Quick Start with Velero

The following steps guide you through installing a sample extension, Velero, into your TCE cluster.

1. Create the `tanzu-extensions` namespace:

    ```shell
    kubectl create ns tanzu-extensions
    ```

2. Deploy kapp controller into the target cluster:

    ```shell
    # TODO: deploy to tanzu-extensions namespace instead?
    kubectl create ns kapp-controller
    kubectl -n kapp-controller apply -f https://github.com/k14s/kapp-controller/releases/latest/download/release.yml
    ```

3. Validate that the kapp-controller started successfully:

    ```shell
    $ kubectl -n kapp-controller get deployment
    NAME              READY   UP-TO-DATE   AVAILABLE   AGE
    kapp-controller   1/1     1            1           8m34s
    ```

4. Deploy Velero extension:

    ```shell
    kubectl -n tanzu-extensions apply -f extensions/velero/extension.yaml
    ```

5. Validate the extension was deployed successfully by checking the App's description.

    ```shell
    $ kubectl -n tanzu-extensions get app velero
    NAME     DESCRIPTION           SINCE-DEPLOY   AGE
    velero   Reconcile succeeded   2m13s          2m47s
    ```

    If the description field shows an error, use the `kubectl describe` command to troubleshoot further:

    ```shell
    kubectl -n tanzu-extensions describe app velero
    ```

See the Velero extension [documentation](./extensions/velero) for more information, including a walkthrough that guides you through a usage example.

### Extension Structure

Each extension lives in its own directory within the `extensions/` directory.

An extension is composed of the following assets:

* `extension.yaml`: Contains the [kapp-controller](https://github.com/vmware-tanzu/carvel-kapp-controller) `App` resource. This file also includes a Service Account and a ClusterRoleBinding for kapp-controller to use to deploy the extension.
* `config/` directory: Contains the Kubernetes deployment manifests necessary to deploy the extension. Manifests can be templatized with [`ytt`](https://github.com/vmware-tanzu/carvel-ytt), when necessary.

### Extension Delivery

The deployment manifests of an extension are bundled into an OCI image using [imgpkg](https://github.com/vmware-tanzu/carvel-imgpkg). The OCI image is pushed to a container registry for consumption.

### Extension Deployment

Assuming you have a TCE cluster up and running, the extension deployment process works as follows:

1. Install the `extension.yaml` file of a given extension. This initiates the App resource reconciliation.
2. `kapp-controller` fetches the OCI image that contains the extension's deployment manifests
3. `kapp-controller` renders the templates using `ytt`
4. `kapp-controller` deploys the extension

### Creating a new Extension

To create a new extension, you must:

1. Create a new directory within the `extensions/` directory. The name of the extension must be the name of the directory.
2. Create an `extension.yaml` file which contains a kapp-controller App custom resource.
3. Create a `config` directory which contains the extension's deployment manifests.
4. Templatize the deployment manifests, when necessary.
5. Bundle the deployment manifests into an OCI image using [imgpkg](https://github.com/vmware-tanzu/carvel-imgpkg) and push to the OCI registry.
6. Test the extension.
