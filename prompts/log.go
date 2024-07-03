package prompts

import (
	"fmt"
	"os"
	"strings"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts/utils"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type MessageLineOptions = core.FormatLineOptions
type MessageOptions = core.FormatLinesOptions

func Message(msg string, options MessageOptions) {
	p := &core.Prompt[string]{}
	formattedMsg := p.FormatLines(strings.Split(msg, "\n"), options)
	os.Stdout.WriteString(fmt.Sprintf("%s\r\n%s\r\n", picocolors.Gray(utils.S_BAR), formattedMsg))
}

func styleMsg(msg string, style func(msg string) string) string {
	parts := strings.Split(msg, "\n")
	styledParts := make([]string, len(parts))
	for i, part := range parts {
		styledParts[i] = style(part)
	}
	return strings.Join(styledParts, "\n")
}

func Intro(msg string) {
	p := &core.Prompt[string]{}
	formattedMsg := p.FormatLines(strings.Split(msg, "\n"), MessageOptions{
		FirstLine: MessageLineOptions{
			Start: picocolors.Gray(utils.S_BAR_START),
		},
		NewLine: MessageLineOptions{
			Start: picocolors.Gray(utils.S_BAR),
		},
	})
	os.Stdout.WriteString(fmt.Sprintf("\r\n%s\r\n%s\r\n", formattedMsg, picocolors.Gray(utils.S_BAR)))
}

func Cancel(msg string) {
	Message(styleMsg(msg, picocolors.Red), MessageOptions{
		Default: MessageLineOptions{
			Start: picocolors.Gray(utils.S_BAR),
		},
		LastLine: MessageLineOptions{
			Start: picocolors.Gray(utils.S_BAR_END),
		},
	})
}

func Outro(msg string) {
	Message("\n"+msg, MessageOptions{
		Default: MessageLineOptions{
			Start: picocolors.Gray(utils.S_BAR),
		},
		LastLine: MessageLineOptions{
			Start: picocolors.Gray(utils.S_BAR_END),
		},
	})
}

func Info(msg string) {
	Message(msg, MessageOptions{
		FirstLine: MessageLineOptions{
			Start: picocolors.Blue(utils.S_INFO),
		},
		NewLine: MessageLineOptions{
			Start: picocolors.Gray(utils.S_BAR),
		},
	})
}

func Success(msg string) {
	Message(msg, MessageOptions{
		FirstLine: MessageLineOptions{
			Start: picocolors.Green(utils.S_SUCCESS),
		},
		NewLine: MessageLineOptions{
			Start: picocolors.Gray(utils.S_BAR),
		},
	})
}

func Step(msg string) {
	Message(msg, MessageOptions{
		FirstLine: MessageLineOptions{
			Start: picocolors.Green(utils.S_STEP_SUBMIT),
		},
		NewLine: MessageLineOptions{
			Start: picocolors.Gray(utils.S_BAR),
		},
	})
}

func Warn(msg string) {
	Message(msg, MessageOptions{
		FirstLine: MessageLineOptions{
			Start: picocolors.Yellow(utils.S_WARN),
		},
		NewLine: MessageLineOptions{
			Start: picocolors.Gray(utils.S_BAR),
		},
	})
}

func Error(msg string) {
	Message(msg, MessageOptions{
		FirstLine: MessageLineOptions{
			Start: picocolors.Red(utils.S_ERROR),
		},
		NewLine: MessageLineOptions{
			Start: picocolors.Gray(utils.S_BAR),
		},
	})
}
