---
title: "Advanced Supply Chain Choreography Now Included in VMware Tanzu® Community Edition"
slug: tanzu-community-edition-011-adds-cartographer-supply-chain-choreography
date: 2022-03-29
author: Kartik Lunkad
image: /img/picture-1.jpg
excerpt: "VMware Tanzu Community Edition now includes advanced software supply chain tooling that helps application teams deliver software more rapidly, securely, and efficiently at scale. The 0.11 release of Tanzu Community Edition, available today, introduces new supply chain choreography capabilities powered by the open source Cartographer project."
tags: ['Kartik Lunkad']
---
!["An abstract image of a computer chip with signal traces coming out of it, accompanied by a cloud and a shield."](/img/picture-1.jpg)

VMware Tanzu Community Edition now includes advanced software supply chain tooling that helps application teams deliver software more rapidly, securely, and efficiently at scale. The 0.11 release of Tanzu Community Edition, available today, introduces new supply chain choreography capabilities powered by the open source [Cartographer](https://cartographer.sh/) project. Application teams use Cartographer to create, evolve, and manage modern software supply chains more easily, and to coordinate flows through those supply chains more flexibly. The release also adds the [Flux CD Source Controller](https://github.com/fluxcd/source-controller) that can be used to monitor source repositories for code changes and make those updates available to supply chains.

## Support for application development teams in Tanzu Community Edition

Today’s release delivers premier functionality to Tanzu Community Edition, the freely available open source distribution of VMware Tanzu that can be installed and configured in minutes on your local workstation or favorite cloud. Already providing capabilities of the VMware Tanzu portfolio that target platform operators, Tanzu Community Edition has been incorporating open source technologies that underpin VMware Tanzu Application Platform to support application development teams. With these latest additions, the community distribution now provides a curated set of open source tools that application teams can use to fully automate and streamline their paths to production. These tools include:

- [Flux CD Source Controller](https://github.com/fluxcd/source-controller) for detecting and delivering updated code

- [Kpack](https://github.com/pivotal/kpack), working with [Paketo Buildpacks](https://paketo.io/), for turning source code into container images

- the [Harbor](https://goharbor.io/) container registry for signing images, scanning them for vulnerabilities (using the embedded [Trivy](https://github.com/aquasecurity/trivy) scanner), and providing secure registry services

- [Knative Serving](https://knative.dev/) for simplifying deployment and running of services

- [Cartographer](https://cartographer.sh/) for building supply chains from these or other components, and choreographing flows across them

## Modern supply chain automation with Cartographer

Cartographer can help application operations teams advance beyond the limits of today’s CI/CD pipelines to deliver more robust, versatile, and adaptable software supply chains that take advantage of a modern supply chain architecture and flow model. Use it to:

- *Improve consistency, compliance, and maintainability of supply chains at scale* by building flows from approved, reusable components. Eliminate the time-consuming and error-prone processes involved in maintaining and extending the proliferation of snowflake pipeline configurations found in most organizations today. With Cartographer, you can instead reuse trusted components across many supply chains to drive increased implementation consistency and maintainability. With declarative deployments that are continuously reconciled, Cartographer can improve the availability, scalability, and fault tolerance of your supply chains too.

- *Enhance supply chain versatility* with event-driven, non-linear flows and asynchronous processing. Eliminate awkward workarounds to initiate supply chain execution that is not naturally driven by code commits; for example, the process of updating images with a new operating system patch can be triggered directly by the availability of the patch, rather than by an artificial code commit. Cartographer also enables you to speed supply chain execution by performing steps asynchronously, and in parallel where appropriate.

- *Adapt more quickly* with loosely coupled supply chains made from independent components that can be updated on the fly. Replace complex, tightly coupled, and brittle pipeline monoliths with flexible, modular supply chains comprising loosely coupled components that you can update independently. Cartographer supply chains are cloud native, so you can introduce updates and evolve your supply chains without impacting their availability.

Software development, security, and operations teams can also take advantage of Cartographer’s clear separation of responsibilities to decrease frustration, work better together, and increase output. Cartographer helps:

- *Boost team productivity* by enabling members to focus their time and efforts in their respective areas of expertise, where they can provide greatest value. Developers can focus on setting application configuration parameters and can use supply chains effectively without deep knowledge of Kubernetes or mastery of operations best practices. Operations teams can create supply chains using predefined tool sets, or allow developers to choose their own preferred tools.

- *Improve security* by enabling teams to control access to the supply chain appropriately. Operations teams can prevent developers from overriding established operational policies and best practices, and they can ensure that compliant tools are used and consistent processes are applied to all applications that move through the supply chain.  

<!-- markdownlint-disable MD026 -->
## Try Tanzu Community Edition now!

Tanzu Community Edition is freely available, community supported, open source software that you can download and use today. Find it, and a welcoming community, at [tanzucommunityedition.io](https://tanzucommunityedition.io/).

<!-- markdownlint-disable MD033 -->
<sub>This article may contain hyperlinks to non-VMware websites that are created and maintained by third parties who are solely responsible for the content on such websites.</sub>
