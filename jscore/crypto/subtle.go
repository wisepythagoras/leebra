package js

import (
	"encoding/json"
	"fmt"

	"rogchap.com/v8go"
)

type SubtleCrypto struct {
	VM          *v8go.Isolate
	ExecContext *v8go.Context
	subtleObj   *v8go.Object
}

// GetGenerateKeyFunction creates the function that generates a key.
func (c *SubtleCrypto) GetGenerateKeyFunction() *v8go.FunctionTemplate {
	return v8go.NewFunctionTemplate(c.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		var algorithm *v8go.Object // *v8go.Value
		var extractable bool
		var keyUsagesObj *v8go.Object

		if len(args) < 3 {
			return nil
		}

		algorithm = args[0].Object()
		extractable = args[1].Boolean()
		keyUsagesObj = args[2].Object()

		if algorithm.IsNullOrUndefined() || keyUsagesObj.IsNullOrUndefined() {
			return nil
		}

		fnContext := info.Context()
		resolver, _ := v8go.NewPromiseResolver(fnContext)

		algoName, err := algorithm.Get("name")

		if algoName.IsNullOrUndefined() || err != nil {
			errorMessage, _ := v8go.NewValue(c.VM, "Invalid algorithm")
			resolver.Reject(errorMessage)
			return nil
		}

		go func() {
			if algoName.String() == "ECDSA" {
				var keyUsages []ECDSAKeyUsage
				var curveType ECDSANamedCurve

				namedCurve, err := algorithm.Get("namedCurve")

				if err != nil {
					errorMessage, _ := v8go.NewValue(c.VM, err.Error())
					resolver.Reject(errorMessage)
					return
				}

				for i := 0; keyUsagesObj.HasIdx(uint32(i)); i++ {
					key, err := keyUsagesObj.GetIdx(uint32(i))

					if err != nil {
						continue
					}

					if key.String() == "verify" {
						keyUsages = append(keyUsages, ECDSAVerify)
					} else if key.String() == "sign" {
						keyUsages = append(keyUsages, ECDSASign)
					}
				}

				if namedCurve.String() == "P-256" {
					curveType = ECDSAP256
				} else if namedCurve.String() == "P-384" {
					curveType = ECDSAP384
				} else if namedCurve.String() == "P-521" {
					curveType = ECDSAP521
				}

				key := &ECDSAKey{
					VM:          c.VM,
					ExecContext: c.ExecContext,
					Extractable: extractable,
					KeyUsages:   keyUsages,
					CurveType:   curveType,
				}

				_, err = key.GenerateKey()

				if err != nil {
					errorMessage, _ := v8go.NewValue(c.VM, err.Error())
					resolver.Reject(errorMessage)
					return
				}

				keyObj, _ := key.GetJSObject()

				resolver.Resolve(keyObj)

				return
			}
		}()

		return resolver.GetPromise().Value
	})
}

// TODO: This is a function testing the ability to work around v8go, which does not have the
// UInt8Array implemented yet. Remove it in the future.
func (c *SubtleCrypto) getTestFunctionCallback() *v8go.FunctionTemplate {
	return v8go.NewFunctionTemplate(c.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		if len(info.Args()) > 0 {
			arg := info.Args()[0]
			m, err := arg.MarshalJSON()

			if err != nil {
				fmt.Println(err)
				return nil
			}

			uint8Arr := make([]uint8, 0)
			err = json.Unmarshal(m, &uint8Arr)

			if err != nil {
				fmt.Println(err)
				return nil
			}

			fmt.Println(uint8Arr)
		}

		return nil
	})
}

func (c *SubtleCrypto) getTestsFunctionCallback() *v8go.FunctionTemplate {
	return v8go.NewFunctionTemplate(c.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		if len(info.Args()) > 0 {
			arg := info.Args()[0].Object()
			messageVal, err := arg.Get("message")

			if err != nil {
				fmt.Println(err)
				return nil
			}

			fmt.Println("Error ", messageVal.String())
		}

		return nil
	})
}

// GetV8Object gets the entire object structure of the V8 Subtle Crypto API.
func (c *SubtleCrypto) GetV8Object() (*v8go.ObjectTemplate, error) {
	cryptoObj := v8go.NewObjectTemplate(c.VM)
	generateKeyFunction := c.GetGenerateKeyFunction()
	testFn := c.getTestFunctionCallback()
	test2Fn := c.getTestsFunctionCallback()

	cryptoObj.Set("generateKey", generateKeyFunction)
	cryptoObj.Set("testFn", testFn)
	cryptoObj.Set("test2Fn", test2Fn)

	return cryptoObj, nil
}

// GetJSObject returns the JS Object that can be mutated.
func (c *SubtleCrypto) GetJSObject() (*v8go.Object, error) {
	crypto, err := c.GetV8Object()

	if err != nil {
		return nil, err
	}

	c.subtleObj, err = crypto.NewInstance(c.ExecContext)

	if err != nil {
		return nil, err
	}

	return c.subtleObj, nil
}
