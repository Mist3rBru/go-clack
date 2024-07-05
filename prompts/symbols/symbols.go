package symbols

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
	STEP_ACTIVE Synbol = s("◆", "*")
	STEP_CANCEL Synbol = s("■", "x")
	STEP_ERROR  Synbol = s("▲", "x")
	STEP_SUBMIT Synbol = s("◇", "o")

	BAR_START Synbol = s("┌", "T")
	BAR       Synbol = s("│", "|")
	BAR_END   Synbol = s("└", "—")

	RADIO_ACTIVE      Synbol = s("●", ">")
	RADIO_INACTIVE    Synbol = s("○", " ")
	CHECKBOX_ACTIVE   Synbol = s("◻", "[•]")
	CHECKBOX_SELECTED Synbol = s("◼", "[+]")
	CHECKBOX_INACTIVE Synbol = s("◻", "[ ]")
	PASSWORD_MASK     Synbol = s("▪", "•")

	BAR_H               Synbol = s("─", "-")
	CORNER_TOP_RIGHT    Synbol = s("╮", "+")
	CONNECT_LEFT        Synbol = s("├", "+")
	CORNER_BOTTOM_RIGHT Synbol = s("╯", "+")

	INFO    Synbol = s("●", "•")
	SUCCESS Synbol = s("◆", "*")
	WARN    Synbol = s("▲", "!")
	ERROR   Synbol = s("■", "x")
)

func State(state core.State) string {
	switch state {
	case core.ErrorState:
		return picocolors.Yellow(STEP_ERROR)
	case core.CancelState:
		return picocolors.Red(STEP_CANCEL)
	case core.SubmitState:
		return picocolors.Green(STEP_SUBMIT)
	default:
		return picocolors.Cyan(STEP_ACTIVE)
	}
}
