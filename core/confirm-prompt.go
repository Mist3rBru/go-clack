package core

import (
	"os"

	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/core/validator"
)

type ConfirmPrompt struct {
	Prompt[bool]
	Active   string
	Inactive string
}

type ConfirmPromptParams struct {
	Input        *os.File
	Output       *os.File
	Active       string
	Inactive     string
	InitialValue bool
	Render       func(p *ConfirmPrompt) string
}

func NewConfirmPrompt(params ConfirmPromptParams) *ConfirmPrompt {
	v := validator.NewValidator("ConfirmPrompt")
	v.ValidateRender(params.Render)

	if params.Active == "" {
		params.Active = "yes"
	}
	if params.Inactive == "" {
		params.Inactive = "no"
	}
	var p ConfirmPrompt
	p = ConfirmPrompt{
		Prompt: *NewPrompt(PromptParams[bool]{
			Input:        params.Input,
			Output:       params.Output,
			InitialValue: params.InitialValue,
			Render:       WrapRender[bool](&p, params.Render),
		}),
		Active:   params.Active,
		Inactive: params.Inactive,
	}
	p.On(KeyEvent, func(args ...any) {
		key := args[0].(*Key)
		switch key.Name {
		case UpKey, DownKey, LeftKey, RightKey:
			p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+1, 2)
			p.Value = !p.Value
		}
	})
	return &p
}
