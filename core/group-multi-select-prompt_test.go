package core_test

import (
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newGroupMultiSelectPrompt() *core.GroupMultiSelectPrompt[string] {
	return core.NewGroupMultiSelectPrompt(core.GroupMultiSelectPromptParams[string]{
		Options: map[string][]core.MultiSelectOption[string]{
			"foo": {
				{Value: "a"},
				{Value: "b"},
				{Value: "c"},
			},
			"bar": {
				{Value: "x"},
				{Value: "y"},
				{Value: "z"},
			},
		},
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

func TestLabelAsGroupMultiSelectValue(t *testing.T) {
	p := core.NewGroupMultiSelectPrompt(core.GroupMultiSelectPromptParams[string]{
		Options: map[string][]core.MultiSelectOption[string]{
			"group": {
				{Label: "foo"},
				{Label: "bar"},
				{Label: "baz"},
			},
		},
	})

	p.PressKey(&core.Key{Name: core.DownKey})
	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, []string{"foo"}, p.Value)
	p.PressKey(&core.Key{Name: core.DownKey})
	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, []string{"foo", "bar"}, p.Value)
}

func TestGroupMultiSelectRequiredValue(t *testing.T) {
	p := core.NewGroupMultiSelectPrompt(core.GroupMultiSelectPromptParams[string]{
		Required: true,
		Options: map[string][]core.MultiSelectOption[string]{
			"foo": {
				{Value: "a"},
				{Value: "b"},
				{Value: "c"},
			},
		},
	})

	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, core.ErrorState, p.State)
}
