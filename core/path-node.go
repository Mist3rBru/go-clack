package core

import (
	"path"
	"regexp"

	"github.com/Mist3rBru/go-clack/core/internals"
)

type PathNode struct {
	Index  int
	Depth  int
	Path   string
	Name   string
	Parent *PathNode

	IsDir    bool
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
	root.Children = root.MapChildren()
	return root
}

func (p *PathNode) MapChildren() []*PathNode {
	if !p.IsDir {
		return nil
	}
	if len(p.Children) > 0 {
		return p.Children
	}
	entries, err := p.FileSystem.ReadDir(p.Path)
	if err != nil {
		return nil
	}
	children := []*PathNode{}
	for _, entry := range entries {
		if p.OnlyShowDir && !entry.IsDir() {
			continue
		}
		child := &PathNode{
			Index:  len(children),
			Depth:  p.Depth + 1,
			Path:   path.Join(p.Path, entry.Name()),
			Name:   entry.Name(),
			Parent: p,

			FileSystem:  p.FileSystem,
			OnlyShowDir: p.OnlyShowDir,
		}
		if entry.IsDir() {
			child.IsDir = true
			child.Children = []*PathNode(nil)
		}
		children = append(children, child)
	}
	return children
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

func (p *PathNode) FilteredFlat(search string, currentNode *PathNode) ([]*PathNode, *PathNode) {
	searchRegex, err := regexp.Compile("(?i)" + search)
	if err != nil || search == "" {
		return p.Flat(), currentNode
	}

	var options []*PathNode
	var firstNode *PathNode
	var hasRemovedCurrentNode bool

	p.TraverseNodes(func(node *PathNode) {
		if search != "" && node.Depth == currentNode.Depth && node.Depth > 0 {
			if matched := searchRegex.MatchString(node.Name); matched {
				options = append(options, node)
				if firstNode == nil && node.Depth == currentNode.Depth {
					firstNode = node
				}
			} else if node.IsEqual(currentNode) {
				hasRemovedCurrentNode = true
			}
		} else {
			options = append(options, node)
		}
	})

	if hasRemovedCurrentNode && firstNode != nil {
		return options, firstNode
	}
	return options, currentNode
}

func (p *PathNode) IsEqual(node *PathNode) bool {
	return node.Path == p.Path
}
