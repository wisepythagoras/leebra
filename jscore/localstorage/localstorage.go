package jscore

import (
	"errors"
	"fmt"

	"github.com/wisepythagoras/leebra/db"
	"rogchap.com/v8go"
)

type LocalStorageOp int8

const (
	LocalStorageInsert LocalStorageOp = iota
	LocalStorageDelete LocalStorageOp = iota
)

// LocalStorage defines the LocalStorage API.
type LocalStorage struct {
	VM              *v8go.Isolate
	DB              *db.KVDatabase
	ExecContext     *v8go.Context
	Context         string
	length          int32
	keys            [][]byte
	localStorageObj *v8go.Object
	onChange        func(LocalStorageOp)
}

// Init creates the key-value database
func (ls *LocalStorage) Init() {
	if ls.DB != nil {
		return
	}

	ls.DB = &db.KVDatabase{
		Name: ls.Context,
	}
	ls.length = 0
}

// ensureDBIsOpen ensures that the key-value DB is open.
func (ls *LocalStorage) ensureDBIsOpen() error {
	if ls.DB == nil {
		ls.Init()
	}

	if !ls.DB.IsOpen() {
		opened, err := ls.DB.Open()

		if err != nil || !opened {
			err := errors.New("Unable to open db")
			fmt.Println(err.Error())
			return err
		}

		// Get the keys in the db.
		ls.keys = ls.DB.GetKeys()
		ls.length = int32(len(ls.keys))

		// If the values are different, then update.
		if ls.localStorageObj != nil {
			val, _ := ls.localStorageObj.Get("length")

			if val.Int32() != ls.length {
				ls.localStorageObj.Delete("length")
				ls.localStorageObj.Set("length", ls.length)
			}
		}
	}

	return nil
}

