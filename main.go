package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/wisepythagoras/leebra/browser"
	"github.com/wisepythagoras/leebra/utils"
	"rogchap.com/v8go"
)

var wg sync.WaitGroup

func main() {
	run := flag.String("run", "", "Runs a JavaScript file")
	url := flag.String("url", "about:blank", "The URL to load")

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

	wg.Add(1)

	// Wrap in its own go routine.
	go func() {
		defer wg.Done()

		frameContext := &browser.FrameContext{}
		err = frameContext.Load(*url)

		if err != nil {
			fmt.Println(err)
			return
		}

		vals, errs := frameContext.RunScript(bin, *run)

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
	}()

	wg.Wait()
}
