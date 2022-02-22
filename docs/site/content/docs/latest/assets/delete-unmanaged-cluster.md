1. List the available clusters.

    ```sh
    tanzu unmanaged-cluster list
    ```

    > **Tip**: You can use the alias `ls`, instead of `list`, in
    > the above command.

1. Review the available clusters.

    ```txt
    NAME      PROVIDER
    beepboop  kind
    ```

1. Delete the cluster.

    ```sh
    tanzu unmanaged-cluster delete beepboop
    ```

    > **Tip**: You can use the alias `rm`, instead of `delete`, in
    > the above command.
