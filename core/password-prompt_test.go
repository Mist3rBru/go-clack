package core_test

import (
	"os"
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newPasswordPrompt() *core.PasswordPrompt {
	return core.NewPasswordPrompt(core.PasswordPromptParams{
		Input:  os.Stdin,
		Output: os.Stdout,
	})
}

func TestChangePasswordValue(t *testing.T) {
	p := newPasswordPrompt()

	assert.Equal(t, "", p.Value)
	p.PressKey(&core.Key{Char: "a"})
	assert.Equal(t, "a", p.Value)
	p.PressKey(&core.Key{Char: "b"})
	assert.Equal(t, "ab", p.Value)
	p.PressKey(&core.Key{Name: "Backspace"})
	assert.Equal(t, "a", p.Value)
}

func TestChangPasswordMask(t *testing.T) {
	p := newPasswordPrompt()

	assert.Equal(t, " ", p.ValueWithCursor())
	p.PressKey(&core.Key{Char: "a"})
	assert.Equal(t, "* ", p.ValueWithCursor())
	p.PressKey(&core.Key{Char: "b"})
	assert.Equal(t, "** ", p.ValueWithCursor())
	p.PressKey(&core.Key{Name: "Left"})
	assert.Equal(t, "**", p.ValueWithCursor())
	p.PressKey(&core.Key{Name: "Backspace"})
	assert.Equal(t, "*", p.ValueWithCursor())
}
