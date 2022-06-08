FROM alpine:latest as builder
ENV GO111MODULE on
ENV GOPATH /go
RUN apk --no-cache --update-cache --force-overwrite --update \
    add vips-dev go make git libc6-compat build-base fftw-dev
WORKDIR /go/src/github.com/hyperscale/hyperpic/
COPY . .
RUN make build/hyperpic

FROM alpine:latest
ENV PORT 8080
HEALTHCHECK --interval=10s --timeout=3s CMD curl -f http://localhost:${PORT}/health > /dev/null 2>&1 || exit 1
EXPOSE ${PORT}
VOLUME /var/lib/hyperpic
RUN apk --no-cache --update-cache --force-overwrite --update \
    add ca-certificates curl vips expat libc6-compat
WORKDIR /root/
COPY --from=builder /go/src/github.com/hyperscale/hyperpic/build/hyperpic .
COPY --from=builder /go/src/github.com/hyperscale/hyperpic/cmd/hyperpic/config.yml.dist /etc/hyperpic/config.yml
CMD ["./hyperpic"]
