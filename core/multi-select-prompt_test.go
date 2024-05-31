package core_test

import (
	"os"
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newMultiSelectPrompt() *core.MultiSelectPrompt {
	return core.NewMultiSelectPrompt(core.MultiSelectPromptParams{
		Input:  os.Stdin,
		Output: os.Stdout,
		Options: []core.SelectOption{
			{Value: "a"},
			{Value: "b"},
			{Value: "c"},
		},
	})
}

func TestChangeMultiSelectCursor(t *testing.T) {
	p := newMultiSelectPrompt()

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

func TestChangeMultiSelectValue(t *testing.T) {
	p := newMultiSelectPrompt()

	assert.Equal(t, []any(nil), p.Value)
	p.PressKey(&core.Key{Name: "Space"})
	assert.Equal(t, []any{p.Options[0].Value}, p.Value)
	p.PressKey(&core.Key{Name: "Space"})
	assert.Equal(t, []any{}, p.Value)

	expected := make([]any, len(p.Options))
	for i, option := range p.Options {
		expected[i] = option.Value
	}
	p.PressKey(&core.Key{Name: "a"})
	assert.Equal(t, expected, p.Value)
	p.PressKey(&core.Key{Name: "a"})
	assert.Equal(t, []any{}, p.Value)

	p.Value = append([]any{}, expected...)
	p.CursorIndex = 1
	p.PressKey(&core.Key{Name: "Space"})
	expected = append([]any{expected[0]}, expected[2:]...)
	assert.Equal(t, expected, p.Value)
}

func TestMultiSelectIsSelected(t *testing.T) {
	p := newMultiSelectPrompt()

	i, isSelected := p.IsSelected(p.Options[0])
	assert.Equal(t, -1, i)
	assert.Equal(t, false, isSelected)

	p.Value = []any{p.Options[0].Value}
	i, isSelected = p.IsSelected(p.Options[0])
	assert.Equal(t, 0, i)
	assert.Equal(t, true, isSelected)

	p.Value = []any{p.Options[0].Value, p.Options[1].Value}
	i, isSelected = p.IsSelected(p.Options[1])
	assert.Equal(t, 1, i)
	assert.Equal(t, true, isSelected)
}
