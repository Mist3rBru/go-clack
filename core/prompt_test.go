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

func newPrompt() *core.Prompt[string] {
	return core.NewPrompt(core.PromptParams[string]{
		Input:  os.Stdin,
		Output: os.Stdout,
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

	p.Value = p.TrackKeyValue(&core.Key{Char: "a"}, p.Value)
	assert.Equal(t, "a", p.Value)
	assert.Equal(t, 1, p.CursorIndex)

	p.Value = p.TrackKeyValue(&core.Key{Char: "b"}, p.Value)
	assert.Equal(t, "ab", p.Value)
	assert.Equal(t, 2, p.CursorIndex)

	p.CursorIndex = 1
	p.Value = p.TrackKeyValue(&core.Key{Char: "c"}, p.Value)
	assert.Equal(t, "acb", p.Value)
	assert.Equal(t, 2, p.CursorIndex)
}

func TestTrackCursor(t *testing.T) {
	p := newPrompt()

	p.Value = "abc"
	p.CursorIndex = 3
	p.TrackKeyValue(&core.Key{Name: "Home"}, p.Value)
	assert.Equal(t, 0, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 0
	p.TrackKeyValue(&core.Key{Name: "End"}, p.Value)
	assert.Equal(t, 3, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 3
	p.TrackKeyValue(&core.Key{Name: "Left"}, p.Value)
	assert.Equal(t, 2, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 0
	p.TrackKeyValue(&core.Key{Name: "Left"}, p.Value)
	assert.Equal(t, 0, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 2
	p.TrackKeyValue(&core.Key{Name: "Right"}, p.Value)
	assert.Equal(t, 3, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 3
	p.TrackKeyValue(&core.Key{Name: "Right"}, p.Value)
	assert.Equal(t, 3, p.CursorIndex)
}

func TestTrackBackspace(t *testing.T) {
	p := newPrompt()

	p.Value = "abc"
	p.CursorIndex = 3
	p.Value = p.TrackKeyValue(&core.Key{Name: "Backspace"}, p.Value)
	assert.Equal(t, "ab", p.Value)
	assert.Equal(t, 2, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 2
	p.Value = p.TrackKeyValue(&core.Key{Name: "Backspace"}, p.Value)
	assert.Equal(t, "ac", p.Value)
	assert.Equal(t, 1, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 1
	p.Value = p.TrackKeyValue(&core.Key{Name: "Backspace"}, p.Value)
	assert.Equal(t, "bc", p.Value)
	assert.Equal(t, 0, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 0
	p.Value = p.TrackKeyValue(&core.Key{Name: "Backspace"}, p.Value)
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

func TestValidateValue(t *testing.T) {
	p := newPrompt()
	p.Validate = func(value string) error {
		return fmt.Errorf("invalid value: %v", value)
	}

	p.Value = "foo"
	p.PressKey(&core.Key{Name: "Enter"})
	assert.Equal(t, "error", p.State)
	assert.Equal(t, "invalid value: foo", p.Error)
}

func TestEmitFinalizeOnSubmit(t *testing.T) {
	p := newPrompt()
	calledTimes := 0
	p.On("finalize", func(args ...any) {
		calledTimes++
	})

	p.PressKey(&core.Key{Name: "Enter"})
	assert.Equal(t, 1, calledTimes)
}
