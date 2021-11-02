package net

import (
	"io"
	"io/ioutil"

	"rogchap.com/v8go"
)

type Response struct {
	VM          *v8go.Isolate
	ExecContext *v8go.Context
	fetchObj    *v8go.Object
	body        io.ReadCloser
}

// TODO: Similar to this one, `.json()` will use `v8go.JSONParse(r.VM, body)`.
func (r *Response) getTextFn() *v8go.FunctionTemplate {
	return v8go.NewFunctionTemplate(r.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		resolver, _ := v8go.NewPromiseResolver(info.Context())

		go func() {
			body, _ := ioutil.ReadAll(r.body)
			val, _ := v8go.NewValue(r.VM, string(body))

			resolver.Resolve(val)
		}()

		// Fetch returns a promise.
		return resolver.GetPromise().Value
	})
}

// GetV8Object gets the entire object structure for the response object.
func (r *Response) GetV8Object() (*v8go.ObjectTemplate, error) {
	respObj := v8go.NewObjectTemplate(r.VM)
	respTextFn := r.getTextFn()

	respObj.Set("text", respTextFn)

	return respObj, nil
}

// GetJSObject returns the JS Object that can be mutated.
func (r *Response) GetJSObject() (*v8go.Object, error) {
	subtle, err := r.GetV8Object()

	if err != nil {
		return nil, err
	}

	r.fetchObj, err = subtle.NewInstance(r.ExecContext)

	if err != nil {
		return nil, err
	}

	return r.fetchObj, nil
}
