package jscore

import (
	"fmt"

	"github.com/wasmerio/wasmer-go/wasmer"
	"rogchap.com/v8go"
)

// Wasm describes the WeAssembly JavaScript API.
type Wasm struct {
	VM          *v8go.Isolate
	ExecContext *v8go.Context
	WasmEngine  *wasmer.Engine
	wasmObj     *v8go.Object
}

// NewEngine creates a new instance of the Wasmer engine.
func (w *Wasm) NewEngine() *wasmer.Engine {
	if w.WasmEngine == nil {
		w.WasmEngine = wasmer.NewEngine()
	}

	return w.WasmEngine
}

// InstantiateFunction allows users to compile and instantiate wasm code.
func (w *Wasm) InstantiateFunction() (*v8go.FunctionTemplate, error) {
	setItemFn, err := v8go.NewFunctionTemplate(w.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()

		if len(args) < 1 {
			// TODO: Figure out how to return errors here.
			return nil
		}

		wasmCode := []byte(args[0].String())
		resolver, _ := v8go.NewPromiseResolver(info.Context())

		go func() {
			wasmEngine := w.NewEngine()
			store := wasmer.NewStore(wasmEngine)

			// Compile the WebAssembly code.
			module, err := wasmer.NewModule(store, wasmCode)

			if err != nil {
				fmt.Println(err)
				fmt.Println("Module:", args)
				errorMessage, _ := v8go.NewValue(w.VM, "Unable to compile module")
				resolver.Reject(errorMessage)
				return
			}

			// instantiate the code.
			importObject := wasmer.NewImportObject()
			instance, _ := wasmer.NewInstance(module, importObject)

			// TODO: Return whatever it is that this function returns.
			// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/WebAssembly/instantiate
			val, _ := v8go.NewValue(w.VM, true)
			resolver.Resolve(val)

			sum, err := instance.Exports.GetFunction("sum")

			fmt.Println(err)

			result, err := sum(1, 2)

			fmt.Println(result, err)
		}()

		return resolver.Value
	})

	if err != nil {
		return nil, err
	}

	return setItemFn, nil
}

// GetV8Object gets the entire object structure of the V8 WebAssembly API.
func (w *Wasm) GetV8Object() (*v8go.ObjectTemplate, error) {
	wasmObj, err := v8go.NewObjectTemplate(w.VM)

	if err != nil {
		return nil, err
	}

	instantiateFn, err := w.InstantiateFunction()

	if err != nil {
		return nil, err
	}

	wasmObj.Set("instantiate", instantiateFn, v8go.ReadOnly)

	return wasmObj, nil
}

// GetJSObject returns the JS Object that can be mutated.
func (w *Wasm) GetJSObject() (*v8go.Object, error) {
	clipboard, err := w.GetV8Object()

	if err != nil {
		return nil, err
	}

	w.wasmObj, err = clipboard.NewInstance(w.ExecContext)

	if err != nil {
		return nil, err
	}

	return w.wasmObj, nil
}
