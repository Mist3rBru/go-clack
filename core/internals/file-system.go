package internals

import "os"

type OSFileSystem struct{}

func (fs OSFileSystem) Getwd() (string, error) {
	return os.Getwd()
}

func (fs OSFileSystem) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}
