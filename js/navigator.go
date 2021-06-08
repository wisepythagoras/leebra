package js

import (
	"strings"

	js "github.com/wisepythagoras/leebra/js/clipboard"
	"github.com/wisepythagoras/leebra/system"
	"rogchap.com/v8go"
)

type Navigator struct {
	VM *v8go.Isolate
}

// GetClipboardObject creates the V8 object for the Clipboard API.
func (nav *Navigator) GetClipboardObject() (*v8go.ObjectTemplate, error) {
	clipboard := &js.Clipboard{
		VM: nav.VM,
	}

	return clipboard.GetV8Object()
}

// GetPlatform returns the platform information.
func (nav *Navigator) GetPlatform() string {
	return system.GetOS() + " " + system.GetKernelArch()
}

// GetUserAgent returns the user agent string.
func (nav *Navigator) GetUserAgent() string {
	platform := nav.GetPlatform()
	details := strings.Join([]string{
		system.GetFlavor(),
		platform,
		"rv:" + system.Version,
	}[:], "; ")
	browser := system.Name + "/" + system.Version

	return "Mozilla/5.0 (" + details + ") Gecko/20100101 " + browser
}

// GetV8Object creates the V8 object.
func (nav *Navigator) GetV8Object() (*v8go.ObjectTemplate, error) {
	navigatorObj, err := v8go.NewObjectTemplate(nav.VM)

	if err != nil {
		return nil, err
	}

	platform := nav.GetPlatform()

	navigatorObj.Set("userAgent", nav.GetUserAgent(), v8go.ReadOnly)
	navigatorObj.Set("cookieEnabled", false, v8go.ReadOnly)
	navigatorObj.Set("doNotTrack", true, v8go.ReadOnly)
	navigatorObj.Set("vendor", "", v8go.ReadOnly)
	navigatorObj.Set("maxTouchPoints", 0, v8go.ReadOnly)
	navigatorObj.Set("webdriver", false, v8go.ReadOnly)
	navigatorObj.Set("javaEnabled", false, v8go.ReadOnly)
	navigatorObj.Set("product", "Leebra", v8go.ReadOnly)
	navigatorObj.Set("platform", platform, v8go.ReadOnly)
	navigatorObj.Set("oscpu", platform, v8go.ReadOnly)
	navigatorObj.Set("appName", "Netscape", v8go.ReadOnly)
	navigatorObj.Set("appCodeName", system.Name, v8go.ReadOnly)
	navigatorObj.Set("appVersion", "5.0 (X11)", v8go.ReadOnly)
	navigatorObj.Set("language", system.Language, v8go.ReadOnly)

	clipboardObj, err := nav.GetClipboardObject()

	if err != nil {
		return nil, err
	}

	navigatorObj.Set("clipboard", clipboardObj, v8go.ReadOnly)

	return navigatorObj, nil
}
