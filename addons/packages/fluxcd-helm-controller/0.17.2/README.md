# Helm Controller

The helm-controller is a Kubernetes operator, allowing one to declaratively manage Helm chart releases. It is part of a composable GitOps toolkit and depends on source-controller to acquire the Helm charts from Helm repositories.

## Configuration

The Helm controller package has following configurable properties.

| Value           | Required/Optional |                           Description                                              |
|-----------------|-------------------|------------------------------------------------------------------------------------|
| `namespace`     | Optional          | Sets namespace for k8s objects where resources of helm-controller with be created. |
| `limits_cpu`    | Optional          | Sets maximum usage of cpu for helm-controller deployment.                          |
| `limits_memory` | Optional          | Sets maximum usage of memory for helm-controller deployment.                       |

```yaml
namespace: helm-system
limits_cpu: 1000m
limits_memory: 1Gi
```

## Installation

To install FluxCD helm-controller from the Tanzu Application Platform package repository:

1. Prerequisites:-

   Install fluxcd-source-controller package:-

   Helm controller package requires source-controller package to be installed for acquiring the Kubernetes manifests from the sources.

    ```shell
    tanzu package install fluxcd-source-controller -p fluxcd-source-controller.community.tanzu.vmware.com -v VERSION-NUMBER
    ```

   here `VERSION-NUMBER` is the version of the package

   For example:

    ```shell
    tanzu package install fluxcd-source-controller -p fluxcd-source-controller.community.tanzu.vmware.com -v 0.21.2
    \ Installing package 'fluxcd-source-controller.community.tanzu.vmware.com'
    | Getting package metadata for 'fluxcd-source-controller.community.tanzu.vmware.com'
    | Creating service account 'fluxcd-source-controller-default-sa'
    | Creating cluster admin role 'fluxcd-source-controller-default-cluster-role'
    | Creating cluster role binding 'fluxcd-source-controller-default-cluster-rolebinding'
    | Creating package resource
    / Waiting for 'PackageInstall' reconciliation for 'fluxcd-source-controller'
    \ 'PackageInstall' resource install status: Reconciling

    Added installed package 'fluxcd-source-controller'
    ```

2. List version information for the package by running:

    ```shell
    tanzu package available list helm-controller.fluxcd.community.tanzu.vmware.com
    ```

   For example:

    ```shell
    $ tanzu package available list helm-controller.fluxcd.community.tanzu.vmware.com
    / Retrieving package versions for helm-controller.fluxcd.community.tanzu.vmware.com...  
    NAME                                                    VERSION  RELEASED-AT  
    helm-controller.fluxcd.community.tanzu.vmware.com        0.17.2   2022-02-23 16:44:08 +0530 IST  
    ```

