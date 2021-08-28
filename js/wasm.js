(async () => {
    const wasmFile = 'js/hello_wasm.wasm';

    try {
        // const response = await fetch(wasmFile);
        // const bytes = await response.arrayBuffer();
        // console.log('Response', response, '1');
        const wasm = await WebAssembly.instantiate(wasmFile, {});

        console.log('The sum is: ', wasm.instance.exports.sum(1, 2));
        console.log('The mul is: ', wasm.instance.exports.mul(3, 2));
        console.log(wasm.instance.exports.__heap_base.value);
    } catch(e) {
        console.log('Error!', e);
    }
})();
