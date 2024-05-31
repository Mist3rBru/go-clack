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
	Prompt
	Root          *PathNode
	CurrentLayer  []*PathNode
	CurrentOption *PathNode
	Value         string
	OnlyShowDir   bool
}

type SelectPathPromptParams struct {
	Input       *os.File
	Output      *os.File
	Value       string
	OnlyShowDir bool
	Render      func(p *SelectPathPrompt) string
}

func NewSelectPathPrompt(params SelectPathPromptParams) *SelectPathPrompt {
	var p *SelectPathPrompt
	p = &SelectPathPrompt{
		Prompt: *NewPrompt(PromptParams{
			Input:  params.Input,
			Output: params.Output,
			Value:  params.Value,
			Track:  false,
			Render: func(_p *Prompt) string {
				return params.Render(p)
			},
		}),
		OnlyShowDir: params.OnlyShowDir,
	}
	if cwd, err := os.Getwd(); err == nil && params.Value == "" {
		params.Value = cwd
	}
	p.Root = p.createRoot(params.Value)
	p.CurrentLayer = p.Root.Children
	p.CurrentOption = p.Root.Children[0]

	p.Prompt.On("key", func(args ...any) {
		key := args[0].(string)
		switch key {
		case "ArrowUp":
			p.CurrentOption = p.CurrentLayer[utils.MinMaxIndex(p.CurrentOption.Index-1, len(p.CurrentLayer))]
		case "ArrowDown":
			p.CurrentOption = p.CurrentLayer[utils.MinMaxIndex(p.CurrentOption.Index+1, len(p.CurrentLayer))]
		case "ArrowLeft":
			p.exitChildren()
		case "ArrowRight":
			p.enterChildren()
		case "Home":
			p.CurrentOption = p.CurrentLayer[0]
		case "End":
			p.CurrentOption = p.CurrentLayer[len(p.CurrentLayer)-1]
		}
		p.Value = p.CurrentOption.Path
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

func (p *SelectPathPrompt) CursorIndex() int {
	for i, option := range p.Options() {
		if option.Path == p.CurrentOption.Path {
			return i
		}
	}
	return -1
}

func (p *SelectPathPrompt) mapNode(node *PathNode) ([]*PathNode, error) {
	entries, err := os.ReadDir(node.Path)
	if err != nil {
		return nil, err
	}
	children := []*PathNode{}
	for i, entry := range entries {
		if p.OnlyShowDir && !entry.IsDir() {
			continue
		}
		child := PathNode{
			Index:  i,
			Depth:  node.Depth + 1,
			Path:   path.Join(node.Path, entry.Name()),
			Name:   entry.Name(),
			Parent: node,
		}
		if entry.IsDir() {
			child.Children = []*PathNode{}
		}
		children = append(children, &child)
	}
	return children, nil
}

func (p *SelectPathPrompt) createRoot(cwd string) *PathNode {
	root := &PathNode{
		Index: 0,
		Depth: 0,
		Path:  cwd,
		Name:  cwd,
	}
	root.Children, _ = p.mapNode(root)
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
	if p.CurrentOption.Children == nil {
		return
	}
	children, err := p.mapNode(p.CurrentOption)
	if err != nil || len(children) == 0 {
		return
	}
	p.CurrentOption.Children = children
	p.CurrentOption = children[0]
	p.CurrentLayer = children
}

func (p *SelectPathPrompt) Run() (string, error) {
	_, err := p.Prompt.Run()
	if err != nil {
		return "", err
	}
	return p.Value, nil
}
