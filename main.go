package main

import (
	"flag"
	"fmt"

	"github.com/wisepythagoras/leebra/jscore"
	ls "github.com/wisepythagoras/leebra/jscore/localstorage"
	"github.com/wisepythagoras/leebra/utils"
	"go.kuoruan.net/v8go-polyfills/console"
	"go.kuoruan.net/v8go-polyfills/fetch"
	"rogchap.com/v8go"
)

func main() {
	run := flag.String("run", "", "Runs a JavaScript file")

	flag.Parse()

	if *run == "" {
		fmt.Println("Invalid or no JS file passed with -run [file.js]")
		return
	}

	if *run == "" {
		fmt.Println("Invalid or no such JS file")
		return
	}

	// Read the JS file.
	bin, err := utils.ReadFile(*run)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Creates the new VM to run all of the code in.
	vm, _ := v8go.NewIsolate()

	// This object will create a new object on which we'll place our overrides.
	obj, _ := v8go.NewObjectTemplate(vm)

	// Here we create a new instance of the Navigator object.
	navigator := &jscore.Navigator{
		VM: vm,
	}

	localStorage := &ls.LocalStorage{
		VM:      vm,
		Context: "about:blank",
	}
	localStorage.Init()

	// This adds the fetch function polyfills.
	if err := fetch.InjectTo(vm, obj); err != nil {
		fmt.Println("Error", err)
	}

	// Create a new context.
	ctx, _ := v8go.NewContext(vm, obj)
	localStorage.ExecContext = ctx
	navigator.ExecContext = ctx

	// Inject the console polyfill.
	if err := console.InjectTo(ctx); err != nil {
		fmt.Println("Error", err)
	}

	global := ctx.Global()
	lsObj, _ := localStorage.GetJSObject()
	navObj, _ := navigator.GetJSObject()
	global.Set("navigator", navObj)
	global.Set("localStorage", lsObj)

	// With this hack we create the window object.
	thisObj, _ := global.Get("this")
	global.Set("window", thisObj)

	vals := make(chan *v8go.Value, 1)
	errs := make(chan error, 1)

	go func() {
		val, err := ctx.RunScript(string(bin), *run)

		if err != nil {
			errs <- err
		}

		vals <- val
	}()

	select {
	case val := <-vals:
		if val.IsPromise() {
			prom, err := val.AsPromise()

			if err != nil {
				fmt.Println(err)
			}

			// Wait for the promise to resolve.
			for prom.State() == v8go.Pending {
				continue
			}
		} else {
			fmt.Println(val)
		}
	case err := <-errs:
		e := err.(*v8go.JSError)
		fmt.Printf("%+v\n", e)
	}
}
