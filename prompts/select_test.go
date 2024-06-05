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

func runSelect() {
	prompts.Select(prompts.SelectParams[string]{
		Message: message,
		Options: []prompts.SelectOption[string]{
			{Label: "a"},
			{Label: "b"},
			{Label: "c"},
		},
	})
}

func TestSelectInitialState(t *testing.T) {
	go runSelect()
	time.Sleep(time.Millisecond)
	p := test.SelectTestingPrompt.(*core.SelectPrompt[string])

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectWithHint(t *testing.T) {
	go prompts.Select(prompts.SelectParams[string]{
		Message: message,
		Options: []prompts.SelectOption[string]{
			{Label: "a", Hint: "b"},
			{Label: "b", Hint: "c"},
			{Label: "c", Hint: "a"},
		},
	})
	time.Sleep(time.Millisecond)
	p := test.SelectTestingPrompt.(*core.SelectPrompt[string])

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectCancelState(t *testing.T) {
	go runSelect()
	time.Sleep(time.Millisecond)

	p := test.SelectTestingPrompt.(*core.SelectPrompt[string])
	p.PressKey(&core.Key{Name: core.CancelKey})

	assert.Equal(t, core.CancelState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectSubmitState(t *testing.T) {
	go runSelect()
	time.Sleep(time.Millisecond)

	p := test.SelectTestingPrompt.(*core.SelectPrompt[string])
	p.PressKey(&core.Key{Name: core.EnterKey})

	assert.Equal(t, core.SubmitState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestSelectWithLongList(t *testing.T) {
	go prompts.Select(prompts.SelectParams[string]{
		Message: message,
		Options: []prompts.SelectOption[string]{
			{Label: "a"},
			{Label: "b"},
			{Label: "c"},
			{Label: "a"},
			{Label: "b"},
			{Label: "c"},
			{Label: "a"},
			{Label: "b"},
			{Label: "c"},
			{Label: "a"},
			{Label: "b"},
			{Label: "c"},
		},
	})
	time.Sleep(time.Millisecond)
	p := test.SelectTestingPrompt.(*core.SelectPrompt[string])

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}
