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
