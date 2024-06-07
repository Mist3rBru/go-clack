package utils_test

import (
	"testing"

	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"

	"github.com/stretchr/testify/assert"
)

func TestDiffLines(t *testing.T) {
	assert.Equal(t, []int{1, 2}, utils.DiffLines("a", "a\nb\nc"))
	assert.Equal(t, []int{1, 2}, utils.DiffLines("a\nb\nc", "a"))
	assert.Equal(t, []int{1, 2}, utils.DiffLines("a\nb\nc", "a\nc\nb"))
	assert.Equal(t, []int{}, utils.DiffLines("a\nb\nc", "a\nb\nc"))
}

func TestStrLength(t *testing.T) {
	assert.Equal(t, 1, utils.StrLength(picocolors.Inverse(" ")))
	assert.Equal(t, 5, utils.StrLength(picocolors.Cyan("| foo")))
	assert.Equal(t, 5, utils.StrLength(picocolors.Gray("|")+" "+picocolors.Dim("foo")))
	assert.Equal(t, 1, utils.StrLength(picocolors.Green("◆")))
	assert.Equal(t, 5, utils.StrLength(picocolors.Green("◇")+" "+"Foo"))
	assert.Equal(t, 5, utils.StrLength(picocolors.Green("o")+" "+"Foo"))
}
