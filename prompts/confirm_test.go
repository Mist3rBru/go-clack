package prompts_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/utils"
	"github.com/stretchr/testify/assert"
)

func TestConfirmInitialState(t *testing.T) {
	go prompts.Confirm(prompts.ConfirmParams{Message: message})
	time.Sleep(time.Millisecond)

	p := test.ConfirmTestingPrompt
	title := utils.SymbolState(core.InitialState) + " " + message
	valueWithCursor := strings.Join([]string{utils.S_BAR, utils.S_RADIO_INACTIVE, p.Active, "/", utils.S_RADIO_ACTIVE, p.Inactive}, " ")
	expected := strings.Join([]string{utils.S_BAR, title, valueWithCursor, utils.S_BAR_END}, "\r\n")
	assert.Equal(t, core.InitialState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestConfirmInitialStateWithInitialValue(t *testing.T) {
	go prompts.Confirm(prompts.ConfirmParams{Message: message, InitialValue: true})
	time.Sleep(time.Millisecond)

	p := test.ConfirmTestingPrompt
	title := utils.SymbolState(core.InitialState) + " " + message
	valueWithCursor := strings.Join([]string{utils.S_BAR, utils.S_RADIO_ACTIVE, p.Active, "/", utils.S_RADIO_INACTIVE, p.Inactive}, " ")
	expected := strings.Join([]string{utils.S_BAR, title, valueWithCursor, utils.S_BAR_END}, "\r\n")
	assert.Equal(t, core.InitialState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestConfirmCancelState(t *testing.T) {
	go prompts.Confirm(prompts.ConfirmParams{Message: message})
	time.Sleep(time.Millisecond)

	p := test.ConfirmTestingPrompt
	p.PressKey(&core.Key{Name: core.CancelKey})

	title := utils.SymbolState(core.CancelState) + " " + message
	value := utils.S_BAR + " " + p.Inactive
	expected := strings.Join([]string{utils.S_BAR, title, value, utils.S_BAR}, "\r\n")
	assert.Equal(t, core.CancelState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestConfirmCancelStateWithValue(t *testing.T) {
	go prompts.Confirm(prompts.ConfirmParams{Message: message, InitialValue: true})
	time.Sleep(time.Millisecond)

	p := test.ConfirmTestingPrompt
	p.PressKey(&core.Key{Name: core.CancelKey})

	title := utils.SymbolState(core.CancelState) + " " + message
	value := utils.S_BAR + " " + p.Active
	expected := strings.Join([]string{utils.S_BAR, title, value, utils.S_BAR}, "\r\n")
	assert.Equal(t, core.CancelState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestConfirmSubmitState(t *testing.T) {
	go prompts.Confirm(prompts.ConfirmParams{Message: message})
	time.Sleep(time.Millisecond)

	p := test.ConfirmTestingPrompt
	p.PressKey(&core.Key{Name: core.EnterKey})

	title := utils.SymbolState(core.SubmitState) + " " + message
	value := utils.S_BAR + " " + p.Inactive
	expected := strings.Join([]string{utils.S_BAR, title, value}, "\r\n")
	assert.Equal(t, core.SubmitState, p.State)
	assert.Equal(t, expected, p.Frame)
}
