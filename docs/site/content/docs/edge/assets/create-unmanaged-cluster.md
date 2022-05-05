1. Create a cluster named `beepboop`.

    ```sh
    tanzu unmanaged-cluster create beepboop
    ```

    > **Tip**: You can use the alias `uc`, instead of `unmanaged-cluster`, in
    > these commands.

1. Wait for the cluster to initialize.

    ```txt
    📁 Created cluster directory

    🔧 Resolving Tanzu Kubernetes Release (TKR)
       projects.registry.vmware.com/tce/tkr:v1.21.5
       TKR exists at /home/josh/.config/tanzu/tkg/unmanaged-cluster/bom/projects.registry.vmware.com_tce_tkr_v1.21.5
       Rendered Config: /home/josh/.config/tanzu/tkg/unmanaged-cluster/beepboop/config.yaml
       Bootstrap Logs: /home/josh/.config/tanzu/tkg/unmanaged-cluster/beepboop/bootstrap.log

    🔧 Processing Tanzu Kubernetes Release

    🎨 Selected base image
       projects.registry.vmware.com/tce/kind/node:v1.21.5

    📦 Selected core package repository
       projects.registry.vmware.com/tkg/packages/core/repo:v1.21.2_vmware.1-tkg.1

    📦 Selected additional package repositories
       projects.registry.vmware.com/tce/main:0.9.1

    📦 Selected kapp-controller image bundle
       projects.registry.vmware.com/tkg/packages/core/kapp-controller:v0.23.0_vmware.1-tkg.1

    🚀 Creating cluster beepboop
       Base image downloaded
       Cluster created
       To troubleshoot, use:
       kubectl ${COMMAND} --kubeconfig /home/josh/.config/tanzu/tkg/unmanaged-cluster/beepboop/kube.conf

    📧 Installing kapp-controller
       kapp-controller status: Running

    📧 Installing package repositories
       Core package repo status: Reconcile succeeded

    🌐 Installing CNI
       antrea.tanzu.vmware.com:0.13.3+vmware.1-tkg.1

    ✅ Cluster created

    🎮 kubectl context set to beepboop

    View available packages:
       tanzu package available list
    View running pods:
       kubectl get po -A
    Delete this cluster:
       tanzu unmanaged-cluster delete beepboop
    ```

    > A container image larger than 1GB is used for running clusters. This may
    > cause your first cluster to take significantly longer to bootstrap than
    > subsequent clusters.

1. If you have [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
   installed, you can now use it to interact with the
   cluster.
