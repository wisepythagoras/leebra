package net

import (
	"net/http"

	"github.com/wisepythagoras/leebra/jscore"
	"rogchap.com/v8go"
)

func responseObject(vm *v8go.Isolate) *v8go.ObjectTemplate {
	respObj := v8go.NewObjectTemplate(vm)

	return respObj
}

// CreateFetchFn creates the V8 function for `fetch(input[, options])`.
func CreateFetchFn(vm *v8go.Isolate, nav *jscore.Navigator) *v8go.FunctionTemplate {
	return v8go.NewFunctionTemplate(vm, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()

		// TODO: This should be a URL or a string or "any other object with a stringifier".
		url := args[0].String()

		resolver, _ := v8go.NewPromiseResolver(info.Context())

		go func() {
			// The request method should come from the options. If there is none defined, then
			// it can default to "GET".
			request, err := http.NewRequest("GET", url, nil)

			if err != nil {
				errVal, _ := v8go.NewValue(vm, err.Error())
				resolver.Reject(errVal)
				return
			}

			// TODO: Here there should be some loop that adds all the headers that are in the
			// options (the second argument).
			request.Header.Set("User-Agent", nav.GetUserAgent())

			client := &http.Client{}
			response, err := client.Do(request)

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
