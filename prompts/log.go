package prompts

import (
	"fmt"
	"os"

	"github.com/Mist3rBru/go-clack/prompts/utils"
)

func write(msg string) {
	os.Stdout.WriteString(msg)
}

func Intro(msg string) {
	write(fmt.Sprintf("%s %s\n", utils.Color["gray"](utils.S_BAR_START), msg))
}

func Cancel(msg string) {
	write(fmt.Sprintf("%s %s\n\n", utils.Color["gray"](utils.S_BAR_END), utils.Color["red"](msg)))
}

func Outro(msg string) {
	write(fmt.Sprintf("%s\n%s %s\n\n", utils.Color["gray"](utils.S_BAR), utils.Color["gray"](utils.S_BAR_END), msg))
}

func Info(msg string) {
	write(fmt.Sprintf("%s %s\n", utils.Color["blue"](utils.S_INFO), msg))
}

func Success(msg string) {
	write(fmt.Sprintf("%s %s\n", utils.Color["green"](utils.S_SUCCESS), msg))
}

func Step(msg string) {
	write(fmt.Sprintf("%s %s\n", utils.Color["green"](utils.S_STEP_SUBMIT), msg))
}

func Warn(msg string) {
	write(fmt.Sprintf("%s %s\n", utils.Color["yellow"](utils.S_WARN), msg))
}

func Error(msg string) {
	write(fmt.Sprintf("%s %s\n", utils.Color["red"](utils.S_ERROR), msg))
}
