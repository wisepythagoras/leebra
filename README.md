# Leebra

Leebra is an experimental browser engine written (mostly) in Go.

## How it works

The browser is built mostly from scratch.

### HTML and Style Engine

These will be built from scratch or I will use Go's [built-in HTML parser](https://pkg.go.dev/golang.org/x/net/html#Parse).

### JavaScript Engine

Leebra is using V8 as its JavaScript engine due to the ease of using it.

## Why?

I've always wanted to learn more about browsers and thought I should write one from scratch (ish).
