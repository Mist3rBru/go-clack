package core

import (
	"fmt"
	"os"
)

type TextPrompt struct {
	Prompt

	Value string
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
		Prompt: *NewPrompt(PromptParams{
			Input:  params.Input,
			Output: params.Output,
			Value:  params.Value,
			Track:  true,
			Validate: func(value any) error {
				return params.Validate(p.Value)
			},
			Render: func(_p *Prompt) string {
				return params.Render(p)
			},
		}),
		Value: params.Value,
	}
	p.On("key", func(args ...any) {
		p.Value = p.Prompt.Value.(string)
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

func (p *TextPrompt) Run() (string, error) {
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
