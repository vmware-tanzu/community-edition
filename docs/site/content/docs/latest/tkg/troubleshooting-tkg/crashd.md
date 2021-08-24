# Troubleshooting Tanzu Kubernetes Clusters with Crash Diagnostics

[Crash Diagnostics (Crashd)](https://github.com/vmware-tanzu/crash-diagnostics) is an open source project that makes it easy to diagnose problems with unstable or even unreponsive Kubernetes clusters.

Crashd uses a script file written in Starlark, a Python-like language, that interacts with your Tanzu Kubernetes clusters to collect infrastructure and cluster information. Crashd takes the output from the commands run by the script and adds the output to a `tar` file. The `tar` file is then saved locally for further analysis.

Tanzu Kubernetes Grid includes signed binaries for Crashd and a diagnostics script file for Photon OS Tanzu Kubernetes clusters.

## <a id="install"></a> Install or Upgrade the Crashd Binary

To install or upgrade `crashd`, follow the instructions below.

1. Go to [the Tanzu Kubernetes Grid downloads page](https://my.vmware.com/en/web/vmware/downloads/info/slug/infrastructure_operations_management/vmware_tanzu_kubernetes_grid/1_x), and log in with your My VMware credentials.
1. Download Crashd for your platform.

   - Linux: **crashd-linux-amd64-v0.3.2-vmware.1.tar.gz**
   - macOS: **crashd-darwin-amd64-v0.3.2-vmware.1.tar.gz**

1. Use the `tar` command to unpack the binary for your platform.

   - Linux:

        ```
        tar -xvf crashd-linux-amd64-v0.3.2-vmware.1.tar.gz
        ```

   - macOS:

        ```
        tar -xvf crashd-darwin-amd64-v0.3.2-vmware.1.tar.gz
        ```

1. The previous step creates a directory named `crashd` with the following files:

  ```
    crashd
    crashd/args
    crashd/diagnostics.crsh
    crashd/crashd-PLATFORM-amd64-v0.3.2+vmware.1
  ```

1. Move the binary into the `/usr/local/bin` folder.

   - Linux:

        ```
        mv ./crashd/crashd-linux-amd64-v0.3.2+vmware.1 /usr/local/bin/crashd
        ```

   - macOS:

        ```
        mv ./crashd/crashd-darwin-amd64-v0.3.2+vmware.1 /usr/local/bin/crashd
        ```

## <a id="tkg-photon"></a> Run Crashd on Photon OS Tanzu Kubernetes Grid Clusters

Crashd for Tanzu Kubernetes Grid provides a script file, `diagnostics.crsh`, along with a script argument file, `args`.  When Crashd runs, Crashd takes the the argument values from the `args` file and passes the values to the script. The script runs commands to extract information that can help diagnose problems on Photon OS Tanzu Kubernetes Grid management clusters and Tanzu Kubernetes workload clusters, which have been deployed on vSphere from Tanzu Kubernetes Grid.

### Prerequisites

Prior to running Crashd script `diagnostics.crsh`, your local machine must have the following programs on its execution path:

- `kind` (v0.7.0 or higher)
- `kubectl`
- `scp`
- `ssh`

Additionally, before you can run Crashd, you must follow these steps:

- Configure Crashd with a SSH private/public key pair.
- Ensure that your Tanzu Kubernetes Grid VMs are configured to use your SSH public key.
- Extract the `kubeconfig` file for the management cluster by using command `tanzu cluster kubeconfig get <management-cluster-name>`.
- For a simpler setup, ensure that the `kubeconfig`, `public-key` file, the `diagnostics.crsh` file, and the `args` file are in the same location.

### Configure Crashd

1. Navigate to the location where you downloaded and unpacked the Crashd bundle.
1. In a text editor, open the argument file `args`.

   For example, use `vi` to edit the file.

   ```
   vi args
   ```

   The file contains a series of named key/value pairs that are passed to the script:

   ```
   # Specifies cluster to target, (supported: bootstrap, mgmt, or workload)
   target=mgmt

   # Underlying infrastructure used by TKG (supported: vsphere, aws)
   infra=vsphere

   # working directory
   workdir=./workdir

   # User and private key for ssh connections to cluster nodes.
   ssh_user=capv
   ssh_pk_file=./capv.pem

   # namespace where mgmt cluster is deployed
   mgmt_cluster_ns=tkg-system

   # kubeconfig file path for management cluster
   mgmt_cluster_config=./tkg_cluster_config

   # Uncomment the following to specify a comma separated 
   # list of workload cluster names
   #workload_clusters=tkg-cluster-wc-498

   # Uncomment the following to specify the namespace
   # associated with the workload cluster names above
   #workload_cluster_ns=default
   ```

1. Configure the collection of diagnostics from a bootstrap cluster.

   If you are troubleshooting the initial setup of your cluster during bootstrap, update the following arguments in the file:

    - `target`: Set this value to `bootstrap`.
    - `workdir`: The location where files are collected.

1. Configure the collection of diagnostics from a management cluster.

   When diagnosing a management cluster failure, update the following arguments in the args file:

    - `target`: Set this value to `mgmt`.
    - `workdir`: The location where files are collected.
    - `ssh_user`: The SSH user used to access cluster machines. For clusters running on vSphere, the user name is `capv`.
    - `ssh_pk_file`: The path to your SSH private key file. For information about creating the SSH key pairs, see [Create an SSH Key Pair](../mgmt-clusters/vsphere.md#ssh-key) in *Deploy a Management Cluster to vSphere*.
    - `mgmt_cluster_ns`: The namespace where the management cluster is deployed.
    - `mgmt_cluster_config` The path of the kubeconfig file for the management cluster.

1. Configure the collection of diagnostics from one or more workload clusters:

   When collecting diagnostics information from workload clusters, you must specify the following arguments:

    - `target`: Set this value to `workload`.
    - `workdir`: The location where files are collected.
    - `mgmt_cluster_ns`: The namespace where the management cluster is deployed.
    - `mgmt_cluster_config` The path of the kubeconfig file for the management cluster.

   In addition to the previous arguments, you must uncomment the following workload cluster values:

    - `workload_clusters`: A comma-separated list of workload cluster names from which to collect diagnostics information.
    - `workload_cluster_ns`: The namespace where the workload clusters are deployed.

### Run Crashd

1. Run the `crashd` command from the location where the script file `diagnostics.crsh` and argument file `args` are located.

   ```
   crashd run --args-file args diagnostics.crsh
   ```

1. Optionally, monitor Crashd output. By default, the `crashd` command runs silently until completion.  However, you can use flag `--debug` to view log messages on the screen similar to the following:

   ```
   crashd run --args-file args --debug diagnostics.crsh

   DEBU[0003] creating working directory ./workdir/tkg-kind-12345
   DEBU[0003] kube_capture(what=objects)
   DEBU[0003] Searching in 20 groups
   ...
   DEBU[0015] Archiving [./workdir/tkg-kind-12345] in bootstrap.tkg-kind-12345.diagnostics.tar.gz
   DEBU[0015] Archived workdir/tkg-kind-12345/kind-logs/docker-info.txt
   DEBU[0015] Archived workdir/tkg-kind-12345/kind-logs/tkg-kind-12345-control-plane/alternatives.log
   DEBU[0015] Archived workdir/tkg-kind-12345/kind-logs/tkg-kind-12345-control-plane/containerd.log
   ```
