package prompts

import (
	"strings"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/core/validator"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/utils"
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
		InitialValue: params.InitialValue,
		Options:      options,
		Render: func(p *core.SelectPrompt[TValue]) string {
			var value string
			switch p.State {
			case core.SubmitState, core.CancelState:
			default:
				radioOptions := make([]string, len(params.Options))
				for i, option := range params.Options {
					var radio, label, hint string
					if i == p.CursorIndex {
						radio = picocolors.Green(utils.S_RADIO_ACTIVE)
						label = option.Label
						if option.Hint != "" {
							hint = picocolors.Dim("(" + option.Hint + ")")
						}
					} else {
						radio = picocolors.Dim(utils.S_RADIO_INACTIVE)
						label = picocolors.Dim(option.Label)
					}
					radioOptions[i] = strings.Join([]string{radio, label, hint}, " ")
				}
				value = p.LimitLines(radioOptions, 3)
			}

			return utils.ApplyTheme(utils.ThemeParams[TValue]{
				Ctx:             p.Prompt,
				Message:         params.Message,
				Value:           params.Options[p.CursorIndex].Label,
				ValueWithCursor: value,
			})
		},
	})
	test.SelectTestingPrompt = p
	return p.Run()
}
