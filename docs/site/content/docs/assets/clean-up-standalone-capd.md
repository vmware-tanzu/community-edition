## Clean-up

Standalone-cluster deletion is currently un-implemented. Until it is, do the
following steps to remove resources created by Docker.

1. Kill all containers.

    ⚠️ : The following command will kill all containers on your system.

    ```sh
    docker kill $(docker ps -q)
    ```

1. Prune your machine of left-over volumes.

    ```sh
    docker system prune --volumes
    ```
