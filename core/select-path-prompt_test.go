package core_test

import (
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newSelectPathPrompt() *core.SelectPathPrompt {
	return core.NewSelectPathPrompt(core.SelectPathPromptParams{})
}

func TestChangeSelectPathCursor(t *testing.T) {
	p := newSelectPathPrompt()

	assert.Equal(t, 1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, 2, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.UpKey})
	assert.Equal(t, 1, p.CursorIndex)

	p.PressKey(&core.Key{Name: core.EndKey})
	assert.Equal(t, len(p.Options())-1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.HomeKey})
	assert.Equal(t, 1, p.CursorIndex)

	p.PressKey(&core.Key{Name: core.HomeKey})
	p.PressKey(&core.Key{Name: core.UpKey})
	assert.Equal(t, len(p.Options())-1, p.CursorIndex)
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, 1, p.CursorIndex)
}

func TestChangeSelectPathValue(t *testing.T) {
	p := newSelectPathPrompt()

	assert.Equal(t, p.CurrentLayer[0].Path, p.Value)
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, p.CurrentLayer[1].Path, p.Value)
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, p.CurrentLayer[2].Path, p.Value)
	p.PressKey(&core.Key{Name: core.UpKey})
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
	p.PressKey(&core.Key{Name: core.RightKey})
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
	p.PressKey(&core.Key{Name: core.RightKey})
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
	p.PressKey(&core.Key{Name: core.RightKey})
	p.PressKey(&core.Key{Name: core.LeftKey})
	assert.Equal(t, 1, p.CurrentOption.Depth)
	assert.Equal(t, pastOption, p.CurrentOption)
}

func TestExitRootDirectory(t *testing.T) {
	p := newSelectPathPrompt()

	p.PressKey(&core.Key{Name: core.LeftKey})
	assert.Equal(t, 0, p.CurrentOption.Depth)
	assert.Equal(t, p.Root, p.CurrentOption)

	pastChildrenLength := len(p.Root.Children)
	p.PressKey(&core.Key{Name: core.LeftKey})
	assert.Equal(t, 0, p.CurrentOption.Depth)
	assert.Equal(t, p.Root, p.CurrentOption)
	assert.NotEqual(t, pastChildrenLength, len(p.Root.Children))
}
