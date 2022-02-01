# Directory Structure

Packages that are added to the Tanzu Community Edition source repository should conform to the following directory structure. If the package is being developed in its own repository, it is free to follow the best practices defined by the Maintainers. For more information about the filesystem structure used for packages, see the [Package Contents Bundle](https://carvel.dev/kapp-controller/docs/latest/packaging-artifact-formats/#package-contents-bundle) topic in the Carvel documentation.

```shell
├── 1.0.0
│   ├── README.md
│   ├── bundle
│   │   ├── .imgpkg
│   │   │   └── images.yml
│   │   ├── config
│   │   │   ├── schema.yml
│   │   │   ├── overlays
│   │   │   │   └── overlay-a.yaml
│   │   │   ├── upstream
│   │   │   │   └── upstream-a.yaml
│   │   │   └── values.yaml
│   │   ├── vendir.lock.yml
│   │   └── vendir.yml
│   ├── package.yaml
│   └── test
│       ├── Makefile
│       ├── README.md
│       ├── e2e
│       │   └── test.go
│       └── unittest
│           └── test.go
├── LICENSE
└── metadata.yaml
```

## Files and Directories

* _version_

  Each package will have a _version_ directory that aligns with the version of the underlying software. In this case, the version is `1.2.3`.

* README.md

  File that contains documentation about this specific version of the package.

* bundle

  Packages are bundled by [imgpkg](https://carvel.dev/imgpkg/docs/latest/basic-workflow/#step-1-creating-the-bundle). This is the directory that contains the source of the package.

* .imgpkg

  A hidden directory required by `imgpkg` for storing bundle related information. Contains the `images.yaml`, which records image references used by the configuration.

* config

  The directory that contains the actual source of the package: upstream manifests, ytt [overlays](https://carvel.dev/ytt/docs/latest/ytt-overlays/), [schema](https://carvel.dev/ytt/docs/latest/how-to-write-schema/) and a default configuration values.

* vendir files

  The [vendir](https://carvel.dev/vendir/docs/latest/vendir-spec/) files reference where to obtain the upstream manifests for the package.

* package.yaml

  The [package.yaml](https://carvel.dev/kapp-controller/docs/latest/packaging/#package) is created for every new version of a package and it carries information about how to fetch, template, and deploy the package.

* test

  Each version of a package should contain unit and/or end-to-end tests. There should also be a README that describes the tests and instructs how to execute them.

* LICENSE

  The license to use the package. See the [Licensing](../licensing/) section for more details.

* _metadata.yaml_

  Contains high level information about the package. See the [PackageMetadata](./cr-files/) section for more details.

## Bootstrap Commands

1. Create a directory for the package with the version

    ```shell
    mkdir -p example/1.0.0
    ```

2. Change into the directory you created in the previous step, and create further directories as follows:

    ```shell
    cd example/1.0.0
    mkdir -p bundle/.imgpkg
    mkdir -p bundle/config/overlays
    mkdir -p bundle/config/upstream
    mkdir -p test
    ```
