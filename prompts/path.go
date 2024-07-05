package prompts

import (
	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/theme"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type PathParams struct {
	Message      string
	InitialValue string
	OnlyShowDir  bool
	Required     bool
	Validate     func(value string) error
}

func Path(params PathParams) (string, error) {
	p := core.NewPathPrompt(core.PathPromptParams{
		InitialValue: params.InitialValue,
		OnlyShowDir:  params.OnlyShowDir,
		Required:     params.Required,
		Validate:     params.Validate,
		Render: func(p *core.PathPrompt) string {
			var hintOptions string
			if len(p.HintOptions) > 0 {
				hintOptions = "\n"
			}
			for i, hintOption := range p.HintOptions {
				if i != 0 {
					hintOptions += " "
				}
				if i == p.HintIndex {
					hintOptions += picocolors.Cyan(hintOption)
				} else {
					hintOptions += picocolors.Dim(hintOption)
				}
			}
			valueWithCursorAndOptions := p.ValueWithCursor() + hintOptions

			return theme.ApplyTheme(theme.ThemeParams[string]{
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
