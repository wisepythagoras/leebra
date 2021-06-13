package browser

import (
	"errors"
	"fmt"

	"github.com/wisepythagoras/leebra/jscore"
	c "github.com/wisepythagoras/leebra/jscore/crypto"
	ls "github.com/wisepythagoras/leebra/jscore/localstorage"
	"go.kuoruan.net/v8go-polyfills/console"
	"go.kuoruan.net/v8go-polyfills/fetch"
	"rogchap.com/v8go"
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
	vm, _ := v8go.NewIsolate()

	// This object will create a new object on which we'll place our overrides.
	obj, _ := v8go.NewObjectTemplate(vm)

	// Here we create a new instance of the Navigator object.
	navigator := &jscore.Navigator{VM: vm}
	crypto := &c.Crypto{VM: vm}

	localStorage := &ls.LocalStorage{
		VM:      vm,
		Context: jsc.DomainContext.GetHost(),
	}
	localStorage.Init()

	// This adds the fetch function polyfills.
	if err := fetch.InjectTo(vm, obj); err != nil {
		fmt.Println("Error", err)
	}

	// Create a new context.
	jsc.ctx, _ = v8go.NewContext(vm, obj)
	localStorage.ExecContext = jsc.ctx
	navigator.ExecContext = jsc.ctx
	crypto.ExecContext = jsc.ctx

	// Inject the console polyfill.
	if err := console.InjectTo(jsc.ctx); err != nil {
		fmt.Println("Error", err)
	}

	global := jsc.ctx.Global()
	lsObj, _ := localStorage.GetJSObject()
	navObj, _ := navigator.GetJSObject()
	cryptoObj, _ := crypto.GetJSObject()
	global.Set("navigator", navObj)
	global.Set("localStorage", lsObj)
	global.Set("crypto", cryptoObj)

	// With this hack we create the window object.
	thisObj, _ := global.Get("this")
	global.Set("window", thisObj)

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
