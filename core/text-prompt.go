package core

import (
	"os"
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
				return params.Render(p)
			},
		}),
		Placeholder: params.Placeholder,
	}
	p.On(EventKey, func(args ...any) {
		key := args[0].(*Key)
		if key.Name == KeyTab && p.Value == "" && p.Placeholder != "" {
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
	inverse := color["inverse"]
	if p.CursorIndex == len(p.Value) {
		return p.Value + inverse(" ")
	}
	return p.Value[0:p.CursorIndex] + inverse(string(p.Value[p.CursorIndex])) + p.Value[p.CursorIndex+1:]
}
