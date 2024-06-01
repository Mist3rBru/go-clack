package core_test

import (
	"os"
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newSelectPrompt() *core.SelectPrompt[string] {
	return core.NewSelectPrompt(core.SelectPromptParams[string]{
		Input:  os.Stdin,
		Output: os.Stdout,
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
	p.PressKey(&core.Key{Name: "Down"})
	assert.Equal(t, 1, p.CursorIndex)
	p.PressKey(&core.Key{Name: "Right"})
	assert.Equal(t, 2, p.CursorIndex)
	p.PressKey(&core.Key{Name: "Up"})
	assert.Equal(t, 1, p.CursorIndex)
	p.PressKey(&core.Key{Name: "Left"})
	assert.Equal(t, 0, p.CursorIndex)

	p.PressKey(&core.Key{Name: "End"})
	assert.Equal(t, len(p.Options)-1, p.CursorIndex)
	p.PressKey(&core.Key{Name: "Home"})
	assert.Equal(t, 0, p.CursorIndex)

	p.CursorIndex = 0
	p.PressKey(&core.Key{Name: "Up"})
	assert.Equal(t, len(p.Options)-1, p.CursorIndex)
	p.PressKey(&core.Key{Name: "Down"})
	assert.Equal(t, 0, p.CursorIndex)
}

func TestChangeSelectValue(t *testing.T) {
	p := newSelectPrompt()

	assert.Equal(t, p.Options[0].Value, p.Value)
	p.PressKey(&core.Key{Name: "Down"})
	assert.Equal(t, p.Options[1].Value, p.Value)
	p.PressKey(&core.Key{Name: "Down"})
	assert.Equal(t, p.Options[2].Value, p.Value)
	p.PressKey(&core.Key{Name: "Up"})
	assert.Equal(t, p.Options[1].Value, p.Value)
}
