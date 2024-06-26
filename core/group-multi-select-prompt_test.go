package core_test

import (
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newGroupMultiSelectPrompt() *core.GroupMultiSelectPrompt[string] {
	options := make(map[string][]core.MultiSelectOption[string])
	options["foo"] = []core.MultiSelectOption[string]{
		{Value: "a"},
		{Value: "b"},
		{Value: "c"},
	}
	options["bar"] = []core.MultiSelectOption[string]{
		{Value: "x"},
		{Value: "y"},
		{Value: "z"},
	}
	return core.NewGroupMultiSelectPrompt(core.GroupMultiSelectPromptParams[string]{
		Options: options,
	})
}

func TestChangeGroupMultiSelectCursor(t *testing.T) {
	p := newGroupMultiSelectPrompt()

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

func TestSelectGroupMultiSelectOption(t *testing.T) {
	p := newGroupMultiSelectPrompt()

	assert.Equal(t, []string(nil), p.Value)

	p.PressKey(&core.Key{Name: core.DownKey})
	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, []string{p.Options[1].Value}, p.Value)

	p.PressKey(&core.Key{Name: core.DownKey})
	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, []string{p.Options[1].Value, p.Options[2].Value}, p.Value)

	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, []string{p.Options[1].Value}, p.Value)
}

func TestSelectGroupMultiSelectGroupOption(t *testing.T) {
	p := newGroupMultiSelectPrompt()

	assert.Equal(t, []string(nil), p.Value)

	p.PressKey(&core.Key{Name: core.SpaceKey})
	for _, option := range p.Options[0].Options {
		assert.Equal(t, true, option.IsSelected, option.Value)
	}
	assert.Equal(t, len(p.Options[0].Options), len(p.Value))

	p.PressKey(&core.Key{Name: core.SpaceKey})
	for _, option := range p.Options[0].Options {
		assert.Equal(t, false, option.IsSelected)
	}
	assert.Equal(t, 0, len(p.Value))
}

func TestGroupMultiSelectIsGroupSelected(t *testing.T) {
	p := newGroupMultiSelectPrompt()
	group := p.Options[0]

	isSelected := p.IsGroupSelected(group)
	assert.Equal(t, false, isSelected)

	for _, option := range group.Options {
		option.IsSelected = true
	}
	isSelected = p.IsGroupSelected(group)
	assert.Equal(t, true, isSelected)
}
