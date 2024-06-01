package core_test

import (
	"os"
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newSelectKeyPrompt() *core.SelectKeyPrompt {
	return core.NewSelectKeyPrompt(core.SelectKeyPromptParams{
		Input:  os.Stdin,
		Output: os.Stdout,
		Options: []core.SelectKeyOption{
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
	assert.Equal(t, "active", p.State)
	assert.Equal(t, nil, p.Value)

	p.PressKey(&core.Key{Name: "Enter"})
	assert.Equal(t, "submit", p.State)
	assert.Equal(t, "enter", p.Value)

	p.State = "active"
	p.PressKey(&core.Key{Name: "a"})
	assert.Equal(t, "submit", p.State)
	assert.Equal(t, "a", p.Value)
}
