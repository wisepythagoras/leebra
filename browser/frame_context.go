package browser

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/wisepythagoras/go-lexbor/html"
	"github.com/wisepythagoras/leebra/jscore/dom"
	"github.com/wisepythagoras/leebra/jscore/net"
	"rogchap.com/v8go"
)

type FrameContext struct {
	URL           string
	Title         string
	domainContext *DomainContext
	jsContext     *JSContext
	resp          *http.Response
	document      *html.Document
	// HistoryManager
	// EventBus
}

// InitJSEngine initializes the JS engine in the appropriate engine and attaches all of
// the APIs.
func (bc *FrameContext) InitJSEngine() error {
	// This should happen when the URL has been parsed.
	if bc.domainContext == nil {
		return errors.New("Can't continue without a domain context")
	}

	// No need to re-initialize if it's already there.
	if bc.jsContext != nil {
		return nil
	}

	// Create the JS engine context and then initialize it.
	bc.jsContext = &JSContext{
		DomainContext: bc.domainContext,
		document:      bc.document,
	}

	// The init function will attach all of the browser APIs on it.
	err := bc.jsContext.Init()

	if err != nil {
		return err
	}

	return nil
}

// Load loads the page.
func (bc *FrameContext) Load(newUrl string) error {
	if bc.domainContext == nil {
		bc.domainContext = &DomainContext{}
	}

	bc.domainContext.URL = newUrl
	err := bc.domainContext.ParseURL()

	if err != nil {
		return err
	}

	// TODO: Expose this to the JS context.
	bc.domainContext.SetTitle(newUrl)

	if newUrl != "about:blank" {
		// Somewhere over here also do the HTTP request to get the page.
		bc.resp, err = net.HTTPRequest(newUrl, nil)

		if err != nil {
			return err
		}

		body, err := ioutil.ReadAll(bc.resp.Body)

		if err != nil {
			return err
		}

		bc.document, err = dom.ParseHTML(body)

		if err != nil {
			return err
		}
	}

	// The JS engine has to be set up before the page load starts.
	err = bc.InitJSEngine()

	if err != nil {
		return err
	}

	return nil
}

// RunScript executes a JS script in the execution context.
func (bc *FrameContext) RunScript(src []byte, name string) (chan *v8go.Value, chan error) {
	return bc.jsContext.RunScript(src, name)
}
