Hyperpic ![Last release](https://img.shields.io/github/release/hyperscale/hyperpic.svg) 
========

![Reactive logo](https://cdn.rawgit.com/hyperscale/hyperpic/master/_resources/hyperpic.svg "Hyperpic logo")


[![Go Report Card](https://goreportcard.com/badge/github.com/hyperscale/hyperpic)](https://goreportcard.com/report/github.com/hyperscale/hyperpic)

| Branch  | Status | Coverage |
|---------|--------|----------|
| master  | [![Build Status](https://img.shields.io/travis/hyperscale/hyperpic/master.svg)](https://travis-ci.org/hyperscale/hyperpic) | [![Coveralls](https://img.shields.io/coveralls/hyperscale/hyperpic/master.svg)](https://coveralls.io/github/hyperscale/hyperpic?branch=master) |
| develop | [![Build Status](https://img.shields.io/travis/hyperscale/hyperpic/develop.svg)](https://travis-ci.org/hyperscale/hyperpic) | [![Coveralls](https://img.shields.io/coveralls/hyperscale/hyperpic/develop.svg)](https://coveralls.io/github/hyperscale/hyperpic?branch=develop) |

Fast HTTP microservice for high-level image processing.

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

Install dependencies with golang dep:
```shell
dep ensure
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

Original: `https://hyperpic.herokuapp.com/kayaks.jpg`

![Original](https://hyperpic.herokuapp.com/kayaks.jpg)

Croped and Resized: `https://hyperpic.herokuapp.com/kayaks.jpg?w=400&h=400&fit=crop`

![Croped and resized](https://hyperpic.herokuapp.com/kayaks.jpg?w=400&h=400&fit=crop)

### Crop on focal point

Original: `https://hyperpic.herokuapp.com/smartcrop.jpg`

![Original](https://hyperpic.herokuapp.com/smartcrop.jpg)

Croped and Resized: `https://hyperpic.herokuapp.com/smartcrop.jpg?w=200&h=200&fit=crop-focal-point`

![Croped and resized](https://hyperpic.herokuapp.com/smartcrop.jpg?w=200&h=200&fit=crop-focal-point)

Documentation
-------------

[Hyperpic API Reference](https://hyperscale.github.io/hyperpic/)

License
-------

hyperpic is licensed under [the MIT license](LICENSE.md).