3. Configure helm-controller package:

   User can optionally provide the configuration parameters with --values-file flag while installing the package. Download the values.yaml file from [values.yaml](https://github.com/vmware-tanzu/package-for-helm-controller/blob/main/0.17.2/bundle/config/values.yaml).

    ```shell
    cat values.yaml  
    #@data/values
    ---
    namespace: helm-system
    limits_cpu: 1000m
    limits_memory: 1Gi
    ```

4. Install the package by running:

    ```shell
    tanzu package install fluxcd-helm-controller -p helm-controller.fluxcd.community.tanzu.vmware.com -v VERSION-NUMBER
    ```

   here: `VERSION-NUMBER` is the version of the package

   For example:

    ```shell
    $ tanzu package install fluxcd-helm-controller -p helm-controller.fluxcd.community.tanzu.vmware.com -v 0.17.2
    / Installing package 'helm-controller.fluxcd.community.tanzu.vmware.com'
    | Getting package metadata for 'helm-controller.fluxcd.community.tanzu.vmware.com'
    | Creating service account 'fluxcd-helm-controller-default-sa'
    | Creating cluster admin role 'fluxcd-helm-controller-default-cluster-role'
    | Creating cluster role binding 'fluxcd-helm-controller-default-cluster-rolebinding'
    | Creating package resource
    / Waiting for 'PackageInstall' reconciliation for 'fluxcd-helm-controller'
    \ 'PackageInstall' resource install status: Reconciling

    Added installed package 'fluxcd-helm-controller'
    ```

5. Verify the package install by running:

    ```shell
    tanzu package installed get fluxcd-helm-controller
    ```

   For example:

    ```shell
    - Retrieving installation details for fluxcd-helm-controller...
    NAME:                    fluxcd-helm-controller
    PACKAGE-NAME:            helm-controller.fluxcd.community.tanzu.vmware.com
    PACKAGE-VERSION:         0.17.2
    STATUS:                  Reconcile succeeded
    CONDITIONS:              [{ReconcileSucceeded True  }]
    ```

   Verify that `STATUS` is `Reconcile succeeded`

    ```shell
    kubectl get pods -n helm-system
    ```

   For example:

    ```shell
    $ kubectl get pods -n helm-system
    NAME                               READY   STATUS    RESTARTS   AGE
    helm-controller-7748954977-cgv59   1/1     Running   2          157m
    ```

   Verify that `STATUS` is `Running`

## Try helm-controller

1. Verify all the objects are installed:

   This package would create a new namespace where all the elements of fluxcd will be hosted called `helm-system` you can verify the main components of `helm-controller` were installed by running:

    ```shell
    $ kubectl get all -n helm-system
    NAME                                   READY   STATUS    RESTARTS   AGE
    pod/helm-controller-7748954977-cgv59   1/1     Running   0          3m28s

    NAME                              READY   UP-TO-DATE   AVAILABLE   AGE
    deployment.apps/helm-controller   1/1     1            1           3m28s

    NAME                                         DESIRED   CURRENT   READY   AGE
    replicaset.apps/helm-controller-7748954977   1         1         1       3m28s
    ```

   you should get something really similar!

2. Verify all the CRD were installed correctly:

   The way you would communicate with `helm-controller` would be through its CRDs, this will be your main action point.

   In order to check all the CRDs were installed you can run:

    ```shell
     $ kubectl get crds | grep ".fluxcd.io"
     buckets.source.toolkit.fluxcd.io                                               2022-04-18T07:29:09Z
     gitrepositories.source.toolkit.fluxcd.io                                       2022-04-18T07:29:09Z
     helmcharts.source.toolkit.fluxcd.io                                            2022-04-18T07:29:09Z
     helmreleases.helm.toolkit.fluxcd.io                                            2022-04-18T07:40:57Z
     helmrepositories.source.toolkit.fluxcd.io                                      2022-04-18T07:29:09Z
    ```

   Apart from `helmreleases.helm.toolkit.fluxcd.io` all the other CRDs are installed as part of the [fluxcd-source-controller](#installation) (which is a pre-requisite) to get helm-controller work as expected.

3. Try one quick example yourself, so you can check everything is working as expected

   - Let's consume a `HelmRelease` and `HelmRepository` object.

   - Create the following `sample.yaml` file:

     Note: `flux-system` namespace is created as part of the pre-requisite - installing fluxcd-source-controller

    ```yaml
    ---
    apiVersion: helm.toolkit.fluxcd.io/v2beta1
    kind: HelmRelease
    metadata:
      name: podinfo
      namespace: flux-system
    spec:
      interval: 5m
      chart:
        spec:
          chart: podinfo
          version: '4.0.x'
          sourceRef:
            kind: HelmRepository
            name: podinfo
            namespace: flux-system
          interval: 1m
      values:
        replicaCount: 2
    ---
    apiVersion: source.toolkit.fluxcd.io/v1beta1
    kind: HelmRepository
    metadata:
      name: podinfo
      namespace: flux-system
    spec:
      interval: 1m
      url: https://stefanprodan.github.io/podinfo
    ```

   - Apply the created conf

    ```shell
    $ kubectl apply -f sample.yaml
    helmrelease.helm.toolkit.fluxcd.io/podinfo created
    helmrepository.source.toolkit.fluxcd.io/podinfo created
    ```

   - Check HelmChart was pulled successfully
   - Check the HelmRepository was fetched successfully
   - Check the HelmRelease is reconciled successfully

    ```shell
    $ kubectl -n flux-system get helmreleases.helm.toolkit.fluxcd.io
    NAME      READY   STATUS                             AGE
    podinfo   True    Release reconciliation succeeded   27s

    $ kubectl -n flux-system get helmcharts.source.toolkit.fluxcd.io
    NAME                  CHART     VERSION   SOURCE KIND      SOURCE NAME   READY   STATUS                                         AGE
    flux-system-podinfo    podinfo   4.0.x     HelmRepository   podinfo       True    Pulled 'podinfo' chart with version '4.0.6'.   42s

    $ kubectl -n flux-system get helmrepositories.source.toolkit.fluxcd.io
    NAME      URL                                      READY   STATUS                                                                               AGE
    podinfo   https://stefanprodan.github.io/podinfo   True    Fetched revision: debd9e30ecc98721b7aa3404eb64db4032786e0572be4629359aa7b6cf0f5fe7   57s

    $ kubectl get all -n flux-system
    NAME                                     READY   STATUS    RESTARTS   AGE
    pod/podinfo-7ddf7dd478-5d7dw             1/1     Running   0          12m
    pod/podinfo-7ddf7dd478-rdk6k             1/1     Running   0          12m
    pod/source-controller-66ccb7ff5c-j4q4l   1/1     Running   2          168m

    NAME                        TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)             AGE
    service/podinfo             ClusterIP   10.96.42.37    <none>        9898/TCP,9999/TCP   12m
    service/source-controller   ClusterIP   10.96.133.58   <none>        80/TCP              168m

    NAME                                READY   UP-TO-DATE   AVAILABLE   AGE
    deployment.apps/podinfo             2/2     2            2           12m
    deployment.apps/source-controller   1/1     1            1           168m

    NAME                                           DESIRED   CURRENT   READY   AGE
    replicaset.apps/podinfo-7ddf7dd478             2         2         2       12m
    replicaset.apps/source-controller-66ccb7ff5c   1         1         1       168m
    ```

   You can find more examples checking out the samples folder on [fluxcd/helm-controller/samples](https://github.com/fluxcd/helm-controller/tree/main/config/samples)

## Documentation

For documentation specific to helm-controller, check out the main repository
[fluxcd/helm-controller](https://github.com/fluxcd/helm-controller).
