package browser

import (
	"net/url"
)

// DomainContext describes the Browser's context.
type DomainContext struct {
	URL    string
	Title  string
	URLObj *url.URL
}

// ParseURL does what the name of the function says it does.
func (d *DomainContext) ParseURL() error {
	var err error
	d.URLObj, err = url.Parse(d.URL)

	return err
}

// GetHost returns the domain. APIs like localStorage, IndexedDB, etc will
// heavily rely on this.
func (d *DomainContext) GetHost() string {
	if d.URLObj == nil {
		return ""
	}

	return d.URLObj.Host
}

// SetTitle updates the title string.
func (d *DomainContext) SetTitle(title string) {
	d.Title = title
}
