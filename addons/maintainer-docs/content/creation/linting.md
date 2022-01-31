---
title: Linting
weight: 3
---

Before publishing a package, run a few checks to ensure that it is free of common mistakes and errors.

* Markdown

   The Tanzu Community Edition source repository has a [linter to check markdown](https://github.com/vmware-tanzu/community-edition/tree/main/hack/check/check-mdlint.sh). You can use this script as a guide to create your own linter in your source, or just run it if the package in the TCE source repo.

* Spelling

  The Tanzu Community Edition source repository has a [linter to check spelling](https://github.com/vmware-tanzu/community-edition/tree/main/hack/check/check-misspell.sh). It is located in `hack/check/check-misspell.sh`. You can use this script as a guide to create your own linter in your source, or just run it if the package in the TCE source repo.

* Tests

   Verify that your unit and end-to-end tests compile, execute and succeed.

* Alpine Image Check

   Due to licensing concerns, packages that contain an Alpine image are [restricted](../restrictions) from use. Use the [`imagelinter` utility](https://github.com/vmware-tanzu/community-edition/tree/main/hack/imagelinter) located in the Tanzu Community Edition source repository to check for Alpine images.
