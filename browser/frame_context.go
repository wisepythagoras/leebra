package browser

import (
	"errors"

	"rogchap.com/v8go"
)

type FrameContext struct {
	URL           string
	domainContext *DomainContext
	jsContext     *JSContext
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

	// The JS engine has to be set up before the page load starts.
	err = bc.InitJSEngine()

	if err != nil {
		return err
	}

	// Somewhere over here also do the HTTP request to get the page.

	return nil
}

// RunScript executes a JS script in the execution context.
func (bc *FrameContext) RunScript(src []byte, name string) (chan *v8go.Value, chan error) {
	return bc.jsContext.RunScript(src, name)
}
