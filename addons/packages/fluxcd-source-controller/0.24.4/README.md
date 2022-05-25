# Fluxcd-source-controller

The source-controller is a Kubernetes operator, specialised in artifacts acquisition
from external sources such as Git, Helm repositories and S3 buckets.
The source-controller implements the
[source.toolkit.fluxcd.io](https://github.com/fluxcd/source-controller/tree/master/docs/spec/v1beta1) API
and is a core component of the [GitOps toolkit](https://toolkit.fluxcd.io).

## Configuration

The fluxcd-source-controller package has following configurable properties.

| Value           | Required/Optional |                           Description                                                |
|-----------------|-------------------|--------------------------------------------------------------------------------------|
| `namespace`     | Optional          | Sets namespace for k8s objects where resources of source-controller will be created. |
| `limits_cpu`    | Optional          | Sets maximum usage of cpu for source-controller deployment.                          |
| `limits_memory` | Optional          | Sets maximum usage of memory for source-controller deployment.                       |
| `no_proxy`      | Optional          | Set domains for which no proxying should be performed                                |
| `https_proxy`   | Optional          | Set secure proxy connection information                                              |
| `http_proxy`    | Optional          | Set unsecure proxy connection information                                            |

```yaml
namespace: source-system
resources:
 limits_cpu: 1050m
 limits_memory: 2Gi
proxy:
 no_proxy: ""
 https_proxy: ""
 http_proxy: ""
service_port: 80
```

## Installation

To install FluxCD source-controller from the Tanzu Application Platform package repository:

1. List version information for the package by running:

    ```shell
    tanzu package available list fluxcd-source-controller.community.tanzu.vmware.com
    ```

    For example:

    ```shell
    \ Retrieving package versions for fluxcd-source-controller.community.tanzu.vmware.com...
      NAME                                                 VERSION  RELEASED-AT
      fluxcd-source-controller.community.tanzu.vmware.com  0.21.2   2022-02-07 06:14:08 -0500 -05
      fluxcd-source-controller.community.tanzu.vmware.com  0.21.3   2022-02-07 06:14:08 -0500 -05
      fluxcd-source-controller.community.tanzu.vmware.com  0.24.4   2022-02-07 06:14:08 -0500 -05
    ```

2. Install the package by running:

    ```shell
    tanzu package install fluxcd-source-controller -p fluxcd-source-controller.community.tanzu.vmware.com -v VERSION-NUMBER
    ```

    Where:

    - `VERSION-NUMBER` is the version of the package listed in step 1.

    For example:

    ```shell
    tanzu package install fluxcd-source-controller -p fluxcd-source-controller.community.tanzu.vmware.com -v 0.24.4
    \ Installing package 'fluxcd-source-controller.community.tanzu.vmware.com'
    | Getting package metadata for 'fluxcd-source-controller.community.tanzu.vmware.com'
    | Creating service account 'fluxcd-source-controller-default-sa'
    | Creating cluster admin role 'fluxcd-source-controller-default-cluster-role'
    | Creating cluster role binding 'fluxcd-source-controller-default-cluster-rolebinding'
    | Creating package resource
    / Waiting for 'PackageInstall' reconciliation for 'fluxcd-source-controller'
    / 'PackageInstall' resource install status: Reconciling


     Added installed package 'fluxcd-source-controller'
    ```

3. Verify the package install by running:

    ```shell
    tanzu package installed get fluxcd-source-controller
    ```

    For example:

    ```shell
    \ Retrieving installation details for fluxcd-source-controller...
    NAME:                    fluxcd-source-controller
    PACKAGE-NAME:            fluxcd-source-controller.community.tanzu.vmware.com
    PACKAGE-VERSION:         0.24.4
    STATUS:                  Reconcile succeeded
    CONDITIONS:              [{ReconcileSucceeded True  }]
    USEFUL-ERROR-MESSAGE:
    ```

    Verify that `STATUS` is `Reconcile succeeded`

    ```shell
    kubectl get pods -n source-system
    ```

    For example:

    ```shell
    $ kubectl get pods -n source-system
    NAME                                 READY   STATUS    RESTARTS   AGE
    source-controller-69859f545d-ll8fj   1/1     Running   0          3m38s
    ```

    Verify that `STATUS` is `Running`

## Try fluxcd-source-controller

1. Verify all the objects are installed:

    This package would create a new namespace where all the elements of fluxcd will be hosted called `source-system`

    you can verify the main components of `fluxcd-source-controller` were installed by running:

    ```shell
    kubectl get all -n source-system
    NAME                                     READY   STATUS    RESTARTS   AGE
    pod/source-controller-7684c85659-2zfxb   1/1     Running   0          40m

    NAME                        TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
    service/source-controller   ClusterIP   10.108.138.74   <none>        80/TCP    40m

    NAME                                READY   UP-TO-DATE   AVAILABLE   AGE
    deployment.apps/source-controller   1/1     1            1           40m

    NAME                                           DESIRED   CURRENT   READY   AGE
    replicaset.apps/source-controller-7684c85659   1         1         1       40m
    ```

    you should get something really similar!

2. Verify all the CRD were installed correctly:

    The way you would communicate with `fluxcd-source-controller` would be through its CRDs, this will be your main action point.

    In order to check all the CRDs were installed you can run:

    ```shell
    kubectl get crds -n source-system | grep ".fluxcd.io"
    buckets.source.toolkit.fluxcd.io                         2022-03-07T19:20:14Z
    gitrepositories.source.toolkit.fluxcd.io                 2022-03-07T19:20:14Z
    helmcharts.source.toolkit.fluxcd.io                      2022-03-07T19:20:14Z
    helmrepositories.source.toolkit.fluxcd.io                2022-03-07T19:20:14Z
    ```

3. Try one simple example:

    Try one quick example yourself, so you can check everything is working as expected

    - Let's consume a `GitRepository` object,

      - Create the following `gitrepository-sample.yaml` file:

          ```yaml
          apiVersion: source.toolkit.fluxcd.io/v1beta1
          kind: GitRepository
          metadata:
            name: gitrepository-sample
          spec:
            interval: 1m
            url: https://github.com/stefanprodan/podinfo
            ref:
              branch: master
          ```

      - Apply the created conf

          ```shell
          kubectl apply -f gitrepository-sample.yaml
          gitrepository.source.toolkit.fluxcd.io/gitrepository-sample created
          ```

      - Check the git-repository was fetched correctly

          ```shell
          kubectl get GitRepository
          NAME                   URL                                       READY   STATUS                                                              AGE
          gitrepository-sample   https://github.com/stefanprodan/podinfo   True    Fetched revision: master/132f4e719209eb10b9485302f8593fc0e680f4fc   4s
          ```

    You can find more examples checking out the samples folder on [fluxcd/source-controller/samples](https://github.com/fluxcd/source-controller/tree/main/config/samples)

## Documentation

For documentation specific to fluxcd-source-controller, check out the main repository
[fluxcd/source-controller](https://github.com/fluxcd/source-controller).
