package prompt_test

import (
	"go-clack/internal/prompt"
	"go-clack/internal/utils"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var color = utils.CreateColors()

func newTextPrompt() *prompt.TextPrompt {
	return prompt.NewTextPrompt(prompt.TextPromptOptions{
		Input:  os.Stdin,
		Output: os.Stdout,
		Value:  "",
		Render: func(p *prompt.TextPrompt) string {
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
