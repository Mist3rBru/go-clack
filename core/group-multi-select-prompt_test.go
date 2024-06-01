package core_test

import (
	"os"
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newGroupMultiSelectPrompt() *core.GroupMultiSelectPrompt[string] {
	options := make(map[string][]core.SelectOption[string])
	options["foo"] = []core.SelectOption[string]{
		{Value: "a"},
		{Value: "b"},
		{Value: "c"},
	}
	options["bar"] = []core.SelectOption[string]{
		{Value: "x"},
		{Value: "y"},
		{Value: "z"},
	}
	return core.NewGroupMultiSelectPrompt(core.GroupMultiSelectPromptParams[string]{
		Input:   os.Stdin,
		Output:  os.Stdout,
		Options: options,
	})
}

func TestChangeGroupMultiSelectCursor(t *testing.T) {
	p := newGroupMultiSelectPrompt()

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

func TestSelectGroupMultiSelectOption(t *testing.T) {
	p := newGroupMultiSelectPrompt()

	assert.Equal(t, []string(nil), p.Value)

	p.PressKey(&core.Key{Name: "Down"})
	p.PressKey(&core.Key{Name: "Space"})
	assert.Equal(t, []string{p.Options[1].Value}, p.Value)

	p.PressKey(&core.Key{Name: "Down"})
	p.PressKey(&core.Key{Name: "Space"})
	assert.Equal(t, []string{p.Options[1].Value, p.Options[2].Value}, p.Value)

	p.PressKey(&core.Key{Name: "Space"})
	assert.Equal(t, []string{p.Options[1].Value}, p.Value)
}

func TestSelectGroupMultiSelectGroupOption(t *testing.T) {
	p := newGroupMultiSelectPrompt()

	assert.Equal(t, []string(nil), p.Value)

	p.PressKey(&core.Key{Name: "Space"})
	expected := []string{}
	for i := 1; i < len(p.Options); i++ {
		if p.Options[i].IsGroup {
			break
		}
		expected = append(expected, p.Options[i].Value)
	}
	assert.Equal(t, expected, p.Value)

	p.PressKey(&core.Key{Name: "Space"})
	assert.Equal(t, []string{}, p.Value)

	expected = []string{}
	p.Value = []string{}
	for i, option := range p.Options {
		if option.IsGroup {
			p.CursorIndex = i
			p.PressKey(&core.Key{Name: "Space"})
		} else {
			expected = append(expected, option.Value)
		}
	}
	assert.Equal(t, expected, p.Value)
}

func TestGroupMultiSelectIsSelected(t *testing.T) {
	p := newGroupMultiSelectPrompt()

	i, isSelected := p.IsSelected(p.Options[0])
	assert.Equal(t, -1, i)
	assert.Equal(t, false, isSelected)

	p.Value = []string{p.Options[1].Value}
	i, isSelected = p.IsSelected(p.Options[1])
	assert.Equal(t, 0, i)
	assert.Equal(t, true, isSelected)

	p.Value = []string{p.Options[1].Value, p.Options[2].Value}
	i, isSelected = p.IsSelected(p.Options[2])
	assert.Equal(t, 1, i)
	assert.Equal(t, true, isSelected)
}

func TestGroupMultiSelectIsGroupSelected(t *testing.T) {
	p := newGroupMultiSelectPrompt()
	group := p.Options[0]

	isSelected := p.IsGroupSelected(group)
	assert.Equal(t, false, isSelected)

	p.Value = []string{p.Options[1].Value}
	isSelected = p.IsGroupSelected(group)
	assert.Equal(t, false, isSelected)

	p.Value = []string{}
	for _, option := range group.Options {
		p.Value = append(p.Value, option.Value)
	}
	isSelected = p.IsGroupSelected(group)
	assert.Equal(t, true, isSelected)
}
