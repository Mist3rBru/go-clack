package prompts

import (
	"strings"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/utils"
)

type ConfirmParams struct {
	Message      string
	InitialValue bool
	Active       string
	Inactive     string
}

func Confirm(params ConfirmParams) (bool, error) {
	p := core.NewConfirmPrompt(core.ConfirmPromptParams{
		InitialValue: params.InitialValue,
		Active:       params.Active,
		Inactive:     params.Inactive,
		Render: func(p *core.ConfirmPrompt) string {
			dim := utils.Color["dim"]
			activeRadio := utils.Color["green"](utils.S_RADIO_ACTIVE)
			inactiveRadio := dim(utils.S_RADIO_INACTIVE)
			slash := dim("/")

			var value, valueWithCursor string
			if p.Value {
				value = p.Active
				valueWithCursor = strings.Join([]string{activeRadio, p.Active, slash, inactiveRadio, dim(p.Inactive)}, " ")
			} else {
				value = p.Inactive
				valueWithCursor = strings.Join([]string{inactiveRadio, dim(p.Active), slash, activeRadio, p.Inactive}, " ")
			}

			return utils.ApplyTheme(utils.ThemeParams[bool]{
				Ctx:             p.Prompt,
				Message:         params.Message,
				Value:           value,
				ValueWithCursor: valueWithCursor,
			})
		},
	})
	test.ConfirmTestingPrompt = p
	return p.Run()
}
