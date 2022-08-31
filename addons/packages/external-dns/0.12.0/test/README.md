# ExternalDNS tests

## End-to-End Tests

End-to-End tests for `external-dns` are located in the `./e2e` directory.

### Prerequisites

To run the `external-dns` end-to-end tests you need:

* A Tanzu Community Edition cluster and the cluster needs to be the
  current-context. See the [Getting Started
  Guide](https://tanzucommunityedition.io/docs/getting-started/) for
  instuctions on how to create a cluster.
* The cluster supports Service type `LoadBalancer`.
* The `external-dns.community.tanzu.vmware.com` Package must exist on the
  cluster so it can be installed by the test.

Clusters built with `docker` or `kind` providers do not support Service type
`LoadBalancer` by default, one method of supporting Service type `LoadBalancer`
is to install MetalLB.

To install MetalLB run:

```bash
$ kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.12.1/manifests/namespace.yaml
namespace/metallb-system created
```

```bash
$ kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.12.1/manifests/metallb.yaml
podsecuritypolicy.policy/controller configured
podsecuritypolicy.policy/speaker configured
serviceaccount/controller created
serviceaccount/speaker created
clusterrole.rbac.authorization.k8s.io/metallb-system:controller unchanged
clusterrole.rbac.authorization.k8s.io/metallb-system:speaker unchanged
role.rbac.authorization.k8s.io/config-watcher created
role.rbac.authorization.k8s.io/pod-lister created
role.rbac.authorization.k8s.io/controller created
clusterrolebinding.rbac.authorization.k8s.io/metallb-system:controller unchanged
clusterrolebinding.rbac.authorization.k8s.io/metallb-system:speaker unchanged
rolebinding.rbac.authorization.k8s.io/config-watcher created
rolebinding.rbac.authorization.k8s.io/pod-lister created
rolebinding.rbac.authorization.k8s.io/controller created
daemonset.apps/speaker created
deployment.apps/controller created
```

To configure LoadBalancer IPs that MetalLB can use you can run the following:

**Note: The addresses you specify must be in the Kind network subnet, but should
not overlap with any node IPs. The range used below is in the default Kind
network subnet and in a high enough IP range that it shouldn't overlap with any
node IPs.**

```bash
$ cat <<EOF | kubectl apply -f -
---
apiVersion: v1
kind: ConfigMap
metadata:
  namespace: metallb-system
  name: config
data:
  config: |
    address-pools:
    - name: default
      protocol: layer2
      addresses:
      - 172.18.0.240-172.18.0.250
EOF
configmap/config created
```

### Test Configuration

Set the `DOCKERHUB_PROXY` environment variable if you would like to override
`docker.io` with a proxy.

### Run Tests

Run the tests from the e2e directory:

```bash
$ cd addons/packages/external-dns/0.12.0/test
$ make e2e-test
...
External-dns Addon E2E Test
...
    Ran 1 of 1 Specs in 108.396 seconds
    SUCCESS! -- 1 Passed | 0 Failed | 0 Pending | 0 Skipped
    PASS

    Ginkgo ran 1 suite in 1m49.44702766s
    Test Suite Passed
```

## Development

The tests have its own Go module. Most tooling for Golang projects (e.g gopls)
require you to be within the directory of the `go.mod` file. It is recommended
that you are in this subdirectory when you are working on this module.

There is also a shared testing library for packages,
[../../test/pkg](../../test/pkg), located outside of this module and it is
required by this module using a replace directive. For Golang tooling to work in
this module you need to be in that subdirectory.
