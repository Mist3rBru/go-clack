package core_test

import "os"

type MockDirEntry struct {
	name  string
	isDir bool
}

func (e MockDirEntry) Name() string {
	return e.name
}

func (e MockDirEntry) IsDir() bool {
	return e.isDir
}

func (e MockDirEntry) Type() os.FileMode {
	if e.isDir {
		return os.ModeDir
	}
	return 0
}

func (e MockDirEntry) Info() (os.FileInfo, error) {
	return nil, nil
}

type MockFileSystem struct{}

func (fs MockFileSystem) Getwd() (string, error) {
	return "/clack", nil
}

func (fs MockFileSystem) ReadDir(name string) ([]os.DirEntry, error) {
	return []os.DirEntry{
		MockDirEntry{name: "dir", isDir: true},
		MockDirEntry{name: "file", isDir: false},
	}, nil
}

func (fs MockFileSystem) UserHomeDir() (string, error) {
	return "/home/clack", nil
}
