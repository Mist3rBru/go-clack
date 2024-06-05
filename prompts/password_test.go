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

func TestPasswordInitialState(t *testing.T) {
	go prompts.Password(prompts.PasswordParams{Message: message})
	time.Sleep(time.Millisecond)

	p := test.PasswordTestingPrompt
	title := utils.SymbolState(core.InitialState) + " " + message
	valueWithCursor := utils.S_BAR + " "
	expected := strings.Join([]string{utils.S_BAR, title, valueWithCursor, utils.S_BAR_END}, "\r\n")
	assert.Equal(t, core.InitialState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestPasswordInitialStateWithInitialValue(t *testing.T) {
	go prompts.Password(prompts.PasswordParams{Message: message, InitialValue: "foo"})
	time.Sleep(time.Millisecond)

	p := test.PasswordTestingPrompt
	title := utils.SymbolState(core.InitialState) + " " + message
	valueWithCursor := utils.S_BAR + " *** "
	expected := strings.Join([]string{utils.S_BAR, title, valueWithCursor, utils.S_BAR_END}, "\r\n")
	assert.Equal(t, core.InitialState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestPasswordErrorState(t *testing.T) {
	go prompts.Password(prompts.PasswordParams{Message: message, InitialValue: "foo", Validate: func(value string) error {
		return fmt.Errorf("invalid value: %s", value)
	}})
	time.Sleep(time.Millisecond)

	p := test.PasswordTestingPrompt
	p.PressKey(&core.Key{Name: core.EnterKey})

	title := utils.SymbolState(core.ErrorState) + " " + message
	valueWithCursor := utils.S_BAR + " *** "
	err := utils.S_BAR_END + " invalid value: foo"
	expected := strings.Join([]string{utils.S_BAR, title, valueWithCursor, err}, "\r\n")
	assert.Equal(t, core.ErrorState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestPasswordCancelState(t *testing.T) {
	go prompts.Password(prompts.PasswordParams{Message: message})
	time.Sleep(time.Millisecond)

	p := test.PasswordTestingPrompt
	p.PressKey(&core.Key{Name: core.CancelKey})

	title := utils.SymbolState(core.CancelState) + " " + message
	value := utils.S_BAR + " "
	expected := strings.Join([]string{utils.S_BAR, title, value}, "\r\n")
	assert.Equal(t, core.CancelState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestPasswordCancelStateWithValue(t *testing.T) {
	go prompts.Password(prompts.PasswordParams{Message: message, InitialValue: "foo"})
	time.Sleep(time.Millisecond)

	p := test.PasswordTestingPrompt
	p.PressKey(&core.Key{Name: core.CancelKey})

	title := utils.SymbolState(core.CancelState) + " " + message
	value := utils.S_BAR + " ***"
	expected := strings.Join([]string{utils.S_BAR, title, value, utils.S_BAR}, "\r\n")
	assert.Equal(t, core.CancelState, p.State)
	assert.Equal(t, expected, p.Frame)
}

func TestPasswordSubmitState(t *testing.T) {
	go prompts.Password(prompts.PasswordParams{Message: message, InitialValue: "foo"})
	time.Sleep(time.Millisecond)

	p := test.PasswordTestingPrompt
	p.PressKey(&core.Key{Name: core.EnterKey})

	title := utils.SymbolState(core.SubmitState) + " " + message
	value := utils.S_BAR + " ***"
	expected := strings.Join([]string{utils.S_BAR, title, value}, "\r\n")
	assert.Equal(t, core.SubmitState, p.State)
	assert.Equal(t, expected, p.Frame)
}
