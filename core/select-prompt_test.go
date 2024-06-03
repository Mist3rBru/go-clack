package core_test

import (
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newSelectPrompt() *core.SelectPrompt[string] {
	return core.NewSelectPrompt(core.SelectPromptParams[string]{
		Options: []core.SelectOption[string]{
			{Value: "a"},
			{Value: "b"},
			{Value: "c"},
		},
	})
}

func TestChangeSelectCursor(t *testing.T) {
	p := newSelectPrompt()

	assert.Equal(t, 0, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.KeyDown})
	assert.Equal(t, 1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.KeyRight})
	assert.Equal(t, 2, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.KeyUp})
	assert.Equal(t, 1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.KeyLeft})
	assert.Equal(t, 0, p.CursorIndex)

	p.PressKey(&core.Key{Name: core.KeyEnd})
	assert.Equal(t, len(p.Options)-1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.KeyHome})
	assert.Equal(t, 0, p.CursorIndex)

	p.CursorIndex = 0
	p.PressKey(&core.Key{Name: core.KeyUp})
	assert.Equal(t, len(p.Options)-1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.KeyDown})
	assert.Equal(t, 0, p.CursorIndex)
}

func TestChangeSelectValue(t *testing.T) {
	p := newSelectPrompt()

	assert.Equal(t, p.Options[0].Value, p.Value)
	p.PressKey(&core.Key{Name: core.KeyDown})
	assert.Equal(t, p.Options[1].Value, p.Value)
	p.PressKey(&core.Key{Name: core.KeyDown})
	assert.Equal(t, p.Options[2].Value, p.Value)
	p.PressKey(&core.Key{Name: core.KeyUp})
	assert.Equal(t, p.Options[1].Value, p.Value)
}
