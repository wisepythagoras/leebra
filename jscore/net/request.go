package net

import (
	"net/http"

	"github.com/wisepythagoras/leebra/jscore"
	"rogchap.com/v8go"
)

// HTTPRequest performs a network request.
// The options follow this format:
// https://developer.mozilla.org/en-US/docs/Web/API/fetch
func HTTPRequest(url string, options *v8go.Value) (*http.Response, error) {
	method := "GET"

	if options != nil && options.IsObject() {
		optObj := options.Object()

		if optObj.Has("method") {
			methodVal, err := optObj.Get("method")

			if err != nil {
				method = methodVal.String()
			}
		}
	}

	// The request method should come from the options. If there is none defined, then
	// it can default to "GET".
	request, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	if options != nil {
		// Loop through the options and apply the headers as well.
	}

	request.Header.Set("User-Agent", jscore.GetUserAgent())

	client := &http.Client{}

	return client.Do(request)
}
