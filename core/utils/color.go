// Imported from https://github.com/alexeyraspopov/picocolors/blob/main/picocolors.js
package utils

import (
	"os"
	"strings"

	"golang.org/x/term"
)

func isColorSupported() bool {
	env := os.Environ()
	argv := os.Args

	noColor := false
	for _, arg := range argv {
		if arg == "--no-color" {
			noColor = true
			break
		}
	}

	if noColor {
		return false
	}

	forceColor := false
	for _, arg := range argv {
		if arg == "--color" {
			forceColor = true
			break
		}
	}

	if forceColor {
		return true
	}

	for _, v := range env {
		if v == "FORCE_COLOR" {
			return true
		}
	}

	if os.Getenv("TERM") != "dumb" && term.IsTerminal(int(os.Stdout.Fd())) {
		return true
	}

	for _, v := range env {
		if v == "CI" {
			return true
		}
	}

	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	return false
}

func formatter(open, close, replace string) func(string) string {
	return func(input string) string {
		index := strings.Index(input, close)
		if index != -1 {
			return open + replaceClose(input, close, replace, index) + close
		}
		return open + input + close
	}
}

func replaceClose(input, close, replace string, index int) string {
	result := ""
	cursor := 0

	for index != -1 {
		result += input[cursor:index] + replace
		cursor = index + len(close)
		index = strings.Index(input[cursor:], close)
		if index != -1 {
			index += cursor
		}
	}

	return result + input[cursor:]
}

func CreateColors() map[string]func(input string) string {
	init := func(open, close, replace string) func(input string) string {
		if isColorSupported() {
			return formatter(open, close, replace)
		}
		return func(input string) string { return input }
	}

	colors := map[string]func(string) string{
		"reset":         init("\x1b[0m", "\x1b[0m", "\x1b[0m"),
		"bold":          init("\x1b[1m", "\x1b[22m", "\x1b[22m\x1b[1m"),
		"dim":           init("\x1b[2m", "\x1b[22m", "\x1b[22m\x1b[2m"),
		"italic":        init("\x1b[3m", "\x1b[23m", "\x1b[3m"),
		"underline":     init("\x1b[4m", "\x1b[24m", "\x1b[4m"),
		"inverse":       init("\x1b[7m", "\x1b[27m", "\x1b[7m"),
		"hidden":        init("\x1b[8m", "\x1b[28m", "\x1b[8m"),
		"strikethrough": init("\x1b[9m", "\x1b[29m", "\x1b[9m"),
		"black":         init("\x1b[30m", "\x1b[39m", "\x1b[30m"),
		"red":           init("\x1b[31m", "\x1b[39m", "\x1b[31m"),
		"green":         init("\x1b[32m", "\x1b[39m", "\x1b[32m"),
		"yellow":        init("\x1b[33m", "\x1b[39m", "\x1b[33m"),
		"blue":          init("\x1b[34m", "\x1b[39m", "\x1b[34m"),
		"magenta":       init("\x1b[35m", "\x1b[39m", "\x1b[35m"),
		"cyan":          init("\x1b[36m", "\x1b[39m", "\x1b[36m"),
		"white":         init("\x1b[37m", "\x1b[39m", "\x1b[37m"),
		"gray":          init("\x1b[90m", "\x1b[39m", "\x1b[90m"),
		"bgBlack":       init("\x1b[40m", "\x1b[49m", "\x1b[40m"),
		"bgRed":         init("\x1b[41m", "\x1b[49m", "\x1b[41m"),
		"bgGreen":       init("\x1b[42m", "\x1b[49m", "\x1b[42m"),
		"bgYellow":      init("\x1b[43m", "\x1b[49m", "\x1b[43m"),
		"bgBlue":        init("\x1b[44m", "\x1b[49m", "\x1b[44m"),
		"bgMagenta":     init("\x1b[45m", "\x1b[49m", "\x1b[45m"),
		"bgCyan":        init("\x1b[46m", "\x1b[49m", "\x1b[46m"),
		"bgWhite":       init("\x1b[47m", "\x1b[49m", "\x1b[47m"),
	}

	return colors
}
