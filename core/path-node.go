package core

import (
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/Mist3rBru/go-clack/core/internals"
)

type PathNode struct {
	Index  int
	Depth  int
	Path   string
	Name   string
	Parent *PathNode

	IsDir    bool
	IsOpen   bool
	Children []*PathNode

	IsSelected bool

	FileSystem  FileSystem
	OnlyShowDir bool
}

type PathNodeOptions struct {
	OnlyShowDir bool
	FileSystem  FileSystem
}

func NewPathNode(rootPath string, options PathNodeOptions) *PathNode {
	if options.FileSystem == nil {
		options.FileSystem = internals.OSFileSystem{}
	}

	root := &PathNode{
		Path:  rootPath,
		Name:  rootPath,
		IsDir: true,

		OnlyShowDir: options.OnlyShowDir,
		FileSystem:  options.FileSystem,
	}
	root.Open()

	return root
}

func (p *PathNode) Open() {
	if !p.IsDir || p.IsOpen {
		return
	}

	entries, err := p.FileSystem.ReadDir(p.Path)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if p.OnlyShowDir && !entry.IsDir() {
			continue
		}
		p.Children = append(p.Children, &PathNode{
			Depth:  p.Depth + 1,
			Path:   path.Join(p.Path, entry.Name()),
			Name:   entry.Name(),
			Parent: p,
			IsDir:  entry.IsDir(),

			FileSystem:  p.FileSystem,
			OnlyShowDir: p.OnlyShowDir,
		})
	}

	sort.SliceStable(p.Children, func(i, j int) bool {
		if p.Children[i].IsDir != p.Children[j].IsDir {
			return p.Children[i].IsDir
		}
		return strings.ToLower(p.Children[i].Name) < strings.ToLower(p.Children[j].Name)
	})

	for i, child := range p.Children {
		child.Index = i
	}

	p.IsOpen = true
}

func (p *PathNode) Close() {
	p.Children = []*PathNode(nil)
	p.IsOpen = false
}

func (p *PathNode) TraverseNodes(visit func(node *PathNode)) {
	var traverse func(node *PathNode)
	traverse = func(node *PathNode) {
		visit(node)
		if !node.IsDir {
			return
		}
		for _, child := range node.Children {
			traverse(child)
		}
	}

	traverse(p)
}

func (p *PathNode) Flat() []*PathNode {
	var options []*PathNode
	p.TraverseNodes(func(node *PathNode) {
		options = append(options, node)
	})
	return options
}

func (p *PathNode) FilteredFlat(search string, currentNode *PathNode) []*PathNode {
	searchRegex, err := regexp.Compile("(?i)" + search)
	if err != nil || search == "" {
		return p.Flat()
	}

	var options []*PathNode
	p.TraverseNodes(func(node *PathNode) {
		if node.Depth == currentNode.Depth && node.Depth > 0 {
			if matched := searchRegex.MatchString(node.Name); matched {
				options = append(options, node)
			}
		} else {
			options = append(options, node)
		}
	})

	return options
}

func (p *PathNode) Layer() []*PathNode {
	if p.IsRoot() {
		return []*PathNode(nil)
	}

	return p.Parent.Children
}

func (p *PathNode) FilteredLayer(search string) []*PathNode {
	searchRegex, err := regexp.Compile("(?i)" + search)
	if err != nil || search == "" {
		return p.Layer()
	}

	var layer []*PathNode
	for _, node := range p.Layer() {
		if matched := searchRegex.MatchString(node.Name); matched {
			layer = append(layer, node)
		}
	}

	return layer
}

func (p *PathNode) FirstChild() *PathNode {
	if len(p.Children) == 0 {
		return nil
	}
	return p.Children[0]
}

func (p *PathNode) LastChild() *PathNode {
	if len(p.Children) == 0 {
		return nil
	}
	return p.Children[len(p.Children)-1]
}

func (p *PathNode) PrevChild(index int) *PathNode {
	if index <= 0 {
		return p.LastChild()
	}
	return p.Children[index-1]
}

func (p *PathNode) NextChild(index int) *PathNode {
	if index+1 >= len(p.Children) {
		return p.FirstChild()
	}
	return p.Children[index+1]
}

func (p *PathNode) IsRoot() bool {
	return p.Parent == nil
}

func (p *PathNode) IsEqual(node *PathNode) bool {
	return node.Path == p.Path
}

func (p *PathNode) IndexOf(node *PathNode, options []*PathNode) int {
	if node != nil {
		for i, option := range options {
			if option.IsEqual(node) {
				return i
			}
		}
	}
	return -1
}
