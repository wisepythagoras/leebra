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

func (n *Node) Parent() *Node {
	parent := n.Element.Node().Parent()

	if parent == nil {
		return nil
	}

	return &Node{
		VM:          n.VM,
		ExecContext: n.ExecContext,
		Document:    n.Document,
		Element:     parent.Element(),
		URL:         n.URL,
	}
}

// GetV8Object gets the entire object structure of the browser Document API.
func (n *Node) GetV8Object(withParent bool) (*v8go.ObjectTemplate, error) {
	nodeObj := v8go.NewObjectTemplate(n.VM)
	children := n.Element.Node().Children()

	nodeObj.Set("accessKey", "", v8go.ReadOnly)
	nodeObj.Set("accessKeyLabel", "", v8go.ReadOnly)
	nodeObj.Set("baseURI", n.URL, v8go.ReadOnly)
	nodeObj.Set("childElementCount", uint32(len(children)), v8go.ReadOnly)
	nodeObj.Set("className", n.ClassName(), v8go.ReadOnly)

	// TODO: This should be done with getters and setters, instead of calling it this way. This
	// code can cause issues if we let it run for every node object that's constructed. Working
	// with getters and setters is not yet supported by v8go.
	if withParent {
		parent := n.Parent()

		if parent != nil {
			parentObj, err := parent.GetV8Object(false)

			if err == nil {
				nodeObj.Set("parentNode", parentObj, v8go.ReadOnly)
			}
		} else {
			nodeObj.Set("parentNode", nil, v8go.ReadOnly)
		}
	}

	return nodeObj, nil
}

// GetJSObject returns the JS Object that can be mutated.
func (n *Node) GetJSObject(withParent bool) (*v8go.Object, error) {
	document, err := n.GetV8Object(withParent)

	if err != nil {
		return nil, err
	}

	n.nodeObj, err = document.NewInstance(n.ExecContext)

	if err != nil {
		return nil, err
	}

	return n.nodeObj, nil
}
