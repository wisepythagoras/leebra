package dom

import (
	"github.com/wisepythagoras/go-lexbor/html"
	"rogchap.com/v8go"
)

// https://developer.mozilla.org/en-US/docs/Web/API/Node

// Node defines what will be the Node object in JavaScript.
type Node struct {
	VM          *v8go.Isolate
	ExecContext *v8go.Context
	Document    *html.Document
	Element     *html.Element
	URL         string
	nodeObj     *v8go.Object
}

func (n *Node) ClassName() string {
	if n.Element == nil {
		return ""
	}

	return n.Element.Attribute("class")
}

// GetV8Object gets the entire object structure of the browser Document API.
func (n *Node) GetV8Object() (*v8go.ObjectTemplate, error) {
	nodeObj := v8go.NewObjectTemplate(n.VM)

	nodeObj.Set("accessKey", "", v8go.ReadOnly)
	nodeObj.Set("accessKeyLabel", "", v8go.ReadOnly)
	nodeObj.Set("baseURI", n.URL, v8go.ReadOnly)
	nodeObj.Set("childElementCount", 1, v8go.ReadOnly)
	nodeObj.Set("className", n.ClassName(), v8go.ReadOnly)

	return nodeObj, nil
}

// GetJSObject returns the JS Object that can be mutated.
func (n *Node) GetJSObject() (*v8go.Object, error) {
	document, err := n.GetV8Object()

	if err != nil {
		return nil, err
	}

	n.nodeObj, err = document.NewInstance(n.ExecContext)

	if err != nil {
		return nil, err
	}

	return n.nodeObj, nil
}
