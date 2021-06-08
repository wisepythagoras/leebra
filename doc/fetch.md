# Fetch

There's no fetch function. I'm using [the polyfill](https://github.com/kuoruan/v8go-polyfills), but here's a simple implementation which is in [the documentation](https://pkg.go.dev/rogchap.com/v8go#FunctionTemplate).

``` go
func FetchFn(vm *v8go.Isolate) (*v8go.FunctionTemplate, error) {
	fetchFn, err := v8go.NewFunctionTemplate(vm, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		url := args[0].String()

		resolver, _ := v8go.NewPromiseResolver(info.Context())

		go func() {
			res, _ := http.Get(url)
			body, _ := ioutil.ReadAll(res.Body)
			val, _ := v8go.NewValue(vm, string(body))
			resolver.Resolve(val)
		}()
		return resolver.GetPromise().Value
	})

	if err != nil {
		return nil, err
	}

	return fetchFn, nil
}
```

And this can be called like this:

``` go
val, err = ctx.RunScript(`
fetch('https://example.com/')
	.then(resp => resp.text());
`, "")

if err != nil {
	fmt.Println(err)
}

prom, err := val.AsPromise()

if err != nil {
	fmt.Println(err)
}

// Wait for the promise to resolve.
for prom.State() == v8go.Pending {
	continue
}

fmt.Println(prom.Result().String())
```