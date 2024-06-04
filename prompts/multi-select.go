package prompts

import (
	"strings"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/utils"
)

type MultiSelectOption[TValue comparable] struct {
	Label      string
	Value      TValue
	IsSelected bool
	Hint       string
}

type MultiSelectParams[TValue comparable] struct {
	Message      string
	Options      []MultiSelectOption[TValue]
	InitialValue []TValue
	Validate     func(value []TValue) error
}

func MultiSelect[TValue comparable](params MultiSelectParams[TValue]) ([]TValue, error) {
	var options []*core.MultiSelectOption[TValue]
	for _, option := range params.Options {
		options = append(options, &core.MultiSelectOption[TValue]{
			Label:      option.Label,
			Value:      option.Value,
			IsSelected: option.IsSelected,
		})
	}

	p := core.NewMultiSelectPrompt(core.MultiSelectPromptParams[TValue]{
		InitialValue: params.InitialValue,
		Options:      options,
		Validate:     params.Validate,
		Render: func(p *core.MultiSelectPrompt[TValue]) string {
			var value string

			switch p.State {
			case core.SubmitState, core.CancelState:
				for _, option := range p.Options {
					if option.IsSelected {
						if value == "" {
							value = option.Label
						} else {
							value += ", " + option.Label
						}
					}
				}

			default:
				var radioOptions []string
				for i, option := range p.Options {
					var radio, label, hint string
					if option.IsSelected && i == p.CursorIndex {
						radio = utils.Color["green"](utils.S_CHECKBOX_SELECTED)
						label = option.Label
						if params.Options[i].Hint != "" {
							hint = utils.Color["dim"]("(" + params.Options[i].Hint + ")")
						}
					} else if i == p.CursorIndex {
						radio = utils.Color["green"](utils.S_CHECKBOX_ACTIVE)
						label = option.Label
						if params.Options[i].Hint != "" {
							hint = utils.Color["dim"]("(" + params.Options[i].Hint + ")")
						}
					} else if option.IsSelected {
						radio = utils.Color["green"](utils.S_CHECKBOX_SELECTED)
						label = utils.Color["dim"](option.Label)
					} else {
						radio = utils.Color["dim"](utils.S_CHECKBOX_INACTIVE)
						label = utils.Color["dim"](option.Label)
					}
					radioOptions = append(radioOptions, strings.Join([]string{radio, label, hint}, " "))
				}
				value = p.LimitLines(radioOptions, 3)
			}

			return utils.ApplyTheme(utils.ThemeParams[[]TValue]{
				Ctx:             p.Prompt,
				Message:         params.Message,
				Value:           value,
				ValueWithCursor: value,
			})
		},
	})
	test.MultiSelectTestingPrompt = p
	return p.Run()
}
