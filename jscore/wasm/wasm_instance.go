package jscore

import (
	"fmt"

	"github.com/wasmerio/wasmer-go/wasmer"
	"rogchap.com/v8go"
)

type WasmInstance struct {
	VM              *v8go.Isolate
	ExecContext     *v8go.Context
	WasmInstance    *wasmer.Instance
	WasmModule      *wasmer.Module
	wasmInstanceObj *v8go.Object
}

func (w *WasmInstance) InstantiateFunction(name string) (*v8go.FunctionTemplate, error) {
	wasmFn, err := w.WasmInstance.Exports.GetFunction(name)

	if err != nil {
		return nil, err
	}

	setItemFn := v8go.NewFunctionTemplate(w.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		nativeArgs := make([]interface{}, 0)

		for _, v := range args {
			if v.IsString() {
				nativeArgs = append(nativeArgs, v.String())
			} else if v.IsNumber() {
				if v.IsInt32() {
					nativeArgs = append(nativeArgs, v.Int32())
				} else if v.IsBigInt() {
					nativeArgs = append(nativeArgs, v.BigInt())
				}
			}
		}

		result, err := wasmFn(nativeArgs...)

		if err != nil {
			errVal, _ := v8go.NewValue(info.Context().Isolate(), err.Error())
			info.Context().Isolate().ThrowException(errVal)
			return nil
		}

		val, _ := v8go.NewValue(w.VM, result)

		return val
	})

	return setItemFn, nil
}

// GetV8Object gets the entire object structure of the V8 WebAssembly.Instance API.
func (w *WasmInstance) GetV8Object() (*v8go.ObjectTemplate, error) {
	wasmInstanceObj := v8go.NewObjectTemplate(w.VM)
	exportsObj := v8go.NewObjectTemplate(w.VM)

	for _, v := range w.WasmModule.Exports() {
		kind := v.Type().Kind().String()

		if kind == "func" {
			wasmFn, err := w.InstantiateFunction(v.Name())

			if err != nil {
				fmt.Println("[DEBUG]", err)
				continue
			}

			// Set the function on the exports object.
			err = exportsObj.Set(v.Name(), wasmFn, v8go.ReadOnly)

			if err != nil {
				fmt.Println("[DEBUG]", err)
			}
		} else if kind == "global" {
			global, err := w.WasmInstance.Exports.GetGlobal(v.Name())

			if err != nil {
				fmt.Println("[DEBUG]", err)
				continue
			}

			obj := v8go.NewObjectTemplate(w.VM)
			value, _ := global.Get()
			obj.Set("value", value)

			err = exportsObj.Set(v.Name(), obj, v8go.ReadOnly)

			if err != nil {
				fmt.Println("[DEBUG]", err)
			}
		}
	}

	wasmInstanceObj.Set("exports", exportsObj, v8go.ReadOnly)

	return wasmInstanceObj, nil
}

// GetJSObject returns the JS Object that can be mutated.
func (w *WasmInstance) GetJSObject() (*v8go.Object, error) {
	clipboard, err := w.GetV8Object()

	if err != nil {
		return nil, err
	}

	w.wasmInstanceObj, err = clipboard.NewInstance(w.ExecContext)

	if err != nil {
		return nil, err
	}

	return w.wasmInstanceObj, nil
}
