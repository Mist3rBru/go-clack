package prompts_test

import (
	"strings"
	"testing"

	"github.com/Mist3rBru/go-clack/prompts"
	"github.com/stretchr/testify/assert"
)

func TestNoteBox(t *testing.T) {
	writer := &MockWriter{}
	prompts.Note("test", prompts.NoteOptions{Output: writer})

	assert.Equal(t, strings.Join([]string{
		"│",
		"│────────╮",
		"│        │",
		"│  test  │",
		"│        │",
		"│────────╯",
		"",
	}, "\r\n"), writer.Data[0])
}

func TestNoteBoxMultiline(t *testing.T) {
	writer := &MockWriter{}
	prompts.Note("test\ntee\nfoooooo", prompts.NoteOptions{Output: writer})

	assert.Equal(t, strings.Join([]string{
		"│",
		"│───────────╮",
		"│           │",
		"│  test     │",
		"│  tee      │",
		"│  foooooo  │",
		"│           │",
		"│───────────╯",
		"",
	}, "\r\n"), writer.Data[0])
}

func TestNoteTitle(t *testing.T) {
	writer := &MockWriter{}
	prompts.Note("test", prompts.NoteOptions{Output: writer, Title: "Title Test"})

	assert.Equal(t, strings.Join([]string{
		"│",
		"│◇ Title Test ────╮",
		"│                 │",
		"│  test           │",
		"│                 │",
		"│─────────────────╯",
		"",
	}, "\r\n"), writer.Data[0])
}

func TestNoteTitleWithMultiline(t *testing.T) {
	writer := &MockWriter{}
	prompts.Note("test\ntee\nfoooooo", prompts.NoteOptions{Output: writer, Title: "Title Test"})

	assert.Equal(t, strings.Join([]string{
		"│",
		"│◇ Title Test ────╮",
		"│                 │",
		"│  test           │",
		"│  tee            │",
		"│  foooooo        │",
		"│                 │",
		"│─────────────────╯",
		"",
	}, "\r\n"), writer.Data[0])
}
