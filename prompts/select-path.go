package prompts

import (
	"fmt"
	"strings"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/utils"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type FileSystem = core.FileSystem

type SelectPathParams struct {
	Message      string
	InitialValue string
	OnlyShowDir  bool
	FileSystem   FileSystem
}

func SelectPath(params SelectPathParams) (string, error) {
	p := core.NewSelectPathPrompt(core.SelectPathPromptParams{
		InitialValue: params.InitialValue,
		OnlyShowDir:  params.OnlyShowDir,
		FileSystem:   params.FileSystem,
		Render: func(p *core.SelectPathPrompt) string {
			var value string
			switch p.State {
			case core.SubmitState, core.CancelState:
			default:
				options := p.Options()
				radioOptions := make([]string, len(options))
				for i, option := range options {
					var radio, label, dir string
					if option.Children != nil {
						if len(option.Children) == 0 {
							dir = ">"
						} else {
							dir = "v"
						}
					}
					if i == p.CursorIndex {
						radio = picocolors.Green(utils.S_RADIO_ACTIVE)
						label = option.Name
					} else {
						radio = picocolors.Dim(utils.S_RADIO_INACTIVE)
						label = picocolors.Dim(option.Name)
						dir = picocolors.Dim(dir)
					}
					depth := strings.Repeat(" ", option.Depth)
					radioOptions[i] = fmt.Sprintf("%s%s %s %s", depth, radio, label, dir)
				}
				value = p.LimitLines(radioOptions, 3)
			}

			return utils.ApplyTheme(utils.ThemeParams[string]{
				Ctx:             p.Prompt,
				Message:         params.Message,
				Value:           p.Value,
				ValueWithCursor: value,
			})
		},
	})
	test.SelectPathTestingPrompt = p
	return p.Run()
}
