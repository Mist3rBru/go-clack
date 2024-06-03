package utils_test

import (
	"testing"

	"github.com/Mist3rBru/go-clack/core/utils"

	"github.com/stretchr/testify/assert"
)

var (
	color = utils.CreateColors()
)

func TestDiffLines(t *testing.T) {
	assert.Equal(t, []int{1, 2}, utils.DiffLines("a", "a\nb\nc"))
	assert.Equal(t, []int{1, 2}, utils.DiffLines("a\nb\nc", "a"))
	assert.Equal(t, []int{1, 2}, utils.DiffLines("a\nb\nc", "a\nc\nb"))
	assert.Equal(t, []int{}, utils.DiffLines("a\nb\nc", "a\nb\nc"))
}

func TestStrLength(t *testing.T) {
	assert.Equal(t, 1, utils.StrLength(color["inverse"](" ")))
	assert.Equal(t, 5, utils.StrLength(color["cyan"]("| foo")))
	assert.Equal(t, 5, utils.StrLength(color["gray"]("|")+" "+color["dim"]("foo")))
	assert.Equal(t, 1, utils.StrLength(color["green"]("◆")))
}
