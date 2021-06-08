package js

import (
	"strings"

	"github.com/shirou/gopsutil/host"
	js "github.com/wisepythagoras/leebra/js/clipboard"
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
func (nav *Navigator) GetPlatform() (string, error) {
	info, err := host.Info()

	if err != nil {
		return "", err
	}

	return strings.Title(info.OS) + " " + info.KernelArch, nil
}

// GetV8Object creates the V8 object.
func (nav *Navigator) GetV8Object() (*v8go.ObjectTemplate, error) {
	navigatorObj, err := v8go.NewObjectTemplate(nav.VM)

	if err != nil {
		return nil, err
	}

	navigatorObj.Set("userAgent", "My UserAgent", v8go.ReadOnly)
	navigatorObj.Set("cookieEnabled", false, v8go.ReadOnly)
	navigatorObj.Set("doNotTrack", true, v8go.ReadOnly)
	navigatorObj.Set("vendor", "", v8go.ReadOnly)
	navigatorObj.Set("maxTouchPoints", 0, v8go.ReadOnly)
	navigatorObj.Set("webdriver", false, v8go.ReadOnly)
	navigatorObj.Set("javaEnabled", false, v8go.ReadOnly)

	platform, _ := nav.GetPlatform()
	navigatorObj.Set("platform", platform)
	navigatorObj.Set("oscpu", platform)

	clipboardObj, err := nav.GetClipboardObject()

	if err != nil {
		return nil, err
	}

	navigatorObj.Set("clipboard", clipboardObj, v8go.ReadOnly)

	return navigatorObj, nil
}
