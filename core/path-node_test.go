package core_test

import (
	"os"
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func TestPathNodeRoot(t *testing.T) {
	cwd, _ := os.Getwd()
	node := core.NewPathNode(cwd, core.PathNodeOptions{
		OnlyShowDir: true,
		FileSystem:  MockFileSystem{},
	})

	assert.Equal(t, cwd, node.Path)
	assert.Equal(t, cwd, node.Name)
	assert.Equal(t, true, node.IsDir)
	assert.Equal(t, 0, node.Index)
	assert.Equal(t, 0, node.Depth)
	assert.Equal(t, (*core.PathNode)(nil), node.Parent)
	assert.Equal(t, true, node.OnlyShowDir)
	assert.Greater(t, len(node.Children), 0)
}

func TestPathNodeIsRoot(t *testing.T) {
	node := core.NewPathNode("/root/go-clack/core", core.PathNodeOptions{FileSystem: MockFileSystem{}})
	other := core.NewPathNode("/root/go-clack/prompts", core.PathNodeOptions{FileSystem: MockFileSystem{}})

	assert.True(t, node.IsRoot())

	node.Parent = other
	assert.False(t, node.IsRoot())
}

func TestPathNodeIsEqual(t *testing.T) {
	node := core.NewPathNode("/root/go-clack/core", core.PathNodeOptions{FileSystem: MockFileSystem{}})
	other := core.NewPathNode("/root/go-clack/prompts", core.PathNodeOptions{FileSystem: MockFileSystem{}})

	assert.True(t, node.IsEqual(node))
	assert.False(t, node.IsEqual(other))
}

func TestPathNodeTraverseNodes(t *testing.T) {
	node := core.NewPathNode("/root/go-clack/core", core.PathNodeOptions{FileSystem: MockFileSystem{}})
	node.MapChildren()

	counter := 0
	node.TraverseNodes(func(node *core.PathNode) {
		counter++
	})

	assert.Equal(t, len(node.Children)+1, counter)
}

func TestPathNodeFlat(t *testing.T) {
	node := core.NewPathNode("/root/go-clack/core", core.PathNodeOptions{FileSystem: MockFileSystem{}})

	assert.Equal(t, 3, len(node.Flat()))
}

func TestPathNodeFilteredFlat(t *testing.T) {
	root := core.NewPathNode("/root/go-clack/core", core.PathNodeOptions{FileSystem: MockFileSystem{}})
	var options []*core.PathNode

	options, _ = root.FilteredFlat("", root)
	assert.Equal(t, 3, len(options))

	options, _ = root.FilteredFlat("d", root)
	assert.Equal(t, 3, len(options))

	options, _ = root.FilteredFlat("", root.Children[0])
	assert.Equal(t, 3, len(options))

	options, _ = root.FilteredFlat("d", root.Children[0])
	assert.Equal(t, 2, len(options))

	options, currentNode := root.FilteredFlat("f", root.Children[0])
	assert.Equal(t, 2, len(options))
	assert.Equal(t, root.Children[1], currentNode)
}

func TestPathNodeIndexOf(t *testing.T) {
	node := core.NewPathNode("/root/go-clack/core", core.PathNodeOptions{FileSystem: MockFileSystem{}})
	node.MapChildren()

	assert.Equal(t, 1, node.IndexOf(node.Children[1], node.Children))
}
