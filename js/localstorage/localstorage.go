package js

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
	VM       *v8go.Isolate
	DB       *db.KVDatabase
	Context  string
	length   int32
	onChange func(LocalStorageOp)
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
			fmt.Errorf(err.Error())
			return err
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

		ls.ensureDBIsOpen()

		// Insert the item to the db.
		inserted, err := ls.DB.Insert([]byte(key), []byte(value))

		if err != nil || !inserted {
			// TODO: Throw error?
			return nil
		}

		ls.onChange(LocalStorageInsert)

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

// GetV8Object returns the entire object structure of the V8 LocalStorage API.
func (ls *LocalStorage) GetV8Object() (*v8go.ObjectTemplate, error) {
	// Just initialize some of the moving parts.
	ls.Init()

	// This will contain the LocalStorage API structure.
	localStorageObj, err := v8go.NewObjectTemplate(ls.VM)

	if err != nil {
		return nil, err
	}

	ls.onChange = func(op LocalStorageOp) {
		if op == LocalStorageInsert {
			ls.length += 1
		} else {
			ls.length -= 1
		}

		// Undate the length.
		err := localStorageObj.Set("length", ls.length)

		fmt.Println("Update", ls.length, err)
	}

	// Set the default length to whatever the amount of items in the database.
	localStorageObj.Set("length", ls.length)

	setItemFn, err := ls.SetItemFunction()

	if err != nil {
		return nil, err
	}

	localStorageObj.Set("setItem", setItemFn, v8go.ReadOnly)

	getItemFn, err := ls.GetItemFunction()

	if err != nil {
		return nil, err
	}

	localStorageObj.Set("getItem", getItemFn, v8go.ReadOnly)

	return localStorageObj, nil
}
