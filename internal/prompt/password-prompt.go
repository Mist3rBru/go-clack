package prompt

import (
	"fmt"
	"os"
	"strings"
)

type PasswordPromptRender func(p *PasswordPrompt) string

type PasswordPrompt struct {
	Prompt

	Value string
}

type PasswordPromptOptions struct {
	Input  *os.File
	Output *os.File
	Value  string
	Render PasswordPromptRender
}

func NewPasswordPrompt(options PasswordPromptOptions) *PasswordPrompt {
	p := &PasswordPrompt{
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
