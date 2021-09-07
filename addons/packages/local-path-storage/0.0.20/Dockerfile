FROM golang:1.16 as builder

RUN go install github.com/rancher/local-path-provisioner@v0.0.20

FROM gcr.io/distroless/base

COPY --from=builder /go/bin/local-path-provisioner /bin/local-path-provisioner

ENTRYPOINT ["local-path-provisioner"]
