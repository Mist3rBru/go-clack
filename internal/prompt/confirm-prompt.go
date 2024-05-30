package prompt

import (
	"go-clack/internal/utils"
	"os"
)

type ConfirmPrompt struct {
	Prompt
	Active   string
	Inactive string
	Value    bool
}

type ConfirmPromptOptions struct {
	Input    *os.File
	Output   *os.File
	Active   string
	Inactive string
	Value    bool
	Render   func(p *ConfirmPrompt) string
}

func NewConfirmPrompt(options ConfirmPromptOptions) *ConfirmPrompt {
	var p *ConfirmPrompt
	p = &ConfirmPrompt{
		Prompt: *NewPrompt(PromptOptions{
			Input:  options.Input,
			Output: options.Output,
			Value:  options.Value,
			Track:  false,
			Render: func(_p *Prompt) string {
				return options.Render(p)
			},
		}),
		Active:   options.Active,
		Inactive: options.Inactive,
		Value:    options.Value,
	}
	p.Prompt.On("key", func(args ...any) {
		key := args[0].(string)
		switch key {
		case "ArrowUp", "ArrowLeft":
			p.CursorIndex = utils.MinMaxIndex(p.CursorIndex-1, 1)
			p.Value = !p.Value
		case "ArrowDown", "ArrowRight":
			p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+1, 1)
			p.Value = !p.Value
		}
	})
	return p
}

func (p *ConfirmPrompt) Run() (bool, error) {
	_, err := p.Prompt.Run()
	if err != nil {
		return false, err
	}
	return p.Value, nil
}
