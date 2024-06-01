package core_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newPasswordPrompt() *core.PasswordPrompt {
	return core.NewPasswordPrompt(core.PasswordPromptParams{
		Input:  os.Stdin,
		Output: os.Stdout,
		Value:  "",
	})
}

func TestPasswordPromptInitialValue(t *testing.T) {
	p := core.NewPasswordPrompt(core.PasswordPromptParams{
		Value: "foo",
	})

	assert.Equal(t, "foo", p.Value)
	assert.Equal(t, 3, p.CursorIndex)
}

func TestChangePasswordValue(t *testing.T) {
	p := newPasswordPrompt()

	assert.Equal(t, "", p.Value)
	p.PressKey(&core.Key{Char: "a"})
	assert.Equal(t, "a", p.Value)
	p.PressKey(&core.Key{Char: "b"})
	assert.Equal(t, "ab", p.Value)
	p.PressKey(&core.Key{Name: core.KeyBackspace})
	assert.Equal(t, "a", p.Value)
}

func TestChangePasswordMask(t *testing.T) {
	p := newPasswordPrompt()

	assert.Equal(t, " ", p.ValueWithCursor())
	p.PressKey(&core.Key{Char: "a"})
	assert.Equal(t, "* ", p.ValueWithCursor())
	p.PressKey(&core.Key{Char: "b"})
	assert.Equal(t, "** ", p.ValueWithCursor())
	p.PressKey(&core.Key{Name: core.KeyLeft})
	assert.Equal(t, "**", p.ValueWithCursor())
	p.PressKey(&core.Key{Name: core.KeyBackspace})
	assert.Equal(t, "*", p.ValueWithCursor())
}

func TestValidatePassword(t *testing.T) {
	p := core.NewPasswordPrompt(core.PasswordPromptParams{
		Value: "123",
		Validate: func(value string) error {
			return fmt.Errorf("invalid password: %s", value)
		},
	})

	p.PressKey(&core.Key{Name: core.KeyEnter})
	assert.Equal(t, "error", p.State)
	assert.Equal(t, "invalid password: 123", p.Error)
}
