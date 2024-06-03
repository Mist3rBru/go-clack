package prompts

import (
	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/utils"
)

type TextParams struct {
	Message      string
	Placeholder  string
	InitialValue string
	Validate     func(value string) error
}

func Text(params TextParams) (string, error) {
	p := core.NewTextPrompt(core.TextPromptParams{
		InitialValue: params.InitialValue,
		Placeholder:  params.Placeholder,
		Validate:     params.Validate,
		Render: func(p *core.TextPrompt) string {
			return utils.ApplyTheme(utils.ThemeParams[string]{
				Ctx:             p.Prompt,
				Message:         params.Message,
				Value:           p.Value,
				ValueWithCursor: p.ValueWithCursor(),
				Placeholder:     p.Placeholder,
			})
		},
	})
	test.TextTestingPrompt = p
	return p.Run()
}
