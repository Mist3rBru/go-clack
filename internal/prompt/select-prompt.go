package prompt

import (
	"go-clack/internal/utils"
	"os"
)

type SelectPromptRender func(p *SelectPrompt) string

type SelectPrompt struct {
	Prompt
	Value   any
	Options []any
	Render  func(p *SelectPrompt) []any
}

type SelectPromptOptions struct {
	Input   *os.File
	Output  *os.File
	Value   any
	Options []any
	Render  SelectPromptRender
}

func NewSelectPrompt(options SelectPromptOptions) *SelectPrompt {
	p := &SelectPrompt{
		Prompt: *NewPrompt(PromptOptions{
			Input:       options.Input,
			Output:      options.Output,
			Value:       options.Value,
			CursorIndex: utils.IndexOf(options.Value, options.Options),
			Track:       false,
		}),
		Value:   options.Value,
		Options: options.Options,
	}
	p.Prompt.Value = []any{}
	p.Prompt.Render = func(_p *Prompt) string {
		return options.Render(p)
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
