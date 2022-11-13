[![Go Report Card](https://goreportcard.com/badge/github.com/bep/gowebp)](https://goreportcard.com/report/github.com/bep/gowebp)
[![libwebp Version](https://img.shields.io/badge/libwebp-v1.2.0-blue)](https://github.com/webmproject/libwebp)
[![codecov](https://codecov.io/gh/bep/gowebp/branch/master/graph/badge.svg)](https://codecov.io/gh/bep/gowebp)
[![GoDoc](https://godoc.org/github.com/bep/gowebp/libwebp?status.svg)](https://godoc.org/github.com/bep/gowebp/libwebp)

This library provides C bindings and an API for **encoding** Webp images using Google's [libwebp](https://github.com/webmproject/libwebp).

It is based on [go-webp](https://github.com/kolesa-team/go-webp), but this includes and builds the libwebp C source from a versioned Git subtree.


## Update libwebp version

1. Pull in the relevant libwebp version, e.g. `./pull-libwebp.sh v1.2.0`
2. Regenerate wrappers with `go generate ./gen`
3. Update the libwebp version badge above.

## Local development

Compiling C code isn' particulary fast; if you install libwebp on your PC you can link against that, useful during development.

On a Mac you may do something like:

```bash
brew install webp
go test ./libwebp -tags dev
```


