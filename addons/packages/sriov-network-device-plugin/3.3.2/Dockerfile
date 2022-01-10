FROM golang:1.13 as builder

# Get source
RUN git clone -b v3.3.2 https://github.com/k8snetworkplumbingwg/sriov-network-device-plugin.git 
WORKDIR /go/sriov-network-device-plugin
RUN make build

# Build related tools
WORKDIR /go/sriov-network-device-plugin/images
RUN tar -xvf ddptool-1.0.1.12.tar.gz -C . && \
    make

FROM debian:bullseye-slim

# Prepare environment
RUN apt-get update -y && \
    apt-get install -y hwdata && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /

# Copy sources from builder
COPY --from=builder /go/sriov-network-device-plugin/ /tmp/sriov-network-device-plugin/

# Now let us make it runable and put them under PATH
RUN chmod a+x /tmp/sriov-network-device-plugin/build/sriovdp && \
    chmod a+x /tmp/sriov-network-device-plugin/images/ddptool && \
    chmod a+x /tmp/sriov-network-device-plugin/images/entrypoint.sh && \
    mv /tmp/sriov-network-device-plugin/build/sriovdp /usr/bin/ && \
    mv /tmp/sriov-network-device-plugin/images/ddptool /usr/bin/ && \
    mv /tmp/sriov-network-device-plugin/images/entrypoint.sh /

## Okay now at least delete the source codes
RUN rm -rf /tmp/sriov-network-device-plugin

LABEL io.k8s.display-name="SRIOV Network Device Plugin"

# Same entrypoint as community
ENTRYPOINT ["/entrypoint.sh"]


