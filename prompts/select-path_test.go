package prompts_test

import (
	"testing"
	"time"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func runSelectPath() {
	prompts.SelectPath(prompts.SelectPathParams{
		Message:    message,
		FileSystem: (prompts.FileSystem)(MockFileSystem{}),
	})
}

func TestSelectPathInitialState(t *testing.T) {
	go runSelectPath()
	time.Sleep(time.Millisecond)
	p := test.SelectPathTestingPrompt

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectPathWithOptionChildren(t *testing.T) {
	go runSelectPath()
	time.Sleep(time.Millisecond)
	p := test.SelectPathTestingPrompt

	p.PressKey(&core.Key{Name: core.RightKey})

	assert.Equal(t, core.ActiveState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectPathCancelState(t *testing.T) {
	go runSelectPath()
	time.Sleep(time.Millisecond)

	p := test.SelectPathTestingPrompt
	p.PressKey(&core.Key{Name: core.CancelKey})

	assert.Equal(t, core.CancelState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectPathSubmitState(t *testing.T) {
	go runSelectPath()
	time.Sleep(time.Millisecond)

	p := test.SelectPathTestingPrompt
	p.PressKey(&core.Key{Name: core.EnterKey})

	assert.Equal(t, core.SubmitState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}
