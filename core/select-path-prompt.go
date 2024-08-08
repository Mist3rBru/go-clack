package core

import (
	"os"
	"path/filepath"

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
	Search        string
	Filter        bool
	FileSystem    FileSystem
}

type SelectPathPromptParams struct {
	Input        *os.File
	Output       *os.File
	InitialValue string
	OnlyShowDir  bool
	Filter       bool
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
		Filter:      params.Filter,
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
	var options []*PathNode
	options, p.CurrentOption = p.Root.FilteredFlat(p.Search, p.CurrentOption)
	return options
}

func (p *SelectPathPrompt) exitChildren() {
	if p.CurrentOption.IsRoot() {
		p.Root = NewPathNode(filepath.Dir(p.Root.Path), PathNodeOptions{
			OnlyShowDir: p.OnlyShowDir,
			FileSystem:  p.FileSystem,
		})
		p.CurrentLayer = []*PathNode{p.Root}
		p.CurrentOption = p.Root
		return
	}

	if p.CurrentOption.Parent.IsRoot() {
		p.CurrentLayer = []*PathNode{p.Root}
		p.CurrentOption = p.Root
		return
	}

	p.CurrentLayer = p.CurrentOption.Parent.Parent.Children
	p.CurrentOption = p.CurrentOption.Parent
	p.CurrentOption.ClearChildren()
}

func (p *SelectPathPrompt) enterChildren() {
	p.CurrentOption.MapChildren()
	if len(p.CurrentOption.Children) == 0 {
		return
	}
	p.CurrentLayer = p.CurrentOption.Children
	p.CurrentOption = p.CurrentOption.Children[0]
}

func (p *SelectPathPrompt) handleKeyPress(key *Key) {
	switch key.Name {
	case UpKey:
		p.CurrentOption = p.CurrentLayer[utils.MinMaxIndex(p.CurrentOption.Index-1, len(p.CurrentLayer))]
	case DownKey:
		p.CurrentOption = p.CurrentLayer[utils.MinMaxIndex(p.CurrentOption.Index+1, len(p.CurrentLayer))]
	case LeftKey:
		p.exitChildren()
		p.Search = ""
	case RightKey:
		p.enterChildren()
		p.Search = ""
	case HomeKey:
		p.CurrentOption = p.CurrentLayer[0]
	case EndKey:
		p.CurrentOption = p.CurrentLayer[len(p.CurrentLayer)-1]
	default:
		if p.Filter {
			p.Search, _ = p.TrackKeyValue(key, p.Search, len(p.Search))
		}
	}
	p.Value = p.CurrentOption.Path
	p.CursorIndex = p.Root.IndexOf(p.CurrentOption, p.Options())
}
