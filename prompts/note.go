package prompts

import (
	"fmt"
	"io"
	"os"
	"strings"

	coreUtils "github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/prompts/utils"
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
		options.Title = fmt.Sprintf("%s %s ", utils.Color["green"](utils.S_STEP_SUBMIT), options.Title)
	}
	titleLength := coreUtils.StrLength(options.Title)
	largestLineLength := titleLength

	lines := strings.Split("\n"+msg+"\n", "\n")
	for _, line := range lines {
		largestLineLength = max(coreUtils.StrLength(line), largestLineLength)
	}

	gray := utils.Color["gray"]
	bar := gray(utils.S_BAR)
	boxTop := gray(strings.Repeat(utils.S_BAR_H, max(largestLineLength+4-titleLength, 1)))
	boxTopTight := gray(utils.S_CORNER_TOP_RIGHT)
	boxHeader := strings.Join([]string{bar, options.Title, boxTop, boxTopTight}, "")

	boxLines := make([]string, len(lines))
	for i, line := range lines {
		whitespace := strings.Repeat(" ", largestLineLength+2-coreUtils.StrLength(line))
		boxLines[i] = fmt.Sprintf("%s  %s%s%s", bar, line, whitespace, bar)
	}
	boxBody := strings.Join(boxLines, "\r\n")

	boxBottom := strings.Repeat(gray(utils.S_BAR_H), largestLineLength+4)
	boxBottomRight := gray(utils.S_CORNER_BOTTOM_RIGHT)
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
