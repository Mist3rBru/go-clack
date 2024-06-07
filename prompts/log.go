package prompts

import (
	"fmt"
	"os"

	"github.com/Mist3rBru/go-clack/prompts/utils"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

func write(msg string) {
	os.Stdout.WriteString(msg)
}

func Intro(msg string) {
	write(fmt.Sprintf("%s %s\n", picocolors.Gray(utils.S_BAR_START), msg))
}

func Cancel(msg string) {
	write(fmt.Sprintf("%s %s\n\n", picocolors.Gray(utils.S_BAR_END), picocolors.Red(msg)))
}

func Outro(msg string) {
	write(fmt.Sprintf("%s\n%s %s\n\n", picocolors.Gray(utils.S_BAR), picocolors.Gray(utils.S_BAR_END), msg))
}

func Info(msg string) {
	write(fmt.Sprintf("%s %s\n", picocolors.Blue(utils.S_INFO), msg))
}

func Success(msg string) {
	write(fmt.Sprintf("%s %s\n", picocolors.Green(utils.S_SUCCESS), msg))
}

func Step(msg string) {
	write(fmt.Sprintf("%s %s\n", picocolors.Green(utils.S_STEP_SUBMIT), msg))
}

func Warn(msg string) {
	write(fmt.Sprintf("%s %s\n", picocolors.Yellow(utils.S_WARN), msg))
}

func Error(msg string) {
	write(fmt.Sprintf("%s %s\n", picocolors.Red(utils.S_ERROR), msg))
}
