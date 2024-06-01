package core

import (
	"os"

	"github.com/Mist3rBru/go-clack/core/utils"
)

type ConfirmPrompt struct {
	Prompt[bool]
	Active   string
	Inactive string
}

type ConfirmPromptParams struct {
	Input    *os.File
	Output   *os.File
	Active   string
	Inactive string
	Value    bool
	Render   func(p *ConfirmPrompt) string
}

func NewConfirmPrompt(params ConfirmPromptParams) *ConfirmPrompt {
	var p *ConfirmPrompt
	p = &ConfirmPrompt{
		Prompt: *NewPrompt(PromptParams[bool]{
			Input:  params.Input,
			Output: params.Output,
			Value:  params.Value,
			Render: func(_p *Prompt[bool]) string {
				return params.Render(p)
			},
		}),
		Active:   params.Active,
		Inactive: params.Inactive,
	}
	p.On(PromptEventKey, func(args ...any) {
		key := args[0].(*Key)
		switch key.Name {
		case KeyUp, KeyDown, KeyLeft, KeyRight:
			p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+1, 2)
			p.Value = !p.Value
		}
	})
	return p
}
