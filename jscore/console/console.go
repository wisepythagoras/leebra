package console

import (
	"fmt"
	"time"

	"rogchap.com/v8go"
)

type Console struct {
	VM          *v8go.Isolate
	ExecContext *v8go.Context
	consoleObj  *v8go.Object
	countMap    map[string]int
	timeMap     map[string]time.Time
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

func (c *Console) getCountFunctionCallback() *v8go.FunctionTemplate {
	c.countMap = make(map[string]int)

	return v8go.NewFunctionTemplate(c.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		label := "default"

		if len(info.Args()) > 0 {
			label = info.Args()[0].String()
		}

		count, exists := c.countMap[label]

		if !exists {
			count = 0
		}

		count += 1
		c.countMap[label] = count

		fmt.Println(label, count)

		return nil
	})
}

func (c *Console) getCountResetFunctionCallback() *v8go.FunctionTemplate {
	c.countMap = make(map[string]int)

	return v8go.NewFunctionTemplate(c.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		label := "default"

		if len(info.Args()) > 0 {
			label = info.Args()[0].String()
		}

		c.countMap[label] = 0
		fmt.Println(label, 0)

		return nil
	})
}

func (c *Console) getTimeFunctionCallback() *v8go.FunctionTemplate {
	c.timeMap = make(map[string]time.Time)

	return v8go.NewFunctionTemplate(c.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		// For any given context there cannot be more than 10 thousand timers running.
		if len(c.timeMap) >= 10000 {
			return nil
		}

		label := "default"

		if len(info.Args()) > 0 {
			label = info.Args()[0].String()
		}

		c.timeMap[label] = time.Now()

		return nil
	})
}

func (c *Console) getTimeEndFunctionCallback() *v8go.FunctionTemplate {
	c.timeMap = make(map[string]time.Time)

	return v8go.NewFunctionTemplate(c.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		label := "default"

		if len(info.Args()) > 0 {
			label = info.Args()[0].String()
		}

		startTime, exists := c.timeMap[label]

		if !exists {
			fmt.Printf("Timer \"%s\" doesn't exist.\n", label)
			return nil
		}

		timeSince := time.Since(startTime).Milliseconds()

		fmt.Printf("%s: %dms - timer end\n", label, timeSince)

		// Delete the entry to free up a timer space.
		delete(c.timeMap, label)

		return nil
	})
}

func (c *Console) getTimeLogFunctionCallback() *v8go.FunctionTemplate {
	c.timeMap = make(map[string]time.Time)

	return v8go.NewFunctionTemplate(c.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		label := "default"

		if len(info.Args()) > 0 {
			label = info.Args()[0].String()
		}

		startTime, exists := c.timeMap[label]

		if !exists {
			fmt.Printf("Timer \"%s\" doesn't exist.\n", label)
			return nil
		}

		timeSince := time.Since(startTime).Milliseconds()

		fmt.Printf("%s: %dms\n", label, timeSince)

		return nil
	})
}

// GetV8Object gets the entire object structure for the console.
func (c *Console) GetV8Object() (*v8go.ObjectTemplate, error) {
	consoleObj := v8go.NewObjectTemplate(c.VM)
	consoleLogFn := c.getLogFunctionCallback()
	consoleCountFn := c.getCountFunctionCallback()
	consoleCountResetFn := c.getCountResetFunctionCallback()
	consoleTimeFn := c.getTimeFunctionCallback()
	consoleTimeEndFn := c.getTimeEndFunctionCallback()
	consoleTimeLogFn := c.getTimeLogFunctionCallback()

	// TODO: The following link may be useful for showing the different log
	// outputs in different colors:
	// https://golangbyexample.com/print-output-text-color-console/

	consoleObj.Set("log", consoleLogFn)
	consoleObj.Set("info", consoleLogFn)
	consoleObj.Set("warn", consoleLogFn)
	consoleObj.Set("error", consoleLogFn)
	consoleObj.Set("debug", consoleLogFn)
	consoleObj.Set("count", consoleCountFn)
	consoleObj.Set("countReset", consoleCountResetFn)
	consoleObj.Set("time", consoleTimeFn)
	consoleObj.Set("timeEnd", consoleTimeEndFn)
	consoleObj.Set("timeLog", consoleTimeLogFn)

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
