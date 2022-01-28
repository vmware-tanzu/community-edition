---
title: Maintainer
---

Maintainers are responsible for the creation, ownership and ongoing maintenance of the package. This could be a person, team or organization. They may or may not be associated the underlying software that is being packaged.

## Responsibilities

The role of being a Package Maintainer falls into 3 categories: Authoring, Owning and Maintaining.

### Authoring

When authroing a package, the Maintainer should be knowledgeable, cognizant and able to satisfy the following.

* Software components that are to be packaged
* Configuration parameters that are to be exposed
* Interactions and dependencies between packages
* Compatible Kubernetes distributions
* Considerations for underlying infrastructure (e.g. AWS, GCP, Docker, vSphere, etc)
* Documentation of the package. This should include a brief overview of the software components contained in the package, a description of configuration parameters, and example usage information.
* Unit and End-to-End tests should be provided along with execution instructions.

### Ownership

* Ensure that the package meets minimum package standards as defined in the Creation documentation.
* Performs as intended
* Is available for consumption

  This means the source of the package needs to be publicly accessible and package images are published to public OCI Registries.

* Has the mechanisms, processes and staff necessary to maintain the package.

  Should a Maintainer no longer be able to effectively own and maintain the package, they should notify the Tanzu Community Edition community at the earliest opportunity of the intent to orphan or retire the package.

### Maintenance

* Updating the package with newer versions of packaged components
* Exposing important configuration parameters
* Updating documentation with new version information, configuration parameters and providing more detailed use cases and examples.
* Community Support - responding to End User and Community questions, comments, concerns
