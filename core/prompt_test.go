package core_test

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func newPrompt() *core.Prompt {
	return core.NewPrompt(core.PromptParams{
		Input:  os.Stdin,
		Output: os.Stdout,
		Track:  true,
		Value:  "",
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

	assert.Equal(t, core.Key{Name: "Enter"}, *p.ParseKey('\n'))
	assert.Equal(t, core.Key{Name: "Enter"}, *p.ParseKey('\r'))
	assert.Equal(t, core.Key{Name: "Space"}, *p.ParseKey(' '))
	assert.Equal(t, core.Key{Name: "Backspace"}, *p.ParseKey('\b'))
	assert.Equal(t, core.Key{Name: "Backspace"}, *p.ParseKey(127))
	assert.Equal(t, core.Key{Name: "Tab"}, *p.ParseKey('\t'))
	assert.Equal(t, core.Key{Name: "Cancel"}, *p.ParseKey(3))
	assert.Equal(t, core.Key{Name: "a", Char: "a"}, *p.ParseKey('a'))
}

func TestTrackValue(t *testing.T) {
	p := newPrompt()

	assert.Equal(t, "", p.Value)
	assert.Equal(t, 0, p.CursorIndex)

	p.PressKey(&core.Key{Char: "a"})
	assert.Equal(t, "a", p.Value)
	assert.Equal(t, 1, p.CursorIndex)

	p.PressKey(&core.Key{Char: "b"})
	assert.Equal(t, "ab", p.Value)
	assert.Equal(t, 2, p.CursorIndex)

	p.CursorIndex = 1
	p.PressKey(&core.Key{Char: "c"})
	assert.Equal(t, "acb", p.Value)
	assert.Equal(t, 2, p.CursorIndex)
}

func TestTrackCursor(t *testing.T) {
	p := newPrompt()

	p.Value = "abc"
	p.CursorIndex = 3
	p.PressKey(&core.Key{Name: "Home"})
	assert.Equal(t, 0, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 0
	p.PressKey(&core.Key{Name: "End"})
	assert.Equal(t, 3, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 3
	p.PressKey(&core.Key{Name: "Left"})
	assert.Equal(t, 2, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 0
	p.PressKey(&core.Key{Name: "Left"})
	assert.Equal(t, 0, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 2
	p.PressKey(&core.Key{Name: "Right"})
	assert.Equal(t, 3, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 3
	p.PressKey(&core.Key{Name: "Right"})
	assert.Equal(t, 3, p.CursorIndex)
}

func TestTrackBackspace(t *testing.T) {
	p := newPrompt()

	p.Value = "abc"
	p.CursorIndex = 3
	p.PressKey(&core.Key{Name: "Backspace"})
	assert.Equal(t, "ab", p.Value)
	assert.Equal(t, 2, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 2
	p.PressKey(&core.Key{Name: "Backspace"})
	assert.Equal(t, "ac", p.Value)
	assert.Equal(t, 1, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 1
	p.PressKey(&core.Key{Name: "Backspace"})
	assert.Equal(t, "bc", p.Value)
	assert.Equal(t, 0, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 0
	p.PressKey(&core.Key{Name: "Backspace"})
	assert.Equal(t, "abc", p.Value)
	assert.Equal(t, 0, p.CursorIndex)
}

func TestTrackState(t *testing.T) {
	p := newPrompt()

	p.PressKey(&core.Key{Name: "Cancel"})
	assert.Equal(t, "cancel", p.State)

	p.PressKey(&core.Key{Name: "Enter"})
	assert.Equal(t, "submit", p.State)
}

func TestLimitLines(t *testing.T) {
	p := newPrompt()
	lines := make([]string, 10)
	for i := range lines {
		lines[i] = fmt.Sprint(i)
	}

	frame := p.LimitLines(core.LimitLinesPamams{
		CursorIndex: 0,
		Lines:       lines,
	})
	startLines := lines[0:5]
	startLines[len(startLines)-1] = "..."
	expected := strings.Join(startLines, "\r\n")
	assert.Equal(t, expected, frame)

	frame = p.LimitLines(core.LimitLinesPamams{
		CursorIndex: 5,
		Lines:       lines,
	})
	midLines := lines[3:8]
	midLines[0] = "..."
	midLines[len(midLines)-1] = "..."
	expected = strings.Join(midLines, "\r\n")
	assert.Equal(t, expected, frame)

	frame = p.LimitLines(core.LimitLinesPamams{
		CursorIndex: 9,
		Lines:       lines,
	})
	lasLines := lines[5:10]
	lasLines[0] = "..."
	expected = strings.Join(lasLines, "\r\n")
	assert.Equal(t, expected, frame)
}
