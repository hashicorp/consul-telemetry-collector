# ===================================
#
#   Non-release images.
#
# ===================================


# devbuild compiles the binary
# -----------------------------------
FROM golang:latest AS devbuild
ARG BIN_NAME
# Escape the GOPATH
WORKDIR /build
COPY . ./
RUN make dev


# ===================================
#
#   Release images.
#
# ===================================


# default release image
# -----------------------------------
FROM alpine:latest AS release-default

ARG BIN_NAME
# Export BIN_NAME for the CMD below, it can't see ARGs directly.
ENV BIN_NAME=$BIN_NAME
ARG PRODUCT_VERSION
ARG PRODUCT_REVISION
ARG PRODUCT_NAME=$BIN_NAME
# TARGETARCH and TARGETOS are set automatically when --platform is provided.
ARG TARGETOS TARGETARCH

LABEL maintainer="Team RelEng <team-rel-eng@hashicorp.com>"
LABEL version=$PRODUCT_VERSION
LABEL revision=$PRODUCT_REVISION

# Create a non-root user to run the software.
RUN addgroup $PRODUCT_NAME && \
    adduser -S -G $PRODUCT_NAME 100

COPY dist/$TARGETOS/$TARGETARCH/$BIN_NAME /bin/

USER 100
COPY .github/docker/entrypoint.sh /usr/bin/entrypoint.sh
ENTRYPOINT [ "entrypoint.sh" ]

# dev runs the binary from devbuild
# -----------------------------------
FROM alpine:latest AS dev
ARG BIN_NAME
# Export BIN_NAME for the CMD below, it can't see ARGs directly.
ENV BIN_NAME=${BIN_NAME}
COPY --from=devbuild /build/${BIN_NAME} /bin/
COPY .github/docker/entrypoint.sh /usr/bin/entrypoint.sh
ENTRYPOINT [ "entrypoint.sh" ]

# alternate release image, just for the sake of example. In this case we're using
# debian as the base image just to make the image different from the default alpine one.
#
# The use cases for alternate images are things like defining an additional UBI compatible
# image, or an image with a different function than the main image, e.g. Waypoint's ODR images.
# -----------------------------------
FROM debian:latest AS release-alternate

ARG BIN_NAME
# Export BIN_NAME for the CMD below, it can't see ARGs directly.
ENV BIN_NAME=$BIN_NAME
ARG PRODUCT_VERSION
ARG PRODUCT_REVISION
ARG PRODUCT_NAME=$BIN_NAME
# TARGETARCH and TARGETOS are set automatically when --platform is provided.
ARG TARGETOS TARGETARCH

LABEL maintainer="Team RelEng <team-rel-eng@hashicorp.com>"
LABEL version=$PRODUCT_VERSION
LABEL revision=$PRODUCT_REVISION

# Create a non-root user to run the software.
RUN addgroup $PRODUCT_NAME && \
    adduser --system --uid 101 --group $PRODUCT_NAME

COPY dist/$TARGETOS/$TARGETARCH/$BIN_NAME /bin/

USER 101
COPY .github/docker/entrypoint.sh /usr/local/bin/entrypoint.sh
ENTRYPOINT [ "entrypoint.sh" ]


# Red Hat UBI-based image
# This image is based on the Red Hat UBI base image, and has the necessary
# labels, license file, and non-root user.
# -----------------------------------
FROM registry.access.redhat.com/ubi8/ubi:latest as release-ubi

ARG BIN_NAME
# Export BIN_NAME for the CMD below, it can't see ARGs directly.
ENV BIN_NAME=$BIN_NAME
ARG PRODUCT_VERSION
ARG PRODUCT_REVISION
ARG PRODUCT_NAME=$BIN_NAME
# TARGETARCH and TARGETOS are set automatically when --platform is provided.
ARG TARGETOS TARGETARCH

LABEL name="CRT Core Hello World" \
	  maintainer="Team RelEng <team-rel-eng@hashicorp.com>" \
      vendor="HashiCorp" \
      version=${PRODUCT_VERSION} \
      release=${PRODUCT_REVISION} \
      revision=${PRODUCT_REVISION} \
      summary="CRT Core Hello World is a demo project." \
      description="Example repository demonstrating CRT."

# Create a non-root user to run the software.
RUN groupadd --gid 1000 $PRODUCT_NAME && \
    adduser --uid 100 --system -g $PRODUCT_NAME $PRODUCT_NAME && \
    usermod -a -G root $PRODUCT_NAME

COPY dist/$TARGETOS/$TARGETARCH/$BIN_NAME /bin/
#COPY LICENSE /licenses/copyright.txt

USER 101
CMD /bin/$BIN_NAME


# ===================================
#
#   Set default target to 'dev'.
#
# ===================================
FROM dev
