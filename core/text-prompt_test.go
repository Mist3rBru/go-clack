package core_test

import (
	"fmt"
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

func TestTextPromptInitialValue(t *testing.T) {
	p := core.NewTextPrompt(core.TextPromptParams{
		Value: "foo",
	})

	assert.Equal(t, "foo", p.Value)
	assert.Equal(t, 3, p.CursorIndex)
}

func TestTextPromptValueTrack(t *testing.T) {
	p := newTextPrompt()

	assert.Equal(t, "", p.Value)

	p.PressKey(&core.Key{Char: "a"})
	assert.Equal(t, "a", p.Value)

	p.PressKey(&core.Key{Char: "b"})
	assert.Equal(t, "ab", p.Value)

	p.PressKey(&core.Key{Name: core.KeyLeft})
	p.PressKey(&core.Key{Name: core.KeyBackspace})
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

func TestValidateText(t *testing.T) {
	p := core.NewTextPrompt(core.TextPromptParams{
		Value: "123",
		Validate: func(value string) error {
			return fmt.Errorf("invalid value: %s", value)
		},
	})

	p.PressKey(&core.Key{Name: core.KeyEnter})
	assert.Equal(t, "error", p.State)
	assert.Equal(t, "invalid value: 123", p.Error)
}
