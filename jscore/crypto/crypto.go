package js

import (
	"fmt"
	"math/rand"

	"rogchap.com/v8go"
)

type Crypto struct {
	VM          *v8go.Isolate
	ExecContext *v8go.Context
	cryptoObj   *v8go.Object
}

// GetGetRandomValuesFunction creates the function that returns random bytes.
func (c *Crypto) GetGetRandomValuesFunction() *v8go.FunctionTemplate {
	return v8go.NewFunctionTemplate(c.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		var randomValuesArr *v8go.Value

		if len(args) > 0 {
			randomValuesArr = args[0]
		}

		rand.Uint32()

		if !randomValuesArr.IsInt8Array() &&
			!randomValuesArr.IsInt16Array() &&
			!randomValuesArr.IsInt32Array() &&
			!randomValuesArr.IsUint8Array() &&
			!randomValuesArr.IsUint16Array() &&
			!randomValuesArr.IsUint32Array() &&
			!randomValuesArr.IsUint8ClampedArray() &&
			!randomValuesArr.IsFloat32Array() &&
			!randomValuesArr.IsFloat64Array() {
			c.ExecContext.RunScript("throw new Error('Invalid input type')", "")
			return nil
		}

		valArray := randomValuesArr.Object()
		for i := 0; valArray.HasIdx(uint32(i)); i++ {
			var err error
			idx := uint32(i)

			if randomValuesArr.IsInt8Array() {
				r := rand.Int31()*(127+128) - 128
				err = valArray.SetIdx(idx, r)
			} else if randomValuesArr.IsInt16Array() {
				r := rand.Int31()*(32767+32768) - 32767
				err = valArray.SetIdx(idx, r)
			} else if randomValuesArr.IsInt32Array() {
				err = valArray.SetIdx(idx, rand.Int31())
			} else if randomValuesArr.IsUint8Array() {
				err = valArray.SetIdx(idx, rand.Uint32()%255)
			} else if randomValuesArr.IsUint16Array() {
				err = valArray.SetIdx(idx, rand.Uint32()%65535)
			} else if randomValuesArr.IsUint32Array() {
				err = valArray.SetIdx(idx, rand.Uint32())
			} else if randomValuesArr.IsFloat32Array() {
				// TODO: Figure out why float32 doesn't work.
				err = valArray.SetIdx(idx, float64(rand.Float32()))
			} else if randomValuesArr.IsFloat64Array() {
				err = valArray.SetIdx(idx, rand.Float64())
			}

			if err != nil {
				fmt.Println(i, err)
				return nil
			}
		}

		return randomValuesArr
	})
}

// GetV8Object gets the entire object structure of the V8 Crypto API.
func (c *Crypto) GetV8Object() (*v8go.ObjectTemplate, error) {
	cryptoObj := v8go.NewObjectTemplate(c.VM)
	getRandomValuesFunction := c.GetGetRandomValuesFunction()
	cryptoObj.Set("getRandomValues", getRandomValuesFunction)

	subtle := &SubtleCrypto{
		VM:          c.VM,
		ExecContext: c.ExecContext,
	}
	subtleObj, err := subtle.GetV8Object()

	if err != nil {
		return nil, err
	}

	cryptoObj.Set("subtle", subtleObj)

	return cryptoObj, nil
}

// GetJSObject returns the JS Object that can be mutated.
func (c *Crypto) GetJSObject() (*v8go.Object, error) {
	crypto, err := c.GetV8Object()

	if err != nil {
		return nil, err
	}

	c.cryptoObj, err = crypto.NewInstance(c.ExecContext)

	if err != nil {
		return nil, err
	}

	return c.cryptoObj, nil
}
