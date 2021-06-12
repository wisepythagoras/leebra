package jscore

import (
	"fmt"

	"golang.design/x/clipboard"
	"rogchap.com/v8go"
)

// Clipboard defines the Clipboard web API.
type Clipboard struct {
	VM           *v8go.Isolate
	ExecContext  *v8go.Context
	clipboardObj *v8go.Object
}

// GetReadTextFunction creates the function that reads text from the clipboard.
func (c *Clipboard) GetReadTextFunction() (*v8go.FunctionTemplate, error) {
	readFn, err := v8go.NewFunctionTemplate(c.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		resolver, _ := v8go.NewPromiseResolver(info.Context())

		go func() {
			bytes := clipboard.Read(clipboard.FmtText)
			textVal, _ := v8go.NewValue(c.VM, string(bytes))
			resolver.Resolve(textVal)
		}()

		return resolver.GetPromise().Value
	})

	if err != nil {
		return nil, err
	}

	return readFn, nil
}

// GetWriteTextFunction creates the function that writes text to the clipboard.
func (c *Clipboard) GetWriteTextFunction() (*v8go.FunctionTemplate, error) {
	readFn, err := v8go.NewFunctionTemplate(c.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		var writeText *v8go.Value

		if len(args) > 0 {
			writeText = args[0]
		}

		resolver, _ := v8go.NewPromiseResolver(info.Context())

		go func() {
			if writeText != nil {
				fmt.Println(writeText)
				clipboard.Write(clipboard.FmtText, []byte(writeText.String()))
				val, _ := v8go.NewValue(c.VM, true)
				resolver.Resolve(val)
			} else {
				errorMessage, _ := v8go.NewValue(c.VM, "Clipboard write failed")
				resolver.Reject(errorMessage)
			}
		}()

		return resolver.GetPromise().Value
	})

	if err != nil {
		return nil, err
	}

	return readFn, nil
}

// GetV8Object gets the entire object structure of the V8 Clipboard API.
func (c *Clipboard) GetV8Object() (*v8go.ObjectTemplate, error) {
	clipboardObj, err := v8go.NewObjectTemplate(c.VM)

	if err != nil {
		return nil, err
	}

	readFn, err := c.GetReadTextFunction()

	if err != nil {
		return nil, err
	}

	clipboardObj.Set("readText", readFn, v8go.ReadOnly)

	writeFn, err := c.GetWriteTextFunction()

	if err != nil {
		return nil, err
	}

	clipboardObj.Set("writeText", writeFn, v8go.ReadOnly)

	return clipboardObj, nil
}

// GetJSObject returns the JS Object that can be mutated.
func (c *Clipboard) GetJSObject() (*v8go.Object, error) {
	clipboard, err := c.GetV8Object()

	if err != nil {
		return nil, err
	}

	c.clipboardObj, err = clipboard.NewInstance(c.ExecContext)

	if err != nil {
		return nil, err
	}

	return c.clipboardObj, nil
}
