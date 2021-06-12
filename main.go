package main

import (
	"flag"
	"fmt"

	"github.com/wisepythagoras/leebra/browser"
	"github.com/wisepythagoras/leebra/utils"
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

	domainContext := &browser.DomainContext{
		URL: "http://127.0.0.1:8000",
	}

	err = domainContext.ParseURL()

	if err != nil {
		fmt.Println(err)
		return
	}

	jsContext := &browser.JSContext{
		DomainContext: domainContext,
	}

	err = jsContext.Init()

	if err != nil {
		fmt.Println(err)
		return
	}

	vals, errs := jsContext.RunScript(bin, *run)

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
