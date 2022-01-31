---
title: Creation
weight: 1
---

Creating a Package is quite straightforward. Here is a high level overview of the steps to create a package. Click through the links in each step for more details and examples.

1. Ensure that you have met all the [prerequisites](./prerequisites/)
2. Create the [directory structure](./directory-structure)
3. Use [vendir](./tooling/#vendir) to synchronize [upstream content](./upstream-content/#example-usage)
4. Use [kbld](./tooling/#kbld) to create [immutable image references](./image-refs/#example-usage)
5. Define [configurable parameters](./configuration/) in the schema
6. Create [overlays](./overlays/) to apply configuration over the upstream content
7. Create [tests](./testing/)
8. Create [documentation](./documentation/)
9. [Linting](./linting/)
10. [Publish](./publish/) the package
11. Create [Package](./cr-files/#package) and [PackageMetadata](./cr-files/#packagemetadata) custom resources
