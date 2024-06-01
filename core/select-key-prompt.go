package core

import (
	"os"
)

type SelectKeyOption struct {
	Label string
	Value any
	Key   string
}

type SelectKeyPrompt struct {
	Prompt
	Options []SelectKeyOption
}

type SelectKeyPromptParams struct {
	Input   *os.File
	Output  *os.File
	Options []SelectKeyOption
	Render  func(p *SelectKeyPrompt) string
}

func NewSelectKeyPrompt(params SelectKeyPromptParams) *SelectKeyPrompt {
	var p *SelectKeyPrompt
	p = &SelectKeyPrompt{
		Prompt: *NewPrompt(PromptParams{
			Input:  params.Input,
			Output: params.Output,
			Track:  false,
			Render: func(_p *Prompt) string {
				return params.Render(p)
			},
		}),
		Options: params.Options,
	}
	p.On("key", func(args ...any) {
		key := args[0].(*Key)
		for _, option := range p.Options {
			if key.Name == option.Key {
				p.State = "submit"
				p.Value = option.Value
				p.Emit("submit", p.Value)
				return
			}
		}
	})
	return p
}
