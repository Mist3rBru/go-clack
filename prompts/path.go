package prompts

import (
	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/utils"
)

type PathParams struct {
	Message      string
	InitialValue string
	OnlyShowDir  bool
	Validate     func(value string) error
}

func Path(params PathParams) (string, error) {
	p := core.NewPathPrompt(core.PathPromptParams{
		InitialValue: params.InitialValue,
		OnlyShowDir:  params.OnlyShowDir,
		Validate:     params.Validate,
		Render: func(p *core.PathPrompt) string {
			var hintOptions string
			if len(p.HintOptions) > 0 {
				hintOptions = "\n"
			}
			for i, hintOption := range p.HintOptions {
				if i == p.HintIndex {
					hintOptions += " " + utils.Color["cyan"](hintOption)
				} else {
					hintOptions += " " + utils.Color["dim"](hintOption)
				}
			}
			valueWithCursorAndOptions := p.ValueWithCursor() + hintOptions

			return utils.ApplyTheme(utils.ThemeParams[string]{
				Ctx:             p.Prompt,
				Message:         params.Message,
				Value:           p.Value,
				ValueWithCursor: valueWithCursorAndOptions,
				Placeholder:     p.Placeholder,
			})
		},
	})
	test.PathTestingPrompt = p
	return p.Run()
}