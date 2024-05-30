package core

import (
	"fmt"
	"os"
)

type TextPrompt struct {
	Prompt

	Value string
}

type TextPromptOptions struct {
	Input  *os.File
	Output *os.File
	Value  string
	Render func(p *TextPrompt) string
}

func NewTextPrompt(options TextPromptOptions) *TextPrompt {
	var p *TextPrompt
	p = &TextPrompt{
		Prompt: *NewPrompt(PromptOptions{
			Input:  options.Input,
			Output: options.Output,
			Value:  options.Value,
			Track:  true,
			Render: func(_p *Prompt) string {
				return options.Render(p)
			},
		}),
		Value: options.Value,
	}
	p.Prompt.On("key", func(args ...any) {
		value, ok := p.Prompt.Value.(string)
		if ok {
			p.Value = value
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
