package jscore

import (
	"fmt"
	"io/ioutil"

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

func (w *Wasm) CreateInstantiateResponse(instance *wasmer.Instance, module *wasmer.Module) (*v8go.Object, error) {
	instanceObj, err := v8go.NewObjectTemplate(w.VM)

	if err != nil {
		return nil, err
	}

	// This could be empty for now.
	moduleObj, _ := v8go.NewObjectTemplate(w.VM)
	instanceObj.Set("module", moduleObj, v8go.ReadOnly)

	// Create the instance object.
	v8Instance := &WasmInstance{
		VM:           w.VM,
		ExecContext:  w.ExecContext,
		WasmInstance: instance,
		WasmModule:   module,
	}
	v8InstanceObj, err := v8Instance.GetV8Object()

	if err != nil {
		return nil, err
	}

	err = instanceObj.Set("instance", v8InstanceObj, v8go.ReadOnly)

	if err != nil {
		fmt.Println("---", err)
	}

	// Now let's create the JS object.
	wasmInstanceObj, err := instanceObj.NewInstance(w.ExecContext)

	if err != nil {
		return nil, err
	}

	return wasmInstanceObj, nil
}

// InstantiateFunction allows users to compile and instantiate wasm code.
func (w *Wasm) InstantiateFunction() (*v8go.FunctionTemplate, error) {
	setItemFn, err := v8go.NewFunctionTemplate(w.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()

		if len(args) < 1 {
			// TODO: Figure out how to return errors here.
			return nil
		}

		// TODO: In production code, this needs to be an array buffer, or []bytes. For now I need to work with a string.
		wasmCode := args[0].String()
		wasmBytes, err := ioutil.ReadFile(wasmCode)

		if err != nil {
			fmt.Println(err)
			return nil
		}

		resolver, _ := v8go.NewPromiseResolver(info.Context())

		go func() {
			wasmEngine := w.NewEngine()
			store := wasmer.NewStore(wasmEngine)

			// Compile the WebAssembly code.
			module, err := wasmer.NewModule(store, wasmBytes)

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

			instantiateResp, err := w.CreateInstantiateResponse(instance, module)

			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(instantiateResp.Get("module"))

			// TODO: Return whatever it is that this function returns.
			// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/WebAssembly/instantiate
			resolver.Resolve(instantiateResp)
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
	wasm, err := w.GetV8Object()

	if err != nil {
		return nil, err
	}

	w.wasmObj, err = wasm.NewInstance(w.ExecContext)

	if err != nil {
		return nil, err
	}

	return w.wasmObj, nil
}
