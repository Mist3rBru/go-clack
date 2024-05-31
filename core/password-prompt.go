package core

import (
	"fmt"
	"os"
	"strings"
)

type PasswordPrompt struct {
	Prompt

	Value string
}

type PasswordPromptParams struct {
	Input  *os.File
	Output *os.File
	Value  string
	Render func(p *PasswordPrompt) string
}

func NewPasswordPrompt(params PasswordPromptParams) *PasswordPrompt {
	var p *PasswordPrompt
	p = &PasswordPrompt{
		Prompt: *NewPrompt(PromptParams{
			Input:  params.Input,
			Output: params.Output,
			Value:  params.Value,
			Track:  true,
			Render: func(_p *Prompt) string {
				return params.Render(p)
			},
		}),
		Value: params.Value,
	}
	p.Prompt.On("key", func(args ...any) {
		value, ok := p.Prompt.Value.(string)
		if ok {
			p.Value = value
		}
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

func (p *PasswordPrompt) Run() (string, error) {
	result, err := p.Prompt.Run()
	if err != nil {
		return "", err
	}
	resultStr, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("unexpected result type")
	}
	return resultStr, nil
}
