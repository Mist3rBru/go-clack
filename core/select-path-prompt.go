package core

import (
	"os"
	"path"

	"github.com/Mist3rBru/go-clack/core/utils"
)

type PathNode struct {
	Index    int
	Depth    int
	Path     string
	Name     string
	Parent   *PathNode
	Children []*PathNode
}

type SelectPathPrompt struct {
	Prompt[string]
	Root          *PathNode
	CurrentLayer  []*PathNode
	CurrentOption *PathNode
	Value         string
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
	if params.FileSystem == nil {
		params.FileSystem = OSFileSystem{}
	}

	var p *SelectPathPrompt
	p = &SelectPathPrompt{
		Prompt: *NewPrompt(PromptParams[string]{
			Input:       params.Input,
			Output:      params.Output,
			CursorIndex: 1,
			Render: func(_p *Prompt[string]) string {
				if params.Render == nil {
					return ErrMissingRender.Error()
				}
				return params.Render(p)
			},
		}),
		OnlyShowDir: params.OnlyShowDir,
		FileSystem:  params.FileSystem,
	}
	if cwd, err := p.FileSystem.Getwd(); err == nil && params.InitialValue == "" {
		params.InitialValue = cwd
	}
	p.Root = p.createRoot(params.InitialValue)
	p.CurrentLayer = p.Root.Children
	p.CurrentOption = p.Root.Children[0]
	p.Value = p.CurrentOption.Path

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
		}
		p.Value = p.CurrentOption.Path
		p.CursorIndex = p.cursorIndex()
	})
	return p
}

func (p *SelectPathPrompt) Options() []*PathNode {
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

func (p *SelectPathPrompt) cursorIndex() int {
	for i, option := range p.Options() {
		if option.Path == p.CurrentOption.Path {
			return i
		}
	}
	return -1
}

func (p *SelectPathPrompt) mapNode(node *PathNode) []*PathNode {
	if node.Children == nil {
		return nil
	}
	entries, err := p.FileSystem.ReadDir(node.Path)
	if err != nil {
		return nil
	}
	children := []*PathNode{}
	for i, entry := range entries {
		if p.OnlyShowDir && !entry.IsDir() {
			continue
		}
		child := &PathNode{
			Index:  i,
			Depth:  node.Depth + 1,
			Path:   path.Join(node.Path, entry.Name()),
			Name:   entry.Name(),
			Parent: node,
		}
		if entry.IsDir() {
			child.Children = []*PathNode{}
		}
		children = append(children, child)
	}
	return children
}

func (p *SelectPathPrompt) createRoot(cwd string) *PathNode {
	root := &PathNode{
		Index:    0,
		Depth:    0,
		Path:     cwd,
		Name:     cwd,
		Children: []*PathNode{},
	}
	root.Children = p.mapNode(root)
	return root
}

func (p *SelectPathPrompt) exitChildren() {
	if p.CurrentOption.Path == p.Root.Path {
		p.Root = p.createRoot(path.Dir(p.Root.Path))
		p.CurrentLayer = []*PathNode{p.Root}
		p.CurrentOption = p.Root
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

func (p *SelectPathPrompt) enterChildren() {
	children := p.mapNode(p.CurrentOption)
	if len(children) == 0 {
		return
	}
	p.CurrentOption.Children = children
	p.CurrentOption = children[0]
	p.CurrentLayer = children
}
