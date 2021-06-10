package db

import (
	"errors"
	"os"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/wisepythagoras/leebra/utils"
)

type IteratorFn func(txn *badger.Txn) error

// DB defines our database object handler.
type KVDatabase struct {
	Name string
	db   *badger.DB
}

// Exists returns `true` or `false` whether the DB exists.
func (d *KVDatabase) Exists() bool {
	basePath := LocalStoragePath + "/" + d.Name

	// Create the base patch for the chain directory if it doesn't exist.
	if !utils.CheckIfFileExists(basePath) {
		return false
	}

	return true
}

// IsOpen checks if the database is open.
func (d *KVDatabase) IsOpen() bool {
	return d.db != nil
}

// Open opens the database.
func (d *KVDatabase) Open() (bool, error) {
	basePath := LocalStoragePath + "/" + d.Name

	// Create the base patch for the chain directory if it doesn't exist.
	if !utils.CheckIfFileExists(basePath) {
		os.Mkdir(basePath, 0777)
	}

	// Here we define the badger options.
	options := badger.DefaultOptions(basePath + "/" + d.Name)
	options.Logger = nil

	// Now try to open the database.
	db, err := badger.Open(options)

	if err != nil {
		return false, err
	}

	// Save our instance here.
	d.db = db

	return true, nil
}

// Insert creates a new entry in the database.
func (d *KVDatabase) Insert(key []byte, value []byte) (bool, error) {
	if d.db == nil {
		return false, errors.New("Uninitialized database")
	}

	// Create a new transaction.
	txn := d.db.NewTransaction(true)
	defer txn.Discard()

	var err error

	// Set the data in the database.
	if err = txn.Set(key, value); err != nil {
		return false, err
	}

	// Commit the changes to the database.
	if err = txn.Commit(); err != nil {
		return false, err
	}

	return true, nil
}

// Get retrieves the contents of a specific key in the database.
func (d *KVDatabase) Get(key []byte) ([]byte, error) {
	if d.db == nil {
		return nil, errors.New("Uninitialized database")
	}

	// Create a new transaction.
	txn := d.db.NewTransaction(true)
	defer txn.Discard()

	// Get the item of the entry with the hash as the key.
	item, err := txn.Get(key)

	if err != nil {
		return nil, err
	}

	// Get the vaue from the item.
	value, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	return value, nil
}

// Delete deletes the key from the database.
func (d *KVDatabase) Delete(key []byte) error {
	if d.db == nil {
		return errors.New("Uninitialized database")
	}

	// Create a new transaction.
	txn := d.db.NewTransaction(true)
	defer txn.Discard()

	// Delete the item of the entry with the hash as the key.
	err := txn.Delete(key)

	if err != nil {
		return err
	}

	// Commit the changes to the database.
	if err = txn.Commit(); err != nil {
		return err
	}

	return nil
}

// GetKeys gets the list of keys.
func (d *KVDatabase) GetKeys() [][]byte {
	// Create the new transaction, that should not allow updates.
	txn := d.db.NewTransaction(false)
	defer txn.Discard()

	opts := badger.DefaultIteratorOptions

	// We don't want any of the values.
	opts.PrefetchValues = false

	// This array will contain all of our keys and we'll just return this.
	var keys [][]byte = [][]byte{}

	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		keys = append(keys, item.Key())
	}

	return keys
}
