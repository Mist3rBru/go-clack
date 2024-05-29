package utils

import (
	"fmt"
)

var (
	ESC = "\x1B"
	CSI = fmt.Sprintf("%s[", ESC)
)

func MoveCursor(row int, col int) string {
	ret := ""

	if col < 0 {
		ret += fmt.Sprintf("%s%dD", CSI, -col)
	} else if col > 0 {
		ret += fmt.Sprintf("%s%dC", CSI, col)
	}

	if row < 0 {
		ret += fmt.Sprintf("%s%dA", CSI, -row)
	} else if row > 0 {
		ret += fmt.Sprintf("%s%dB", CSI, row)
	}

	return ret
}

func MoveCursorUp(row int) string {
	return fmt.Sprintf("%s%dA", CSI, row)
}

func MoveCursorDown(row int) string {
	return fmt.Sprintf("%s%dB", CSI, row)
}

func HideCursor() string {
	return fmt.Sprintf("%s?251", CSI)
}

func ShowCursor() string {
	return fmt.Sprintf("%s?25h", CSI)
}
