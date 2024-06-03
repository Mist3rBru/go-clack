package utils

import (
	"strings"

	"github.com/Mist3rBru/go-clack/core"
)

type ThemeValue interface {
	string | any | []any
}

type ThemeParams[TValue ThemeValue] struct {
	Ctx             core.Prompt[TValue]
	Message         string
	Value           string
	ValueWithCursor string
	Placeholder     string
}

func ApplyTheme[TValue ThemeValue](params ThemeParams[TValue]) string {
	ctx := params.Ctx

	title := strings.Join([]string{
		Color["gray"](S_BAR),
		ctx.FormatLines(strings.Split(params.Message, "\n"), core.FormatLinesOptions{
			FirstLine: core.FormatLineOptions{
				Start: SymbolState(ctx.State),
			},
			Default: core.FormatLineOptions{
				Start: Color["gray"](S_BAR),
			},
		}),
	}, "\r\n")

	switch ctx.State {
	case core.StateError:
		value := ctx.FormatLines(strings.Split(params.Value, "\n"), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: Color["yellow"](S_BAR),
			},
		})
		if ctx.Error == "" {
			return strings.Join([]string{title, value}, "\r\n")
		}
		err := ctx.FormatLines(strings.Split(ctx.Error, "\n"), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: Color["yellow"](S_BAR),
				Style: Color["yellow"],
			},
			LastLine: core.FormatLineOptions{
				Start: Color["yellow"](S_BAR_END),
			},
		})
		return strings.Join([]string{title, value, err}, "\r\n")

	case core.StateCancel:
		value := ctx.FormatLines(strings.Split(params.Value, "\n"), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: Color["gray"](S_BAR),
				Style: func(line string) string {
					return Color["strikethrough"](Color["dim"](line))
				},
			},
		})
		if params.Value == "" {
			return strings.Join([]string{title, value}, "\r\n")
		}
		end := Color["gray"](S_BAR)
		return strings.Join([]string{title, value, end}, "\r\n")

	case core.StateSubmit:
		value := ctx.FormatLines(strings.Split(params.Value, "\n"), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: Color["gray"](S_BAR),
				Style: Color["dim"],
			},
		})
		return strings.Join([]string{title, value}, "\r\n")

	default:
		start := Color["gray"](S_BAR)
		title = ctx.FormatLines(strings.Split(params.Message, "\n"), core.FormatLinesOptions{
			FirstLine: core.FormatLineOptions{
				Start: SymbolState(ctx.State),
			},
			Default: core.FormatLineOptions{
				Start: Color["cyan"](S_BAR),
			},
		})

		var valueWithCursor string
		if params.Placeholder != "" && params.Value == "" {
			valueWithCursor = Color["inverse"](string(params.Placeholder[0])) + Color["dim"](params.Placeholder[1:])
		} else if params.ValueWithCursor != "" {
			valueWithCursor = params.ValueWithCursor
		} else {
			valueWithCursor = params.Value
		}
		value := ctx.FormatLines(strings.Split(valueWithCursor, "\n"), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: Color["cyan"](S_BAR),
			},
		})
		end := Color["cyan"](S_BAR_END)

		return strings.Join([]string{start, title, value, end}, "\r\n")
	}
}
