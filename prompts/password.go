package prompts

import (
	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/utils"
)

type PasswordParams struct {
	Message      string
	InitialValue string
	Validate     func(value string) error
}

func Password(params PasswordParams) (string, error) {
	p := core.NewPasswordPrompt(core.PasswordPromptParams{
		InitialValue: params.InitialValue,
		Validate:     params.Validate,
		Render: func(p *core.PasswordPrompt) string {
			return utils.ApplyTheme(utils.ThemeParams[string]{
				Ctx:             p.Prompt,
				Message:         params.Message,
				Value:           p.ValueWithMask(),
				ValueWithCursor: p.ValueWithMaskAndCursor(),
			})
		},
	})
	test.PasswordTestingPrompt = p
	return p.Run()
}
