# Package Maintainer Docs

This page offers an index of resources for potential package maintainers looking
for guidance on how to get their package into TCE.

At this time, packages in our official repository are maintained exclusively by
VMware. This documentation is intended for VMware employees who will be
maintaining official packages.

> ⚠️ This README is a temporary solution until we create a proper package
maintainer website. We recommend opening each link in a new tab so that you can
refer back to this index as you progress with your package. ⚠️

## Considerations

Details considerations before creating a package.

1. [Proposing a new package](considerations/proposal.md)
1. [Git repository](considerations/git.md)
1. [Versioning](considerations/versioning.md)

## Creation

Details on the creation and contribution of a package.

1. [Prerequisites](creation/prerequisites.md)
1. [Directory structure](creation/directory-structure.md)
1. [Syncronize upstream assets](creation/upstream-content.md)
1. [Lock container image references](creation/image-refs.md)
1. [Add configurable parameters](creation/configuration.md)
1. [Add overlays to upstream content](creation/overlays.md)
1. [Add tests](creation/testing.md)
1. [Add documentation](creation/documentation.md)
1. [Lint the package](creation/linting.md)
1. [Publish the package](creation/publish.md)
1. [Create custom resources](creation/cr-files.md)

## Maintenance

Details expectations of package maintenance over time.

* [Community Support](maintenance/community-support.md)
* [Deprecation](maintenance/deprecation.md)
* [Maintenance](maintenance/maintenance.md)
* [Security](maintenance/security.md)

## Submitting

Details on submitting a created package into the TCE distribution.

* [Submitting](submitting.md)

## Roles

Details on the types of roles described in the above.

* [Maintainers](roles/maintainer.md)
* [End User](roles/end-user.md)
