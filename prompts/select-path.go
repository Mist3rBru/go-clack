package prompts

import (
	"fmt"
	"strings"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/symbols"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/theme"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type FileSystem = core.FileSystem

type SelectPathParams struct {
	Message      string
	InitialValue string
	OnlyShowDir  bool
	Filter       bool
	FileSystem   FileSystem
}

func SelectPath(params SelectPathParams) (string, error) {
	p := core.NewSelectPathPrompt(core.SelectPathPromptParams{
		InitialValue: params.InitialValue,
		OnlyShowDir:  params.OnlyShowDir,
		Filter:       params.Filter,
		FileSystem:   params.FileSystem,
		Render: func(p *core.SelectPathPrompt) string {
			message := params.Message
			var value string

			switch p.State {
			case core.SubmitState, core.CancelState:
			default:
				options := p.Options()
				radioOptions := make([]string, len(options))
				for i, option := range options {
					var radio, label, dir string
					if option.IsDir {
						if len(option.Children) == 0 {
							dir = ">"
						} else {
							dir = "v"
						}
					}
					if option.IsEqual(p.CurrentOption) {
						radio = picocolors.Green(symbols.RADIO_ACTIVE)
						label = option.Name
					} else {
						radio = picocolors.Dim(symbols.RADIO_INACTIVE)
						label = picocolors.Dim(option.Name)
						dir = picocolors.Dim(dir)
					}
					depth := strings.Repeat(" ", option.Depth)
					radioOptions[i] = fmt.Sprintf("%s%s %s %s", depth, radio, label, dir)
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

			return theme.ApplyTheme(theme.ThemeParams[string]{
				Ctx:             p.Prompt,
				Message:         message,
				Value:           p.Value,
				ValueWithCursor: value,
			})
		},
	})
	test.SelectPathTestingPrompt = p
	return p.Run()
}
