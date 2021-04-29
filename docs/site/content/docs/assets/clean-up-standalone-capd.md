## Clean-up

Standalone-cluster deletion is currently un-implemented. Until it is, do the
following steps to remove resources created by CAPD.

1. Kill all containers.

    ```sh
    docker kill $(docker ps -q)
    ```

1. Prune your machine of left-over volumes.

    ```sh
    docker system prune --volumes
    ```
