package prompt_test

import (
	"go-clack/internal/prompt"
	"go-clack/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

var color = utils.CreateColors()

func render(p *prompt.TextPrompt) string {
	return p.Value
}

func TestValeuWithCursor(t *testing.T) {
	p := prompt.DefaultTextPrompt(render)
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
