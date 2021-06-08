package utils

import (
	"io/ioutil"
	"os"
)

// CheckIfFileExists checks if the given path exists.
func CheckIfFileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

// WriteToFile writes the contents to the file at the given path.
func WriteToFile(path string, contents []byte) error {
	return ioutil.WriteFile(path, contents, 0644)
}

// ReadFile reads the contents of the file at the given path.
func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// WriteReadOnlyFile writes the contents to the read only file at the given path.
func WriteReadOnlyFile(path string, contents []byte) error {
	return ioutil.WriteFile(path, contents, 0444)
}
