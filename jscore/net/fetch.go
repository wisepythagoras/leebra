package net

import (
	"rogchap.com/v8go"
)

func responseObject(vm *v8go.Isolate) *v8go.ObjectTemplate {
	respObj := v8go.NewObjectTemplate(vm)

	return respObj
}

// CreateFetchFn creates the V8 function for `fetch(input[, options])`.
func CreateFetchFn(vm *v8go.Isolate) *v8go.FunctionTemplate {
	return v8go.NewFunctionTemplate(vm, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()

		// TODO: This should be a URL or a string or "any other object with a stringifier".
		url := args[0].String()

		resolver, _ := v8go.NewPromiseResolver(info.Context())
		var optArg *v8go.Value

		if len(args) > 1 {
			optArg = args[1]
		}

		go func() {
			response, err := HTTPRequest(url, optArg)

			if err != nil {
				errVal, _ := v8go.NewValue(vm, err.Error())
				resolver.Reject(errVal)
				return
			}

			respObj := &Response{
				VM:          vm,
				ExecContext: info.Context(),
				body:        response.Body,
			}
			resp, _ := respObj.GetJSObject()

			resolver.Resolve(resp)
		}()

		// Fetch returns a promise.
		return resolver.GetPromise().Value
	})
}
