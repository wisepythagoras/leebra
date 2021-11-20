package dom

import (
	"github.com/wisepythagoras/go-lexbor/html"
)

// https://developer.mozilla.org/en-US/docs/Web/API/Node

// Node defines what will be the Node object in JavaScript.
type Node struct {
	raw      *html.Node
	document *html.Document
}

// Init will initialize all the maps and arrays that are needed.
func (n *Node) Init() {
	//
}
