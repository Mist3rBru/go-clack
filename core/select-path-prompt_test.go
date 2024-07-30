package core_test

import (
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newSelectPathPrompt() *core.SelectPathPrompt {
	return core.NewSelectPathPrompt(core.SelectPathPromptParams{
		Render: func(p *core.SelectPathPrompt) string { return "" },
	})
}

func TestSelectPathChangeCursor(t *testing.T) {
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

func TestSelectPathChangeValue(t *testing.T) {
	p := newSelectPathPrompt()

	assert.Equal(t, p.CurrentLayer[0].Path, p.Value)
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, p.CurrentLayer[1].Path, p.Value)
	p.PressKey(&core.Key{Name: core.DownKey})
	assert.Equal(t, p.CurrentLayer[2].Path, p.Value)
	p.PressKey(&core.Key{Name: core.UpKey})
	assert.Equal(t, p.CurrentLayer[1].Path, p.Value)
}

func TestSelectPathEnterDirectory(t *testing.T) {
	p := newSelectPathPrompt()

	for _, node := range p.CurrentLayer {
		if node.IsDir {
			p.CurrentOption = node
			break
		}
	}
	pastOption := p.CurrentOption
	p.PressKey(&core.Key{Name: core.RightKey})
	assert.Equal(t, 2, p.CurrentOption.Depth)
	assert.Equal(t, pastOption, p.CurrentOption.Parent)
}

func TestSelectPathEnterNonDirectory(t *testing.T) {
	p := newSelectPathPrompt()

	for _, node := range p.CurrentLayer {
		if !node.IsDir {
			p.CurrentOption = node
			break
		}
	}
	pastOption := p.CurrentOption
	p.PressKey(&core.Key{Name: core.RightKey})
	assert.Equal(t, 1, p.CurrentOption.Depth)
	assert.Equal(t, pastOption, p.CurrentOption)
}

func TestSelectPathExitDirectory(t *testing.T) {
	p := newSelectPathPrompt()

	for _, node := range p.CurrentLayer {
		if node.IsDir {
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

func TestSelectPathExitRootDirectory(t *testing.T) {
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

func TestSelectPathFilter(t *testing.T) {
	p1 := core.NewSelectPathPrompt(core.SelectPathPromptParams{
		Render: func(p *core.SelectPathPrompt) string { return "" },
	})
	p2 := core.NewSelectPathPrompt(core.SelectPathPromptParams{
		Filter:     true,
		FileSystem: MockFileSystem{},
		Render:     func(p *core.SelectPathPrompt) string { return "" },
	})

	p1.PressKey(&core.Key{Char: "d"})
	p2.PressKey(&core.Key{Char: "d"})

	assert.Greater(t, len(p1.Options()), 0)
	assert.Greater(t, len(p2.Options()), 0)
	assert.Greater(t, len(p1.Options()), len(p2.Options()))
}
