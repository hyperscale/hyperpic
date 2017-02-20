FROM alpine:latest
MAINTAINER Axel Etcheverry <axel@etcheverry.biz>
ENV PORT 8080
RUN apk add --update curl && rm -rf /var/cache/apk/*
HEALTHCHECK --interval=1m --timeout=3s CMD curl -f http://localhost:${PORT}/health > /dev/null 2>&1 || exit 1
EXPOSE ${PORT}
VOLUME /var/lib/image-service
ADD image-service /opt/image-service/
ENTRYPOINT ["/opt/image-service/image-service"]
