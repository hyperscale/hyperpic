FROM alpine:edge as builder
ARG VERSION
ARG VCS_URL
ARG VCS_REF
ARG BUILD_DATE
ENV GO111MODULE on
ENV GOPATH /go
RUN echo "http://nl.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories
RUN apk --no-cache --update add vips-dev go make git musl-dev fftw-dev
RUN go version
WORKDIR /go/src/github.com/hyperscale/hyperpic/
COPY . .
RUN go build -ldflags "-X github.com/hyperscale/hyperpic/version.Version=${VERSION} -X github.com/hyperscale/hyperpic/version.Revision=${VCS_REF} -X github.com/hyperscale/hyperpic/version.BuildAt=${BUILD_DATE}" ./cmd/hyperpic/

FROM alpine:edge
ARG VERSION
ARG VCS_URL
ARG VCS_REF
ARG BUILD_DATE
ENV PORT 8080
HEALTHCHECK --interval=10s --timeout=3s CMD curl -f http://localhost:${PORT}/health > /dev/null 2>&1 || exit 1
EXPOSE ${PORT}
VOLUME /var/lib/hyperpic
RUN echo "http://nl.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories
RUN apk --no-cache --update add ca-certificates curl vips
WORKDIR /root/
COPY --from=builder /go/src/github.com/hyperscale/hyperpic/hyperpic .
COPY --from=builder /go/src/github.com/hyperscale/hyperpic/config.yml.dist /etc/hyperpic/config.yml
CMD ["./hyperpic"]

# Metadata
LABEL org.label-schema.vendor="Hyperscale" \
      org.label-schema.url="https://github.com/hyperscale/hyperpic" \
      org.label-schema.name="Hyperpic" \
      org.label-schema.description="Fast HTTP microservice for high-level image processing." \
      org.label-schema.version="v${VERSION}" \
      org.label-schema.vcs-url=${VCS_URL} \
      org.label-schema.vcs-ref=${VCS_REF} \
      org.label-schema.build-date=${BUILD_DATE} \
      org.label-schema.docker.schema-version="1.0"
