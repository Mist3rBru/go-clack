package core

import (
	"os"
)

type SelectKeyOption[TValue any] struct {
	Label string
	Value TValue
	Key   string
}

type SelectKeyPrompt[TValue any] struct {
	Prompt[TValue]
	Options []SelectKeyOption[TValue]
}

type SelectKeyPromptParams[TValue any] struct {
	Input   *os.File
	Output  *os.File
	Options []SelectKeyOption[TValue]
	Render  func(p *SelectKeyPrompt[TValue]) string
}

func NewSelectKeyPrompt[TValue any](params SelectKeyPromptParams[TValue]) *SelectKeyPrompt[TValue] {
	var p *SelectKeyPrompt[TValue]
	p = &SelectKeyPrompt[TValue]{
		Prompt: *NewPrompt(PromptParams[TValue]{
			Input:  params.Input,
			Output: params.Output,
			Render: func(_p *Prompt[TValue]) string {
				return params.Render(p)
			},
		}),
		Options: params.Options,
	}
	p.On(PromptEventKey, func(args ...any) {
		key := args[0].(*Key)
		for _, option := range p.Options {
			if key.Name == KeyName(option.Key) {
				p.State = "submit"
				p.Value = option.Value
				p.Emit(PromptEventSubmit, p.Value)
				return
			}
		}
	})
	return p
}
