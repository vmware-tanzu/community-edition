# test

## This directory is intended to hold the orchestration, framework and utilities needed to test Tanzu Community Edition

### contents

1. Deployment Automation Test:

    - [aws](https://github.com/vmware-tanzu/community-edition/tree/main/test/aws), [azure](https://github.com/vmware-tanzu/community-edition/tree/main/test/azure), [docker](https://github.com/vmware-tanzu/community-edition/tree/main/test/docker) and [vsphere](https://github.com/vmware-tanzu/community-edition/tree/main/test/vsphere) directories are home to deployment automation of Tanzu Community Edition on the respective providers.

    - Usage:
        - Environment variables:
            - Make sure that the environment variables needed for different providers are exported.
            - Here are the environment variables needed for [aws](https://github.com/vmware-tanzu/community-edition/blob/f2e47713ec0da0e9f68649f6e8f678c77bd26c24/test/aws/deploy-tce-managed.sh#L11), [azure](https://github.com/vmware-tanzu/community-edition/blob/f2e47713ec0da0e9f68649f6e8f678c77bd26c24/test/azure/deploy-management-and-workload-cluster.sh#L11) and [vsphere](https://github.com/vmware-tanzu/community-edition/blob/f2e47713ec0da0e9f68649f6e8f678c77bd26c24/test/vsphere/run-tce-vsphere-standalone-cluster.sh#L10).

        - Run the following commands from the root of the repository to trigger deployment automation:
            - Management Cluster deployment automation on:
                - aws:

                    ```shell
                    make aws-management-and-workload-cluster-e2e-test
                    ```

                - azure:

                    ```shell
                    make azure-management-and-workload-cluster-e2e-test
                    ```

                - docker:

                    ```shell
                    make docker-management-and-cluster-e2e-test:
                    ```

            - Standalone Cluster deployment automation on:
                - aws:
  
                    ```shell
                    make aws-standalone-cluster-e2e-test
                    ```

                - azure:

                    ```shell
                    make azure-standalone-cluster-e2e-test
                    ```

                - docker:

                    ```shell
                    make docker-standalone-cluster-e2e-test
                    ```

                - vSphere:

                    ```shell
                    make tce-vsphere-standalone-cluster-e2e-test
                    ```

1. E2E Test Framework:

    - Framework to build and test Tanzu Community Edition packages lives in [e2e](https://github.com/vmware-tanzu/community-edition/tree/main/test/e2e) directory.

1. Test util:

    - [gatekeeper](https://github.com/vmware-tanzu/community-edition/tree/main/test/gatekeeper) directory contains the e2e test for Gatekeeper package along with the required configurations.

    - [utils](https://github.com/vmware-tanzu/community-edition/tree/main/test/util) directory contains the various utilities to test Tanzu Community Edition.
