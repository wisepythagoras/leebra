package net

import (
	"rogchap.com/v8go"
)

type URL struct {
	VM          *v8go.Isolate
	ExecContext *v8go.Context
	urlObj      *v8go.Object
}

func URLFn(vm *v8go.Isolate) *v8go.FunctionTemplate {
	urlConstructor := v8go.NewFunctionTemplate(vm, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		return nil
	})

	// TODO: Set a prototype.

	return urlConstructor
}
