package core

import (
	"os"

	"github.com/Mist3rBru/go-clack/core/utils"
)

type SelectPrompt[TValue comparable] struct {
	Prompt[TValue]
	Options []SelectOption[TValue]
}

type SelectPromptParams[TValue comparable] struct {
	Input        *os.File
	Output       *os.File
	InitialValue TValue
	Options      []SelectOption[TValue]
	Render       func(p *SelectPrompt[TValue]) string
}

func NewSelectPrompt[TValue comparable](params SelectPromptParams[TValue]) *SelectPrompt[TValue] {
	startIndex := 0
	for i, option := range params.Options {
		if option.Value == params.InitialValue {
			startIndex = i
			break
		}
	}
	var p *SelectPrompt[TValue]
	p = &SelectPrompt[TValue]{
		Prompt: *NewPrompt(PromptParams[TValue]{
			Input:        params.Input,
			Output:       params.Output,
			InitialValue: params.Options[startIndex].Value,
			CursorIndex:  startIndex,
			Render: func(_p *Prompt[TValue]) string {
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
		}
		p.Value = p.Options[p.CursorIndex].Value
	})
	return p
}
