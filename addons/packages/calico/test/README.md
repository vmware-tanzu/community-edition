# Calico tests

## Unit Tests

The unit tests for Calico test manifest generation of the package given some set
of data values.

### Prerequisites

To run the unit tests you need:

* [ginkgo](https://onsi.github.io/ginkgo/)
* [ytt](https://carvel.dev/ytt/)

### Run Tests

To run the unit tests you can run from this directory:

```bash
make test
```

## Development

The tests have its own Go module. Most tooling for Golang projects (e.g gopls)
require you to be within the directory of the `go.mod` file. It is recommended
that you are in this subdirectory when you are working on this module.

There is also a shared testing library for packages
[../../test/pkg](../../test/pkg), located outside of this module and it is
required by this module using a replace directive For Golang tooling to work in
this module you need to be in that subdirectory.
