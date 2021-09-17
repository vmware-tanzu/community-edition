## Clean-up

1. Run the `delete` command.

    ```sh
    tanzu standalone-cluster delete <STANDALONE-CLUSTER-NAME>
    ```

    > This may take several minutes to complete!

1. _Note:_ If you configured a proxy, you may need to provide the following environment variables `TKG_NO_PROXY`, `TKG_HTTP_PROXY`, `TKG_HTTPS_PROXY`

    ```sh
    TKG_HTTP_PROXY="127.0.0.1" tanzu standalone-cluster delete <STANDALONE-CLUSTER-NAME>
    ```
