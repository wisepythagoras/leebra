package console

import (
	"fmt"

	"rogchap.com/v8go"
)

type Console struct {
	VM          *v8go.Isolate
	ExecContext *v8go.Context
	consoleObj  *v8go.Object
}

func (c *Console) getLogFunctionCallback() *v8go.FunctionTemplate {
	return v8go.NewFunctionTemplate(c.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		argv := info.Args()

		if len(argv) > 0 {
			args := make([]interface{}, len(argv))

			for i, input := range argv {
				args[i] = input
			}

			fmt.Println(args...)
		}

		return nil
	})
}

// GetV8Object gets the entire object structure for the console.
func (c *Console) GetV8Object() (*v8go.ObjectTemplate, error) {
	consoleObj := v8go.NewObjectTemplate(c.VM)
	consoleLogObj := c.getLogFunctionCallback()

	consoleObj.Set("log", consoleLogObj)
	consoleObj.Set("info", consoleLogObj)
	consoleObj.Set("warn", consoleLogObj)
	consoleObj.Set("error", consoleLogObj)

	return consoleObj, nil
}

// GetJSObject returns the JS Object that can be mutated.
func (c *Console) GetJSObject() (*v8go.Object, error) {
	subtle, err := c.GetV8Object()

	if err != nil {
		return nil, err
	}

	c.consoleObj, err = subtle.NewInstance(c.ExecContext)

	if err != nil {
		return nil, err
	}

	return c.consoleObj, nil
}
