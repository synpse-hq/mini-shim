FROM golang:buster AS go-build-env

WORKDIR /app

COPY . .
RUN go mod download

ARG TARGETARCH
RUN CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /app/app ./cmd/mini-shim

FROM quay.io/synpse/alpine:3.9
RUN apk --update add ca-certificates

COPY --from=go-build-env /app/app /bin/
RUN chmod +x /bin/app

ENTRYPOINT ["/bin/app"]