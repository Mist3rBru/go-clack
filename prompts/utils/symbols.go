package utils

import (
	"github.com/Mist3rBru/go-clack/core"
	isunicodesupported "github.com/Mist3rBru/go-clack/third_party/is-unicode-supported"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

func s(c, fallback string) string {
	if isunicodesupported.IsUnicodeSupported() {
		return c
	}
	return fallback
}

type Synbol = string

var (
	S_STEP_ACTIVE Synbol = s("◆", "*")
	S_STEP_CANCEL Synbol = s("■", "x")
	S_STEP_ERROR  Synbol = s("▲", "x")
	S_STEP_SUBMIT Synbol = s("◇", "o")

	S_BAR_START Synbol = s("┌", "T")
	S_BAR       Synbol = s("│", "|")
	S_BAR_END   Synbol = s("└", "—")

	S_RADIO_ACTIVE      Synbol = s("●", ">")
	S_RADIO_INACTIVE    Synbol = s("○", " ")
	S_CHECKBOX_ACTIVE   Synbol = s("◻", "[•]")
	S_CHECKBOX_SELECTED Synbol = s("◼", "[+]")
	S_CHECKBOX_INACTIVE Synbol = s("◻", "[ ]")
	S_PASSWORD_MASK     Synbol = s("▪", "•")

	S_BAR_H               Synbol = s("─", "-")
	S_CORNER_TOP_RIGHT    Synbol = s("╮", "+")
	S_CONNECT_LEFT        Synbol = s("├", "+")
	S_CORNER_BOTTOM_RIGHT Synbol = s("╯", "+")

	S_INFO    Synbol = s("●", "•")
	S_SUCCESS Synbol = s("◆", "*")
	S_WARN    Synbol = s("▲", "!")
	S_ERROR   Synbol = s("■", "x")
)

func SymbolState(state core.State) string {
	switch state {
	case "error":
		return picocolors.Yellow(S_STEP_ERROR)
	case "cancel":
		return picocolors.Red(S_STEP_CANCEL)
	case "submit":
		return picocolors.Green(S_STEP_SUBMIT)
	default:
		return picocolors.Cyan(S_STEP_ACTIVE)
	}
}
