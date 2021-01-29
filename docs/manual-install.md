### Manual installation example

> This page exists to capture manual steps used during development. It will
likely be deleted in the future. TODO(joshrosso)

The following steps guide you through installing a sample extension, Velero, into your TCE cluster.

1. Create the `tanzu-extensions` namespace:

    ```shell
    kubectl create namespace tanzu-extensions
    ```

2. Deploy kapp-controller into the target cluster:

    ```shell
    kubectl create namespace kapp-controller
    kubectl --namespace kapp-controller apply --file https://github.com/vmware-tanzu/kapp-controller/releases/latest/download/release.yml
    ```

3. Validate that the kapp-controller started successfully:

    ```shell
    $ kubectl --namespace kapp-controller get deployment
    NAME              READY   UP-TO-DATE   AVAILABLE   AGE
    kapp-controller   1/1     1            1           8m34s
    ```

4. Deploy Velero extension:

    ```shell
    kubectl --namespace tanzu-extensions apply --file extensions/velero/extension.yaml
    ```

5. Validate the extension was deployed successfully by checking the App's description.

    ```shell
    $ kubectl --namespace tanzu-extensions get app velero
    NAME     DESCRIPTION           SINCE-DEPLOY   AGE
    velero   Reconcile succeeded   2m13s          2m47s
    ```

    If the description field shows an error, use the `kubectl describe` command to troubleshoot further:

    ```shell
    kubectl --namespace tanzu-extensions describe app velero
    ```

See the Velero extension [documentation](./extensions/velero) for more information, including a walkthrough that guides you through a usage example.
