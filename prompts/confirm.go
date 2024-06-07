package prompts

import (
	"strings"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/utils"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
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
			activeRadio := picocolors.Green(utils.S_RADIO_ACTIVE)
			inactiveRadio := picocolors.Dim(utils.S_RADIO_INACTIVE)
			slash := picocolors.Dim("/")

			var value, valueWithCursor string
			if p.Value {
				value = p.Active
				valueWithCursor = strings.Join([]string{activeRadio, p.Active, slash, inactiveRadio, picocolors.Dim(p.Inactive)}, " ")
			} else {
				value = p.Inactive
				valueWithCursor = strings.Join([]string{inactiveRadio, picocolors.Dim(p.Active), slash, activeRadio, p.Inactive}, " ")
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
