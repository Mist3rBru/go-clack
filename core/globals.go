package core

import (
	"errors"
	"os"
	"reflect"
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
	InitialState  State = "initial"
	ActiveState   State = "active"
	ValidateState State = "validate"
	ErrorState    State = "error"
	CancelState   State = "cancel"
	SubmitState   State = "submit"
)

var (
	ErrCancelPrompt error = errors.New("prompt canceled")
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
	UserHomeDir() (string, error)
}

func WrapRender[T any, TPrompt any](p TPrompt, render func(p TPrompt) string) func(_ *Prompt[T]) string {
	return func(_ *Prompt[T]) string {
		return render(p)
	}
}

func WrapValidate[TValue any](validate func(value TValue) error, isRequired *bool, msg string) func(value TValue) error {
	return func(value TValue) error {
		var err error
		if validate != nil {
			err = validate(value)
		}
		if err == nil && *isRequired {
			v := reflect.ValueOf(value)
			if !v.IsValid() {
				return errors.New(msg)
			}

			k := v.Kind()
			if (k == reflect.Ptr || k == reflect.Interface) && v.IsNil() {
				err = errors.New(msg)
			} else if k != reflect.Bool &&
				((k == reflect.Slice && v.Len() == 0) ||
					(k == reflect.Array && v.Len() == 0) ||
					(k == reflect.Map && v.Len() == 0) ||
					(k == reflect.Struct && v.IsZero()) ||
					v.IsZero()) {
				err = errors.New(msg)
			}
		}
		return err
	}
}
