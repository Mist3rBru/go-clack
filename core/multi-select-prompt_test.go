package core_test

import (
	"os"
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newMultiSelectPrompt() *core.MultiSelectPrompt[string] {
	return core.NewMultiSelectPrompt(core.MultiSelectPromptParams[string]{
		Input:  os.Stdin,
		Output: os.Stdout,
		Options: []core.SelectOption[string]{
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

	assert.Equal(t, []string(nil), p.Value)
	p.PressKey(&core.Key{Name: "Space"})
	assert.Equal(t, []string{p.Options[0].Value}, p.Value)
	p.PressKey(&core.Key{Name: "Space"})
	assert.Equal(t, []string{}, p.Value)

	expected := make([]string, len(p.Options))
	for i, option := range p.Options {
		expected[i] = option.Value
	}
	p.PressKey(&core.Key{Name: "a"})
	assert.Equal(t, expected, p.Value)
	p.PressKey(&core.Key{Name: "a"})
	assert.Equal(t, []string{}, p.Value)

	p.Value = append([]string{}, expected...)
	p.CursorIndex = 1
	p.PressKey(&core.Key{Name: "Space"})
	expected = append([]string{expected[0]}, expected[2:]...)
	assert.Equal(t, expected, p.Value)
}

func TestMultiSelectIsSelected(t *testing.T) {
	p := newMultiSelectPrompt()

	i, isSelected := p.IsSelected(p.Options[0])
	assert.Equal(t, -1, i)
	assert.Equal(t, false, isSelected)

	p.Value = []string{p.Options[0].Value}
	i, isSelected = p.IsSelected(p.Options[0])
	assert.Equal(t, 0, i)
	assert.Equal(t, true, isSelected)

	p.Value = []string{p.Options[0].Value, p.Options[1].Value}
	i, isSelected = p.IsSelected(p.Options[1])
	assert.Equal(t, 1, i)
	assert.Equal(t, true, isSelected)
}
