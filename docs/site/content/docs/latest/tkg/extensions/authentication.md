# Implementing User Authentication

This release of Tanzu Kubernetes Grid introduces user authentication with Pinniped, that runs automatically in management clusters if you enable identity management during deployment. Previously, user authentication was implemented by deploying the Dex and Gangway extensions. The extensions for Dex and Gangway are included in this release for historical reasons, but are deprecated. For new deployments, you must enable Pinniped in your management clusters. Do not to use the Dex and Gangway extensions.

- For information about identity management with Pinniped, see [Enabling Identity Management in Tanzu Kubernetes Grid](../mgmt-clusters/enabling-id-mgmt.md).
- For information about migrating existing Dex and Gangway deployments to Pinniped, see [Register Core Add-ons](../upgrade-tkg/addons.md).
