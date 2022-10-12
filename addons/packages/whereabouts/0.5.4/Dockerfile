FROM golang:1.16
RUN mkdir -p $GOPATH/src/github.com/k8snetworkplumbingwg/whereabouts
WORKDIR $GOPATH/src/github.com/k8snetworkplumbingwg/whereabouts
RUN git clone -b v0.5.4 --single-branch https://github.com/k8snetworkplumbingwg/whereabouts.git .
RUN CGO_ENABLED=0 go build -o bin/whereabouts cmd/whereabouts.go
RUN CGO_ENABLED=0 go build -o bin/ip-reconciler cmd/reconciler/*.go
RUN CGO_ENABLED=0 go build -o bin/ip-control-loop cmd/controlloop/*.go

# Use distroless base image with debug tag, as we need to run install-cni.sh shell scripts here.
FROM gcr.io/distroless/base:debug
COPY --from=0 /go/src/github.com/k8snetworkplumbingwg/whereabouts/bin/whereabouts .
COPY --from=0 /go/src/github.com/k8snetworkplumbingwg/whereabouts/bin/ip-reconciler .
COPY --from=0 /go/src/github.com/k8snetworkplumbingwg/whereabouts/script/install-cni.sh .
COPY --from=0 /go/src/github.com/k8snetworkplumbingwg/whereabouts/bin/ip-control-loop .
CMD ["/install-cni.sh"]
