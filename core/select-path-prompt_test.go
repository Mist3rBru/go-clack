package core_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newSelectPathPrompt() *core.SelectPathPrompt {
	return core.NewSelectPathPrompt(core.SelectPathPromptParams{
		Input:  os.Stdin,
		Output: os.Stdout,
	})
}

func TestChangeSelectPathCursor(t *testing.T) {
	p := newSelectPathPrompt()

	assert.Equal(t, 1, p.CursorIndex)
	p.PressKey(&core.Key{Name: "Down"})
	assert.Equal(t, 2, p.CursorIndex)
	p.PressKey(&core.Key{Name: "Up"})
	assert.Equal(t, 1, p.CursorIndex)

	p.PressKey(&core.Key{Name: "End"})
	assert.Equal(t, len(p.Options())-1, p.CursorIndex)
	p.PressKey(&core.Key{Name: "Home"})
	assert.Equal(t, 1, p.CursorIndex)

	p.PressKey(&core.Key{Name: "Home"})
	p.PressKey(&core.Key{Name: "Up"})
	assert.Equal(t, len(p.Options())-1, p.CursorIndex)
	p.PressKey(&core.Key{Name: "Down"})
	assert.Equal(t, 1, p.CursorIndex)
}

func TestChangeSelectPathValue(t *testing.T) {
	p := newSelectPathPrompt()

	assert.Equal(t, p.CurrentLayer[0].Path, p.Value)
	p.PressKey(&core.Key{Name: "Down"})
	assert.Equal(t, p.CurrentLayer[1].Path, p.Value)
	p.PressKey(&core.Key{Name: "Down"})
	assert.Equal(t, p.CurrentLayer[2].Path, p.Value)
	p.PressKey(&core.Key{Name: "Up"})
	assert.Equal(t, p.CurrentLayer[1].Path, p.Value)
}

func TestEnterDirectory(t *testing.T) {
	p := newSelectPathPrompt()

	for _, node := range p.CurrentLayer {
		if node.Children != nil {
			p.CurrentOption = node
			break
		}
	}
	pastOption := p.CurrentOption
	p.PressKey(&core.Key{Name: "Right"})
	assert.Equal(t, 2, p.CurrentOption.Depth)
	assert.Equal(t, pastOption, p.CurrentOption.Parent)
}

func TestEnterNonDirectory(t *testing.T) {
	p := newSelectPathPrompt()

	for _, node := range p.CurrentLayer {
		if node.Children == nil {
			p.CurrentOption = node
			break
		}
	}
	pastOption := p.CurrentOption
	p.PressKey(&core.Key{Name: "Right"})
	assert.Equal(t, 1, p.CurrentOption.Depth)
	assert.Equal(t, pastOption, p.CurrentOption)
}

func TestExitDirectory(t *testing.T) {
	p := newSelectPathPrompt()

	for _, node := range p.CurrentLayer {
		if node.Children != nil {
			p.CurrentOption = node
			break
		}
	}
	pastOption := p.CurrentOption
	p.PressKey(&core.Key{Name: "Right"})
	p.PressKey(&core.Key{Name: "Left"})
	assert.Equal(t, 1, p.CurrentOption.Depth)
	assert.Equal(t, pastOption, p.CurrentOption)
}

func TestExitRootDirectory(t *testing.T) {
	p := newSelectPathPrompt()

	p.PressKey(&core.Key{Name: "Left"})
	assert.Equal(t, 0, p.CurrentOption.Depth)
	assert.Equal(t, p.Root, p.CurrentOption)

	pastChildrenLength := len(p.Root.Children)
	p.PressKey(&core.Key{Name: "Left"})
	assert.Equal(t, 0, p.CurrentOption.Depth)
	assert.Equal(t, p.Root, p.CurrentOption)
	assert.NotEqual(t, pastChildrenLength, len(p.Root.Children))
}

func TestValidateSelectPathValue(t *testing.T) {
	p := core.NewSelectPathPrompt(core.SelectPathPromptParams{
		Validate: func(path string) error {
			return fmt.Errorf("invalid path: %s", path)
		},
	})

	p.PressKey(&core.Key{Name: "Enter"})
	assert.Equal(t, "error", p.State)
	assert.Equal(t, fmt.Sprintf("invalid path: %s", p.CurrentOption.Path), p.Error)
}
