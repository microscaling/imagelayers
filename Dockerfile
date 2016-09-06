FROM alpine:3.4

MAINTAINER Ross Fairbanks <ross@microscaling.com>

ENV BUILD_PACKAGES ca-certificates

RUN apk update && \
    apk upgrade && \
    apk add $BUILD_PACKAGES && \
    rm -rf /var/cache/apk/*

EXPOSE 8888
COPY imagelayers Dockerfile /

# Metadata params
ARG BUILD_DATE
ARG VERSION
ARG VCS_REF

# Metadata
LABEL org.label-schema.url="https://imagelayers.io" \
      org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.version=$VERSION \
      org.label-schema.vcs-url="https://github.com/microscaling/imagelayers.git" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.docker.dockerfile="/Dockerfile" \
      org.label-schema.description="This utility provides a browser-based visualization of user-specified Docker Images and their layers." \
      org.label-schema.schema-version="1.0"

ENTRYPOINT ["/imagelayers"]
