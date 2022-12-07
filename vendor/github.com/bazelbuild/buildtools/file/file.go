// Package file provides utility file operations.
package file

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// ReadFile is like ioutil.ReadFile.
func ReadFile(name string) ([]byte, os.FileInfo, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return nil, nil, err
	}

	data, err := ioutil.ReadFile(name)
	return data, fi, err
}

// WriteFile is like ioutil.WriteFile
func WriteFile(name string, data []byte) error {
	return ioutil.WriteFile(name, data, 0644)
}

// OpenReadFile is like os.Open.
func OpenReadFile(name string) io.ReadCloser {
	f, err := os.Open(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open %s\n", name)
		os.Exit(1)
	}
	return f
}
