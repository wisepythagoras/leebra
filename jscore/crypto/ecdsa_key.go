package js

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"rogchap.com/v8go"
)

// ECDSAKey defines the
type ECDSAKey struct {
	PrivateKey  *ecdsa.PrivateKey
	PublicKey   *ecdsa.PublicKey
	VM          *v8go.Isolate
	ExecContext *v8go.Context
	Extractable bool
	KeyUsages   []ECDSAKeyUsage
	CurveType   ECDSANamedCurve
	subtleObj   *v8go.Object
}

// GenerateKey generates a key pair.
func (e *ECDSAKey) GenerateKey() (*ecdsa.PrivateKey, error) {
	var curve elliptic.Curve
	var err error

	if e.CurveType == ECDSAP256 {
		curve = elliptic.P256()
	} else if e.CurveType == ECDSAP384 {
		curve = elliptic.P384()
	} else if e.CurveType == ECDSAP521 {
		curve = elliptic.P521()
	}

	e.PrivateKey, err = ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		return nil, err
	}

	e.PublicKey = &e.PrivateKey.PublicKey

	return e.PrivateKey, nil
}

func (e *ECDSAKey) getPrivateKeyObj() (*v8go.ObjectTemplate, error) {
	pkObj := v8go.NewObjectTemplate(e.VM)
	algoObj := v8go.NewObjectTemplate(e.VM)
	algoObj.Set("name", "ECDSA")

	if e.CurveType == ECDSAP256 {
		algoObj.Set("type", "P-256")
	} else if e.CurveType == ECDSAP384 {
		algoObj.Set("type", "P-384")
	} else if e.CurveType == ECDSAP521 {
		algoObj.Set("type", "P-521")
	}

	usagesObj := v8go.NewObjectTemplate(e.VM)
	usagesArr, _ := usagesObj.NewInstance(e.ExecContext)
	usagesArr.SetIdx(0, "sign")

	pkObj.Set("extractable", false, v8go.DontEnum)
	pkObj.Set("usages", usagesArr, v8go.DontEnum)
	pkObj.Set("type", "private", v8go.DontEnum)
	pkObj.Set("algorithm", algoObj, v8go.DontEnum)
	pkObj.Set("Symbol(Symbol.toStringTag)", "CryptoKey")

	return pkObj, nil
}

// GetV8Object gets the entire object structure for this ECDSA key.
func (e *ECDSAKey) GetV8Object() (*v8go.ObjectTemplate, error) {
	ecdsaObj := v8go.NewObjectTemplate(e.VM)
	privateKeyObj, err := e.getPrivateKeyObj()

	if err != nil {
		return nil, err
	}

	ecdsaObj.SetInternalFieldCount(1)
	ecdsaObj.Set("privateKey", privateKeyObj)

	// TODO: In order for us to be able to sign and verify input, we need support
	// for UInt8Array and TextEncoder.
	// https://developer.mozilla.org/en-US/docs/Web/API/TextEncoder
	//
	// Example:
	// https://github.com/mdn/dom-examples/blob/master/web-crypto/sign-verify/ecdsa.js
	//
	// More on this:
	// https://developer.mozilla.org/en-US/docs/Web/API/SubtleCrypto/generateKey

	return ecdsaObj, nil
}

// GetJSObject returns the JS Object that can be mutated.
func (e *ECDSAKey) GetJSObject() (*v8go.Object, error) {
	subtle, err := e.GetV8Object()

	if err != nil {
		return nil, err
	}

	e.subtleObj, err = subtle.NewInstance(e.ExecContext)

	if err != nil {
		return nil, err
	}

	e.subtleObj.SetInternalField(0, ECDSAEncodePrivate(e.PrivateKey))

	return e.subtleObj, nil
}
