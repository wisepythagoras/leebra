# WebAssembly Support

Leebra uses [wasmer-go](https://github.com/wasmerio/wasmer-go) as its WebAssembly runtime engine. As of now, [v8go](https://github.com/rogchap/v8go) does not support array buffers, and therefore a path to the WebAssembly binary needs to be passed, instead of getting the contents with a `fetch`.

## WASM Module

The way I built the WebAssembly module, was this:

I created a simple project with the help of `cargo`.

``` sh
cargo new --lib hello-wasm
cd hello-wasm
```

Then I opened `src/lib.rs` and pasted the following code inside.

``` rust
#[no_mangle]
pub extern "C" fn sum(x: i32, y: i32) -> i32 {
    x + y
}

#[no_mangle]
pub extern "C" fn mul(x: i32, y: i32) -> i64 {
    (x as i64) * (y as i64)
}
```

The above code can be built with the following command:

```sh
rustup target add wasm32-unknown-unknown
cargo build --target wasm32-unknown-unknown --release
```

If everything went well, your module will be available at the following path:

```
target/wasm32-unknown-unknown/release/hello_wasm.wasm
```

## JavaScript Example

There's a dedicated JavaScript example that lives [here](https://github.com/wisepythagoras/leebra/blob/main/js/wasm.js).

### Example break-down

``` js
// The Go code will load the wasm module at the specific path.
const wasm = await WebAssembly.instantiate('js/hello_wasm.wasm', {});

// Call the `sum` function with 1 and 2 as the arguments.
wasm.instance.exports.sum(1, 2); // This should yield 3.

// Call the `mul` function with 3 and 2 as the arguments.
wasm.instance.exports.mul(3, 2); // This should yield 6.

// This gets the value of the `__heap_base` global.
wasm.instance.exports.__heap_base.value
```

## What's missing

1. The memory export,
2. Native array buffer support so we can load through a `fetch`
