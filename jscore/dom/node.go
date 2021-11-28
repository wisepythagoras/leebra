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

func (n *Node) Children() []*Node {
	childNodes := make([]*Node, 0)
	children := n.Element.Node().Children()

	for _, child := range children {
		childNode := &Node{
			VM:          n.VM,
			ExecContext: n.ExecContext,
			Document:    n.Document,
			Element:     child.Element(),
			URL:         n.URL,
		}
		childNodes = append(childNodes, childNode)
	}

	return childNodes
}

func (n *Node) FirstChild() *Node {
	firstChild := n.Element.Node().FirstChild()

	if firstChild == nil {
		return nil
	}

	return &Node{
		VM:          n.VM,
		ExecContext: n.ExecContext,
		Document:    n.Document,
		Element:     firstChild.Element(),
		URL:         n.URL,
	}
}

// GetV8Object gets the entire object structure of the browser Document API.
func (n *Node) GetV8Object(withParent bool) (*v8go.ObjectTemplate, error) {
	nodeObj := v8go.NewObjectTemplate(n.VM)
	childObjs := v8go.NewObjectTemplate(n.VM)
	childJSObjs, _ := childObjs.NewInstance(n.ExecContext)
	children := n.Children()
	firstChild := n.FirstChild()

	for i, child := range children {
		childObj, _ := child.GetV8Object(false)
		childJSObjs.SetIdx(uint32(i), childObj)
	}

	nodeObj.Set("accessKey", "", v8go.DontDelete)
	nodeObj.Set("accessKeyLabel", "", v8go.DontDelete)
	nodeObj.Set("baseURI", n.URL, v8go.DontDelete)
	nodeObj.Set("childElementCount", uint32(len(children)), v8go.DontDelete)
	nodeObj.Set("children", childObjs, v8go.DontDelete)
	nodeObj.Set("className", n.ClassName(), v8go.DontDelete)
	nodeObj.Set("firstChild", v8go.Null(n.VM), v8go.DontDelete)

	if firstChild != nil {
		firstChildObj, _ := firstChild.GetV8Object(false)
		nodeObj.Set("firstChild", firstChildObj, v8go.DontDelete)
	}

	// TODO: This should be done with getters and setters, instead of calling it this way. This
	// code can cause issues if we let it run for every node object that's constructed. Working
	// with getters and setters is not yet supported by v8go.
	if withParent {
		parent := n.Parent()

		if parent != nil {
			parentObj, err := parent.GetV8Object(false)

			if err == nil {
				nodeObj.Set("parentNode", parentObj, v8go.DontDelete)
			}
		} else {
			nodeObj.Set("parentNode", nil, v8go.DontDelete)
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
