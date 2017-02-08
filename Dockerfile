FROM docker:1.12-dind
MAINTAINER Axel Etcheverry <axel@etcheverry.biz>
ENV PORT 8080
ENV DEBUG true
HEALTHCHECK --interval=1m --timeout=3s CMD curl -f http://localhost:${PORT}/health > /dev/null 2>&1 || exit 1
EXPOSE ${PORT}
VOLUME /var/lib/docker-manager
ADD ui/dist/* /opt/docker-manager/ui/
ADD docker-manager /opt/docker-manager/
ENTRYPOINT ["/usr/local/bin/dockerd-entrypoint.sh", "/opt/docker-manager/docker-manager"]
