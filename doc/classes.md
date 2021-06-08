# Classes

Inspired from [this](https://github.com/rogchap/v8go/issues/122) issue.

``` go
func PersonClassObj(vm *v8go.Isolate) (*v8go.FunctionTemplate, error) {
	type Person struct {
		Age  int64
		Name string
	}

	var person *Person = nil

	personFn, err := v8go.NewFunctionTemplate(vm, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		name := args[0].String()
		age := args[1].Integer()

		person = &Person{
			Age:  age,
			Name: name,
		}

		getNameFn, _ := v8go.NewFunctionTemplate(vm, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if person == nil {
				return nil
			}

			personName, _ := v8go.NewValue(vm, person.Name)
			return personName
		})

		getAgeFn, _ := v8go.NewFunctionTemplate(vm, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if person == nil {
				return nil
			}

			personAge, _ := v8go.NewValue(vm, person.Age)
			return personAge
		})

		prototype, _ := v8go.NewObjectTemplate(vm)
		prototype.Set("getName", getNameFn)
		prototype.Set("getAge", getAgeFn)

		obj, _ := prototype.NewInstance(info.Context())

		return obj.Value
	})

	if err != nil {
		return nil, err
	}

	return personFn, nil
}
```