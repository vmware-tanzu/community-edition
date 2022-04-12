# Test Execution Prerequisites

- A `app-toolkit-values.yaml` file containing the following information:

  ```yaml
  contour:
    envoy:
      service:
        type: ClusterIP
      hostPorts:
        enable: true

  knative_serving:
    domain:
      type: real
      name: 127-0-0-1.sslip.io

  kpack:
    # name of registry secret where build artifacts are stored
    kp_default_repository: [DEFAULT_REGISTRY_URL]
    kp_default_repository_username: [DEFAULT_REGISTRY_USERNAME]
    kp_default_repository_password: [DEFAULT_REGISTRY_PASSWORD]
  ```

  Where:
  - `DEFAULT_REGISTRY_URL` is a valid OCI registry to store kpack images, like `https://index.docker.io/v1/`
  - `DEFAULT_REGISTRY_USERNAME` and `DEFAULT_REGISTRY_PASSWORD` are the credentials for the specified registry.

- A `supplychain-example-values.yaml` file containing the following information:

  ```yaml
  kpack:
    registry:
      url: [REGISTRY_URL]
      username: [REGISTRY_USERNAME]
      password: [REGISTRY_PASSWORD]
    builder:
      # path to the container repository where kpack build artifacts are stored
      tag: [REGISTRY_TAG]
    # A comma-separated list of languages e.g. [java,nodejs] that will be supported for development
    # Allowed values are:
    # - java
    # - nodejs
    # - dotnet-core
    # - go
    # - ruby
    # - php
    languages: [java]
    image_prefix: [REGISTRY_PREFIX]

  ```

  Where:
  - `REGISTRY_URL` is a valid OCI registry to store kpack images, like `https://index.docker.io/v1/`
  - `REGISTRY_USERNAME` and `REGISTRY_PASSWORD` are the credentials for the specified registry.
  - `REGISTRY_TAG` is the path to the container repository where kpack build artifacts are stored. For Docker, it is username/tag, e.g. csamp/builder
  - `REGISTRY_PREFIX` is prefix for your images as they reside on the registry. For Docker, it is the username with a trailing slash, e.g. csamp/

After creating two files with the required fields, you can start the actual test execution with this command: `go run app-toolkit-test.go`
