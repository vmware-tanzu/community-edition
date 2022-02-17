# Upgrade kapp-controller

Tanzu Community Edition provides a built in version of kapp-controller that is dictated
by the Bill of Materials (BOM) of the Tanzu Kubernetes Release (TKR) in use.

It is possible to upgrade the kapp-controller version provided.

## Managed clusters

On the management cluster, there's a controller named `tanzu-addons-manager` that reconciles core packages defined
in the BOM to the downstream cluster. One of these packages is `kapp-controller`, so we need to pause reconciliation
of this package to allow us to replace the version of the package being used.

To do that, **on the management cluster**, you need to:

```shell
kubectl patch secret/<WORKLOAD_CLUSTER_NAME>-kapp-controller-addon -n default -p '{"metadata":{"annotations":{"tkg.tanzu.vmware.com/addon-paused": ""}}}' --type=merge
```

Example:

```shell
kubectl patch secret/my-tce-workload-kapp-controller-addon -n default -p '{"metadata":{"annotations":{"tkg.tanzu.vmware.com/addon-paused": ""}}}' --type=merge
```

Once this is done, you can change the image of the kapp-controller application installed on your workload cluster to the one you want.

You will need to edit the `app/<WORKLOAD_CLUSTER_NAME>-kapp-controller` on the `default` namespace (or the namespace where your workload cluster has been defined) of the management cluster.

You can verify the current image in use with:

```shell
kubectl get app/<WORKLOAD_CLUSTER_NAME>-kapp-controller -n default -o jsonpath='{.spec.fetch[0].imgpkgBundle.image}'
```

Example:

```shell
$ kubectl get app/my-tce-workload-kapp-controller -n default -o jsonpath='{.spec.fetch[0].imgpkgBundle.image}'
projects.registry.vmware.com/tkg/packages/core/kapp-controller:v0.23.0_vmware.1-tkg.1
```

Create a patch file with the version of kapp-controller you want to use:

```yaml
spec:
  fetch:
  - imgpkgBundle:
      image: projects-stg.registry.vmware.com/tkg/packages/core/kapp-controller:v0.25.0_vmware.1-tkg.1-zshippable
```

Now patch the Application:

```shell
kubectl patch app/<WORKLOAD_CLUSTER_NAME>-kapp-controller --type merge --patch "$(cat /tmp/patch.yaml)"
```

Example:

```shell
kubectl patch app/my-tce-workload-kapp-controller --type merge --patch "$(cat /tmp/patch.yaml)"
```

Your kapp-controller will be updated in a few seconds.

You should be able to verify the version of kapp-controller in use in the workload cluster by looking
at the `kapp-controller.carvel.dev/version` annotation on the kapp-controller deployment:

```shell
kubectl get deploy kapp-controller -n tkg-system -ojsonpath='{.metadata.annotations.kapp-controller\.carvel\.dev/version}'
```

**NOTE**: Remember this command needs to be run on the workload cluster.

Example:

```shell
$ kubectl get deploy kapp-controller -n tkg-system -ojsonpath='{.metadata.annotations.kapp-controller\.carvel\.dev/version}'
v0.25.0
```

## Existing kapp-controller imgpkg bundle versions

To verify the existing versions of the `kapp-controller` imgpkg bundle you can list the existing tags on
the VMware OCI registries.

To list released versions of the package:

```shell
imgpkg tag list -i projects.registry.vmware.com/tkg/packages/core/kapp-controller
```

To list pre-release versions of the package:

```shell
imgpkg tag list -i projects-stg.registry.vmware.com/tkg/packages/core/kapp-controller
```

Find one that matches the version you want and use it in the patch commands above.

Example:

```shell
$ imgpkg tag list -i projects-stg.registry.vmware.com/tkg/packages/core/kapp-controller
Tags

Name
v0.20.0_vmware.1-tkg.1-rc.1
v0.20.0_vmware.1-tkg.1-rc.2
v0.20.0_vmware.1-tkg.1-zshippable
v0.20.0_vmware.1-tkg.1-zshippablerelease
v0.23.0_vmware.1-tkg.1
v0.23.0_vmware.1-tkg.1-rc.3
v0.23.0_vmware.1-tkg.1-rc.4
v0.23.0_vmware.1-tkg.1-rc.5
v0.23.0_vmware.1-tkg.1-zshippable
v0.23.0_vmware.1-tkg.2-20210924-539f8b15
v0.23.0_vmware.1-tkg.2-20210930-5b764f3e
v0.23.0_vmware.1-tkg.2-framework-v0.3.0
v0.23.0_vmware.1-tkg.2-framework-v0.4.0
v0.23.0_vmware.1-tkg.2-framework-v0.5.0
v0.23.0_vmware.1-tkg.2-zshippable
v0.25.0_vmware.1-tkg.1-20211007-6d459b1c
v0.25.0_vmware.1-tkg.1-framework-v0.6.0
v0.25.0_vmware.1-tkg.1-zshippable

18 tags

Succeeded
```

## Existing kapp-controller image versions

To verify the existing versions of the `kapp-controller` imgpkg bundle you can list the existing tags on
the VMware OCI registries.

To list released versions of the package:

```shell
imgpkg tag list -i projects.registry.vmware.com/tkg/kapp-controller
```

To list pre-release versions of the package:

```shell
imgpkg tag list -i projects-stg.registry.vmware.com/tkg/kapp-controller
```

Find one that matches the version you want and use it in the patch commands above.

Example:

```shell
$ imgpkg tag list -i projects-stg.registry.vmware.com/tkg/kapp-controller
Tags

Name
v0.20.0_vmware.1-tkg.1-rc.1
v0.20.0_vmware.1-tkg.1-rc.2
v0.20.0_vmware.1-tkg.1-zshippable
v0.20.0_vmware.1-tkg.1-zshippablerelease
v0.23.0_vmware.1-tkg.1
v0.23.0_vmware.1-tkg.1-rc.3
v0.23.0_vmware.1-tkg.1-rc.4
v0.23.0_vmware.1-tkg.1-rc.5
v0.23.0_vmware.1-tkg.1-zshippable
v0.23.0_vmware.1-tkg.2-20210924-539f8b15
v0.23.0_vmware.1-tkg.2-20210930-5b764f3e
v0.23.0_vmware.1-tkg.2-framework-v0.3.0
v0.23.0_vmware.1-tkg.2-framework-v0.4.0
v0.23.0_vmware.1-tkg.2-framework-v0.5.0
v0.23.0_vmware.1-tkg.2-zshippable
v0.25.0_vmware.1-tkg.1-20211007-6d459b1c
v0.25.0_vmware.1-tkg.1-framework-v0.6.0
v0.25.0_vmware.1-tkg.1-zshippable

18 tags

Succeeded
```
