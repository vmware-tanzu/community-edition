# syntax=docker/dockerfile:1
FROM golang:1.17.6 AS build-env
WORKDIR /go/src/app
COPY . /go/src/app
RUN go build -o /tagger .

RUN go get -d -v ./...

FROM gcr.io/distroless/base
COPY --from=build-env /tagger /
CMD ["/tagger"]
