---
title: "VMware Tanzu Community Edition Deepens Appeal to Developers with Container Build Automation"
slug: tanzu-community-edition-adds-kpack
date: 2022-01-25
author: Roger Klorese
image: /img/path-to-production.png
excerpt: "VMware Tanzu Community Edition is today expanding its support for cloud native application development and deployment with the addition of
kpack, a Kubernetes-native automated container build platform that saves time for developers, while also enabling operations teams to deliver containers more securely, reliably,
and consistently at scale."
tags: ['Roger Klorese']
---

[VMware Tanzu Community Edition](https://tanzucommunityedition.io/) is
today expanding its support for cloud native application development and
deployment with the addition of
[kpack](https://github.com/pivotal/kpack), a Kubernetes-native automated
container build platform that saves time for developers, while also
enabling operations teams to deliver containers more securely, reliably,
and consistently at scale. Kpack joins [Harbor](https://goharbor.io/),
[Knative](https://knative.dev/docs/serving/), and other open source
components of Tanzu Community Edition to further accelerate delivery of
containerized workloads.

Today's addition marks an important step toward simplifying application
development and deployment for users of Tanzu Community Edition, the
freely available open source distribution of VMware Tanzu that can be
installed and configured in minutes on your local workstation or
favorite cloud. Already delivering capabilities of the VMware Tanzu
portfolio that target platform operators as an open source platform,
Tanzu Community Edition is now adding technologies that underpin [VMware
Tanzu Application Platform](https://tanzu.vmware.com/application-platform) to support
application development teams. In time, Tanzu Community Edition will
provide an automated and streamlined path to production, offering much
of the commercial platform's technology base and serving as the upstream
source for core capabilities of Tanzu Application Platform.

## Kpack build automation serves developers and operators

Kpack is a build automation platform that implements the powerful [Cloud
Native Buildpacks](https://buildpacks.io/) specification. Driving the
automation process, kpack leverages language-specific buildpacks, such
as those available from the open source [Paketo](https://paketo.io/)
project, to transform source code into container images.

By using kpack to manage the container build process, developers can
save time, eliminate worry, and focus on writing application code. Kpack
builds OCI-compliant containers directly from source code, eliminating
the need for developers to build and maintain dockerfiles. It analyzes
the code to identify dependencies, and builds containers using approved
versions of language libraries and operating systems (OSs), ensuring
that internal compliance standards are satisfied. Further enhancing
productivity, kpack also detects changes to source code, dependency, or
OS components and automatically updates containers, using software that
is approved by the developer's organization.

!["Four Phases of Build Automation: 1. Detect, 2. Restore and Analyze, 3. Build, 4. Export and Cache"](/img/four-phases.png)

Application operators and DevOps teams benefit too. Automated container
re-builds update images rapidly and efficiently as container components
are updated, saving time while also enabling operations teams to reduce
the risks posed by Common Vulnerabilities and Exposures (CVEs).
Centralized management of container metadata -- including information
about approved software versions for example -- makes policy enforcement
easier so that ops teams can make informed image promotion and
deployment decisions at scale. In addition, as part of the kpack build
process, bill of materials (BOM) metadata is surfaced for every
container built, enabling full-stack auditing, tracking, and container
patching.

Learn more about using kpack, the Buildpacks model for build automation,
and open source Paketo buildpacks [here](https://tanzu.vmware.com/content/blog/tanzu-community-edition-kpack-kubernetes-container-build-platform).

## Growing developer support in Tanzu Community Edition

With the addition of kpack, Tanzu Community Edition now provides a set
of open source tooling that supports the core steps of service delivery:
kpack turns source code into container images; the Harbor container
registry signs images, scans them for vulnerabilities (using the
embedded Trivy scanner), and provides secure registry services; Knative
Serving simplifies deployment and running of services.

!["Building a Path to Production in Tanzu Community Edition now including Open Source software components from Tanzu Application Platform"](/img/path-to-production.png)

Because all of these capabilities are packaged as part of the Tanzu
Community Edition platform, it is easy for platform operators to provide
them to their users. In Tanzu Community Edition, kpack, Harbor, and
Knative, as well as many other platform components, are all provided as
part of a curated, open source platform distribution that is easy to
install, can be flexibly configured, and consistently managed.

Reaching this stage of functional support for application development
and deployment marks an exciting milestone for Tanzu Community Edition.
And there is more to come. Looking ahead, the Tanzu Community Edition
project is expected to add supply chain choreography via the
Cartographer open source project, as well as additional
developer-focused capabilities that will enable developers to further
simplify, accelerate, and secure the creation and delivery of cloud
native workloads.

## Try it out now

Tanzu Community Edition is freely available, community supported, open
source software that you can download and use today. Find it, and a
welcoming community, at
[tanzucommunityedition.io](https://tanzucommunityedition.io/)
