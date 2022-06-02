#!/bin/bash

VIPS_VERSION=$1

if [ ! -d "$HOME/libvips" ]; then
    wget https://github.com/libvips/libvips/archive/refs/tags/v${VIPS_VERSION}.zip
    unzip v${VIPS_VERSION}
    cd libvips-${VIPS_VERSION}
    test -f autogen.sh && ./autogen.sh || ./bootstrap.sh
    CXXFLAGS=-D_GLIBCXX_USE_CXX11_ABI=0
    ./configure \
        --disable-debug \
        --disable-dependency-tracking \
        --disable-introspection \
        --disable-static \
        --enable-gtk-doc-html=no \
        --enable-gtk-doc=no \
        --enable-pyvips8=no \
        --without-orc \
        --without-python \
        --prefix=$HOME/libvips
    make
    sudo make install
    sudo ldconfig

    cd ..
fi

echo "PATH=$PATH:$HOME/libvips/bin" >> $GITHUB_ENV
echo "PKG_CONFIG_PATH=$PKG_CONFIG_PATH:$HOME/libvips/lib/pkgconfig" >> $GITHUB_ENV
echo "LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$HOME/libvips/lib" >> $GITHUB_ENV

$HOME/libvips/bin/vips --vips-version
