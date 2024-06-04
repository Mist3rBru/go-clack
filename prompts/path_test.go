package prompts_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/utils"
	"github.com/stretchr/testify/assert"
)

func TestPathInitialState(t *testing.T) {
	go prompts.Path(prompts.PathParams{Message: message})
	time.Sleep(1 * time.Millisecond)

	p := test.PathTestingPrompt
	title := utils.SymbolState(core.InitialState) + " " + message
	valueWithCursor := utils.S_BAR + " " + p.ValueWithCursor()
	expected := strings.Join([]string{utils.S_BAR, title, valueWithCursor, utils.S_BAR_END}, "\r\n")
	assert.Equal(t, core.InitialState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestPathInitialStateWithInitialValue(t *testing.T) {
	go prompts.Path(prompts.PathParams{Message: message, InitialValue: "/foo"})
	time.Sleep(1 * time.Millisecond)

	p := test.PathTestingPrompt
	title := utils.SymbolState(core.InitialState) + " " + message
	valueWithCursor := utils.S_BAR + " /foo "
	expected := strings.Join([]string{utils.S_BAR, title, valueWithCursor, utils.S_BAR_END}, "\r\n")
	assert.Equal(t, core.InitialState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestPathErrorState(t *testing.T) {
	go prompts.Path(prompts.PathParams{Message: message, InitialValue: "/foo", Validate: func(value string) error {
		return fmt.Errorf("invalid value: %s", value)
	}})
	time.Sleep(1 * time.Millisecond)

	p := test.PathTestingPrompt
	p.PressKey(&core.Key{Name: core.EnterKey})

	title := utils.SymbolState(core.ErrorState) + " " + message
	valueWithCursor := utils.S_BAR + " /foo "
	err := utils.S_BAR_END + " invalid value: /foo"
	expected := strings.Join([]string{utils.S_BAR, title, valueWithCursor, err}, "\r\n")
	assert.Equal(t, core.ErrorState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestPathCancelState(t *testing.T) {
	go prompts.Path(prompts.PathParams{Message: message})
	time.Sleep(1 * time.Millisecond)

	p := test.PathTestingPrompt
	p.Value = ""
	p.PressKey(&core.Key{Name: core.CancelKey})

	title := utils.SymbolState(core.CancelState) + " " + message
	value := utils.S_BAR + " "
	expected := strings.Join([]string{utils.S_BAR, title, value}, "\r\n")
	assert.Equal(t, core.CancelState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestPathCancelStateWithValue(t *testing.T) {
	go prompts.Path(prompts.PathParams{Message: message})
	time.Sleep(1 * time.Millisecond)

	p := test.PathTestingPrompt
	p.PressKey(&core.Key{Name: core.CancelKey})

	title := utils.SymbolState(core.CancelState) + " " + message
	value := utils.S_BAR + " " + p.Value
	expected := strings.Join([]string{utils.S_BAR, title, value, utils.S_BAR}, "\r\n")
	assert.Equal(t, core.CancelState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestPathSubmitState(t *testing.T) {
	go prompts.Path(prompts.PathParams{Message: message})
	time.Sleep(1 * time.Millisecond)

	p := test.PathTestingPrompt
	p.PressKey(&core.Key{Name: core.EnterKey})

	title := utils.SymbolState(core.SubmitState) + " " + message
	value := utils.S_BAR + " " + p.Value
	expected := strings.Join([]string{utils.S_BAR, title, value}, "\r\n")
	assert.Equal(t, core.SubmitState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestPathValueWithOptions(t *testing.T) {
	go prompts.Path(prompts.PathParams{Message: message})
	time.Sleep(1 * time.Millisecond)

	p := test.PathTestingPrompt
	p.PressKey(&core.Key{Name: core.TabKey})
	p.PressKey(&core.Key{Name: core.TabKey})

	title := utils.SymbolState(core.ActiveState) + " " + message
	value := utils.S_BAR + " " + p.ValueWithCursor()
	options := p.FormatLines([]string{strings.Join(p.HintOptions, " ")}, core.FormatLinesOptions{
		Default: core.FormatLineOptions{
			Start: utils.S_BAR,
		},
	})
	expected := strings.Join([]string{utils.S_BAR, title, value, options, utils.S_BAR_END}, "\r\n")
	assert.Equal(t, core.ActiveState, p.State)
	assert.Equal(t, expected, p.Frame)
}
