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

type MultiSelectPathParams struct {
	Message      string
	InitialValue []string
	InitialPath  string
	Required     bool
	Validate     func(value []string) error
	OnlyShowDir  bool
	FileSystem   FileSystem
}

func MultiSelectPath(params MultiSelectPathParams) ([]string, error) {
	p := core.NewMultiSelectPathPrompt(core.MultiSelectPathPromptParams{
		InitialValue: params.InitialValue,
		InitialPath:  params.InitialPath,
		OnlyShowDir:  params.OnlyShowDir,
		FileSystem:   params.FileSystem,
		Required:     params.Required,
		Validate:     params.Validate,
		Render: func(p *core.MultiSelectPathPrompt) string {
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
					if option.IsSelected && i == p.CursorIndex {
						radio = picocolors.Green(symbols.CHECKBOX_SELECTED)
						label = option.Name
					} else if option.IsSelected {
						radio = picocolors.Green(symbols.CHECKBOX_SELECTED)
						label = picocolors.Dim(option.Name)
						dir = picocolors.Dim(dir)
					} else if i == p.CursorIndex {
						radio = picocolors.Green(symbols.CHECKBOX_ACTIVE)
						label = option.Name
					} else {
						radio = picocolors.Dim(symbols.CHECKBOX_INACTIVE)
						label = picocolors.Dim(option.Name)
						dir = picocolors.Dim(dir)
					}
					depth := strings.Repeat(" ", option.Depth)
					radioOptions[i] = fmt.Sprintf("%s%s %s %s", depth, radio, label, dir)
				}
				value = p.LimitLines(radioOptions, 3)
			}

			return theme.ApplyTheme(theme.ThemeParams[[]string]{
				Ctx:             p.Prompt,
				Message:         params.Message,
				Value:           strings.Join(p.Value, "\n"),
				ValueWithCursor: value,
			})
		},
	})
	test.MultiSelectPathTestingPrompt = p
	return p.Run()
}
