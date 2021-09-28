# Tagger

Tagger is a hack to prevent gc from breaking packages

Packages are referenced from repositories via their digest. This means
that if a tag on a configuration bundle mutates, the previous artifact
may no longer have a tag associated with it. These are considered
"dangling" images and will be collected (deleted) when Harbor runs its
garbage collection cycle.

The issue with this is, if a package repository references that
no longer existent digest, it will fail to retrieve the bundle, breaking
the package.

There are two changes that must take place:

1. imgpkg should set a default (additional) tag on bundles it pushes[0]
2. we must scan our package bundles and add an arbitrary and unique tag to ensure
   that they are never garbage collected.

This commit solves 2. It scans all the projects that represent package
bundles in TCE. It then resolves the artifacts for each of those
bundles. For each artifact it looks at the tags and sees if it has the
expected, unique, tag. The tag we've got with is the first 10 characters
of the digest's sha256 value. As an example:

```txt
sha256:0c3d0f33c171e437268e57bdbe0d83feb1606362d8235f1b656556da8c944e18
```

would have a required tag value of `c3d0f33c17`.

If missing, this tag is added to the artifact. If it exists, no
operation is triggered.

## Building

```sh
TAG=1.0.0 make build
```

## Building and Pushing

```sh
TAG=1.0.0 make build-and-push
```

## Running it

1. Run docker login, pointed at `projects.registry.vmware.com`
1. Run Tagger

  ```sh
  docker run \
    -v ${HOME}/.docker/config.json:/root/.docker/config.json \
    projects.registry.vmware.com/tce/tagger:0.1.0
  ```
