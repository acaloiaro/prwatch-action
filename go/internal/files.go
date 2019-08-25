package internal

import "os"

// utilities for working with files

type fileProvider interface {
	Exists(path string) bool
}

type posixFileProvider struct{}

func (p *posixFileProvider) Exists(path string) bool {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
