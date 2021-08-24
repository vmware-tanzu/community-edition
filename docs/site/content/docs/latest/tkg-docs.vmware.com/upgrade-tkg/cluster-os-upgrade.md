# Select an OS During Cluster Upgrade

If your IaaS account has multiple base VM images with the same version of Kubernetes that you are upgrading to, your `tanzu management-cluster upgrade` or `tanzu cluster upgrade` command can specify which OS version to use.

You specify the OS version with the `--os-arch`, `--os-name`, or `--os-version` options to the `upgrade` command.
Possible values and defaults for these options include:

* `--os-name` value depends on cloud infrastructure:
  - **vSphere**: `ubuntu` (default) or `photon` for [Photon OS](https://vmware.github.io/photon/assets/files/html/3.0/)
  - **Amazon EC2**: `ubuntu` (default) or `amazon` for [Amazon Linux](https://aws.amazon.com/amazon-linux-2/)
  - **Azure**: `ubuntu`
* `--os-version` value depends on `os-name`:
  - `ubuntu` values include: `20.04` (default), `18.04`
  - `photon` values include: `3` (default)
  - `amazon` values include: `2` (default)
* `--os-arch` value: `amd64` (default)

If you do not specify an `--os-name` when upgrading a cluster, its nodes retain their existing `--os-name` setting.
