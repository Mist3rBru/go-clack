package core_test

import (
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newConfirmPrompt() *core.ConfirmPrompt {
	return core.NewConfirmPrompt(core.ConfirmPromptParams{
		Active:   "Yes",
		Inactive: "No",
	})
}

func TestChangeConfirmValue(t *testing.T) {
	p := newConfirmPrompt()

	assert.Equal(t, false, p.Value)
	assert.Equal(t, 0, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.KeyUp})
	assert.Equal(t, true, p.Value)
	assert.Equal(t, 1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.KeyDown})
	assert.Equal(t, false, p.Value)
	assert.Equal(t, 0, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.KeyLeft})
	assert.Equal(t, true, p.Value)
	assert.Equal(t, 1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.KeyRight})
	assert.Equal(t, false, p.Value)
	assert.Equal(t, 0, p.CursorIndex)
}
