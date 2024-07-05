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
	if options.Title != "" {
		options.Title = fmt.Sprintf("%s %s ", picocolors.Green(symbols.STEP_SUBMIT), options.Title)
	}
	titleLength := coreUtils.StrLength(options.Title)
	largestLineLength := titleLength

	lines := strings.Split("\n"+msg+"\n", "\n")
	for _, line := range lines {
		largestLineLength = max(coreUtils.StrLength(line), largestLineLength)
	}

	bar := picocolors.Gray(symbols.BAR)
	boxTop := picocolors.Gray(strings.Repeat(symbols.BAR_H, max(largestLineLength+4-titleLength, 1)))
	boxTopTight := picocolors.Gray(symbols.CORNER_TOP_RIGHT)
	boxHeader := strings.Join([]string{bar, options.Title, boxTop, boxTopTight}, "")

	boxLines := make([]string, len(lines))
	for i, line := range lines {
		whitespace := strings.Repeat(" ", largestLineLength+2-coreUtils.StrLength(line))
		boxLines[i] = fmt.Sprintf("%s  %s%s%s", bar, line, whitespace, bar)
	}
	boxBody := strings.Join(boxLines, "\r\n")

	boxBottom := strings.Repeat(picocolors.Gray(symbols.BAR_H), largestLineLength+4)
	boxBottomRight := picocolors.Gray(symbols.CORNER_BOTTOM_RIGHT)
	boxFooter := strings.Join([]string{bar, boxBottom, boxBottomRight}, "")

	box := strings.Join([]string{
		bar,
		boxHeader,
		boxBody,
		boxFooter,
		"",
	}, "\r\n")

	options.Output.Write([]byte(box))
}
