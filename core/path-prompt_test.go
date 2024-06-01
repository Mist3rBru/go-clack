package core_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newPathPrompt() *core.PathPrompt {
	return core.NewPathPrompt(core.PathPromptParams{
		Input:  os.Stdin,
		Output: os.Stdout,
	})
}

func TestPathDefaultValue(t *testing.T) {
	p := newPathPrompt()

	cwd, _ := os.Getwd()
	assert.Equal(t, cwd, p.Value)
	assert.Equal(t, len(cwd), p.CursorIndex)
}

func TestPathChangeValue(t *testing.T) {
	p := newPathPrompt()
	p.Prompt.Value = "/go-clack"
	p.Prompt.CursorIndex = len("/go-clack")
	p.Value = "/go-clack"

	assert.Equal(t, "/go-clack", p.Value)
	p.PressKey(&core.Key{Char: "/"})
	assert.Equal(t, "/go-clack/", p.Value)
	p.PressKey(&core.Key{Char: "a"})
	assert.Equal(t, "/go-clack/a", p.Value)
	p.PressKey(&core.Key{Name: "Backspace"})
	assert.Equal(t, "/go-clack/", p.Value)
}

func TestPathChangeHint(t *testing.T) {
	p := newPathPrompt()

	assert.Equal(t, "/", p.Hint)
	p.PressKey(&core.Key{Char: "/"})
	assert.Equal(t, "confirm-prompt.go", p.Hint)
	p.PressKey(&core.Key{Name: "Backspace"})
	assert.Equal(t, "/", p.Hint)
}

func TestPathTabComplete(t *testing.T) {
	p := newPathPrompt()

	assert.Equal(t, 0, len(p.HintOptions))
	assert.Equal(t, -1, p.HintIndex)
	assert.Equal(t, "/", p.Hint)

	p.PressKey(&core.Key{Name: "Tab"})
	assert.Equal(t, 0, len(p.HintOptions))
	assert.Equal(t, -1, p.HintIndex)
	assert.NotEqual(t, "", p.Hint)
	assert.NotEqual(t, "/", p.Hint)

	p.PressKey(&core.Key{Name: "Tab"})
	assert.GreaterOrEqual(t, len(p.HintOptions), 1)
	assert.Equal(t, 0, p.HintIndex)

	p.PressKey(&core.Key{Name: "Tab"})
	assert.GreaterOrEqual(t, len(p.HintOptions), 1)
	assert.Equal(t, 1, p.HintIndex)

	p.PressKey(&core.Key{Name: "Tab"})
	assert.GreaterOrEqual(t, len(p.HintOptions), 1)
	assert.Equal(t, 2, p.HintIndex)
}

func TestPathComplete(t *testing.T) {
	p := newPathPrompt()

	assert.Equal(t, "/", p.Hint)

	expected := p.Value + p.Hint
	p.PressKey(&core.Key{Name: "Right"})
	assert.Equal(t, expected, p.Value)
	assert.NotEqual(t, "", p.Hint)
	assert.NotEqual(t, "/", p.Hint)

	expected = p.Value + p.Hint
	p.PressKey(&core.Key{Name: "Right"})
	assert.Equal(t, expected, p.Value)
}

func TestPathValueWithHint(t *testing.T) {
	p := newPathPrompt()

	assert.Equal(t, p.Value+p.Hint, p.ValueWithHint())

	p.Hint = ""
	p.CursorIndex = len(p.Value)
	assert.Equal(t, p.Value+" ", p.ValueWithHint())

	p.Hint = ""
	p.CursorIndex = len(p.Value) - 1
	assert.Equal(t, p.Value, p.ValueWithHint())

	p.Hint = "/"
	p.CursorIndex = len(p.Value) - 1
	assert.Equal(t, p.Value+p.Hint, p.ValueWithHint())
}

func TestValidatePath(t *testing.T) {
	p := core.NewPathPrompt(core.PathPromptParams{
		Value: "/folder",
		Validate: func(value string) error {
			return fmt.Errorf("invalid path: %s", value)
		},
	})

	p.PressKey(&core.Key{Name: "Enter"})
	assert.Equal(t, "error", p.State)
	assert.Equal(t, "invalid path: /folder", p.Error)
}
