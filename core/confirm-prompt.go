package core

import (
	"os"

	"github.com/Mist3rBru/go-clack/core/utils"
)

type ConfirmPrompt struct {
	Prompt
	Active   string
	Inactive string
	Value    bool
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
		Prompt: *NewPrompt(PromptParams{
			Input:  params.Input,
			Output: params.Output,
			Value:  params.Value,
			Track:  false,
			Render: func(_p *Prompt) string {
				return params.Render(p)
			},
		}),
		Active:   params.Active,
		Inactive: params.Inactive,
		Value:    params.Value,
	}
	p.On("key", func(args ...any) {
		key := args[0].(*Key)
		switch key.Name {
		case "Up", "Left":
			p.CursorIndex = utils.MinMaxIndex(p.CursorIndex-1, 1)
			p.Value = !p.Value
		case "Down", "Right":
			p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+1, 1)
			p.Value = !p.Value
		}
	})
	return p
}

func (p *ConfirmPrompt) Run() (bool, error) {
	_, err := p.Prompt.Run()
	if err != nil {
		return false, err
	}
	return p.Value, nil
}
