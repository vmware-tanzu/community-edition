FROM projects-stg.registry.vmware.com/tkg/bazel:3.5.0
ARG KUBECTL_VERSION=v1.21.0
ARG CERT_MANAGER_VERSION=v1.1.0

RUN apt-get update -y && apt-get install vim jq git -y
RUN mkdir -p /workspace
RUN git clone https://github.com/jetstack/cert-manager.git workspace
WORKDIR /workspace
RUN git checkout $CERT_MANAGER_VERSION

# Copy
COPY patch-registry.sh patch-registry.sh
#Replace modified scripts
COPY create-kind.sh devel/cluster/create-kind.sh
COPY setup-e2e-deps.sh devel/setup-e2e-deps.sh

RUN curl -sL https://dl.k8s.io/release/$KUBECTL_VERSION/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl
RUN chmod +x /usr/local/bin/kubectl

#TODO goes away with packages
RUN curl -L -o cert-manager.yaml https://gist.githubusercontent.com/vijaykatam/0fae1dda55a460fe6e14258d7676a771/raw/08bc7e63bd6beec42af4088050887d3e8d3158c3/gistfile1.txt
ENTRYPOINT /workspace/devel/ci-run-e2e.sh --ginkgo.focus="Certificates|CA Certificate|Self Signed Certificate|RBAC" --ginkgo.skip=ACME