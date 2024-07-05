package core_test

import (
	"path/filepath"
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newMultiSelectPathPrompt() *core.MultiSelectPathPrompt {
	return core.NewMultiSelectPathPrompt(core.MultiSelectPathPromptParams{
		Render: func(p *core.MultiSelectPathPrompt) string { return "" },
	})
}

func TestMultiSelectPathChangeCursor(t *testing.T) {
	p := newMultiSelectPathPrompt()

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

func TestMultiSelectPathChangeValue(t *testing.T) {
	p := newMultiSelectPathPrompt()
	options := p.Options()

	assert.Equal(t, []string(nil), p.Value)
	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, []string{options[1].Path}, p.Value)
	assert.Equal(t, true, options[1].IsSelected)
	p.PressKey(&core.Key{Name: core.SpaceKey})
	assert.Equal(t, []string{}, p.Value)
	assert.Equal(t, false, options[1].IsSelected)

	p.CurrentOption = options[1]
	for i := 0; i < 3; i++ {
		p.PressKey(&core.Key{Name: core.SpaceKey})
		p.PressKey(&core.Key{Name: core.DownKey})
	}
	p.CurrentOption = options[2]
	p.PressKey(&core.Key{Name: core.SpaceKey})
	expected := append([]string{options[1].Path}, options[3].Path)
	assert.Equal(t, expected, p.Value)
}

func TestMultiSelectPathEnterDirectory(t *testing.T) {
	p := newMultiSelectPathPrompt()

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

func TestMultiSelectPathEnterNonDirectory(t *testing.T) {
	p := newMultiSelectPathPrompt()

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

func TestMultiSelectPathExitDirectory(t *testing.T) {
	p := newMultiSelectPathPrompt()

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

func TestMultiSelectPathExitRootDirectory(t *testing.T) {
	p := newMultiSelectPathPrompt()

	p.PressKey(&core.Key{Name: core.LeftKey})
	assert.Equal(t, 0, p.CurrentOption.Depth)
	assert.Equal(t, p.Root, p.CurrentOption)

	pastChildrenLength := len(p.Root.Children)
	p.PressKey(&core.Key{Name: core.LeftKey})
	assert.Equal(t, 0, p.CurrentOption.Depth)
	assert.Equal(t, p.Root, p.CurrentOption)
	assert.NotEqual(t, pastChildrenLength, len(p.Root.Children))
}

func TestMultiSelectPathExitToSelectedDirectory(t *testing.T) {
	p := newMultiSelectPathPrompt()
	p.Value = []string{filepath.Dir(p.Root.Path)}

	p.PressKey(&core.Key{Name: core.LeftKey})
	p.PressKey(&core.Key{Name: core.LeftKey})
	assert.True(t, p.CurrentOption.IsSelected)
}

func TestMultiSelectPathRequiredValue(t *testing.T) {
	p := newMultiSelectPathPrompt()
	p.Required = true

	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, core.ErrorState, p.State)
}
