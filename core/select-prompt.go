package core

import (
	"os"

	"github.com/Mist3rBru/go-clack/core/utils"
)

type SelectPrompt struct {
	Prompt
	Options []any
}

type SelectPromptOptions struct {
	Input   *os.File
	Output  *os.File
	Value   any
	Options []any
	Render  func(p *SelectPrompt) string
}

func NewSelectPrompt(options SelectPromptOptions) *SelectPrompt {
	var p *SelectPrompt
	p = &SelectPrompt{
		Prompt: *NewPrompt(PromptOptions{
			Input:       options.Input,
			Output:      options.Output,
			Value:       options.Value,
			CursorIndex: utils.IndexOf(options.Value, options.Options),
			Track:       false,
			Render: func(_p *Prompt) string {
				return options.Render(p)
			},
		}),
		Options: options.Options,
	}
	p.Prompt.On("key", func(args ...any) {
		key := args[0].(string)
		switch key {
		case "ArrowUp", "ArrowLeft":
			p.CursorIndex = utils.MinMaxIndex(p.CursorIndex-1, len(p.Options))
		case "ArrowDown", "ArrowRight":
			p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+1, len(p.Options))
		case "Home":
			p.CursorIndex = 0
		case "End":
			p.CursorIndex = len(p.Options) - 1
		}
		p.Value = p.Options[p.CursorIndex]
	})
	return p
}

func (p *SelectPrompt) Run() (any, error) {
	_, err := p.Prompt.Run()
	if err != nil {
		return nil, err
	}
	return p.Value, nil
}
