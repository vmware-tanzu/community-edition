# kustomize-controller

The kustomize-controller is a Kubernetes operator, specialized in running continuous delivery pipelines for
infrastructure and workloads defined with Kubernetes manifests and assembled with Kustomize.

The kustomize-controller is part of a composable [GitOps](https://fluxcd.io/docs/components/) toolkit and depends on [source-controller](https://github.com/fluxcd/source-controller) to acquire the Kubernetes manifests from Git repositories and S3 compatible storage buckets.

The kustomize-controller implements the [kustomize.toolkit.fluxcd.io](https://github.com/fluxcd/kustomize-controller/tree/main/docs/spec/v1beta2) API and is one of the component of GitOps toolkit.

## Configuration  
  
The Kustomize controller package has following configurable properties.

| Value           | Required/Optional | Description                     |
|-----------------|-------------------|---------------------------------|
| `namespace`     | Optional          | Sets namespace for k8s objects. |
| `limits_cpu`    | Optional          | Sets maximum usage of cpu.      |
| `limits_memory` | Optional          | Sets maximum usage of memory.   |

## Installation

To install FluxCD kustomize-controller from the Tanzu Application Platform package repository:

1. Prerequisite:-
  
   - Install source-controller package:-

     Kustomize-controller package requires source-controller package to be installed for acquiring the Kubernetes manifests from the sources.
  
      ```shell
      tanzu package install fluxcd-source-controller -p fluxcd-source-controller.community.tanzu.vmware.com -v VERSION-NUMBER
      ```
  
      Where:

      - VERSION-NUMBER is the version of the package listed in step 1.  
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
      \ 'PackageInstall' resource install status: Reconciling
  
       Added installed package 'fluxcd-source-controller'

      ```

2. List version information for the package by running:

    ```shell
    tanzu package available list kustomize-controller.fluxcd.community.tanzu.vmware.com
    ```

    For example:

    ```shell
    $ tanzu package available list kustomize-controller.fluxcd.community.tanzu.vmware.com
      / Retrieving package versions for kustomize-controller.fluxcd.community.tanzu.vmware.com...  
      NAME                                                    VERSION  RELEASED-AT  
      kustomize-controller.fluxcd.community.tanzu.vmware.com  0.21.1   2022-02-23 16:44:08 +0530 IST  
      kustomize-controller.fluxcd.community.tanzu.vmware.com  0.24.4   2022-02-23 16:44:08 +0530 IST  
    ```
  
3. Configure kustomize-controller package:  
  
   User can optionally provide the configuration parameters with --values-file flag while installing the package. Download the values.yaml file from [values.yaml](https://github.com/vmware-tanzu/package-for-kustomize-controller/blob/main/0.24.4/bundle/config/values.yaml).
  
   ```shell
   namespace: kustomize-system
   limits_cpu: 1000m
   limits_memory: 1Gi
   ```

4. Install the package by running:

   ```shell
   tanzu package install fluxcd-kustomize-controller -p kustomize-controller.fluxcd.community.tanzu.vmware.com -v VERSION-NUMBER --values-file VALUES-FILE
   ```

   Where:

   - `VERSION-NUMBER` is the version of the package listed in step 1.
   - `VALUES-FILE` is the configuration file for the package listed in step1.

   For example:

   ```shell
   tanzu package install fluxcd-kustomize-controller -p kustomize-controller.fluxcd.community.tanzu.vmware.com -v 0.24.4 --values-file values.yaml  
   | Installing package 'kustomize-controller.fluxcd.community.tanzu.vmware.com'  
   | Getting package metadata for 'kustomize-controller.fluxcd.community.tanzu.vmware.com'  
   | Creating service account 'fluxcd-kustomize-controller-default-sa'  
   | Creating cluster admin role 'fluxcd-kustomize-controller-default-cluster-role'  
   | Creating cluster role binding 'fluxcd-kustomize-controller-default-cluster-rolebinding'  
   | Creating secret 'fluxcd-kustomize-controller-default-values'  
   | Creating package resource  
   / Waiting for 'PackageInstall' reconciliation for 'fluxcd-kustomize-controller'  
   \ 'PackageInstall' resource install status: Reconciling  

   Added installed package 'fluxcd-kustomize-controller'
   ```

5. Verify the package install by running:

   ```shell
   tanzu package installed get fluxcd-kustomize-controller
   ```

   For example:

   ```shell
   - Retrieving installation details for fluxcd-kustomize-controller...
     NAME:                    fluxcd-kustomize-controller
     PACKAGE-NAME:            kustomize-controller.fluxcd.community.tanzu.vmware.com
     PACKAGE-VERSION:         0.24.4
     STATUS:                  Reconcile succeeded
     CONDITIONS:              [{ReconcileSucceeded True  }]
   ```
  
   Verify that `STATUS` is `Reconcile succeeded`

   ```shell
   kubectl get pods -n kustomize-system
   ```

   For example:

   ```shell
   kubectl get pods -n kustomize-system
   NAME                                   READY   STATUS    RESTARTS   AGE
   kustomize-controller-fbdbbdfd8-8w9r7   1/1     Running   0          2m47s
   ```

   Verify that `STATUS` is `Running`

## Try kustomize-controller

1. Verify all the objects are installed:

    This package would create a new namespace where all the elements of fluxcd will be hosted called `kustomize-system`

    you can verify the main components of `kustomize-controller` were installed by running:

    ```shell
    kubectl get all -n kustomize-system  
    NAME                                       READY   STATUS    RESTARTS   AGE
    pod/kustomize-controller-fbdbbdfd8-8w9r7   1/1     Running   0          3m54s

    NAME                                   READY   UP-TO-DATE   AVAILABLE   AGE
    deployment.apps/kustomize-controller   1/1     1            1           3m54s

    NAME                                             DESIRED   CURRENT   READY   AGE
    replicaset.apps/kustomize-controller-fbdbbdfd8   1         1         1       3m54s
    ```

    you should get something really similar!

2. Verify all the CRD were installed correctly:

    The way you would communicate with `kustomize-controller` would be through its CRDs, this will be your main action point.

    In order to check all the CRDs were installed you can run:

    ```shell
    kubectl get crds -n flux-system | grep ".fluxcd.io"  
    buckets.source.toolkit.fluxcd.io                                               2022-04-08T11:28:55Z
    gitrepositories.source.toolkit.fluxcd.io                                       2022-04-08T11:28:55Z
    helmcharts.source.toolkit.fluxcd.io                                            2022-04-08T11:28:55Z
    helmrepositories.source.toolkit.fluxcd.io                                      2022-04-08T11:28:55Z
    kustomizations.kustomize.toolkit.fluxcd.io                                     2022-04-10T16:12:13Z
    ```

3. Try one simple example:

    Try one quick example yourself, so you can check everything is working as expected

    - Let's consume a `GitRepository` and `Kustomization` object,

      - Create the following `sample.yaml` file:

       ```yaml  
       apiVersion: source.toolkit.fluxcd.io/v1beta1  
       kind: GitRepository
       metadata:
         name: podinfo
         namespace: flux-system
       spec:
         interval: 1m
         url: https://github.com/stefanprodan/podinfo
       ---
       apiVersion: kustomize.toolkit.fluxcd.io/v1beta2
       kind: Kustomization
       metadata:
         name: podinfo
         namespace: flux-system
       spec:
         interval: 5m0s
         path: ./kustomize
         prune: true
         sourceRef:
           kind: GitRepository
           name: podinfo
         targetNamespace: flux-system
       ```

      - Apply the created conf

       ```shell
       kubectl apply -f sample.yaml
       gitrepository.source.toolkit.fluxcd.io/podinfo created
       kustomization.kustomize.toolkit.fluxcd.io/podinfo created
       ```

      - Check the git-repository was fetched correctly
      - Check the Kustomization is created

       ```shell

       kubectl get GitRepository -A
       NAMESPACE     NAME      URL                                       READY   STATUS                          6m12s
       flux-system   podinfo   https://github.com/stefanprodan/podinfo   True    Fetched revision: master/b8910253653d375bce5b71178518fa50a9cab23f   3s
  
       kubectl get Kustomization -A
       NAMESPACE     NAME      READY   STATUS                                                              AGE
       flux-system   podinfo   True    Applied revision: master/b8910253653d375bce5b71178518fa50a9cab23f   8s
  
       k get all -n flux-system
       NAME                                    READY   STATUS    RESTARTS   AGE
       pod/podinfo-694f589bf6-5p8bc            1/1     Running   0          7m59s
       pod/podinfo-694f589bf6-tll6k            1/1     Running   0          8m14s
  
       NAME                        TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)             AGE
       service/podinfo             ClusterIP   10.96.169.222   <none>        9898/TCP,9999/TCP   8m14s
  
       NAME                                READY   UP-TO-DATE   AVAILABLE   AGE
       deployment.apps/podinfo             2/2     2            2           8m14s
  
       NAME                                          DESIRED   CURRENT   READY   AGE
       replicaset.apps/podinfo-694f589bf6            2         2         2       8m14s
  
       NAME                                          REFERENCE            TARGETS         MINPODS   MAXPODS   REPLICAS   AGE
       horizontalpodautoscaler.autoscaling/podinfo   Deployment/podinfo   <unknown>/99%   2         4         2          8m14s

       ```

    You can find more examples checking out the samples folder on [fluxcd/kustomize-controller/samples](https://github.com/fluxcd/kustomize-controller/tree/main/config/samples)

## Documentation

For documentation specific to kustomize-controller, check out the main repository
[fluxcd/kustomize-controller](https://github.com/fluxcd/kustomize-controller).
