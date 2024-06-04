package prompts

import (
	"strings"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/utils"
)

type SelectOption[TValue comparable] struct {
	Label string
	Value TValue
	Hint  string
}

type SelectParams[TValue comparable] struct {
	Message      string
	InitialValue TValue
	Options      []SelectOption[TValue]
	Validate     func(value string) error
}

func Select[TValue comparable](params SelectParams[TValue]) (TValue, error) {
	var options []core.SelectOption[TValue]
	for _, option := range params.Options {
		options = append(options, core.SelectOption[TValue]{
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
				var lines []string
				for i, option := range params.Options {
					if i == p.CursorIndex {
						lines = append(lines, opt[TValue](option, "active"))
					} else {
						lines = append(lines, opt[TValue](option, "inactive"))
					}
				}
				value = p.LimitLines(lines, 3)
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

func opt[TValue comparable](option SelectOption[TValue], state string) string {
	if state == "active" {
		radio := utils.Color["green"](utils.S_RADIO_ACTIVE)
		if option.Hint == "" {
			return strings.Join([]string{radio, option.Label}, " ")
		}
		hint := utils.Color["dim"]("(" + option.Hint + ")")
		return strings.Join([]string{radio, option.Label, hint}, " ")
	}

	radio := utils.Color["dim"](utils.S_RADIO_INACTIVE)
	label := utils.Color["dim"](option.Label)
	return strings.Join([]string{radio, label}, " ")
}
