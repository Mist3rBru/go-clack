package core

import (
	"os"
	"regexp"

	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/core/validator"
)

type SelectPrompt[TValue comparable] struct {
	Prompt[TValue]
	Filter         bool
	Search         string
	initialOptions []*SelectOption[TValue]
	Options        []*SelectOption[TValue]
}

type SelectPromptParams[TValue comparable] struct {
	Input        *os.File
	Output       *os.File
	InitialValue TValue
	Filter       bool
	Options      []*SelectOption[TValue]
	Render       func(p *SelectPrompt[TValue]) string
}

func NewSelectPrompt[TValue comparable](params SelectPromptParams[TValue]) *SelectPrompt[TValue] {
	v := validator.NewValidator("SelectPrompt")
	v.ValidateRender(params.Render)
	v.ValidateOptions(len(params.Options))

	startIndex := 0
	for i, option := range params.Options {
		if value, ok := any(option.Value).(string); ok && value == "" {
			option.Value = any(option.Label).(TValue)
		}
		if option.Value == params.InitialValue {
			startIndex = i
		}
	}

	var p SelectPrompt[TValue]
	p = SelectPrompt[TValue]{
		Prompt: *NewPrompt(PromptParams[TValue]{
			Input:        params.Input,
			Output:       params.Output,
			InitialValue: params.Options[startIndex].Value,
			CursorIndex:  startIndex,
			Render:       WrapRender[TValue](&p, params.Render),
		}),
		Filter:         params.Filter,
		initialOptions: params.Options,
		Options:        params.Options,
	}

	p.On(KeyEvent, func(args ...any) {
		p.handleKeyPress(args[0].(*Key))
	})

	return &p
}

func (p *SelectPrompt[TValue]) handleKeyPress(key *Key) {
	switch key.Name {
	case UpKey, LeftKey:
		p.CursorIndex = utils.MinMaxIndex(p.CursorIndex-1, len(p.Options))
	case DownKey, RightKey:
		p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+1, len(p.Options))
	case HomeKey:
		p.CursorIndex = 0
	case EndKey:
		p.CursorIndex = len(p.Options) - 1
	case EnterKey, CancelKey:
	default:
		if p.Filter {
			p.Search, _ = p.TrackKeyValue(key, p.Search, len(p.Search))
			p.CursorIndex = 0
			if p.Search == "" {
				p.Options = p.initialOptions
				for i, option := range p.Options {
					if option.Value == p.Value {
						p.CursorIndex = i
						break
					}
				}
			} else {
				p.Options = []*SelectOption[TValue]{}
				for _, option := range p.initialOptions {
					if matched, _ := regexp.MatchString("(?i)"+p.Search, option.Label); matched {
						p.Options = append(p.Options, option)
						if option.Value == p.Value {
							p.CursorIndex = len(p.Options) - 1
						}
					}
				}
			}
		}
	}
	if p.CursorIndex >= 0 && p.CursorIndex < len(p.Options) {
		p.Value = p.Options[p.CursorIndex].Value
	} else {
		p.Value = *new(TValue)
	}
}
