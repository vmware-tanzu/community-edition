# Test Execution Prerequisites

- A `app-toolkit-values.yaml` file containing the following information:

  ```yaml
  contour:
    envoy:
      service:
        type: ClusterIP
      hostPorts:
        enable: true
  
  cartographer_catalog:
    registry:
      server: [REGISTRY_SERVER]
      repository: [REGISTRY_REPOSITORY]

  developer_namespace: dev-test
  
  knative_serving:
    domain:
      type: real
      name: 127-0-0-1.sslip.io

  kpack:
    # name of registry secret where build artifacts are stored
    kp_default_repository: [DEFAULT_REGISTRY_URL]
    kp_default_repository_username: [DEFAULT_REGISTRY_USERNAME]
    kp_default_repository_password: [DEFAULT_REGISTRY_PASSWORD]

  # Below is used for testing, but is not part of the intended App Toolkit flow
  #
  # if using DockerHub, provide it in the format "https://index.docker.io/v1/"
  # registry.server: [REGISTRY_SERVER_SECRET]
  # registry.username: [REGISTRY_USERNAME_SECRET]
  # registry.password: [REGISTRY_PASSWORD_SECRET]
  ```

  Where:
  - `REGISTRY_SERVER` and `DEFAULT_REGISTRY_URL` are valid OCI registries to store kpack images, like `https://index.docker.io/v1/`
  - `REGISTRY_REPOSITORY` is the repository name (i.e., on Dockerhub, this is likely your username)
  - `DEFAULT_REGISTRY_USERNAME` and `DEFAULT_REGISTRY_PASSWORD` are the credentials for the specified registry.
  - `REGISTRY_<FOOBAR>_SECRET` are the values that appear in your `~/.docker/config.json` for your registry login

After creating the file with the required fields, you can start the actual test execution with this command: `./app-toolkit-test.sh`
You can also provide a `PackageRepo` url to the script and it will use that repository for the tests.
