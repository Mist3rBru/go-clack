package theme

import (
	"strings"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/symbols"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
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
		picocolors.Gray(symbols.BAR),
		ctx.FormatLines(strings.Split(params.Message, "\n"), core.FormatLinesOptions{
			FirstLine: core.FormatLineOptions{
				Start: symbols.State(ctx.State),
			},
			NewLine: core.FormatLineOptions{
				Start: picocolors.Gray(symbols.BAR),
			},
		}),
	}, "\r\n")

	switch ctx.State {
	case core.ErrorState:
		value := ctx.FormatLines(strings.Split(params.ValueWithCursor, "\n"), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: picocolors.Yellow(symbols.BAR),
			},
		})
		if ctx.Error == "" {
			return strings.Join([]string{title, value}, "\r\n")
		}
		err := ctx.FormatLines(strings.Split(ctx.Error, "\n"), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: picocolors.Yellow(symbols.BAR),
				Style: picocolors.Yellow,
			},
			LastLine: core.FormatLineOptions{
				Start: picocolors.Yellow(symbols.BAR_END),
			},
		})
		return strings.Join([]string{title, value, err}, "\r\n")

	case core.CancelState:
		value := ctx.FormatLines(strings.Split(params.Value, "\n"), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: picocolors.Gray(symbols.BAR),
				Style: func(line string) string {
					return picocolors.Strikethrough(picocolors.Dim(line))
				},
			},
		})
		if params.Value == "" {
			return strings.Join([]string{title, value}, "\r\n")
		}
		end := picocolors.Gray(symbols.BAR)
		return strings.Join([]string{title, value, end}, "\r\n")

	case core.SubmitState:
		value := ctx.FormatLines(strings.Split(params.Value, "\n"), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: picocolors.Gray(symbols.BAR),
				Style: picocolors.Dim,
			},
		})
		return strings.Join([]string{title, value}, "\r\n")

	default:
		start := picocolors.Gray(symbols.BAR)
		title = ctx.FormatLines(strings.Split(params.Message, "\n"), core.FormatLinesOptions{
			FirstLine: core.FormatLineOptions{
				Start: symbols.State(ctx.State),
			},
			NewLine: core.FormatLineOptions{
				Start: picocolors.Cyan(symbols.BAR),
			},
		})

		var valueWithCursor string
		if params.Placeholder != "" && params.Value == "" {
			valueWithCursor = picocolors.Inverse(string(params.Placeholder[0])) + picocolors.Dim(params.Placeholder[1:])
		} else {
			valueWithCursor = params.ValueWithCursor
		}
		value := ctx.FormatLines(strings.Split(valueWithCursor, "\n"), core.FormatLinesOptions{
			Default: core.FormatLineOptions{
				Start: picocolors.Cyan(symbols.BAR),
			},
		})
		end := picocolors.Cyan(symbols.BAR_END)

		return strings.Join([]string{start, title, value, end}, "\r\n")
	}
}
