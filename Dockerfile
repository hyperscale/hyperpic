FROM alpine:edge as builder
ENV GOPATH /go
RUN echo "http://nl.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories
RUN apk --no-cache --update add vips-dev go make git musl-dev fftw-dev
WORKDIR /go/src/github.com/hyperscale/hyperpic/
RUN go get -u github.com/golang/dep/cmd/dep
COPY . .
RUN /go/bin/dep ensure
RUN make build

FROM alpine:edge
LABEL maintainer "axel@etcheverry.biz"
ENV PORT 8080
HEALTHCHECK --interval=1m --timeout=3s CMD curl -f http://localhost:${PORT}/health > /dev/null 2>&1 || exit 1
EXPOSE ${PORT}
VOLUME /var/lib/hyperpic
RUN echo "http://nl.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories
RUN apk --no-cache --update add ca-certificates curl vips
WORKDIR /root/
COPY --from=builder /go/src/github.com/hyperscale/hyperpic/hyperpic .
COPY --from=builder /go/src/github.com/hyperscale/hyperpic/config.yml.dist /etc/hyperpic/config.yml
CMD ["./hyperpic"]
