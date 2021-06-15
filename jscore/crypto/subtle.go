package js

import (
	"rogchap.com/v8go"
)

type SubtleCrypto struct {
	VM          *v8go.Isolate
	ExecContext *v8go.Context
	subtleObj   *v8go.Object
}

// GetGenerateKeyFunction creates the function that generates a key.
func (c *SubtleCrypto) GetGenerateKeyFunction() (*v8go.FunctionTemplate, error) {
	generateKeyFn, err := v8go.NewFunctionTemplate(c.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
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

	if err != nil {
		return nil, err
	}

	return generateKeyFn, nil
}

// GetV8Object gets the entire object structure of the V8 Subtle Crypto API.
func (c *SubtleCrypto) GetV8Object() (*v8go.ObjectTemplate, error) {
	cryptoObj, err := v8go.NewObjectTemplate(c.VM)

	if err != nil {
		return nil, err
	}

	generateKeyFunction, err := c.GetGenerateKeyFunction()

	if err != nil {
		return nil, err
	}

	cryptoObj.Set("generateKey", generateKeyFunction)

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
