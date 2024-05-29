package utils_test

import (
	"go-clack/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffLines(t *testing.T) {
	assert.Equal(t, []int{1, 2}, utils.DiffLines("a", "a\nb\nc"))
	assert.Equal(t, []int{1, 2}, utils.DiffLines("a\nb\nc", "a"))
	assert.Equal(t, []int{1, 2}, utils.DiffLines("a\nb\nc", "a\nc\nb"))
	assert.Equal(t, []int{}, utils.DiffLines("a\nb\nc", "a\nb\nc"))
}
