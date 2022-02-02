# Changing Editions for Existing Binaries

If you have an existing set of `tanzu` binaries on your system, you can change between Tanzu Community Edition (TCE) and Tanzu Kubernetes Grid (TKG) modes without uninstalling anything.

To do so, run the following command:

```shell
tanzu config set cli.edition (tce or tkg)
```

Doing this will have the following implications:

* All binaries, including plugins, on your local system will remain the same. However, some plugins may not be exposed depending on the edition you are using.
* Using the `tce` edition will _not_ install a user package repository into the management clusters; this will have to be added.
  Using the `tkg` edition _will_ install a user package repository into the management clusters, which is the default package repository for TKG users.
* Packages for the `tkg` edition will include VMware-built and validated components, some of which may not be open source.
