---
title: Upstream Content
weight: 3
---

[vendir](https://carvel.dev/vendir) is used to synchronize the upstream content for the package. This is done by defining the location of the upstream content in a [`vendir.yaml`](https://carvel.dev/vendir/docs/latest/vendir-spec/) file and running `vendir sync`. Vendir supports a large number of upstream content [sources](https://carvel.dev/vendir/docs/latest/), such as source control systems, HTTP sources and Helm Charts.

## Sample vendir File

This is an example `vendir.yaml` file with a `githubRelease` as the source of the upstream content. The vendir source repository provides an extensive set of [examples](https://github.com/vmware-tanzu/carvel-vendir/tree/develop/examples).

```yaml
apiVersion: vendir.k14s.io/v1alpha1
  kind: Config
  minimumRequiredVersion: 0.12.0
  directories:
    - path: config/upstream
      contents:
        - path: .
          githubRelease:
            slug: example-repo/example
            tag: 1.0.0
            disableAutoChecksumValidation: true
          includePaths:
            - example.yaml
```

## Example Usage

For the purposes of this example, we'll use cert-manager.

1. Create a `vendir.yml` file that indicates where to find the remote, upstream configuration for cert-manager.  It indicates to vendir to synchronize the `config/upstream` directory created in the previous step with the contents of the cert-manager v1.5.3 GitHub release located in the `jetstack/cert-manager` repository. From that release, we want the `cert-manager.yaml` file.

    ```shell
    cat > bundle/vendir.yml <<EOF
    apiVersion: vendir.k14s.io/v1alpha1
    kind: Config
    minimumRequiredVersion: 0.12.0
    directories:
    - path: config/upstream
      contents:
        - path: .
          githubRelease:
            slug: jetstack/cert-manager
            tag: v1.5.3
            disableAutoChecksumValidation: true
          includePaths:
            - cert-manager.yaml
    EOF
    ```

2. Run the vendir sync command to pull down the cert-manager manifest.

    ```shell
    vendir sync --chdir bundle
    ```

3. Run the following command to inspect your local `bundle/config/upstream` directory.

    ```shell
    ls -l bundle/config/upstream
    ```

    You should see the `cert-manager.yaml` file from the `v1.5.3` cert-manger release is present.

    ```shell
    -rw-r--r--  1 seemiller  staff  1442034 Oct 18 12:39 cert-manager.yaml
     ```

    You should also see the `bundle/vendir.lock.yml` file has been created. This lock file resolves the `v1.5.3` release tag to the specific GitHub release and declares that the `config/upstream` is the synchronization target path. If you inspect the file, the contents should look like this:

    ```yaml
    apiVersion: vendir.k14s.io/v1alpha1
    directories:
    - contents:
      - githubRelease:
          url: https://api.github.com/repos/jetstack/cert-manager/releases/48370396
          path: .
      path: config/upstream
    kind: LockConfig
    ```
