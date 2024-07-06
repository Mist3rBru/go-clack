package prompts

import (
	"fmt"
	"io"
	"os"
	"strings"

	coreUtils "github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/prompts/symbols"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type NoteOptions struct {
	Title  string
	Output io.Writer
}

func Note(msg string, options NoteOptions) {
	if options.Output == nil {
		options.Output = os.Stdout
	}

	lineLength := coreUtils.StrLength(options.Title) + 7
	for _, line := range strings.Split(msg, "\n") {
		lineLength = max(coreUtils.StrLength(line)+4, lineLength)
	}

	header := noteHeader(options.Title, lineLength)
	body := noteBody(msg, lineLength)
	footer := noteFooter(lineLength)

	box := strings.Join([]string{
		picocolors.Gray(symbols.BAR),
		header,
		body,
		footer,
		"",
	}, "\r\n")

	options.Output.Write([]byte(box))
}

func noteHeader(title string, lineLength int) string {
	var left, header, top, right string

	if title == "" {
		left = symbols.CONNECT_LEFT
		top = strings.Repeat(symbols.BAR_H, lineLength)
		right = symbols.CORNER_TOP_RIGHT
		header = picocolors.Gray(fmt.Sprint(left, top, right))
	} else {
		left = picocolors.Green(symbols.STEP_SUBMIT)
		topLength := max(lineLength-coreUtils.StrLength(title)-2, 0)
		top = picocolors.Gray(strings.Repeat(symbols.BAR_H, topLength))
		right = picocolors.Gray(symbols.CORNER_TOP_RIGHT)
		header = fmt.Sprintf("%s %s %s%s", left, title, top, right)
	}

	return header
}

func noteBody(msg string, lineLength int) string {
	bar := picocolors.Gray(symbols.BAR)

	lines := strings.Split("\n"+msg+"\n", "\n")
	body := make([]string, len(lines))

	for i, line := range lines {
		whitespace := strings.Repeat(" ", max(lineLength-2-coreUtils.StrLength(line), 1))
		body[i] = fmt.Sprintf("%s  %s%s%s", bar, line, whitespace, bar)
	}

	return strings.Join(body, "\r\n")
}

func noteFooter(lineLength int) string {
	left := symbols.CONNECT_LEFT
	bottom := strings.Repeat(symbols.BAR_H, lineLength)
	right := symbols.CORNER_BOTTOM_RIGHT

	return picocolors.Gray(fmt.Sprint(left, bottom, right))
}
