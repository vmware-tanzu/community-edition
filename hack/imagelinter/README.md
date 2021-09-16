# Imagelinter

Imagelinter is used to verify and validate docker images.Docker images can be from dockerfiles or yaml | yml files.

## The problem

In order to adhere to open source compliance, some repos should not contain Alpine based images in any Docker or Image files. Hence, we make sure that all docker images (from Dockerfiles or from yaml|yml files) should not be Alpine-based images.

## Solution

There is no straight command or a way to identify Linux OS from Docker perspective. After a deep analysis, we have found a few solutions that can apparently identify Linux OS of a given image.The below are steps that are executed when the linter is run.

1. Read Image metadata from the registry/Image libraries
2. Pull image and analyze Image meta-data (History)
3. Create a container and read /etc/os-release file
4. Copy etc/os-release to the local path and then analyze
5. Copy /usr/lib/os-release to the local path and then analyze
6. Check if there is any License file inside the container
7. Copy a simple binary to the container
8. Nothing determines the OS then reject the image

## Configuration

imagelinter supports few configurations.The below is the default configuration file.

```yaml
---
includeExts:
- ".yaml"
- ".yml"
includeLines:
- 'image:'
- FROM
matchPattern:
- "packages/*/*/bundle/.imgpkg/*"
ignoreImages:
- "index.docker.io/grafana/grafana@sha256:4e5835bcfd55cf72563a06932f10c75d9d92a0e1334a4c83eaa9c5b897370b25"
- "index.docker.io/rancher/local-path-provisioner@sha256:9666b1635fec95d4e2251661e135c90678b8f45fd0f8324c55db99c80e2a958c"
- "k8s.gcr.io/external-dns/external-dns@sha256:e49f63e07498ce8484c9e9050c1dfbc2584f4c9c262433d80387e855725e6bce"
- "index.docker.io/kiwigrid/k8s-sidecar@sha256:444be8cef8b25b4aaddea692ae09e3883d9064c0d31b43c4ba388a83c920552f"
- "index.docker.io/minio/minio@sha256:e7e9a563f52bf95f614e8017c2da9bd5d9f2f1ae0dc1127767fa341b3ae22088"
- "gcr.io/cadvisor/cadvisor@sha256:10638ceca79c01f4045f4a645242e763fe62eeb71d859ff93b09b0854a0d2220"
succesValidators:
- apt-get
- apt
- yum
- "/lib/x86_64-linux-gnu"
- "/usr/lib/x86_64-linux-gnu"
- "imgpkg"
failureValidators:
- Alpine

```

## How to run imagelinter

-To manually run imagelinter, user must download(clone) the source code.

- Navigate to the root directory and run the following make command ```make imagelint```  or

- ```cd hack/imagelinter && go run main.go --path <any valid path or it takes only current working directory> -- config <provide config path or default path will be taken> --summary=true --details=fail```
