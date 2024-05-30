package core

import (
	"os"

	"github.com/Mist3rBru/go-clack/core/utils"
)

type MultiSelectPrompt struct {
	Prompt
	Value   []any
	Options []any
}

type MultiSelectPromptOptions struct {
	Input   *os.File
	Output  *os.File
	Value   []any
	Options []any
	Render  func(p *MultiSelectPrompt) string
}

func NewMultiSelectPrompt(options MultiSelectPromptOptions) *MultiSelectPrompt {
	var p *MultiSelectPrompt
	p = &MultiSelectPrompt{
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
		Value:   options.Value,
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
		case "Space":
			option := p.Options[p.CursorIndex]
			if i, isSelected := p.IsSelected(option); isSelected {
				p.Value = append(p.Value[0:i], p.Value[i+1:]...)
				break
			}
			p.Value = append(p.Value, option)
		case "a":
			if len(p.Value) == len(p.Options) {
				p.Value = []any{}
				break
			}
			p.Value = append([]any{}, p.Options...)
		}
	})
	return p
}

func (p *MultiSelectPrompt) IsSelected(option any) (int, bool) {
	for i, value := range p.Value {
		if value == option {
			return i, true
		}
	}
	return -1, false
}

func (p *MultiSelectPrompt) Run() ([]any, error) {
	_, err := p.Prompt.Run()
	if err != nil {
		return nil, err
	}
	return p.Value, nil
}
