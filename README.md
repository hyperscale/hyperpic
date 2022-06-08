Hyperpic [![Last release](https://img.shields.io/github/release/hyperscale/hyperpic.svg)](https://github.com/hyperscale/hyperpic/releases/latest) [![](https://img.shields.io/docker/pulls/hyperscale/hyperpic.svg)](https://hub.docker.com/r/hyperscale/hyperpic)
========

![Reactive logo](https://cdn.rawgit.com/hyperscale/hyperpic/master/_resources/hyperpic.svg "Hyperpic logo")


[![Go Report Card](https://goreportcard.com/badge/github.com/hyperscale/hyperpic)](https://goreportcard.com/report/github.com/hyperscale/hyperpic)

| Branch  | Status | Coverage | Docker |
|---------|--------|----------|--------|
| master  | [![Go](https://github.com/hyperscale/hyperpic/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/hyperscale/hyperpic/actions/workflows/go.yml) | [![Coveralls](https://img.shields.io/coveralls/hyperscale/hyperpic/master.svg)](https://coveralls.io/github/hyperscale/hyperpic?branch=master) | [![](https://img.shields.io/microbadger/image-size/hyperscale/hyperpic/latest.svg)](https://hub.docker.com/r/hyperscale/hyperpic) |
| develop | [![Go](https://github.com/hyperscale/hyperpic/actions/workflows/go.yml/badge.svg)](https://github.com/hyperscale/hyperpic/actions/workflows/go.yml) | [![Coveralls](https://img.shields.io/coveralls/hyperscale/hyperpic/develop.svg)](https://coveralls.io/github/hyperscale/hyperpic?branch=develop) | [![](https://img.shields.io/microbadger/image-size/hyperscale/hyperpic/dev.svg)](https://hub.docker.com/r/hyperscale/hyperpic) |

Fast HTTP microservice for high-level image processing.

[![Run on Google Cloud](https://storage.googleapis.com/cloudrun/button.svg)](https://deploy.cloud.run)

Install
-------

### Docker

```shell
docker pull hyperscale/hyperpic
```

### MacOS

Install libvips:
```shell
brew install homebrew/science/vips --with-webp
```

Build hyperpic:
```shell
make build
```

Run hyperpic
```shell
./hyperpic
```

Example
-------

### Crop

Original: `https://hyperpic-euskadi31.koyeb.app/kayaks.jpg`

![Original](https://hyperpic-euskadi31.koyeb.app/kayaks.jpg)

Croped and Resized: `https://hyperpic-euskadi31.koyeb.app/kayaks.jpg?w=400&h=400&fit=crop`

![Croped and resized](https://hyperpic-euskadi31.koyeb.app/kayaks.jpg?w=400&h=400&fit=crop)

### Crop on focal point

Original: `https://hyperpic-euskadi31.koyeb.app/smartcrop.jpg`

![Original](https://hyperpic-euskadi31.koyeb.app/smartcrop.jpg)

Croped and Resized: `https://hyperpic-euskadi31.koyeb.app/smartcrop.jpg?w=200&h=200&fit=crop-focal-point`

![Croped and resized](https://hyperpic-euskadi31.koyeb.app/smartcrop.jpg?w=200&h=200&fit=crop-focal-point)

Documentation
-------------

[Hyperpic API Reference](https://hyperscale.github.io/hyperpic/)

License
-------

hyperpic is licensed under [the MIT license](LICENSE.md).
