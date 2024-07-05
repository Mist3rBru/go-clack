package core

import (
	"os"
	"path"
	"sort"

	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/core/validator"
)

type MultiSelectPathPrompt struct {
	Prompt[[]string]
	Root          *PathNode
	CurrentLayer  []*PathNode
	CurrentOption *PathNode
	OnlyShowDir   bool
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
	FileSystem   FileSystem
	Validate     func(value []string) error
	Render       func(p *MultiSelectPathPrompt) string
}

func NewMultiSelectPathPrompt(params MultiSelectPathPromptParams) *MultiSelectPathPrompt {
	v := validator.NewValidator("MultiSelectPathPrompt")
	v.ValidateRender(params.Render)

	if params.FileSystem == nil {
		params.FileSystem = OSFileSystem{}
	}

	var p MultiSelectPathPrompt
	p = MultiSelectPathPrompt{
		Prompt: *NewPrompt(PromptParams[[]string]{
			Input:        params.Input,
			Output:       params.Output,
			InitialValue: params.InitialValue,
			CursorIndex:  1,
			Validate:     WrapValidateSlice(params.Validate, &p.Required, "Please select at least one option. Press `space` to select"),
			Render:       WrapRender[[]string](&p, params.Render),
		}),
		OnlyShowDir: params.OnlyShowDir,
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
		key := args[0].(*Key)
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
		}
		if key.Name != SpaceKey {
			p.CursorIndex = p.cursorIndex()
		}
	})
	p.On(FinalizeEvent, func(args ...any) {
		sort.Slice(p.Value, func(i, j int) bool {
			return p.Value[i] < p.Value[j]
		})
	})
	return &p
}

func (p *MultiSelectPathPrompt) Options() []*PathNode {
	options := []*PathNode{}

	var traverse func(node *PathNode)
	traverse = func(node *PathNode) {
		options = append(options, node)
		if node.Children == nil {
			return
		}
		for _, child := range node.Children {
			traverse(child)
		}
	}

	traverse(p.Root)
	return options
}

func (p *MultiSelectPathPrompt) cursorIndex() int {
	for i, option := range p.Options() {
		if option.Path == p.CurrentOption.Path {
			return i
		}
	}
	return -1
}

func (p *MultiSelectPathPrompt) exitChildren() {
	if p.CurrentOption.Path == p.Root.Path {
		p.Root = NewPathNode(path.Dir(p.Root.Path), PathNodeOptions{
			OnlyShowDir: p.OnlyShowDir,
			FileSystem:  p.FileSystem,
		})
		p.CurrentLayer = []*PathNode{p.Root}
		p.CurrentOption = p.Root
		p.mapSelectedOptions(p.Root)
		return
	}
	if p.CurrentOption.Parent.Path == p.Root.Path {
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

func (p *MultiSelectPathPrompt) enterChildren() {
	children := p.CurrentOption.MapChildren()
	if len(children) == 0 {
		return
	}
	p.CurrentOption.Children = children
	p.mapSelectedOptions(p.CurrentOption)
	p.CurrentOption = children[0]
	p.CurrentLayer = children
}

func (p *MultiSelectPathPrompt) mapSelectedOptions(node *PathNode) {
	var traverse func(node *PathNode)
	traverse = func(node *PathNode) {
		for _, path := range p.Value {
			if path == node.Path {
				node.IsSelected = true
				break
			}
		}
		if node.Children == nil {
			return
		}
		for _, child := range node.Children {
			traverse(child)
		}
	}

	traverse(node)
}
