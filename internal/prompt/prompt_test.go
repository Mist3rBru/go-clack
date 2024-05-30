package prompt_test

import (
	"go-clack/internal/prompt"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newPrompt() *prompt.Prompt {
	return prompt.NewPrompt(prompt.PromptOptions{
		Input:  os.Stdin,
		Output: os.Stdout,
		Track:  true,
	})
}

func TestEmitEvent(t *testing.T) {
	p := newPrompt()
	event := "test"
	arg := rand.Int()

	p.On(event, func(args ...any) {
		assert.Equal(t, args[0], arg)
	})
	p.Emit(event, arg)
}

func TestEmitOtherEvent(t *testing.T) {
	p := newPrompt()
	event := "test"
	calledTimes := 0

	p.On(event, func(args ...any) {
		calledTimes++
	})
	p.Emit("other" + event)
	assert.Equal(t, 0, calledTimes)
}

func TestEmitEventWithMultiArgs(t *testing.T) {
	p := newPrompt()
	event := "test"
	args := []any{rand.Int(), rand.Int()}

	p.On(event, func(_args ...any) {
		assert.Equal(t, _args, args)
	})
	p.Emit(event, args...)
	p.Emit(event, args[0], args[1])
}

func TestEmitEventTwice(t *testing.T) {
	p := newPrompt()
	event := "test"
	calledTimes := 0

	p.On(event, func(args ...any) {
		calledTimes++
	})
	p.Emit(event)
	p.Emit(event)
	assert.Equal(t, 2, calledTimes)
}

func TestEmitEventOnce(t *testing.T) {
	p := newPrompt()
	event := "test"
	calledTimes := 0

	p.Once(event, func(args ...any) {
		calledTimes++
	})
	p.Emit(event)
	p.Emit(event)
	assert.Equal(t, 1, calledTimes)
}

func TestEmitUnsubscribedEvent(t *testing.T) {
	p := newPrompt()
	event := "test"
	calledTimes := 0
	listener := func(args ...any) {
		calledTimes++
	}

	p.On(event, listener)
	p.Off(event, listener)
	p.Emit(event)
	assert.Equal(t, 0, calledTimes)
}

func TestSetValue(t *testing.T) {
	p := newPrompt()
	p.Value = ""
	p.SetValue("test")
	assert.Equal(t, "test", p.Value)
}

func TestParseKey(t *testing.T) {
	p := newPrompt()

	key, char := p.ParseKey('\n')
	assert.Equal(t, "Enter", key)
	assert.Equal(t, "\n", char)

	key, char = p.ParseKey('a')
	assert.Equal(t, "a", key)
	assert.Equal(t, "a", char)

	key, char = p.ParseKey(3)
	assert.Equal(t, "Cancel", key)
	assert.Equal(t, "", char)

	key, char = p.ParseKey(27)
	assert.Equal(t, "", key)
	assert.Equal(t, "", char)
}
