# Deploy Clusters with Different Kubernetes Versions

Tanzu Community Edition manages Kubernetes versions with Tanzu Kubernetes release (TKr) objects.

Complete the following steps to deploy a cluster with a non default version of Kubernetes:

1. The Tanzu Community Edition project keeps TKrs here:`projects.registry.vmware.com/tce/tkr`.
Use [imgpkg](https://carvel.dev/imgpkg) to query what TKrs are available:

```sh
imgpkg tag list -i projects.registry.vmware.com/tce/tkr
```

Output:

```txt
Name
sha256-2fd337282cf17357c6329f441dc970ec900145faef9e2ec6122f98fa75d529c3.imgpkg
sha256-33f63314fb72ead645715f6ac85128c0fe0fd380d14f0a79eddba3dd361b73dd.imgpkg
sha256-ac6566268e0f113a4b91bab870a34353685e886f97e248633bb2c2fcf6490dc8.imgpkg
v1.21.5
v1.21.5-1
v1.21.5-2
v1.21.5-3
v1.22.2
```

1. To create a cluster with an alternative TKr, specify the TKr in the --tkr option:

```sh
tanzu cluster create <WORKLOAD-CLUSTER-NAME> --tkr projects.registry.vmware.com/tce/tkr:<TKr-VERSION>
```

For example

```sh
tanzu cluster create tce-wl-cluster --tkr projects.registry.vmware.com/tce/tkr:v1.22.2
```

1. (Optional) To customize a TKr, you can pull an existing one down using `imgpkg`. For example,

```sh
imgpkg pull -i projects.registry.vmware.com/tce/tkr:v1.22.2 -o tkr
Pulling image 'projects.registry.vmware.com/tce/tkr@sha256:7c1a241dc57fe94f02be4dd6d7e4b29f159415417164abc4b5ab6bb10cf4cbaa'
Extracting layer 'sha256:e17e901811682a2c8c91c8865f3344a21fdf8f83f012de167c15d2ab06cc494a' (1/1)

Succeeded
```

You can then edit the TKr in the `tkr/tkr-bom-v1.22.2.yaml`. After modifying it, you may also wish to rename the YAML file. Once you have made your modifications, you can repush it using:

```sh
imgpkg push -f ./tkr/tkr-bom-CUSTOM.yaml -i ${YOUR_REGISTRY}:${YOUR_TAG}
```

Once pushed, you can reference this repo using the `--tkr` flag.
