package core

import (
	"os"
	"path"
	"sort"

	"github.com/Mist3rBru/go-clack/core/internals"
	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/core/validator"
)

type MultiSelectPathPrompt struct {
	Prompt[[]string]
	Root          *PathNode
	CurrentLayer  []*PathNode
	CurrentOption *PathNode
	OnlyShowDir   bool
	Filter        bool
	Search        string
	Required      bool
	FileSystem    FileSystem
}

type MultiSelectPathPromptParams struct {
	Input        *os.File
	Output       *os.File
	InitialValue []string
	InitialPath  string
	OnlyShowDir  bool
	Required     bool
	Filter       bool
	FileSystem   FileSystem
	Validate     func(value []string) error
	Render       func(p *MultiSelectPathPrompt) string
}

func NewMultiSelectPathPrompt(params MultiSelectPathPromptParams) *MultiSelectPathPrompt {
	v := validator.NewValidator("MultiSelectPathPrompt")
	v.ValidateRender(params.Render)

	if params.FileSystem == nil {
		params.FileSystem = internals.OSFileSystem{}
	}

	var p MultiSelectPathPrompt
	p = MultiSelectPathPrompt{
		Prompt: *NewPrompt(PromptParams[[]string]{
			Input:        params.Input,
			Output:       params.Output,
			InitialValue: params.InitialValue,
			CursorIndex:  1,
			Validate:     WrapValidate(params.Validate, &p.Required, "Please select at least one option. Press `space` to select"),
			Render:       WrapRender[[]string](&p, params.Render),
		}),
		OnlyShowDir: params.OnlyShowDir,
		Filter:      params.Filter,
		Required:    params.Required,
		FileSystem:  params.FileSystem,
	}

	if cwd, err := p.FileSystem.Getwd(); err == nil && params.InitialPath == "" {
		params.InitialPath = cwd
	}
	p.Root = NewPathNode(params.InitialPath, PathNodeOptions{
		OnlyShowDir: p.OnlyShowDir,
		FileSystem:  p.FileSystem,
	})
	p.CurrentLayer = p.Root.Children
	p.CurrentOption = p.Root.Children[0]
	p.mapSelectedOptions(p.Root)

	p.On(KeyEvent, func(args ...any) {
		p.handleKeyPress(args[0].(*Key))
	})
	p.On(FinalizeEvent, func(args ...any) {
		sort.Slice(p.Value, func(i, j int) bool {
			return p.Value[i] < p.Value[j]
		})
	})

	return &p
}

func (p *MultiSelectPathPrompt) Options() []*PathNode {
	var options []*PathNode
	options, p.CurrentOption = p.Root.FilteredFlat(p.Search, p.CurrentOption)
	return options
}

func (p *MultiSelectPathPrompt) exitChildren() {
	if p.CurrentOption.IsEqual(p.Root) {
		p.Root = NewPathNode(path.Dir(p.Root.Path), PathNodeOptions{
			OnlyShowDir: p.OnlyShowDir,
			FileSystem:  p.FileSystem,
		})
		p.CurrentLayer = []*PathNode{p.Root}
		p.CurrentOption = p.Root
		p.mapSelectedOptions(p.Root)
		return
	}
	if p.CurrentOption.Parent.IsEqual(p.Root) {
		p.CurrentLayer = []*PathNode{p.Root}
		p.CurrentOption = p.Root
		return
	}
	p.CurrentLayer = p.CurrentOption.Parent.Parent.Children
	p.CurrentOption = p.CurrentOption.Parent
	p.CurrentOption.ClearChildren()
}

func (p *MultiSelectPathPrompt) enterChildren() {
	p.CurrentOption.MapChildren()
	if len(p.CurrentOption.Children) == 0 {
		return
	}
	p.mapSelectedOptions(p.CurrentOption)
	p.CurrentLayer = p.CurrentOption.Children
	p.CurrentOption = p.CurrentOption.Children[0]
}

func (p *MultiSelectPathPrompt) handleKeyPress(key *Key) {
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
	case SpaceKey:
		if p.CurrentOption.IsSelected {
			p.CurrentOption.IsSelected = false
			value := []string{}
			for _, v := range p.Value {
				if v != p.CurrentOption.Path {
					value = append(value, v)
				}
			}
			p.Value = value
		} else {
			p.CurrentOption.IsSelected = true
			p.Value = append(p.Value, p.CurrentOption.Path)
		}
	default:
		if p.Filter {
			p.Search, _ = p.TrackKeyValue(key, p.Search, len(p.Search))
		}
	}
	if key.Name != SpaceKey {
		p.CursorIndex = p.Root.IndexOf(p.CurrentOption, p.Options())
	}
}

func (p *MultiSelectPathPrompt) mapSelectedOptions(node *PathNode) {
	node.TraverseNodes(func(node *PathNode) {
		for _, path := range p.Value {
			if path == node.Path {
				node.IsSelected = true
				return
			}
		}
		node.IsSelected = false
	})
}
