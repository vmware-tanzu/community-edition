---
title: Submitting a Package to TCE
---

To submit a package for consideration into Tanzu Community Edition's package repository, you will need to create a GitHub Issue. The issue should provide details about:

* Functionality it provides
* Software involved
* TODO (joshrosso)

**highly recommended:** Wait for `status/approved` on proposal before doing work. You're welcome to begin work immediately, but if the proposal is `status/declined`, the work may go to waste.

### Open Pull Request in community-edition repo

When your package is complete and ready for acceptance, create a Pull Request in the [community-edition](https://github.com/vmware-tanzu/community-edition/pulls) GitHub repository. The PR should reference the GitHub issue created for introducing the package. The PR will be reviewed by the Tanzu Community Edition Engineering Team. Upon approval, the package will be added to the main package repository.
