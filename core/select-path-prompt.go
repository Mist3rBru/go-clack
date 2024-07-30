package core

import (
	"os"
	"path"

	"github.com/Mist3rBru/go-clack/core/internals"
	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/core/validator"
)

type SelectPathPrompt struct {
	Prompt[string]
	Root          *PathNode
	CurrentLayer  []*PathNode
	CurrentOption *PathNode
	OnlyShowDir   bool
	FileSystem    FileSystem
}

type SelectPathPromptParams struct {
	Input        *os.File
	Output       *os.File
	InitialValue string
	OnlyShowDir  bool
	FileSystem   FileSystem
	Render       func(p *SelectPathPrompt) string
}

func NewSelectPathPrompt(params SelectPathPromptParams) *SelectPathPrompt {
	v := validator.NewValidator("SelectPathPrompt")
	v.ValidateRender(params.Render)

	if params.FileSystem == nil {
		params.FileSystem = internals.OSFileSystem{}
	}

	var p SelectPathPrompt
	p = SelectPathPrompt{
		Prompt: *NewPrompt(PromptParams[string]{
			Input:       params.Input,
			Output:      params.Output,
			CursorIndex: 1,
			Render:      WrapRender[string](&p, params.Render),
		}),
		OnlyShowDir: params.OnlyShowDir,
		FileSystem:  params.FileSystem,
	}

	if cwd, err := p.FileSystem.Getwd(); err == nil && params.InitialValue == "" {
		params.InitialValue = cwd
	}
	p.Root = NewPathNode(params.InitialValue, PathNodeOptions{
		OnlyShowDir: p.OnlyShowDir,
		FileSystem:  p.FileSystem,
	})
	p.CurrentLayer = p.Root.Children
	p.CurrentOption = p.Root.Children[0]
	p.Value = p.CurrentOption.Path

	p.On(KeyEvent, func(args ...any) {
		p.handleKeyPress(args[0].(*Key))
	})

	return &p
}

func (p *SelectPathPrompt) Options() []*PathNode {
	return p.Root.Flat()
}

func (p *SelectPathPrompt) cursorIndex() int {
	for i, option := range p.Options() {
		if option.IsEqual(p.CurrentOption) {
			return i
		}
	}
	return -1
}

func (p *SelectPathPrompt) exitChildren() {
	if p.CurrentOption.IsEqual(p.Root) {
		p.Root = NewPathNode(path.Dir(p.Root.Path), PathNodeOptions{
			OnlyShowDir: p.OnlyShowDir,
			FileSystem:  p.FileSystem,
		})
		p.CurrentLayer = []*PathNode{p.Root}
		p.CurrentOption = p.Root
		return
	}
	if p.CurrentOption.Parent.IsEqual(p.Root) {
		p.CurrentLayer = []*PathNode{p.Root}
		p.CurrentOption = p.Root
		return
	}
	p.CurrentLayer = p.CurrentOption.Parent.Parent.Children
	p.CurrentOption = p.CurrentOption.Parent
	if p.CurrentOption.Children != nil {
		p.CurrentOption.Children = []*PathNode{}
	}
}

func (p *SelectPathPrompt) enterChildren() {
	children := p.CurrentOption.MapChildren()
	if len(children) == 0 {
		return
	}
	p.CurrentOption.Children = children
	p.CurrentOption = children[0]
	p.CurrentLayer = children
}

func (p *SelectPathPrompt) handleKeyPress(key *Key) {
	switch key.Name {
	case UpKey:
		p.CurrentOption = p.CurrentLayer[utils.MinMaxIndex(p.CurrentOption.Index-1, len(p.CurrentLayer))]
	case DownKey:
		p.CurrentOption = p.CurrentLayer[utils.MinMaxIndex(p.CurrentOption.Index+1, len(p.CurrentLayer))]
	case LeftKey:
		p.exitChildren()
	case RightKey:
		p.enterChildren()
	case HomeKey:
		p.CurrentOption = p.CurrentLayer[0]
	case EndKey:
		p.CurrentOption = p.CurrentLayer[len(p.CurrentLayer)-1]
	}
	p.Value = p.CurrentOption.Path
	p.CursorIndex = p.cursorIndex()
}