// SetItemFunction sets an item to the DB.
func (ls *LocalStorage) SetItemFunction() (*v8go.FunctionTemplate, error) {
	setItemFn, err := v8go.NewFunctionTemplate(ls.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()

		if len(args) < 2 {
			// TODO: Figure out how to return errors here.
			return nil
		}

		// Here we convert the arguments, whatever they may be, to a string.
		key := args[0].Object().String()
		value := args[1].Object().String()

		err := ls.ensureDBIsOpen()

		if err != nil {
			fmt.Println(err)
			info.Context().RunScript("throw new Error('Unable to access index')", "")
			return nil
		}

		// Check if the key is already in the DB.
		existingData, _ := ls.DB.Get([]byte(key))

		// Insert the item to the db.
		inserted, err := ls.DB.Insert([]byte(key), []byte(value))

		if err != nil || !inserted {
			info.Context().RunScript("throw new Error('Unable to create record for key')", "")
			return nil
		}

		if existingData == nil {
			ls.onChange(LocalStorageInsert)
			ls.keys = append(ls.keys, []byte(key))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return setItemFn, nil
}

// GetItemFunction sets an item to the DB.
func (ls *LocalStorage) GetItemFunction() (*v8go.FunctionTemplate, error) {
	getItemFn, err := v8go.NewFunctionTemplate(ls.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()

		if len(args) < 1 {
			// TODO: Figure out how to return errors here.
			return nil
		}

		// Here we convert the key to a string.
		key := args[0].Object().String()

		ls.ensureDBIsOpen()

		// Get the item to the db.
		value, _ := ls.DB.Get([]byte(key))

		if value == nil {
			return nil
		}

		// Create a new V8 value for the value in the db.
		val, _ := v8go.NewValue(ls.VM, string(value))

		return val
	})

	if err != nil {
		return nil, err
	}

	return getItemFn, nil
}

// RemoveItemFunction defines the function that removes an item from the DB.
func (ls *LocalStorage) RemoveItemFunction() (*v8go.FunctionTemplate, error) {
	removeItemFn, err := v8go.NewFunctionTemplate(ls.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()

		if len(args) < 1 {
			// TODO: Figure out how to return errors here.
			return nil
		}

		// Here we convert the key to a string.
		key := args[0].String()

		ls.ensureDBIsOpen()

		// Get the item to the db.
		err := ls.DB.Delete([]byte(key))

		if err != nil {
			ls.ExecContext.RunScript("throw new Error('Unable to remove key')", "")
			return nil
		}

		ls.onChange(LocalStorageDelete)

		indexToRemove := -1

		for i, lsKey := range ls.keys {
			if string(lsKey) == key {
				indexToRemove = i
				break
			}
		}

		if indexToRemove > -1 {
			ls.keys = append(ls.keys[:indexToRemove], ls.keys[indexToRemove+1:]...)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return removeItemFn, nil
}

// KeyFunction returns the key, given an index.
func (ls *LocalStorage) KeyFunction() (*v8go.FunctionTemplate, error) {
	keyFn, err := v8go.NewFunctionTemplate(ls.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()

		if len(args) < 1 {
			// TODO: Figure out how to return errors here.
			return nil
		}

		// Here we convert the idx to a string.
		idx := args[0]

		if idx.IsNullOrUndefined() {
			return nil
		}

		intIndex := idx.Int32()

		if intIndex > ls.length {
			return nil
		}

		ls.ensureDBIsOpen()

		// Create a new V8 value for the value in the db.
		val, _ := v8go.NewValue(ls.VM, string(ls.keys[intIndex]))

		return val
	})

	if err != nil {
		return nil, err
	}

	return keyFn, nil
}

// ClearFunction clears all keys from the database.
func (ls *LocalStorage) ClearFunction() (*v8go.FunctionTemplate, error) {
	keyFn, err := v8go.NewFunctionTemplate(ls.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		if ls.length > 0 {
			ls.ensureDBIsOpen()

			for _, key := range ls.keys {
				// Get the item to the db.
				err := ls.DB.Delete([]byte(key))

				if err != nil {
					ls.ExecContext.RunScript("throw new Error('Unable to remove key')", "")
				} else if ls.onChange != nil {
					ls.onChange(LocalStorageDelete)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return keyFn, nil
}

// GetV8Object returns the entire object structure of the V8 LocalStorage API.
func (ls *LocalStorage) GetV8Object() (*v8go.ObjectTemplate, error) {
	// Just initialize some of the moving parts.
	ls.Init()

	// This will contain the LocalStorage API structure.
	localStorage, err := v8go.NewObjectTemplate(ls.VM)

	if err != nil {
		return nil, err
	}

	// Set the default length to whatever the amount of items in the database.
	localStorage.Set("length", ls.length)

	setItemFn, err := ls.SetItemFunction()

	if err != nil {
		return nil, err
	}

	localStorage.Set("setItem", setItemFn, v8go.ReadOnly)

	getItemFn, err := ls.GetItemFunction()

	if err != nil {
		return nil, err
	}

	localStorage.Set("getItem", getItemFn, v8go.ReadOnly)

	removeItemFn, err := ls.RemoveItemFunction()

	if err != nil {
		return nil, err
	}

	localStorage.Set("removeItem", removeItemFn, v8go.ReadOnly)

	keyFn, err := ls.KeyFunction()

	if err != nil {
		return nil, err
	}

	localStorage.Set("key", keyFn, v8go.ReadOnly)

	clearFn, err := ls.ClearFunction()

	if err != nil {
		return nil, err
	}

	localStorage.Set("clear", clearFn, v8go.ReadOnly)

	return localStorage, nil
}

// GetJSObject returns the JS Object that can be mutated.
func (ls *LocalStorage) GetJSObject() (*v8go.Object, error) {
	localStorage, err := ls.GetV8Object()

	if err != nil {
		return nil, err
	}

	ls.localStorageObj, err = localStorage.NewInstance(ls.ExecContext)

	if err != nil {
		return nil, err
	}

	// Define the internal onChange event handler.
	ls.onChange = func(op LocalStorageOp) {
		if op == LocalStorageInsert {
			ls.length += 1
		} else {
			ls.length -= 1
		}

		ls.localStorageObj.Delete("length")

		ls.localStorageObj.Set("length", ls.length)
	}

	return ls.localStorageObj, nil
}
