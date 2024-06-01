package core

import (
	"os"
	"strings"
)

type PasswordPrompt struct {
	Prompt[string]
}

type PasswordPromptParams struct {
	Input    *os.File
	Output   *os.File
	Value    string
	Validate func(value string) error
	Render   func(p *PasswordPrompt) string
}

func NewPasswordPrompt(params PasswordPromptParams) *PasswordPrompt {
	var p *PasswordPrompt
	p = &PasswordPrompt{
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
	p.On(EventKey, func(args ...any) {
		p.Value = p.TrackKeyValue(args[0].(*Key), p.Value)
	})
	return p
}

func (p *PasswordPrompt) ValueWithCursor() string {
	inverse := color["inverse"]
	maskedValue := strings.Repeat("*", len(p.Value))
	if p.CursorIndex == len(p.Value) {
		return maskedValue + inverse(" ")
	}
	return maskedValue[0:p.CursorIndex] + inverse(string(maskedValue[p.CursorIndex])) + maskedValue[p.CursorIndex+1:]
}
