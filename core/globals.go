package core

import (
	"errors"
	"os"
	"path"

	"github.com/Mist3rBru/go-clack/core/internals"
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
	ErrCancelPrompt error = errors.New("prompt canceled error")
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

type PathNode struct {
	Index    int
	Depth    int
	Path     string
	Name     string
	Parent   *PathNode
	Children []*PathNode

	IsSelected bool

	FileSystem  FileSystem
	OnlyShowDir bool
}

type PathNodeOptions struct {
	OnlyShowDir bool
	FileSystem  FileSystem
}

func NewPathNode(rootPath string, options PathNodeOptions) *PathNode {
	if options.FileSystem == nil {
		options.FileSystem = internals.OSFileSystem{}
	}

	root := &PathNode{
		Index:    0,
		Depth:    0,
		Path:     rootPath,
		Name:     rootPath,
		Children: []*PathNode{},

		OnlyShowDir: options.OnlyShowDir,
		FileSystem:  options.FileSystem,
	}
	root.Children = root.MapChildren()
	return root
}

func (p *PathNode) MapChildren() []*PathNode {
	if p.Children == nil {
		return nil
	}
	if len(p.Children) > 0 {
		return p.Children
	}
	entries, err := p.FileSystem.ReadDir(p.Path)
	if err != nil {
		return nil
	}
	children := []*PathNode{}
	for _, entry := range entries {
		if p.OnlyShowDir && !entry.IsDir() {
			continue
		}
		child := &PathNode{
			Index:  len(children),
			Depth:  p.Depth + 1,
			Path:   path.Join(p.Path, entry.Name()),
			Name:   entry.Name(),
			Parent: p,

			FileSystem:  p.FileSystem,
			OnlyShowDir: p.OnlyShowDir,
		}
		if entry.IsDir() {
			child.Children = []*PathNode{}
		}
		children = append(children, child)
	}
	return children
}

func WrapRender[T any, TPrompt any](p TPrompt, render func(p TPrompt) string) func(_ *Prompt[T]) string {
	return func(_ *Prompt[T]) string {
		return render(p)
	}
}

func WrapValidateSlice[T []E, E any](validate func(value T) error, isRequired *bool, msg string) func(value T) error {
	return func(value T) error {
		var err error
		if validate != nil {
			err = validate(value)
		}
		if err == nil && *isRequired && len(value) == 0 {
			err = errors.New(msg)
		}
		return err
	}
}

func WrapValidateString(validate func(value string) error, isRequired *bool, msg string) func(value string) error {
	return func(value string) error {
		var err error
		if validate != nil {
			err = validate(value)
		}
		if err == nil && *isRequired && value == "" {
			err = errors.New(msg)
		}
		return err
	}
}
