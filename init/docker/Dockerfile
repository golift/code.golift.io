#
# This is part of Application Builder.
# https://github.com/golift/application-builder
#

FROM golang:1-bookworm as builder

# Build speed magic! 
# https://www.docker.com/blog/faster-multi-platform-builds-dockerfile-cross-compilation-guide/
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
COPY pkg pkg

ARG TARGETOS
ARG TARGETARCH
ARG VERSION

# Build the app.
RUN GOFLAGS="-trimpath -mod=readonly -modcacherw" \
 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /tmp/image -tags osusergo,netgo \
 -ldflags "-w -s  -X \"main.Version=${VERSION}\""

FROM scratch

# Make sure we have an ssl cert chain and timezone data.
COPY --from=builder /etc/ssl /etc/ssl
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=builder /tmp/image /image

ARG TARGETOS
ARG TARGETARCH
ARG BUILD_DATE
ARG COMMIT
ARG VERSION
ARG LICENSE=Apache-2
ARG SOURCE_URL=http://github.com/golift/turbovanityurls
ARG DESC=application-builder
ARG VENDOR=golift
ARG AUTHOR=golift

# Build-time metadata as defined at https://github.com/opencontainers/image-spec/blob/master/annotations.md
LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="turbovanityurls" \
      org.opencontainers.image.documentation="${SOURCE_URL}/wiki/Docker" \
      org.opencontainers.image.description="${DESC}" \
      org.opencontainers.image.url="${SOURCE_URL}" \
      org.opencontainers.image.revision="${COMMIT}" \
      org.opencontainers.image.source="${SOURCE_URL}" \
      org.opencontainers.image.vendor="${VENDOR}" \
      org.opencontainers.image.authors="${AUTHOR}" \
      org.opencontainers.image.architecture="${TARGETOS} ${TARGETARCH}" \
      org.opencontainers.image.licenses="${LICENSE}" \
      org.opencontainers.image.version="${VERSION}"

ENV TZ=UTC

EXPOSE 8080
ENTRYPOINT [ "/image" ]