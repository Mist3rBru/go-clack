package core

import (
	"os"

	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type TextPrompt struct {
	Prompt[string]
	Placeholder string
}

type TextPromptParams struct {
	Input        *os.File
	Output       *os.File
	InitialValue string
	Placeholder  string
	Validate     func(value string) error
	Render       func(p *TextPrompt) string
}

func NewTextPrompt(params TextPromptParams) *TextPrompt {
	var p *TextPrompt
	p = &TextPrompt{
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
		Placeholder: params.Placeholder,
	}
	p.On(KeyEvent, func(args ...any) {
		key := args[0].(*Key)
		if key.Name == TabKey && p.Value == "" && p.Placeholder != "" {
			p.Value = p.Placeholder
			p.CursorIndex = len(p.Placeholder)
			return
		} else {
			p.TrackKeyValue(key, &p.Value)
		}
	})
	return p
}

func (p *TextPrompt) ValueWithCursor() string {
	if p.CursorIndex == len(p.Value) {
		return p.Value + picocolors.Inverse(" ")
	}
	return p.Value[0:p.CursorIndex] + picocolors.Inverse(string(p.Value[p.CursorIndex])) + p.Value[p.CursorIndex+1:]
}
