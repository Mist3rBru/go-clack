package core_test

import (
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newMultiSelectPrompt() *core.MultiSelectPrompt[string] {
	return core.NewMultiSelectPrompt(core.MultiSelectPromptParams[string]{
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
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, 1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.RightKey})
	assert.Equal(t, 2, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.UpKey})
	assert.Equal(t, 1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.LeftKey})
	assert.Equal(t, 0, p.CursorIndex)

	p.PressKey(&core.Key{Name: core.EndKey})
	assert.Equal(t, len(p.Options)-1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.HomeKey})
	assert.Equal(t, 0, p.CursorIndex)

	p.CursorIndex = 0
	p.PressKey(&core.Key{Name: core.UpKey})
	assert.Equal(t, len(p.Options)-1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, 0, p.CursorIndex)
}

func TestChangeMultiSelectValue(t *testing.T) {
	p := newMultiSelectPrompt()

	assert.Equal(t, []string(nil), p.Value)
	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, []string{p.Options[0].Value}, p.Value)
	p.PressKey(&core.Key{Name: core.SpaceKey})
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
	p.PressKey(&core.Key{Name: core.SpaceKey})
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
