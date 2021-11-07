LDFLAGS		+= -s -w

buildx-image:
	docker buildx create --use && \
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t quay.io/synpse/mini-shim --push -f Dockerfile .