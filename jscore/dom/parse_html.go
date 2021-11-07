package dom

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"
)

// ParseHTML was adapted from here: https://pkg.go.dev/golang.org/x/net/html#Parse
func ParseHTML(htmlStr []byte) {
	doc, err := html.Parse(strings.NewReader(string(htmlStr)))

	if err != nil {
		log.Fatal(err)
	}

	var f func(*html.Node, string)

	f = func(n *html.Node, indent string) {
		if n.Type == html.ElementNode {
			// Create the DOM object here.
			fmt.Print(indent + n.Data)

			if n.FirstChild != nil && n.FirstChild.Type != html.ElementNode {
				childData := strings.Trim(strings.Trim(n.FirstChild.Data, "\n"), " ")

				if childData != "" {
					fmt.Print(" => '", childData, "'")
				}
			}

			fmt.Println()
		}

		newIndent := indent + "    "

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, newIndent)
		}
	}

	f(doc, "")

	fmt.Println("--------------")
}
