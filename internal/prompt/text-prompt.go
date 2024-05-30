package prompt

import (
	"fmt"
	"os"
)

type TextRender func(p *TextPrompt) string

type TextPrompt struct {
	Prompt

	Value string
}

type TextPromptOptions struct {
	Input  *os.File
	Output *os.File
	Value  string
	Render TextRender
}

func NewTextPrompt(options TextPromptOptions) *TextPrompt {
	p := &TextPrompt{
		Prompt: *NewPrompt(PromptOptions{
			Input:  options.Input,
			Output: options.Output,
			Value:  options.Value,
			Track:  true,
		}),
		Value: options.Value,
	}
	p.Prompt.Value = ""
	p.Prompt.Render = func(_p *Prompt) string {
		return options.Render(p)
	}
	p.Prompt.On("key", func(args ...any) {
		value, ok := p.Prompt.Value.(string)
		if ok {
			p.Value = value
		}
	})
	return p
}

func DefaultTextPrompt(render TextRender) *TextPrompt {
	return NewTextPrompt(TextPromptOptions{
		Input:  os.Stdin,
		Output: os.Stdout,
		Render: render,
	})
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
