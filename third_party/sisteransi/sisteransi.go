// Forked from https://github.com/terkelg/sisteransi/blob/master/src/index.js
package sisteransi

import (
	"fmt"
)

var (
	ESC = "\x1B"
	CSI = fmt.Sprintf("%s[", ESC)
)

func MoveCursorUp(row int) string {
	return fmt.Sprintf("%s%dA", CSI, row)
}

func MoveCursorDown(row int) string {
	return fmt.Sprintf("%s%dB", CSI, row)
}

func MoveCursor(row int, col int) string {
	ret := ""

	if col < 0 {
		ret += fmt.Sprintf("%s%dD", CSI, -col)
	} else if col > 0 {
		ret += fmt.Sprintf("%s%dC", CSI, col)
	}

	if row < 0 {
		ret += MoveCursorUp(-row)
	} else if row > 0 {
		ret += MoveCursorDown(row)
	}

	return ret
}

func HideCursor() string {
	return fmt.Sprintf("%s?25l", CSI)
}

func ShowCursor() string {
	return fmt.Sprintf("%s?25h", CSI)
}

func Restore() string {
	return fmt.Sprintf("%s8", CSI)
}

func EraseCurrentLine() string {
	return "\x1b[K"
}

func EraseDown() string {
	return "\x1b[J"
}
