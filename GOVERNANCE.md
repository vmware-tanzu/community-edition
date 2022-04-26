
# Tanzu Community Edition Governance

This document defines the project governance for Tanzu Community Edition, an open source project by VMware.

## Overview

VMware Tanzu Community Edition is a full-featured, easy to manage Kubernetes platform for learners and users. It is a freely available, community supported, open source distribution of VMware Tanzu that can be installed and configured in minutes on your local workstation or your favorite cloud.

We are committed to the project not only delivering the distribution, but also building an open, inclusive, and productive vendor driven open source community; together, we will advance a reliable, nimble, and extensible foundation for modern applications.

## Community

Users: Members who engage with the Tanzu Community Edition community, providing valuable feedback and unique perspectives.

Contributors: Members who contribute to the project through documentation, code reviews, responding to issues, participation in proposal discussions, contributing code, etc. The project welcomes code contributions to the Tanzu Community Edition project, as well as contributor-originated packages that add capabilities from other projects. These contributed packages will conform to the Tanzu Community Edition packaging requirements and lifecycle management.

Maintainers: Given the nature of this project and its relationship to VMware’s Tanzu product offerings, the Tanzu Community Edition project leaders are current employees of VMware. They are responsible for the overall health and direction of the project; final reviewers of PRs and responsible for releases. Some maintainers are responsible for one or more components within a project, acting as technical leads for that component. Maintainers are expected to contribute code and documentation, review PRs including ensuring quality of code, triage issues, proactively fix bugs, and perform maintenance tasks for these components. If a maintainer leaves VMware, they will also leave their maintainer position.

## Proposal Process

One of the most important aspects in any open source community is the concept of proposals. All large changes to the codebase and/or new features, including ones proposed by maintainers, should be preceded by a proposal in this repository. This process allows for all members of the community to weigh in on the concept (including the technical details), share their comments and ideas, and offer to help. It also ensures that members are not duplicating work or inadvertently stepping on toes by making large conflicting changes.

To understand our proposal process, please review [docs/designs](docs/designs/README.md).

## Lazy Consensus

To maintain velocity in a project as busy as Tanzu Community Edition, the concept of Lazy Consensus is practiced. Ideas and / or proposals should be shared by maintainers via GitHub with the appropriate maintainer groups (e.g., @vmware-tanzu/tce-maintainers) tagged. Out of respect for other contributors, major changes should also be accompanied by a ping on Slack, and a note on the Tanzu Community Edition mailing list as appropriate. Author(s) of proposals, pull requests, issues, etc. will specify a time period of no less than five (5) working days for comment and remain cognizant of popular observed world holidays.
Other maintainers may request additional time for review, but should avoid blocking progress and abstain from delaying progress unless absolutely needed. The expectation is that blocking progress is accompanied by a guarantee to review and respond to the relevant action(s) (proposals, PRs, issues, etc.) in short order. All pull requests need to be approved by at least 1 maintainer.

Lazy Consensus is practiced for the main project repository and the additional repositories listed above.
