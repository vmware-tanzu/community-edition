
FROM golang:1.13 as builder

# Get source
RUN git clone -b v2.6.1 https://github.com/k8snetworkplumbingwg/sriov-cni.git 
WORKDIR /go/sriov-cni
RUN apt-get update -y && \
    apt-get install -y apt-utils build-essential && \
    make clean && \
    make build

# Prepare environment
FROM debian:bullseye-slim
WORKDIR /

# Copy sources from builder
COPY --from=builder /go/sriov-cni/ /tmp/sriov-cni/

RUN chmod a+x /tmp/sriov-cni/build/sriov && \
    chmod a+x /tmp/sriov-cni/images/entrypoint.sh && \
    mv /tmp/sriov-cni/build/sriov /usr/bin/ && \
    mv /tmp/sriov-cni/images/entrypoint.sh /

LABEL io.k8s.display-name="SR-IOV CNI"

# Same entrypoint as community
ENTRYPOINT ["/entrypoint.sh"]