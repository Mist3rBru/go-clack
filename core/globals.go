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

type Event int

const (
	// KeyEvent is emitted after each user's input
	KeyEvent Event = iota
	// ValidateEvent is emitted when the input is being validated
	ValidateEvent
	// ErrorEvent is emitted if an error occurs during the validation process
	ErrorEvent
	// FinalizeEvent is emitted on user's submit or cancel, and before rendering the related state
	FinalizeEvent
	// CancelEvent is emitted after the user cancels the prompt, and after rendering the cancel state
	CancelEvent
	// SubmitEvent is emitted after the user submits the input, and after rendering the submit state
	SubmitEvent
)

type State int

const (
	// InitialState is the initial state of the prompt
	InitialState State = iota
	// ActiveState is set after the user's first action
	ActiveState
	// ValidateState is set after 400ms of validation (e.g., checking user input)
	ValidateState
	// ErrorState is set if there is an error during validation
	ErrorState
	// CancelState is set after the user cancels the prompt
	CancelState
	// SubmitState is set after the user submits the input
	SubmitState
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

func WrapValidate[TValue any](validate func(value TValue) error, isRequired *bool, errMsg string) func(value TValue) error {
	return func(value TValue) error {
		if validate == nil && !*isRequired {
			return nil
		}

		if validate != nil {
			if err := validate(value); err != nil {
				return err
			}
		}

		if *isRequired {
			v := reflect.ValueOf(value)
			errRequired := errors.New(errMsg)

			if !v.IsValid() {
				return errRequired
			}

			k := v.Kind()
			if (k == reflect.Ptr || k == reflect.Interface) && v.IsNil() {
				return errRequired
			}

			if k != reflect.Bool &&
				((k == reflect.Slice && v.Len() == 0) ||
					(k == reflect.Array && v.Len() == 0) ||
					(k == reflect.Map && v.Len() == 0) ||
					(k == reflect.Struct && v.IsZero()) ||
					v.IsZero()) {
				return errRequired
			}
		}

		return nil
	}
}
