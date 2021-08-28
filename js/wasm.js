(async () => {
    const wasmFile = 'js/hello_wasm.wasm';

    try {
        // const response = await fetch(wasmFile);
        // const bytes = await response.arrayBuffer();
        // console.log('Response', response, '1');
        const wasm = await WebAssembly.instantiate(wasmFile, {});

        console.log('----');
        for (let k in wasm) {
            console.log(k);
        }
        console.log('----', wasm.module);

        console.log('The sum is: ', wasm.instance.exports.sum(1, 2));
        console.log('The mul is: ', wasm.instance.exports.mul(3, 2));
    } catch(e) {
        console.log('Error!', e);
    }
})();
