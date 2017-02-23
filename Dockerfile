FROM golang:1.8-alpine
MAINTAINER Axel Etcheverry <axel@etcheverry.biz>
ENV PORT 8080
# Environment Variables
ARG LIBVIPS_VERSION_MAJOR_MINOR=8.4
ARG LIBVIPS_VERSION_PATCH=5
ARG MOZJPEG_VERSION="v3.1"

# Install dependencies
RUN echo "http://dl-cdn.alpinelinux.org/alpine/v3.5/community" >> /etc/apk/repositories && \
    apk update && \
    apk upgrade && \
    apk add \
    zlib libxml2 libxslt glib libexif lcms2 fftw ca-certificates curl git \
    giflib libpng libwebp orc tiff poppler-glib librsvg && \

    apk add --no-cache --virtual .build-dependencies autoconf automake build-base \
    git libtool nasm zlib-dev libxml2-dev libxslt-dev glib-dev \
    libexif-dev lcms2-dev fftw-dev giflib-dev libpng-dev libwebp-dev orc-dev tiff-dev \
    poppler-dev librsvg-dev && \

# Install mozjpeg
    cd /tmp && \
    git clone git://github.com/mozilla/mozjpeg.git && \
    cd /tmp/mozjpeg && \
    git checkout ${MOZJPEG_VERSION} && \
    autoreconf -fiv && ./configure --prefix=/usr && make install && \

# Install libvips
    wget -O- http://www.vips.ecs.soton.ac.uk/supported/${LIBVIPS_VERSION_MAJOR_MINOR}/vips-${LIBVIPS_VERSION_MAJOR_MINOR}.${LIBVIPS_VERSION_PATCH}.tar.gz | tar xzC /tmp && \
    cd /tmp/vips-${LIBVIPS_VERSION_MAJOR_MINOR}.${LIBVIPS_VERSION_PATCH} && \
    ./configure --prefix=/usr \
                --without-python \
                --without-gsf \
                --enable-debug=no \
                --disable-dependency-tracking \
                --disable-static \
                --enable-silent-rules && \
    make -s install-strip && \
    cd $OLDPWD && \

    go get -u github.com/euskadi31/image-service && \

# Cleanup
    rm -rf /tmp/vips-${LIBVIPS_VERSION_MAJOR_MINOR}.${LIBVIPS_VERSION_PATCH} && \
    rm -rf /tmp/mozjpeg && \
    apk del --purge .build-dependencies && \
    rm -rf /var/cache/apk/*

HEALTHCHECK --interval=1m --timeout=3s CMD curl -f http://localhost:${PORT}/health > /dev/null 2>&1 || exit 1
EXPOSE ${PORT}
VOLUME /var/lib/image-service
ADD https://raw.githubusercontent.com/euskadi31/image-service/master/config.yml.dist /etc/image-service/config.yml
CMD ["/go/bin/image-service"]
