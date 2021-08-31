## Create vSphere Clusters

This section describes setting up standalone clusters for
vSphere. These clusters are not managed by a management cluster.

1. Download the machine image that matches the version of the Kubernetes you plan on deploying.

    At this time, we cannot guarantee the plugin versions that will be
    used for cluster management. While running the installer interface to bootstrap a
    cluster, you are required to add an `ova` to your vSphere environment.

    The official OVA publishing location is still to be determined. In the meantime, to get access to the necessary OVAs for the current build, please ask on the `#tanzu-community-edition` Slack channel.

    Please note, validation work so far has focused on the **Photon** based
    images.

2. In vCenter, right-click on your datacenter and import the OVF template.

3. After importing, right-click and covert to a template.

4. Initialize the Tanzu Community Edition installer interface.

    ```sh
    tanzu standalone-cluster create --ui
    ```

5.  Complete the configuration steps in the installer interface and create the standalone cluster. The following configuration settings are recommended:

      * Set all instance profile to large or larger. In our testing, we found resource constraints caused bootstrapping issues. Choosing a large profile will give a better chance for successful bootstrapping.
      * Set your control plane IP. The control plane IP is a virtual IP that fronts the Kubernetes API server. You **must** set an IP that is routable and won't be taken by another system (e.g. DHCP).
      * Disable **Enable Identity Management Settings**. You can disable identity management for proof-of-concept/development deployments, but it is strongly recommended to implement identity management in production deployments.

6. Set your kubectl context to the cluster.

    ```sh
    kubectl config use-context <STANDALONE-CLUSTER-NAME>-admin@<STANDALONE-CLUSTER-NAME>
    ```
    Where `<STANDALONE-CLUSTER-NAME>` is the name of the standalone cluster that you specified or if you didn't specify a name, it's the randomly generated name.

7. Validate you can access the cluster's API server.

    ```sh
    kubectl get nodes

    NAME         STATUS   ROLES                  AGE    VERSION
    10-0-1-133   Ready    <none>                 123m   v1.20.1+vmware.2
    10-0-1-76    Ready    control-plane,master   125m   v1.20.1+vmware.2
    ```
