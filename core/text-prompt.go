package core

import (
	"os"
)

type TextPrompt struct {
	Prompt[string]
}

type TextPromptParams struct {
	Input    *os.File
	Output   *os.File
	Value    string
	Validate func(value string) error
	Render   func(p *TextPrompt) string
}

func NewTextPrompt(params TextPromptParams) *TextPrompt {
	var p *TextPrompt
	p = &TextPrompt{
		Prompt: *NewPrompt(PromptParams[string]{
			Input:       params.Input,
			Output:      params.Output,
			Value:       params.Value,
			CursorIndex: len(params.Value),
			Validate:    params.Validate,
			Render: func(_p *Prompt[string]) string {
				return params.Render(p)
			},
		}),
	}
	p.On(PromptEventKey, func(args ...any) {
		p.Value = p.TrackKeyValue(args[0].(*Key), p.Value)
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
