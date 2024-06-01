package core

import (
	"os"

	"github.com/Mist3rBru/go-clack/core/utils"
)

type MultiSelectPrompt[TValue comparable] struct {
	Prompt[[]TValue]
	Options []SelectOption[TValue]
}

type MultiSelectPromptParams[TValue comparable] struct {
	Input   *os.File
	Output  *os.File
	Value   []TValue
	Options []SelectOption[TValue]
	Render  func(p *MultiSelectPrompt[TValue]) string
}

func NewMultiSelectPrompt[TValue comparable](params MultiSelectPromptParams[TValue]) *MultiSelectPrompt[TValue] {
	var p *MultiSelectPrompt[TValue]
	p = &MultiSelectPrompt[TValue]{
		Prompt: *NewPrompt(PromptParams[[]TValue]{
			Input:  params.Input,
			Output: params.Output,
			Value:  params.Value,
			Render: func(_p *Prompt[[]TValue]) string {
				return params.Render(p)
			},
		}),
		Options: params.Options,
	}
	p.On(EventKey, func(args ...any) {
		key := args[0].(*Key)
		switch key.Name {
		case KeyUp, KeyLeft:
			p.CursorIndex = utils.MinMaxIndex(p.CursorIndex-1, len(p.Options))
		case KeyDown, KeyRight:
			p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+1, len(p.Options))
		case KeyHome:
			p.CursorIndex = 0
		case KeyEnd:
			p.CursorIndex = len(p.Options) - 1
		case KeySpace:
			option := p.Options[p.CursorIndex]
			if i, isSelected := p.IsSelected(option); isSelected {
				p.Value = append(p.Value[0:i], p.Value[i+1:]...)
			} else {
				p.Value = append(p.Value, option.Value)
			}
		case "a":
			if len(p.Value) == len(p.Options) {
				p.Value = []TValue{}
				break
			}
			p.Value = []TValue{}
			for _, option := range p.Options {
				p.Value = append(p.Value, option.Value)
			}
		}
	})
	return p
}

func (p *MultiSelectPrompt[TValue]) IsSelected(option SelectOption[TValue]) (int, bool) {
	for i, value := range p.Value {
		if value == option.Value {
			return i, true
		}
	}
	return -1, false
}
