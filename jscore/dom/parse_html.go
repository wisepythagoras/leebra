package dom

import (
	"errors"

	"github.com/wisepythagoras/go-lexbor/html"
)

// ParseHTML was adapted from here: https://pkg.go.dev/golang.org/x/net/html#Parse
func ParseHTML(htmlStr []byte) (*html.Document, error) {
	doc := &html.Document{}
	doc.Create()
	success := doc.Parse(string(htmlStr))

	if !success {
		return nil, errors.New("Unable to parse the HTML")
	}

	html.Serialize(doc.BodyElement().Element().Node())

	return doc, nil
}
