# Testing `cartographer-catalog`

As the supply chain provided by this package is designed to take source code
from a git repository, build a container image, push that image to a registry,
and then submit to the cluster the necessary objects for running the
application

```text
source-provider  -------- image-builder ----------- deployer
       .                       .                       .
fluxucd/GitRepository     kpack/Image         kapp-ctrl/App
```

there are a couple of pre-requisites that must be met before executing the
test.

## Prerequisites

1. Package and PackageMetadata of the following packages in the cluster

* cartographer
* cartographer-catalog
* fluxcd-source-controller
* kpack
* kpack-dependencies
* knative-serving

1. secretgen-controller installed

1. credentials for the registry exported to all namespaces for a secret named
   `cartographer-catalog-test`

  ```bash
  tanzu secret registry add cartographer-catalog-test \
    --server REGISTRY-SERVER \
    --username USERNAME \
    --password PASSWORD \
    --export-to-all-namespaces \
    --yes
  ```

  this Secret is consumed in the `kpack` namespace for the ClusterBuilder
  setup, as well as the dynamically-generated namespace for the Workload.

1. `kpack.yaml` configured with credentials and registry where images can be
   pushed to

1. `kpack-dependencies.yaml` configured with the same `kp_default_repository`
   as `kpack.yaml`

1. `cartographer-catalog.yaml` configured with the details about the registry
   where the application container images should be pushed to

## Running the test

With the prerequisites met, we can go ahead with running the test:

```bash
./cartographer-catalog-test.sh
```

Under the hood, it'll take care of:

1. ensuring that all necessary Package dependencies are installed in the
   cluster

2. submitting a Workload with `kapp` and waiting for it to complete (`kapp`
   knows how to wait for a Workload to be Ready, in which case all supply chain
   resources got satisfied).
