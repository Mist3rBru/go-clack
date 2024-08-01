package prompts

import (
	"fmt"
	"strings"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/core/validator"
	"github.com/Mist3rBru/go-clack/prompts/symbols"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/theme"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type SelectOption[TValue comparable] struct {
	Label string
	Value TValue
	Hint  string
}

type SelectParams[TValue comparable] struct {
	Message      string
	InitialValue TValue
	Filter       bool
	Options      []*SelectOption[TValue]
}

func Select[TValue comparable](params SelectParams[TValue]) (TValue, error) {
	v := validator.NewValidator("Select")
	v.ValidateOptions(len(params.Options))

	var options []*core.SelectOption[TValue]
	for _, option := range params.Options {
		options = append(options, &core.SelectOption[TValue]{
			Label: option.Label,
			Value: option.Value,
		})
	}

	p := core.NewSelectPrompt(core.SelectPromptParams[TValue]{
		Filter:       params.Filter,
		InitialValue: params.InitialValue,
		Options:      options,
		Render: func(p *core.SelectPrompt[TValue]) string {
			message := params.Message
			var value string

			switch p.State {
			case core.SubmitState, core.CancelState:
				if p.CursorIndex >= 0 && p.CursorIndex < len(p.Options) {
					value = p.Options[p.CursorIndex].Label
				}
			default:
				radioOptions := make([]string, len(p.Options))
				for _, option := range params.Options {
					for i, _option := range p.Options {
						if option.Label != _option.Label {
							continue
						}

						var radio, label, hint string
						if i == p.CursorIndex {
							radio = picocolors.Green(symbols.RADIO_ACTIVE)
							label = option.Label
							if option.Hint != "" {
								hint = picocolors.Dim("(" + option.Hint + ")")
							}
						} else {
							radio = picocolors.Dim(symbols.RADIO_INACTIVE)
							label = picocolors.Dim(option.Label)
						}

						radioOptions[i] = strings.Join([]string{radio, label, hint}, " ")
						break
					}
				}

				if p.Filter {
					if p.Search == "" {
						message = fmt.Sprintf("%s\n> %s", message, picocolors.Inverse("T")+picocolors.Dim("ype to filter..."))
					} else {
						message = fmt.Sprintf("%s\n> %s", message, p.Search+picocolors.Inverse(" "))
					}

					value = p.LimitLines(radioOptions, 4)
				} else {
					value = p.LimitLines(radioOptions, 3)
				}
			}

			return theme.ApplyTheme(theme.ThemeParams[TValue]{
				Ctx:             p.Prompt,
				Message:         message,
				Value:           value,
				ValueWithCursor: value,
			})
		},
	})
	test.SelectTestingPrompt = p
	return p.Run()
}
