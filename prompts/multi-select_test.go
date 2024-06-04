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

func runMultiSelect() {
	prompts.MultiSelect(prompts.MultiSelectParams[string]{
		Message: message,
		Options: []prompts.MultiSelectOption[string]{
			{Label: "a", IsSelected: true},
			{Label: "b", IsSelected: true},
			{Label: "c"},
		},
	})
}

func TestMultiSelectInitialState(t *testing.T) {
	go prompts.MultiSelect(prompts.MultiSelectParams[string]{
		Message: message,
		Options: []prompts.MultiSelectOption[string]{
			{Label: "a"},
			{Label: "b"},
			{Label: "c"},
		},
	})
	time.Sleep(1 * time.Millisecond)
	p := test.MultiSelectTestingPrompt.(*core.MultiSelectPrompt[string])

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestMultiSelectWithHint(t *testing.T) {
	go prompts.MultiSelect(prompts.MultiSelectParams[string]{
		Message: message,
		Options: []prompts.MultiSelectOption[string]{
			{Label: "a", Hint: "b"},
			{Label: "b", Hint: "c"},
			{Label: "c", Hint: "a"},
		},
	})
	time.Sleep(1 * time.Millisecond)
	p := test.MultiSelectTestingPrompt.(*core.MultiSelectPrompt[string])

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestMultiSelectCancelState(t *testing.T) {
	go runMultiSelect()
	time.Sleep(1 * time.Millisecond)

	p := test.MultiSelectTestingPrompt.(*core.MultiSelectPrompt[string])
	p.PressKey(&core.Key{Name: core.CancelKey})

	assert.Equal(t, core.CancelState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestMultiSelectSubmitState(t *testing.T) {
	go runMultiSelect()
	time.Sleep(1 * time.Millisecond)

	p := test.MultiSelectTestingPrompt.(*core.MultiSelectPrompt[string])
	p.PressKey(&core.Key{Name: core.EnterKey})

	assert.Equal(t, core.SubmitState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestMultiSelectWithLongList(t *testing.T) {
	go prompts.MultiSelect(prompts.MultiSelectParams[string]{
		Message: message,
		Options: []prompts.MultiSelectOption[string]{
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
	time.Sleep(1 * time.Millisecond)
	p := test.MultiSelectTestingPrompt.(*core.MultiSelectPrompt[string])

	assert.Equal(t, core.InitialState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}

func TestMultiSelectMultiValue(t *testing.T) {
	go prompts.MultiSelect(prompts.MultiSelectParams[string]{
		Message: message,
		Options: []prompts.MultiSelectOption[string]{
			{Label: "a", IsSelected: true},
			{Label: "b", IsSelected: true},
			{Label: "c"},
			{Label: "a", IsSelected: true},
			{Label: "b", IsSelected: true},
			{Label: "c"},
		},
	})
	time.Sleep(1 * time.Millisecond)
	p := test.MultiSelectTestingPrompt.(*core.MultiSelectPrompt[string])
	p.CursorIndex = 1
	p.PressKey(&core.Key{Name: core.DownKey})

	assert.Equal(t, core.ActiveState, p.State)
	cupaloy.SnapshotT(t, p.Frame)
}
