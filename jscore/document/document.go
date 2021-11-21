package document

import (
	"github.com/wisepythagoras/go-lexbor/html"
	"github.com/wisepythagoras/leebra/jscore/dom"
	"rogchap.com/v8go"
)

// Document defines the Document web API.
type Document struct {
	VM          *v8go.Isolate
	ExecContext *v8go.Context
	Document    *html.Document
	URL         string
	documentObj *v8go.Object
}

// GetElementByIdFunction gets an element by its id.
func (c *Document) GetElementByIdFunction() *v8go.FunctionTemplate {
	return v8go.NewFunctionTemplate(c.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()

		if len(args) < 1 {
			val, _ := v8go.NewValue(c.VM, "Expected argument")
			c.VM.ThrowException(val)
			return nil
		}

		element, err := c.Document.GetElementById(args[0].String())

		if err != nil {
			errVal, _ := v8go.NewValue(c.VM, err.Error())
			c.VM.ThrowException(errVal)
			return nil
		}

		node := &dom.Node{
			VM:          c.VM,
			ExecContext: c.ExecContext,
			Document:    c.Document,
			Element:     element,
			URL:         c.URL,
		}

		nodeObj, _ := node.GetJSObject()

		// html.Serialize(element.Node())

		return nodeObj.Value
	})
}

// GetV8Object gets the entire object structure of the browser Document API.
func (c *Document) GetV8Object() (*v8go.ObjectTemplate, error) {
	documentObj := v8go.NewObjectTemplate(c.VM)
	getElementByIdFn := c.GetElementByIdFunction()

	// TODO: Get the title from the domain or frame context.
	documentObj.Set("title", "This is a test", v8go.ReadOnly)
	documentObj.Set("getElementById", getElementByIdFn, v8go.ReadOnly)

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
