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
		Value:  "",
		Render: func(p *core.TextPrompt) string {
			return p.Value
		},
	})
}

func TestValeuWithCursor(t *testing.T) {
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
