# TCE Release Process

## Summary

This document describes the TCE release process. The release process is functionally complete at this point based on our current needs, but improvements can be made as the project evolves.

There are many different layers and capabilities that the release process offers us, which we will get into more detail later on this document, and they include:

- Building only TCE specific components for development
- Building a installable TCE distribution for development, test, and self-signing (unsigned Tanzu CLI + TKG plugins + TCE plugins)

## Key Terms and Concepts

- **unit test**: it is a type of test that exercises discrete sections of code at varying levels of integration. In Go(lang) this is typically done using `go test`.
- **integration test**: it is a type of test that exercises interactions (ie input and output) between those components. Tests for integration can happen at many different levels. It could be between 2 lower level libraries (ie third party and TCE. example, gRPC and a TCE library), between 2 TCE produced libraries, and even an e2e test (mentioned below) is a kind of integration test between the TCE or the software and a human.
- **e2e test**: end-to-end test, it is a type of test that exercises typical user tasks and user scenarios, it tends to focus on integration points and workflow.
- **Design Verification**: this is the collective set of unit tests, integration tests, and e2e tests that should be run in an automated fashion to ensure the reliability and confidence of a given TCE release.
- **CI - Continuous Integration**: Think Jenkins, TravisCI, GitHub Actions, etc. It’s typically a scheduling framework, application or server that facilitates building, testing, and deploying code, a project, bits, and other code components in a repeatable fashion or cadence.
- **Core APIs**:  Core APIs are APIs that are needed across multiple products and programs.  Having them in a central location enables teams to have a single vendored point that can be leveraged in other projects. The APIs that TCE will heavily depend on are available here: [https://github.com/vmware-tanzu/core](https://github.com/vmware-tanzu-private/core)
- **Tanzu CLI**:  Is a community wide effort to consolidate multiple disparate CLIs and unify them in a single user experience for the products they manage. As of writing this document, the code/implementation for the Tanzu CLI resides in the aforementioned Core APIs repo: [https://github.com/vmware-tanzu/core](https://github.com/vmware-tanzu-private/core )
- **Addons**: Now called a package (see below)
- **Package**: Installable packages that can be deployed on top of Kubernetes clusters.
- **TCE**: Tanzu Community Edition
- **TKG**: Tanzu Kubernetes Grid
- **Carvel**: Is a toolchain or kit composed of many different open source projects used primarily in conjunction with Kubernetes or container-based offerings. More information can be found here: [https://github.com/vmware-tanzu/carvel](https://github.com/vmware-tanzu/carvel)
- **Imgpkg**: Is a tool (in the Carvel toolchain) that allows users to store a set of arbitrary files as an OCI image. One of the driving use cases is to store Kubernetes configuration (plain YAML, ytt templates, Helm templates, etc.) in an OCI registry as an image (link).

## Motivation

This document will attempt to describe the release process in order to facilitate:

- Understanding how TCE is built today and the needs it satisfies
- What constitutes a TCE release
- What capabilities are available for developers or others would might want to build TCE

This document builds upon the TCE Test/Release Requirements which is broken down into 2 major deliverables for TCE:

- the release process
- the design verification process

This document attempts to jump into and explore the implementation details with respect to the first item: the release process.

## Goals

- Describe the current state/implementation for the release process today and define all the components, interactions, boxes, moving parts, and pieces that make a TCE release.
- Provide a contextual reference for folks who are onboarding into the space

## Non-Goals

- Discussion on design or implementation for the design verification deliverables. This document should only expose where/how those tests might integrate into those processes.
- This document being the final word on an implementation. As we have dependencies on other projects/teams and the Release Engineering team in India ramps up, the team will collectively have a unique skill set that should be leveraged in the final implementation.

## Roles/Personas

The following personas are available when writing user stories:

- Community user/member
- TCE or community contributor
- A Developer - someone who contributes in the form of a PR
- Release engineer
- Software lifecycle owner

## User Stories

This section lists the typical user stories as they relate to the TCE release process. It is not complete and subject to change as TCE evolves.

### Community user/member Use-Cases

As community user/member, I would like to:

- Download, install, and use the capabilities provided by TCE
- Through an easy to consume method, like a badge, percent or dashboard (To be implement - design verification) I would like to:
  - See the health of the code with in the repo
  - See if TCE specific components are building successfully

### TCE or community contributor Use-Cases

As TCE or community contributor, I would like to:

- Submit a contribution in the form of a PR and have a set of checks (completed) and tests (To be implemented) run against my PR for design verification
- Once a contribution/PR is merged, I want to be able to consume binaries generated by the CI system
- For local developer verification, I want to be able to:
  - Build a TCE specific binaries with changes I’m looking to contribute
  - Build an entire TCE release with changes I’m looking to contribute

### Release engineer Use-Cases

As a release engineer, I would like to understand:

- What components are composed in a TCE release
- The integration points for where the release process and design verification process intersects
- The different software and systems that are involved in a TCE release for:
  - Making improvements in the process
  - Fixing breaks in the process

### Software lifecycle owner Use-Cases

As a software lifecycle owner, I would like to:

- Have a defined software development lifecycle for doing things like:
  - Defining milestones
  - Scheduling
- Produce or generate a release

## How a Contributor Creates a TCE Release

We are currently consuming the CLI from the Core repo using a TCE specific branch, this is a temporary stopgap until the next Core release stabilizes and our TCE specific PRs can be merged into Core. As a result, we temporarily need to build the Core repo (and some others).

The side effect or consequence for temporarily needing to rebuild all components is that it happens to also be useful for developers to build their own TCE release locally for testing out features they would like to introduce or contribute. This can simply be done by running `make release`.

![makefile](/docs/img/makerelease.png)

This builds everything required and generates the final packaging. This also happens to be the best integration point for building a release to be consumed for design verification (specifically e2e testing).

Other build targets that developers might find useful are `build-all` which just builds and installs all components locally on the file system and `build-plugin` which assumes Tanzu CLI is already installed and configured and only builds the TCE components and installs them locally. Another useful target is `release-docker` which can do everything that make release does (ie generate a TCE release), but does not require a development environment, and instead uses a standard Go Docker container.

## Automation for Pull Requests

When a PR is opened, updated, etc by a community contributor, a GitHub Action kicks off:

- **make check** which runs about 35 checks to make sure what is being introduced falls within the best practices for code quality, for example, static code analyzers, and markdown verification.
- **make build** which makes sure that TCE related components build successfully. This is a TCE only build because Tanzu CLI code doesn’t exist in the repo.

The GitHub Actions that will be invoked are (and depends on what is changed): check-build.yaml, check-lint.yaml, check-mdlint.yaml, and check-shell.yaml. For smoke tests, we probably should have a simple light-weight check using CAPD to do some quick checks for plugins and packages within the repo.

When a PR is merged, we perform a due diligence check and run make check and make build just to verify the expected result against main. Additionally, a GitHub Action runs in release-staging.yaml which builds all the plugins within the TCE repo and uploads them to the staging TCE GCP bucket for doing dynamic installation of plugins (i.e. running `tanzu plugin upgrade [plugin name]`). For contributors, this is a simple way for downloading a CI generated build of a plugin with your changes for test. To do that, you can override your production TCE GCP bucket with the staging GCP bucket by deleting the production bucket, if it exists:

```bash
tanzu plugin repo delete tce
```

Then add in the staging repo and perform the upgrade:

```bash
tanzu plugin repo add --name tce --gcp-bucket-name tce-cli-plugins-staging --gcp-root-path artifacts
tanzu plugin upgrade <plugin name>
```

That should cover the “how” for a community member and contributor looking to interact with TCE to fulfill their use-cases.

## The Release Process

I started with how a developer will build and interact with the TCE repo because we take that information and leverage it. CI automation is complex because it can be complicated because of all the moving pieces. Below is a block diagram that contains a more detailed view of all of the component/system interactions:

![release process](/docs/img/publicrelease.png)

The starting point is:

1. a human, actually someone responsible for the Software Lifecycle for the project (ie the tech or community lead), pushes either a Release Candidate (RC - vX.Y.Z-rc.B) or General Availability (GA - vX.Y.Z) tag to the repo to trigger the automation. This all takes place on the main TCE repo: [https://github.com/vmware-tanzu/tce/](https://github.com/vmware-tanzu/tce/)
2. A GitHub Action then runs a full build via make release (the same process a developer would do) to verify we can build everything end-to-end. This GitHub Action does the following:
  i. If successful, the unsigned binaries get created
  ii. A draft release is created on the Releases Page
  iii. The unsigned binaries are uploaded to the draft release as verification the build was successful. These binaries will be removed upon publishing the release.
3. A trigger is invoked to sign the TCE binaries using VMware's signing keys. This is done behind the VMware VPN/Firewall. Those binaries are attached to the draft release.
4. Then a human should delete the unsigned tarballs attached to the draft release, fill out the release notes and info on the draft, and then hit publish to make the release go live!

## What Constitutes a TCE release

Now that we know how the release process works, we should define what a TCE release is. TCE should be thought of as a flavor or distribution of a software package similar to how Kubernetes distros work. The TCE release should contain all the bits needed to run TCE in an airgapped fashion and allow the user to easily place the assets to their desired location through some form of installation.

The current layout of the TCE release, which is contained in a tarball, looks like the following:

![tarball layout](/docs/img/tarballlayout.png)

Ideally, if Tanzu CLI/Core really is thought of as a quick, nimble, open source platform to be consumed like Kubernetes, we should be able to leverage previously built versions of Tanzu CLI/Core instead of having to rebuild and re-sign them.
