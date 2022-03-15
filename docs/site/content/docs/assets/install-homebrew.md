1. Install using [Homebrew](https://brew.sh).

    ```sh
    brew install vmware-tanzu/tanzu/tanzu-community-edition
    ```

1. Run the configure command displayed after Homebrew completes installation.

    ```sh
    {HOMEBREW-INSTALL-LOCATION}/configure-tce.sh
    ```

    > This puts all the Tanzu plugins in the correct location. The first time
    > you run the `tanzu` command the installed plugins and plugin repositories
    > are initialized. This action might take a minute.
