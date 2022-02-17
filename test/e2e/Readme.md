# Build Deploy and E2E test Automation Framework

This is the folder in which we can invoke e2e tests for tce addons packages.

## Prerequisites

1. Dependent packages or tools (like ytt) would be handled in respective addons package suites.

## Supported flags

1. --kubeconfig : Use to set the kube config path
2. --kube-context : Use to set cluster context
3. --packages : Use to filter out the addon package tests to run
4. --version: Provide respective version for packages mentioned.
5. --provider : The environment in which test should run (eg: docker, aws, vsphere etc)
6. --cluster-type : Set type of cluster (eg: management)
7. --create-cluster : Provide true if cluster has to provision. Provide false if the cluster is already exist. If true automation creates cluster based on the provider and cluster type.
8. --tce-version : Provide tce release version to install(eg: "v0.7.0"). If not provided then build it from source code.
9. --guest-cluster-name: Provide cluster name (for workload cluster incase of managed)
10. --management-cluster-name: Provide cluster name for management cluster.
11. --cluster-plan: Provide Cluster Plan (eg: dev, prod etc). By default it will be set to "dev".
12. --cleanup-cluster: Provide true for tearing down the cluster.

## How to run framework to install TCE release from github page, provision cluster and test

   ```ginkgo -v -- --kubeconfig=$KUBECONFIG --packages="external-dns" --version="0.8.0" --provider="docker" --cluster-type="management" --guest-cluster-name="tce-mycluster" --create-cluster --tce-version="v0.7.0"```

## How to create management cluster on docker

   ```ginkgo -v -- --kubeconfig=$KUBECONFIG --packages="external-dns" --version="0.8.0" --provider="docker" --cluster-type="managed" --tce-version="v0.7.0" --management-cluster-name="tce-management-cluster" --create-cluster --guest-cluster-name="tce-workload-cluster"```

## How to run framework to install TCE by building release from source code, provision cluster and test

   ```ginkgo -v -- --kubeconfig=$KUBECONFIG --packages="all" --provider="docker" --cluster-type="management" --guest-cluster-name="tce-mycluster" --create-cluster```

## How to run addons package test if cluster is already available

1. to run individual packages test

   ```ginkgo -v -- --packages="calico,contour" --version="3.11.3,1.17.1" --guest-cluster-name="tce-mycluster"```
    or
    ```ginkgo -v -- --packages="contour" --version="1.17.1" --guest-cluster-name="tce-mycluster"```

2. to run all addons package tests

    ```ginkgo -v -- --kubeconfig=$KUBECONFIG --packages="all" --guest-cluster-name="tce-mycluster"```

3. how to pass cluster context

    ```ginkgo -v -- --kubeconfig=$KUBECONFIG --packages="antrea" --version="0.11.3" --kube-context="cluster-dev-test-admin@cluster-dev-test"```

### Notes

1. If all packages tests need to run then --version is not required along with --packages="all"
2. By default --packages="all" will run the package testing for the latest version of respective packages. One can modify the required package version in addons_config.yaml if they wanted to test lower version of supported packages.
3. Default cluster plan is dev.
4. Provide release version of TCE in --tce-version to fetch it from github and install. If value is not provided then automation code checks if the tanzu cli already available in the environment, if not it will create a release from source code and installs tce.
5. In case of any failures, created cluster will be tear down.
6. Cluster tearing down is optional. One can delete cluster passing --cleanup-cluster.
