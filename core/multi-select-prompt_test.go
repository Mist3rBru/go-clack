package core_test

import (
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newMultiSelectPrompt() *core.MultiSelectPrompt[string] {
	return core.NewMultiSelectPrompt(core.MultiSelectPromptParams[string]{
		Options: []*core.MultiSelectOption[string]{
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
	assert.Equal(t, true, p.Options[0].IsSelected)
	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, []string{}, p.Value)
	assert.Equal(t, false, p.Options[0].IsSelected)

	expected := make([]string, len(p.Options))
	for i, option := range p.Options {
		expected[i] = option.Value
	}
	p.PressKey(&core.Key{Name: "a"})
	assert.Equal(t, expected, p.Value)
	p.PressKey(&core.Key{Name: "a"})
	assert.Equal(t, []string{}, p.Value)

	for _, option := range p.Options {
		option.IsSelected = true
		p.Value = append(p.Value, option.Value)
	}
	p.CursorIndex = 1
	p.PressKey(&core.Key{Name: core.SpaceKey})
	expected = append([]string{expected[0]}, expected[2:]...)
	assert.Equal(t, expected, p.Value)
}

func TestMultiSelectInitialValue(t *testing.T) {
	initialValue := []string{"a", "c"}
	p := core.NewMultiSelectPrompt(core.MultiSelectPromptParams[string]{
		Options: []*core.MultiSelectOption[string]{
			{Value: "a"},
			{Value: "b"},
			{Value: "c"},
		},
		InitialValue: initialValue,
	})

	assert.Equal(t, initialValue, p.Value)

	for _, value := range p.Value {
		for _, option := range p.Options {
			if option.Value == value {
				assert.True(t, option.IsSelected, option.Value)
			}
		}
	}
}

func TestMultiSelectInitialSelectedOptions(t *testing.T) {
	initialValue := []string{"a", "c"}
	p := core.NewMultiSelectPrompt(core.MultiSelectPromptParams[string]{
		Options: []*core.MultiSelectOption[string]{
			{Value: "a", IsSelected: true},
			{Value: "b"},
			{Value: "c", IsSelected: true},
		},
	})

	assert.Equal(t, initialValue, p.Value)

	for _, value := range p.Value {
		for _, option := range p.Options {
			if option.Value == value {
				assert.True(t, option.IsSelected, option.Value)
			}
		}
	}
}
