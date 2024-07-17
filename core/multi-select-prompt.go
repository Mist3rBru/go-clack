package core

import (
	"os"

	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/core/validator"
)

type MultiSelectPrompt[TValue comparable] struct {
	Prompt[[]TValue]
	Options  []*MultiSelectOption[TValue]
	Required bool
}

type MultiSelectPromptParams[TValue comparable] struct {
	Input        *os.File
	Output       *os.File
	InitialValue []TValue
	Options      []*MultiSelectOption[TValue]
	Required     bool
	Validate     func(value []TValue) error
	Render       func(p *MultiSelectPrompt[TValue]) string
}

func NewMultiSelectPrompt[TValue comparable](params MultiSelectPromptParams[TValue]) *MultiSelectPrompt[TValue] {
	v := validator.NewValidator("MultiSelectPrompt")
	v.ValidateRender(params.Render)
	v.ValidateOptions(len(params.Options))

	for _, option := range params.Options {
		if value, ok := any(option.Value).(string); ok && value == "" {
			option.Value = any(option.Label).(TValue)
		}
	}

	var p MultiSelectPrompt[TValue]
	p = MultiSelectPrompt[TValue]{
		Prompt: *NewPrompt(PromptParams[[]TValue]{
			Input:        params.Input,
			Output:       params.Output,
			InitialValue: mapMultiSelectInitialValue(params.InitialValue, params.Options),
			Validate:     WrapValidate(params.Validate, &p.Required, "Please select at least one option. Press `space` to select"),
			Render:       WrapRender[[]TValue](&p, params.Render),
		}),
		Options:  params.Options,
		Required: params.Required,
	}

	p.On(KeyEvent, func(args ...any) {
		p.handleKeyPress(args[0].(*Key))
	})

	return &p
}

func (p *MultiSelectPrompt[TValue]) handleKeyPress(key *Key) {
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
}

func mapMultiSelectInitialValue[TValue comparable](value []TValue, options []*MultiSelectOption[TValue]) []TValue {
	var initialValue []TValue

	if len(value) > 0 {
		initialValue = value
		for _, value := range value {
			for _, option := range options {
				if option.Value == value {
					option.IsSelected = true
				}
			}
		}
	} else {
		for _, option := range options {
			if option.IsSelected {
				initialValue = append(initialValue, option.Value)
			}
		}
	}

	return initialValue
}
