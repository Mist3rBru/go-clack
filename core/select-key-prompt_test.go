package core_test

import (
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newSelectKeyPrompt() *core.SelectKeyPrompt[string] {
	return core.NewSelectKeyPrompt(core.SelectKeyPromptParams[string]{
		Options: []*core.SelectKeyOption[string]{
			{
				Key:   "a",
				Value: "a",
			},
			{
				Key:   "Enter",
				Value: "enter",
			},
		},
	})
}

func TestSelectKeyPromptKey(t *testing.T) {
	p := newSelectKeyPrompt()

	p.PressKey(&core.Key{Name: "invalid-key"})
	assert.Equal(t, core.ActiveState, p.State)
	assert.Equal(t, "", p.Value)

	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, core.SubmitState, p.State)
	assert.Equal(t, "enter", p.Value)

	p.State = core.ActiveState
	p.PressKey(&core.Key{Name: "a"})
	assert.Equal(t, core.SubmitState, p.State)
	assert.Equal(t, "a", p.Value)
}

func TestKeyAsSelectValue(t *testing.T) {
	p := core.NewSelectKeyPrompt(core.SelectKeyPromptParams[string]{
		Options: []*core.SelectKeyOption[string]{
			{Key: "a", Label: "foo"},
			{Key: "b", Label: "bar"},
			{Key: "c", Label: "baz"},
		},
	})

	p.PressKey(&core.Key{Name: "a"})
	assert.Equal(t, "a", p.Value)
}
