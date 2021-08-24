# Upgrade vSphere Deployments in an Internet-Restricted Environment

If you deployed the previous version of Tanzu Kubernetes Grid in an Internet-restricted environment, do the following steps on a machine with an Internet connection.

1. Download and install the new version of the Tanzu CLI. See [Download and Install the Tanzu CLI and Other Tools](index.md#download-cli).

1. Perform the steps in [Prepare to Upgrade Clusters on vSphere](index.md#vsphere) to deploy the new base OS image OVA files.

1. Perform the steps in [Deploying Tanzu Kubernetes Grid in an Internet-Restricted Environment](../mgmt-clusters/airgapped-environments.md#procedure) to run the `gen-publish-images.sh` and `publish-images.sh` scripts.

If you still have the `publish-images.sh` script from when you deployed the previous version of Tanzu Kubernetes Grid, you must regenerate it by running `gen-publish-images.sh` before you run `publish-images.sh`.

Running `gen-publish-images.sh` updates `publish-images.sh` so that it pulls the correct versions of the components for the new version of Tanzu Kubernetes Grid and pushes them into your local private Docker registry.

The `gen-publish-images.sh` script obtains the correct versions of the components from the YAML files that are created in the `~/.tanzu/tkg/bom` folder when you first run a `tanzu` CLI command with a new version of Tanzu Kubernetes Grid.
