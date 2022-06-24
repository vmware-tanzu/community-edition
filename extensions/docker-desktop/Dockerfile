#syntax=docker/dockerfile:1.4

# Copyright 2022 VMware, Inc.
# SPDX-License-Identifier: Apache-2.0

ARG REGISTRY=ghcr.io/
ARG ORG=vmware-tanzu/
ARG IMAGE_NAME=community-edition
ARG TAG=dev

ARG DOWNLOADER_IMAGE=${REGISTRY}${ORG}${IMAGE_NAME}-downloader:${TAG}
ARG TANZU_CLI_IMAGE=${REGISTRY}${ORG}${IMAGE_NAME}-tanzu-cli:${TAG}
ARG APPS_BIN_IMAGE=${REGISTRY}${ORG}${IMAGE_NAME}-apps-bin:${TAG}
ARG CLIENT_IMAGE=${REGISTRY}${ORG}${IMAGE_NAME}-client:${TAG}

#-----------------------------
FROM ${DOWNLOADER_IMAGE} as downloader
FROM ${TANZU_CLI_IMAGE} as tanzu-cli
FROM ${APPS_BIN_IMAGE}  as apps-bin
FROM ${CLIENT_IMAGE} as client-builder
# See latest tags at https://github.com/kubernetes/k8s.io/blob/main/k8s.gcr.io/images/k8s-staging-kubernetes/images.yaml
FROM k8s.gcr.io/pause:3.7 as pause


#-----------------------------
FROM bitnami/minideb:bullseye

ARG EXTENSION_NAME='VMware Tanzu Community Edition'
ARG DESCRIPTION='Spin up a local Kubernetes cluster that delivers VMware Tanzu integrated application management and tooling.'
ARG DESCRIPTION_LONG='Spin up a local Kubernetes cluster that delivers VMware Tanzu integrated application management and tooling.'
ARG VENDOR='VMware, Inc'
ARG LICENSE='Apache-2.0'
ARG ICON_URL='https://raw.githubusercontent.com/vmware-tanzu/community-edition/main/docs/images/logos/tce-logo-only.png'
ARG SCREENSHOTS_URLS='[ { "alt": "Intro", "url": "https://raw.githubusercontent.com/vmware-tanzu/community-edition/main/docs/images/ddextension-intro.png" }, { "alt": "Cluster view", "url": "https://raw.githubusercontent.com/vmware-tanzu/community-edition/main/docs/images/ddextension-cluster.png" }, { "alt": "Kubeconfig view", "url": "https://raw.githubusercontent.com/vmware-tanzu/community-edition/main/docs/images/ddextension-kubeconfig.png" } ]'
ARG PUBLISHER_URL='https://github.com/vmware-tanzu/community-edition'
ARG ADDITIONAL_URLS='[ { "title": "VMware Tanzu Community Edition", "url": "https://tanzucommunityedition.io/" }, { "title": "Documentation", "url": "https://tanzucommunityedition.io/docs/latest/" }, { "title": "Support", "url": "https://github.com/vmware-tanzu/community-edition/issues/new/choose" }, { "title": "Terms of Service", "url": "https://github.com/vmware-tanzu/community-edition/blob/main/LICENSE" }, { "title": "Privacy policy", "url": "https://www.vmware.com/help/privacy.html" }]'

ARG CHANGELOG='<p>Extension changelog:</p> <ul> <li>Initial release</li> </ul>'

ARG DD_VERSION='>=0.2.3'

# Note that some "org.opencontainers.xxx" labels are overridden in the publish gh action
# See github\workflows\build-publish.yml
LABEL org.opencontainers.image.title="${EXTENSION_NAME}" \
    org.opencontainers.image.description="${DESCRIPTION}"\
    org.opencontainers.image.vendor="${VENDOR}" \
    org.opencontainers.image.licenses="${LICENSE}" \
    com.docker.desktop.extension.icon="${ICON_URL}" \
    com.docker.desktop.extension.api.version="${DD_VERSION}" \
    com.docker.extension.screenshots="${SCREENSHOTS_URLS}" \
    com.docker.extension.detailed-description="${DESCRIPTION_LONG}" \
    com.docker.extension.publisher-url="${PUBLISHER_URL}" \
    com.docker.extension.additional-urls="${ADDITIONAL_URLS}" \
    com.docker.extension.changelog="${CHANGELOG}"

ARG TARGETARCH

RUN <<EOT
    apt-get update
    apt -y install ca-certificates
EOT

RUN useradd -u 1001 -g users -G root -M -d /home/tanzu tanzu && \
    mkdir -p /home/tanzu && \
    chown -R 1001:100 /home/tanzu

COPY --from=downloader /usr/bin/docker /usr/local/bin/docker
COPY --from=downloader /usr/bin/docker-compose /usr/local/lib/docker/cli-plugins/docker-compose
COPY --from=downloader /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --from=downloader /usr/local/bin/ytt /usr/local/bin/ytt

COPY --from=apps-bin /clustermgr-${TARGETARCH} /backend/clustermgr

COPY --from=client-builder /ui/dist ui

COPY --from=pause /pause /usr/local/bin/pause

COPY logo.svg .
COPY metadata.json .
COPY backend backend
COPY docker-compose.yaml .
COPY open_source_licenses.txt .

RUN mkdir /opt/kubeconfig && \
    chown 1001:100 /opt/kubeconfig

USER 1001

WORKDIR /home/tanzu

ENV PATH=/backend:$PATH

CMD ["/usr/local/bin/pause"]
