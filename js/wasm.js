console.log(WebAssembly.instantiate);

(async () => {
    const wasmFile = 'http://localhost:8000/target/wasm32-unknown-unknown/release/hello_wasm.wasm';

    try {
        const response = await fetch(wasmFile);
        const bytes = await response.arrayBuffer();
        console.log('Response', response, '1');
        await WebAssembly.instantiate(bytes, {});

        // console.log('The answer is: ', instance.exports.sum(1, 2));
    } catch(e) {
        console.log('Error!', e);
    }
})();
