package prompts_test

import (
	"os"
	"testing"
	"time"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

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

func runSelectPath() {
	prompts.SelectPath(prompts.SelectPathParams{
		Message:    message,
		FileSystem: (prompts.FileSystem)(MockFileSystem{}),
	})
}

func TestSelectPathInitialState(t *testing.T) {
	go runSelectPath()
	time.Sleep(time.Millisecond)
	p := test.SelectPathTestingPrompt

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectPathWithOptionChildren(t *testing.T) {
	go runSelectPath()
	time.Sleep(time.Millisecond)
	p := test.SelectPathTestingPrompt

	p.PressKey(&core.Key{Name: core.RightKey})

	assert.Equal(t, core.ActiveState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectPathCancelState(t *testing.T) {
	go runSelectPath()
	time.Sleep(time.Millisecond)

	p := test.SelectPathTestingPrompt
	p.PressKey(&core.Key{Name: core.CancelKey})

	assert.Equal(t, core.CancelState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectPathSubmitState(t *testing.T) {
	go runSelectPath()
	time.Sleep(time.Millisecond)

	p := test.SelectPathTestingPrompt
	p.PressKey(&core.Key{Name: core.EnterKey})

	assert.Equal(t, core.SubmitState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}
