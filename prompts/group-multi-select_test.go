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

func runGroupMultiSelect() {
	options := make(map[string][]prompts.MultiSelectOption[string])
	options["1"] = []prompts.MultiSelectOption[string]{
		{Label: "a"},
		{Label: "b"},
		{Label: "c"},
	}
	options["2"] = []prompts.MultiSelectOption[string]{
		{Label: "x"},
		{Label: "y"},
		{Label: "z"},
	}
	prompts.GroupMultiSelect(prompts.GroupMultiSelectParams[string]{
		Message: message,
		Options: options,
	})
}

func TestGroupMultiSelectInitialState(t *testing.T) {
	go runGroupMultiSelect()
	time.Sleep(time.Millisecond)
	p := test.GroupMultiSelectTestingPrompt.(*core.GroupMultiSelectPrompt[string])

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestGroupMultiSelectCancelState(t *testing.T) {
	go runGroupMultiSelect()
	time.Sleep(time.Millisecond)

	p := test.GroupMultiSelectTestingPrompt.(*core.GroupMultiSelectPrompt[string])
	p.PressKey(&core.Key{Name: core.SpaceKey})
	p.PressKey(&core.Key{Name: core.CancelKey})

	assert.Equal(t, core.CancelState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestGroupMultiSelectSubmitState(t *testing.T) {
	go runGroupMultiSelect()
	time.Sleep(time.Millisecond)

	p := test.GroupMultiSelectTestingPrompt.(*core.GroupMultiSelectPrompt[string])
	p.PressKey(&core.Key{Name: core.SpaceKey})
	p.PressKey(&core.Key{Name: core.EnterKey})

	assert.Equal(t, core.SubmitState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestGroupMultiSelectWithLongList(t *testing.T) {
	options := make(map[string][]prompts.MultiSelectOption[string])
	options["1"] = []prompts.MultiSelectOption[string]{
		{Label: "a"},
		{Label: "b"},
		{Label: "c"},
		{Label: "d"},
		{Label: "e"},
		{Label: "f"},
	}
	options["2"] = []prompts.MultiSelectOption[string]{
		{Label: "u"},
		{Label: "v"},
		{Label: "w"},
		{Label: "x"},
		{Label: "y"},
		{Label: "z"},
	}
	go prompts.GroupMultiSelect(prompts.GroupMultiSelectParams[string]{
		Message: message,
		Options: options,
	})
	time.Sleep(time.Millisecond)
	p := test.GroupMultiSelectTestingPrompt.(*core.GroupMultiSelectPrompt[string])

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestGroupMultiSelectMultiValue(t *testing.T) {
	go runGroupMultiSelect()
	time.Sleep(time.Millisecond)

	p := test.GroupMultiSelectTestingPrompt.(*core.GroupMultiSelectPrompt[string])
	p.PressKey(&core.Key{Name: core.SpaceKey})
	p.CursorIndex = 5
	p.PressKey(&core.Key{Name: core.SpaceKey})
	p.PressKey(&core.Key{Name: core.DownKey})
	p.PressKey(&core.Key{Name: core.SpaceKey})

	assert.Equal(t, core.ActiveState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}
