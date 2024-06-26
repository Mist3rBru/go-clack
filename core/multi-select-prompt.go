package core

import (
	"os"

	"github.com/Mist3rBru/go-clack/core/utils"
)

type MultiSelectPrompt[TValue comparable] struct {
	Prompt[[]TValue]
	Options []*MultiSelectOption[TValue]
}

type MultiSelectPromptParams[TValue comparable] struct {
	Input        *os.File
	Output       *os.File
	InitialValue []TValue
	Options      []*MultiSelectOption[TValue]
	Validate     func(value []TValue) error
	Render       func(p *MultiSelectPrompt[TValue]) string
}

func NewMultiSelectPrompt[TValue comparable](params MultiSelectPromptParams[TValue]) *MultiSelectPrompt[TValue] {
	var initialValue []TValue
	if len(params.InitialValue) > 0 {
		initialValue = params.InitialValue
		for _, value := range params.InitialValue {
			for _, option := range params.Options {
				if option.Value == value {
					option.IsSelected = true
				}
			}
		}
	} else {
		for _, option := range params.Options {
			if option.IsSelected {
				initialValue = append(initialValue, option.Value)
			}
		}
	}

	var p *MultiSelectPrompt[TValue]
	p = &MultiSelectPrompt[TValue]{
		Prompt: *NewPrompt(PromptParams[[]TValue]{
			Input:        params.Input,
			Output:       params.Output,
			InitialValue: initialValue,
			Validate:     params.Validate,
			Render: func(_p *Prompt[[]TValue]) string {
				if params.Render == nil {
					return ErrMissingRender.Error()
				}
				return params.Render(p)
			},
		}),
		Options: params.Options,
	}
	p.On(KeyEvent, func(args ...any) {
		key := args[0].(*Key)
		switch key.Name {
		case UpKey, LeftKey:
			p.CursorIndex = utils.MinMaxIndex(p.CursorIndex-1, len(p.Options))
		case DownKey, RightKey:
			p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+1, len(p.Options))
		case HomeKey:
			p.CursorIndex = 0
		case EndKey:
			p.CursorIndex = len(p.Options) - 1
		case SpaceKey:
			option := p.Options[p.CursorIndex]
			if option.IsSelected {
				option.IsSelected = false
				value := []TValue{}
				for _, v := range p.Value {
					if v != option.Value {
						value = append(value, v)
					}
				}
				p.Value = value
			} else {
				option.IsSelected = true
				p.Value = append(p.Value, option.Value)
			}
		case "a":
			if len(p.Value) == len(p.Options) {
				p.Value = []TValue{}
				for _, option := range p.Options {
					option.IsSelected = false
				}
			} else {
				p.Value = make([]TValue, len(p.Options))
				for i, option := range p.Options {
					option.IsSelected = true
					p.Value[i] = option.Value
				}
			}
		}
	})
	return p
}
