#
# This is part of Application Builder.
# https://github.com/golift/application-builder
#

FROM golang:1-buster as builder

# Build speed magic! 
# https://www.docker.com/blog/faster-multi-platform-builds-dockerfile-cross-compilation-guide/
WORKDIR /src
COPY go.mod go.sum .
RUN go mod download
COPY . .

ARG TARGETOS
ARG BINARY
ARG TARGETARCH

# Build the app.
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=1 make clean ${BINARY}.${TARGETARCH}.${TARGETOS}

FROM scratch
ARG TARGETOS
ARG TARGETARCH
ARG BUILD_DATE
ARG COMMIT
ARG VERSION
ARG LICENSE=MIT
ARG BINARY
ARG SOURCE_URL=http://github.com/golift/application-builder
ARG DESC=application-builder
ARG VENDOR=golift
ARG AUTHOR=golift

# Build-time metadata as defined at https://github.com/opencontainers/image-spec/blob/master/annotations.md
LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="${BINARY}" \
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

COPY --from=builder /src/${BINARY}.${TARGETARCH}.${TARGETOS} /image


# Make sure we have an ssl cert chain and timezone data.
COPY --from=builder /etc/ssl /etc/ssl
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

ENV TZ=UTC

EXPOSE 8080
ENTRYPOINT [ "/image" ]