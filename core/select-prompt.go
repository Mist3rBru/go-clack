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
		}
		p.Value = p.Options[p.CursorIndex].Value
	})
	return p
}
