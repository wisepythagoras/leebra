package dom

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// ParseHTML was adapted from here: https://pkg.go.dev/golang.org/x/net/html#Parse
func ParseHTML(htmlStr []byte) (*Node, error) {
	doc, err := html.Parse(strings.NewReader(string(htmlStr)))
	// nodes := map[string]*Node

	if err != nil {
		return nil, err
	}

	var f func(*html.Node, string) *Node

	f = func(n *html.Node, indent string) *Node {
		if n.Type == html.ElementNode {
			// Create the DOM object here.
			fmt.Print(indent+n.Data, n.Attr)

			if n.FirstChild != nil && n.FirstChild.Type != html.ElementNode {
				childData := strings.Trim(strings.Trim(n.FirstChild.Data, "\n"), " ")

				if childData != "" {
					fmt.Print(" => '", childData, "'")
				}
			}

			fmt.Println()
		}

		newIndent := indent + "    "
		newNode := &Node{raw: n}
		newNode.Init()

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			newNode.AddChildNode(f(c, newIndent))
		}

		newNode.HTMLAttributesToAttributes(n.Attr)
		return newNode
	}

	return f(doc, ""), nil
}
