# Deploy Clusters with Different Kubernetes Versions

Tanzu Community Edition manages Kubernetes versions with Tanzu Kubernetes
release (TKr) objects.

The following steps describe how to view available TKr versions and use them
when creating a workload cluster.

1. List the available TKrs.

    ```sh
    tanzu kubernetes-release
    ```

    Ouptut:

    ```txt
    NAME                                  VERSION                             COMPATIBLE  ACTIVE  UPDATES AVAILABLE
    v1.20.14---vmware.1-tkg.4-tf-v0.11.2  v1.20.14+vmware.1-tkg.4-tf-v0.11.2  True        True    True
    v1.21.8---vmware.1-tkg.4-tf-v0.11.2   v1.21.8+vmware.1-tkg.4-tf-v0.11.2   True        True    True
    v1.22.5---vmware.1-tkg.3-tf-v0.11.2   v1.22.5+vmware.1-tkg.3-tf-v0.11.2   True        True    False
    ```

1. When creating a cluster, specify the TKr using the `--tkr` flag.

    ```sh
    tanzu cluster create --tkr v1.21.8---vmware.1-tkg.4-tf-v0.11.2
    ```

    > This will deploy the TKr with Kubernetes v1.21.8. This cluster will be
    > upgradable to the TKr featuring Kubernetes v1.22.5.
