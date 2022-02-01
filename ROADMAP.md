# VMware Tanzu Community Edition High-Level Program Roadmap

## About this Document

This document provides an overview of the major themes driving Tanzu Community Edition program development, as well as constituent features and capabilities. Roadmap items are tracked by their implementation stage within the enhancements proposal process: Awaiting Proposal, Proposal in Review, and Implementation in Progress. The stages are documented below. Most items are gathered from the community or include a feedback loop with the community.

The scope of this roadmap is the entire Tanzu Community Edition program, which includes this GitHub project repository, the [Tanzu Framework](https://github.com/vmware-tanzu/tanzu-framework) project frepository, the project repositories for packages delivered as part of Tanzu Comm unity Edition, and non-code program aspects such as videos, training materials, and community activities.

This roadmap should serve as a reference point for Tanzu Community Edition users to help us prioritize existing features, provide input on unmet needs, and for contributors to understand where the project is heading. Contributors can also consult the roadmap when thinking of proposing new ideas and determine if their idea conflicts with the roadmap. If it does conflict, the team and community will need to determine whether to adjust the roadmap or to recommend changes to the idea.

Enhancement stages help everyone understand if something is in the early stages of an idea or if it has wider community buy-in and committed resources. This is important because it helps limit waste, and focus work on the most important and timely activities. The stages are described below:

* Awaiting Proposal - these enhancements are in the idea mode, and the project maintainers are awaiting a detailed proposal. This state is useful to explore different solutions to a known problem and encourage discussion.
* Proposal in Review - maintainers are currently reviewing the proposal, which includes resource requirements and commitments.
* Implementation in Progress - engineering work has begun and will be tracking towards a release target.

## How to help

If you'd like to contribute but don't have anything to propose, check out our [issues](https://github.com/vmware-tanzu/community-edition/issues) for issues tagged
 [help-wanted](https://github.com/vmware-tanzu/community-edition/issues?q=is%3Aopen+is%3Aissue+label%3Ahelp-wanted) or
 [good-first-issue](https://github.com/vmware-tanzu/community-edition/labels/good-first-issue).
 
 You can participate in discussion on the roadmap in threads under Issues and in sdcheduled roadmap serssions in community meetings. Questions may also be brought up in our Google Group or our Slack channel. Community members are encouraged to be actively involved, and also stay informed so contributions can be made with the most positive effect and limited duplication of effort.

## How to add an item to the roadmap

Please create an issue using the Feature Request template to propose a feature for the project. We will work with and rely on our community to help us flesh out and prioritize our efforts to improve Tanzu Community Edition. Or feel free to submit a pull request against your issue for assessment by the team and community.

## Current High-Level Roadmap

The following table includes the current roadmap for Tanzu Community Edition. Please take timelines and dates as proposals and goals. Priorities and requirements change based on community feedback, roadblocks encountered, community contributions, etc. If you depend on a specific item, we encourage you to attend community meetings to get updated status information, or help us deliver that feature by contributing.

Last Updated:  January 2022

| Theme | Feature | Status (Phase) | Targeted Release |
| ----- | ------- | -------------- | ---------------- |
| Build user community | Reward community support | n/a | ongoing |
| |Community roadmap prioritization | n/a | ongoing |
| Enhance cloud-native user experience for skilled and new users | Unmanaged-cluster model introduced to enable minimal-cluster deployments in under 4 minutes that can run on consumer hardware. | Implementation in Progress | v0.11 |
| | Work to reduce time to cluster availability on vSphere and public clouds (Cluster API v1.1 Proposal) | Awaiting Proposal | TBD |
| | Consider more efficient and resource-saving methods to bootstrap | Awaiting proposal | TBD |
| | Make management clusters registerable in Tanzu Mission Control | Implementation in progress | v0.11 |
| | Introduce a UI that takes a user from installation through creating clusters to installing packages on clusters. Embed documentation and guidance throughout so that users can be successful exclusively through this tooling. | Awaiting Proposal | TBD |
| | Expose Cluster API (capi)-provider bootstrapping details, client-side, to users when creating management and workload clusters. Provides users an understanding of where bootstrapping failures occurred. | Awaiting Proposal [(issue)](https://github.com/vmware-tanzu/community-edition/issues/2730 ) | TBD |
|Robust Kubernetes platform | Limited-internet (including airgapped) support | Awaiting Proposal | TBD |
| | Validate and improve end-to-end tests across all supported infrastructure providers. | Awaiting proposal | TBD |
| Rich package library | Work with community to solicit and prioritize desired packages| n/a |ongoing |
| | Pursue community and ISV package lifecycle and configuration integrations via Carvel and Tanzu Framework | n/a | ongoing |
| Positioning Community Edition as leading-edge upstream for Tanzu products | Community Edition produces its own Bill of Materials (BOM), Tanzu Kubernetes Release (TKr), OCI container images, host images, and packages for CLI plugins and in-cluster software | Awaiting Proposals | TBD |
| Application Platform | Installable package is made available containing pre-configured, minimally-viable, platform services that enable developers to run their tools and apps (contour, kpack, cert-manager, knative serving, Cartographer, etc.) | Implementation in Progress (TODO: tracker proposal for all components) | Beginning in v0.10; please see issue (in dev) |
| | Tools and processes for package contribution, maintenance, and installation | Awaiting Proposal | TBD |
| | Package maintainer documentation and guidance is live such that everyone can understand the responsibilities of being a package maintainer | Awaiting Proposal | TBD |
| |Provide guidance on bringing new provider packages to TCE; enabling the inclusion of new infrastructure providers | Awaiting Proposal (dependent on Tanzu Kubernetes Provider implementation) | TBD |
| |Introduce kubeapps to TCE for discovering and installation packages in Bitnami library [(issue)](https://github.com/vmware-tanzu/community-edition/issues/2418) | Awaiting Proposal | TBD |
| Better defined lifecycle management Cluster API extensibility points |Build the Runtime Extension SDK in Cluster API for ClusterClass patching, pre/post control plane upgrade hooks | Awaiting proposal | TBD |
