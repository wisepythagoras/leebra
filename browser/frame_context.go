package browser

import (
	"errors"
)

type FrameContext struct {
	URL           string
	domainContext *DomainContext
	// JSEngine
	// HistoryManager
	// EventBus
}

func (bc *FrameContext) InitJSEngine() error {
	// This should happen when the URL has been parsed.
	if bc.domainContext == nil {
		return errors.New("Can't continue without a domain context")
	}

	return nil
}

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
