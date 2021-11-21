FROM golang:buster AS go-build-env
WORKDIR /app

RUN apt-get update
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get install -y ca-certificates

COPY . .
ENV GOPROXY="direct"
RUN  go mod download

ARG TARGETARCH
RUN GOARCH=$TARGETARCH go build -o /release/mini-shim ./cmd/mini-shim

FROM quay.io/synpse/alpine:3.9
RUN apk --update add ca-certificates

COPY --from=go-build-env /release/mini-shim /mini-shim
ENTRYPOINT ["/mini-shim"]