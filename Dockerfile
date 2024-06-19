# golang version 1 = latest 1.x.x
ARG GOLANG_VERSION="1"

###############################################################################
FROM golang:${GOLANG_VERSION} as builder

COPY ./ /go/build
WORKDIR /go/build/

ARG APP=observer
RUN set -eux  \
    && mkdir -p /root-out/  \
    && cp ./cmd/${APP}/.env /root-out/.env

ARG TARGETOS
ARG TARGETARCH
RUN set -eux  \
    && go mod download -x  \
    && go test  \
    && CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build  \
        -tags=notrace  \
        -ldflags="-s -w"  \
        -o="/root-out/main"  \
        ./cmd/${APP}/...

###############################################################################
# create actual image
FROM gcr.io/distroless/static
COPY --from=builder --chown=nonroot:nonroot /root-out/ /

EXPOSE 2512
USER nonroot
ENTRYPOINT ["/main"]
