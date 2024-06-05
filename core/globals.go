package core

import (
	"errors"
	"os"

	"github.com/Mist3rBru/go-clack/third_party"
)

var (
	color = thirdparty.CreateColors()
)

type SelectOption[TValue comparable] struct {
	Label string
	Value TValue
}

type MultiSelectOption[TValue comparable] struct {
	Label      string
	Value      TValue
	IsSelected bool
}

type Event string

const (
	KeyEvent      Event = "key"
	FinalizeEvent Event = "finalize"
	CancelEvent   Event = "cancel"
	SubmitEvent   Event = "submit"
)

type State string

const (
	InitialState State = "initial"
	ActiveState  State = "active"
	ErrorState   State = "error"
	CancelState  State = "cancel"
	SubmitState  State = "submit"
)

var (
	ErrMissingRender error = errors.New("missing render function error")
	ErrCancelPrompt  error = errors.New("prompt canceled error")
)

type KeyName string

type Key struct {
	Name  KeyName
	Char  string
	Shift bool
	Ctrl  bool
}

const (
	EnterKey     KeyName = "Enter"
	SpaceKey     KeyName = "Space"
	TabKey       KeyName = "Tab"
	UpKey        KeyName = "Up"
	DownKey      KeyName = "Down"
	LeftKey      KeyName = "Left"
	RightKey     KeyName = "Right"
	CancelKey    KeyName = "Cancel"
	HomeKey      KeyName = "Home"
	EndKey       KeyName = "End"
	BackspaceKey KeyName = "Backspace"
)

type FileSystem interface {
	Getwd() (string, error)
	ReadDir(name string) ([]os.DirEntry, error)
}

type OSFileSystem struct{}

func (fs OSFileSystem) Getwd() (string, error) {
	return os.Getwd()
}

func (fs OSFileSystem) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}
