package dom

import "golang.org/x/net/html"

// Node defines what will be the Node object in JavaScript.
type Node struct {
	raw            *html.Node
	ChildNodes     []*Node
	attributes     map[string]string
	htmlAttributes []html.Attribute
}

// Init will initialize all the maps and arrays that are needed.
func (n *Node) Init() {
	n.attributes = make(map[string]string)
	n.ChildNodes = make([]*Node, 0)
}

// HTMLAttributesToAttributes imports Go attributes to the map.
func (n *Node) HTMLAttributesToAttributes(htmlAttributes []html.Attribute) {
	n.htmlAttributes = htmlAttributes

	for _, attr := range htmlAttributes {
		n.attributes[attr.Key] = attr.Val
	}
}

// AddChildNode appends a node to the child list.
func (n *Node) AddChildNode(node *Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// GetAttribute returns the value of an attribute.
func (n *Node) GetAttribute(key string) string {
	return n.attributes[key]
}

// LocalName returns the tag name.
func (n *Node) LocalName() string {
	return n.raw.Data
}

// FirstChild returns the first child of this node.
func (n *Node) FirstChild() *Node {
	if len(n.ChildNodes) == 0 {
		return nil
	}

	return n.ChildNodes[0]
}
