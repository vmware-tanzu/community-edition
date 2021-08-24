# Tanzu Kubernetes Cluster Secrets

This topic explains how to configure and manage secrets used by management and Tanzu Kubernetes (workload) clusters in Tanzu Kubernetes Grid, including:

* Credentials that clusters use to access cloud infrastructure APIs and resources.
* Certificate Authority (CA) certificates that clusters use to access private container registries.

## <a id="mgmt-creds-update"></a> Update Management and Workload Cluster Credentials (vSphere)

To update the vSphere credentials used by the current management cluster and by all of its workload clusters, use the `tanzu management-cluster credentials update --cascading` command:

1. Run `tanzu login` to log in to the management cluster that you are updating.

1. Run `tanzu management-cluster credentials update`

  You can pass values to the following command options, or let the CLI prompt you for them:
  
  - `--vsphere-user`: Name for the vSphere account.
  - `--vsphere-password`: Password the vSphere account.

### <a id="mgmt-creds-only-update"></a> Update Management Cluster Credentials Only

To update a management cluster's vSphere credentials without also updating them for its workload clusters, use the `tanzu management-cluster credentials update` command as above, but without the `--cascading` option.

## <a id="creds-update"></a> Update Workload Cluster Credentials (vSphere)

To update the credentials that a single workload cluster uses to access vSphere, use the `tanzu cluster credentials update` command:

1. Run `tanzu login` to log in to the management cluster that created the workload cluster that you are updating.

1. Run `tanzu cluster credentials update CLUSTER_NAME`

  You can pass values to the following command options, or let the CLI prompt you for them:
  
  - `--namespace`: The namespace of the cluster you are updating credentials for, such as `default`.
  - `--vsphere-user`: Name for the vSphere account.
  - `--vsphere-password`: Password the vSphere account.

You can also use `tanzu management-cluster credentials update --cascading` to update vSphere credentials for a management cluster and all of the workload clusters it manages.

## <a id="custom-ca"></a> Trust Custom CA Certificates on Cluster Nodes

You can add custom CA certificates in Tanzu Kubernetes cluster nodes by using a `ytt` overlay file to enable the cluster nodes to pull images from a container registry that uses self signed certificates.

The overlay code below adds custom CA certificates to all nodes in a new cluster, so that `containerd` and other tools trust the certificate.
The code works on all cloud infrastructure providers, for clusters based on Photon or Ubuntu VM image templates.

For overlays that customize clusters and create a new cluster plan, see [`ytt` Overlays](../tanzu-k8s-clusters/config-plans.md#examples) in the _Configure Tanzu Kubernetes Plans and Clusters_ topic.

1. Choose whether to apply the custom CA to all new clusters, just the ones created on one cloud infrastructure, or ones created with a specific Cluster API provider version, such as [Cluster API Provider vSphere v0.7.4](https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/v0.7.4).

1. In your local `~/.tanzu/tkg/providers/` directory, find the `ytt` directory that covers your chosen scope.
For example, `/ytt/03_customizations/` applies to all clusters, and `/infrastructure-vsphere/ytt/` applies to all vSphere clusters.

1. In your chosen `ytt` directory, create a new `.yaml` file or augment an existing overlay file with the following code:

    ```
    #@ load("@ytt:overlay", "overlay")
    #@ load("@ytt:data", "data")

    #! This ytt overlay adds additional custom CA certificates on TKG cluster nodes, so containerd and other tools trust these CA certificates.
    #! It works when using Photon or Ubuntu as the TKG node template on all TKG infrastructure providers.

    #! Trust your custom CA certificates on all Control Plane nodes.
    #@overlay/match by=overlay.subset({"kind":"KubeadmControlPlane"})
    ---
    spec:
      kubeadmConfigSpec:
        #@overlay/match missing_ok=True
        files:
          #@overlay/append
          - content: #@ data.read("tkg-custom-ca.pem")
            owner: root:root
            permissions: "0644"
            path: /etc/ssl/certs/tkg-custom-ca.pem
        #@overlay/match missing_ok=True
        preKubeadmCommands:
          #! For Photon OS
          #@overlay/append
          - '! which rehash_ca_certificates.sh 2>/dev/null || rehash_ca_certificates.sh'
          #! For Ubuntu
          #@overlay/append
          - '! which update-ca-certificates 2>/dev/null || (mv /etc/ssl/certs/tkg-custom-ca.pem /usr/local/share/ca-certificates/tkg-custom-ca.crt && update-ca-certificates)'

    #! Trust your custom CA certificates on all worker nodes.
    #@overlay/match by=overlay.subset({"kind":"KubeadmConfigTemplate"})
    ---
    spec:
      template:
        spec:
          #@overlay/match missing_ok=True
          files:
            #@overlay/append
            - content: #@ data.read("tkg-custom-ca.pem")
              owner: root:root
              permissions: "0644"
              path: /etc/ssl/certs/tkg-custom-ca.pem
          #@overlay/match missing_ok=True
          preKubeadmCommands:
            #! For Photon OS
            #@overlay/append
            - '! which rehash_ca_certificates.sh 2>/dev/null || rehash_ca_certificates.sh'
            #! For Ubuntu
            #@overlay/append
            - '! which update-ca-certificates 2>/dev/null || (mv /etc/ssl/certs/tkg-custom-ca.pem /usr/local/share/ca-certificates/tkg-custom-ca.crt && update-ca-certificates)'
    ```

1. In the same `ytt` directory, add the Certificate Authority to a new or existing `tkg-custom-ca.pem` file.
