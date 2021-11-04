package browser

import (
	"errors"

	"github.com/wisepythagoras/leebra/jscore"
	"github.com/wisepythagoras/leebra/jscore/console"
	c "github.com/wisepythagoras/leebra/jscore/crypto"
	doc "github.com/wisepythagoras/leebra/jscore/document"
	ls "github.com/wisepythagoras/leebra/jscore/localstorage"
	"github.com/wisepythagoras/leebra/jscore/net"
	w "github.com/wisepythagoras/leebra/jscore/wasm"

	"rogchap.com/v8go"
	// "go.kuoruan.net/v8go-polyfills/fetch"
)

// JSContext manages the JavaScript execution context.
type JSContext struct {
	DomainContext *DomainContext
	ctx           *v8go.Context
}

// Init must be run right after the creation of the JSContext instance so that scripts
// can be run.
func (jsc *JSContext) Init() error {
	if jsc.DomainContext == nil {
		return errors.New("Invalid domain context")
	}

	// Creates the new VM to run all of the code in.
	vm := v8go.NewIsolate()

	// This object will create a new object on which we'll place our overrides.
	obj := v8go.NewObjectTemplate(vm)

	// Here we create a new instance of the Navigator object.
	navigator := &jscore.Navigator{VM: vm}
	crypto := &c.Crypto{VM: vm}
	consoleInstance := &console.Console{VM: vm}

	localStorage := &ls.LocalStorage{
		VM:      vm,
		Context: jsc.DomainContext.GetHost(),
	}
	localStorage.Init()

	document := &doc.Document{VM: vm}
	wasm := &w.Wasm{VM: vm}
	wasm.NewEngine()

	fetchFn := net.CreateFetchFn(vm)
	obj.Set("fetch", fetchFn)

	// Create a new context.
	jsc.ctx = v8go.NewContext(vm, obj)
	localStorage.ExecContext = jsc.ctx
	navigator.ExecContext = jsc.ctx
	crypto.ExecContext = jsc.ctx
	wasm.ExecContext = jsc.ctx
	document.ExecContext = jsc.ctx
	consoleInstance.ExecContext = jsc.ctx

	global := jsc.ctx.Global()
	lsObj, _ := localStorage.GetJSObject()
	navObj, _ := navigator.GetJSObject()
	cryptoObj, _ := crypto.GetJSObject()
	wasmObj, _ := wasm.GetJSObject()
	documentObj, _ := document.GetJSObject()
	consoleObj, _ := consoleInstance.GetJSObject()

	global.Set("navigator", navObj)
	global.Set("clientInformation", navObj)
	global.Set("localStorage", lsObj)
	global.Set("crypto", cryptoObj)
	global.Set("fullscreen", false)
	global.Set("frames", []*JSContext{})
	global.Set("length", 0)
	global.Set("isSecureContext", false)
	global.Set("innerHeight", 1024)
	global.Set("innerWidth", 768)
	global.Set("WebAssembly", wasmObj)
	global.Set("document", documentObj)
	global.Set("console", consoleObj)

	// With this hack we create the window object.
	global.Set("window", global)

	return nil
}

// RunScript executes a JS script.
func (jsc *JSContext) RunScript(src []byte, name string) (chan *v8go.Value, chan error) {
	vals := make(chan *v8go.Value, 1)
	errs := make(chan error, 1)

	go func() {
		val, err := jsc.ctx.RunScript(string(src), name)

		if err != nil {
			errs <- err
		}

		vals <- val
	}()

	return vals, errs
}
