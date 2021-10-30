package jscore

import (
	"rogchap.com/v8go"
)

// Document defines the Document web API.
type Document struct {
	VM          *v8go.Isolate
	ExecContext *v8go.Context
	documentObj *v8go.Object
}

// GetV8Object gets the entire object structure of the browser Document API.
func (c *Document) GetV8Object() (*v8go.ObjectTemplate, error) {
	documentObj := v8go.NewObjectTemplate(c.VM)

	// TODO: Get the title from the domain or frame context.
	documentObj.Set("title", "This is a test", v8go.ReadOnly)

	return documentObj, nil
}

// GetJSObject returns the JS Object that can be mutated.
func (c *Document) GetJSObject() (*v8go.Object, error) {
	document, err := c.GetV8Object()

	if err != nil {
		return nil, err
	}

	c.documentObj, err = document.NewInstance(c.ExecContext)

	if err != nil {
		return nil, err
	}

	return c.documentObj, nil
}
