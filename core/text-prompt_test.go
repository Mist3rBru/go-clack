package core_test

import (
	"os"
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/core/utils"

	"github.com/stretchr/testify/assert"
)

var color = utils.CreateColors()

func newTextPrompt() *core.TextPrompt {
	return core.NewTextPrompt(core.TextPromptParams{
		Input:  os.Stdin,
		Output: os.Stdout,
		Render: func(p *core.TextPrompt) string {
			return p.Value
		},
	})
}

func TestTextPromptValueTrack(t *testing.T) {
	p := newTextPrompt()

	assert.Equal(t, "", p.Value)

	p.PressKey(&core.Key{Char: "a"})
	assert.Equal(t, "a", p.Value)

	p.PressKey(&core.Key{Char: "b"})
	assert.Equal(t, "ab", p.Value)

	p.PressKey(&core.Key{Name: "Left"})
	p.PressKey(&core.Key{Name: "Backspace"})
	assert.Equal(t, "b", p.Value)
}

func TestTextPromptValueWithCursor(t *testing.T) {
	p := newTextPrompt()
	inverse := color["inverse"]
	cursor := inverse(" ")

	expected := cursor
	assert.Equal(t, expected, p.ValueWithCursor())

	p.Value = "foo"
	p.CursorIndex = len(p.Value)
	expected = p.Value + cursor
	assert.Equal(t, expected, p.ValueWithCursor())

	p.CursorIndex = len(p.Value) - 1
	expected = "fo" + inverse("o")
	assert.Equal(t, expected, p.ValueWithCursor())

	p.CursorIndex = 0
	expected = inverse("f") + "oo"
	assert.Equal(t, expected, p.ValueWithCursor())
}
