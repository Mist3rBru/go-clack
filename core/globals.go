package core

import (
	"errors"
	"os"
	"reflect"
)

var (
	ErrCancelPrompt error = errors.New("prompt canceled")
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
