package core

import (
	"os"
	"strings"

	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type PasswordPrompt struct {
	Prompt[string]
}

type PasswordPromptParams struct {
	Input        *os.File
	Output       *os.File
	InitialValue string
	Validate     func(value string) error
	Render       func(p *PasswordPrompt) string
}

func NewPasswordPrompt(params PasswordPromptParams) *PasswordPrompt {
	var p *PasswordPrompt
	p = &PasswordPrompt{
		Prompt: *NewPrompt(PromptParams[string]{
			Input:        params.Input,
			Output:       params.Output,
			InitialValue: params.InitialValue,
			CursorIndex:  len(params.InitialValue),
			Validate:     params.Validate,
			Render: func(_p *Prompt[string]) string {
				if params.Render == nil {
					return ErrMissingRender.Error()
				}
				return params.Render(p)
			},
		}),
	}
	p.On(KeyEvent, func(args ...any) {
		p.TrackKeyValue(args[0].(*Key), &p.Value)
	})
	return p
}

func (p *PasswordPrompt) ValueWithMask() string {
	return strings.Repeat("*", len(p.Value))
}

func (p *PasswordPrompt) ValueWithMaskAndCursor() string {
	maskedValue := strings.Repeat("*", len(p.Value))
	if p.CursorIndex == len(p.Value) {
		return maskedValue + picocolors.Inverse(" ")
	}
	return maskedValue[0:p.CursorIndex] + picocolors.Inverse(string(maskedValue[p.CursorIndex])) + maskedValue[p.CursorIndex+1:]
}
