---
title: "Tanzu Community Edition 0.10 is Now Available"
slug: tanzu-community-edition-0.10-available
date: 2022-02-17
author: Roger Klorese
image: /img/laptop-user.jpg
excerpt: "The Tanzu Community Edition team is pleased to announce the availability of Tanzu Community Edition 0.10. The new version
enhances usability on your local system, makes it easier to customize
your platform, and simplifies IP address management across all nodes in
a cluster."
tags: ['Roger Klorese']
---

The Tanzu Community Edition team is pleased to announce the availability
of [Tanzu Community Edition 0.10](https://tanzucommunityedition.io/download/). The new version
enhances usability on your local system, makes it easier to customize
your platform, and simplifies IP address management across all nodes in
a cluster. Today's release follows the recent [addition of build automation](https://tanzucommunityedition.io/posts/tanzu-community-edition-adds-kpack/)
capabilities into Tanzu Community Edition as the team continues to
deliver the power of Tanzu to you in a free, open source distribution.

Here is a quick summary of the highlights of the new release:

- The 0.10 release introduces a new 'unmanaged' cluster type that
    significantly improves the efficiency of Tanzu Community Edition
    clusters in most local environments. Offered today in alpha stage,
    unmanaged clusters will replace the existing 'standalone' cluster
    type, cutting cluster setup time by half or more, shortening
    tear-down time to seconds, and significantly reducing overall
    resource requirements. Depending on network configuration, they can
    also survive system reboots. Unmanaged clusters are ideal for use in
    your local development environment; in fact, they're a great
    solution whenever you're working locally with more limited
    resources, need only one cluster at a time, and especially if you're
    frequently creating and destroying those clusters. Standalone
    clusters are also included in this release but are now deprecated
    and will be removed from Tanzu Community Edition in the [0.11 release](https://github.com/vmware-tanzu/community-edition/milestone/12). We've posted a [cool demo](https://youtu.be/VH6OWtxvzbM) for it.

- Responding to users, the 0.10 release also provides requested
    platform management enhancements (keep those requests coming!). You
    can now configure most packages more flexibly with expanded YTT
    template parameters that give you more options to tweak when
    packages are installed. And the addition of the optional
    Whereabouts IP Address Management (IPAM) CNI plugin provides a Pod
    IP space that is not split on a per-node basis. This supports an IP
    assignment pattern that is consistent with more traditional
    enterprise networking models, making it easier to integrate clusters
    into environments that use those long-standing models.

- Continuous improvement enhancements are offered as part of the 0.10
    release as well. Among them are updates to documentation, to
    underlying Tanzu components, and to many of the packages included in
    the distribution, as well as expanded test coverage.

For more details about the release and changelog, check out the Tanzu
Community Edition 0.10 [release notes](https://github.com/vmware-tanzu/community-edition/releases/tag/v0.10.0).

## Get involved: the community welcomes you

It is wonderful to see our community growing -- thank you to all who
have been providing ideas, feedback, and other contributions! There is
much more to do, and lots of opportunity for us to help each other:

- If you haven't done so already at the end of the installation
    process, please take a few minutes to complete a [short survey](https://tb3xduryx4x.typeform.com/to/RVmhMHwR) about
    your experience installing Tanzu Community Edition.

- Looking for help? Please send your questions to the [Tanzu Community Edition Google Group](https://groups.google.com/g/tanzu-community-edition?pli=1).
    We'll respond to you there so that others with the same question can
    more easily find the answer.

- You can also get help at [Office Hours](https://tanzucommunityedition.io/community/). Take advantage
    of an open floor with live participants to discuss what's on your
    mind and find the expert support you need.

- Learn more about the technology and your community at [Community Meetings](https://tanzucommunityedition.io/community/). We review
    upcoming roadmap and development work, dive deep into technical
    discussions, and generally enjoy spending time together

- Our [Slack](https://kubernetes.slack.com/archives/C02GY94A8KT)
    channel is open to anyone interested in watching or participating in
    up-to-the-minute community conversation. It's a great place to
    engage directly with the development team

- Make your wishes known! [Open an issue](https://github.com/vmware-tanzu/community-edition/issues/new/choose)
    to file an enhancement request or report a defect.

Enjoy the [new release](https://tanzucommunityedition.io/download/)!
