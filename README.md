# Leebra

Leebra is an experimental browser engine written (mostly) in Go. To get it running you'll need to have at least Go 1.17.1 installed. Then you need to run the following:

``` sh
cd leebra
go build .
```

If everything went well, you should get a binary called `leebra`.

``` sh
./leebra -h
Usage of ./leebra:
  -run string
        Runs a JavaScript file
  -url string
        The URL to load (default "about:blank")
```

To just run a JavaScript file, all you need to do is run the following:

``` sh
./leebra -run path/to/js/file.js
```

There are plenty of JavaScript examples in the [js folder](js).

Also, you can supply a url with the `-url` command line argument. For now, this doesn't do much; it will only download the page and attempt to - at some level - parse the HTML and create DOM objects from it.

## How it works

The browser is built mostly from scratch.

### HTML and Style Engine

These will be built from scratch or I will use Go's [built-in HTML parser](https://pkg.go.dev/golang.org/x/net/html#Parse).

### JavaScript Engine

Leebra is using V8 as its JavaScript engine due to the ease of using it.

### WebAssembly Engine

This part is documented in a separate [README](jscore/wasm/README.md).

## Why?

I've always wanted to learn more about browsers and thought I should write one from scratch (ish).
