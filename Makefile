LDFLAGS		+= -s -w

fetch-certs:
	curl --remote-name --time-cond cacert.pem https://curl.se/ca/cacert.pem
	cp cacert.pem ca-certificates.crt

build-arm: fetch-certs
	CGO_ENABLED=0 GOARCH=arm GOOS=linux go build \
			-ldflags "$(LDFLAGS)" \
			-o release/arm/mini-shim ./cmd/mini-shim

docker-arm: build-arm	
	docker build -t quay.io/synpse/mini-shim -f Dockerfile.armhf .
	docker push quay.io/synpse/mini-shim:latest