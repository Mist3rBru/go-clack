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

const testEvent = core.Event("test")

func TestEmitEvent(t *testing.T) {
	p := newPrompt()
	arg := rand.Int()

	p.On(testEvent, func(args ...any) {
		assert.Equal(t, args[0], arg)
	})
	p.Emit(testEvent, arg)
}

func TestEmitOtherEvent(t *testing.T) {
	p := newPrompt()
	calledTimes := 0

	p.On(testEvent, func(args ...any) {
		calledTimes++
	})
	p.Emit(core.Event("other") + testEvent)
	assert.Equal(t, 0, calledTimes)
}

func TestEmitEventWithMultiArgs(t *testing.T) {
	p := newPrompt()
	args := []any{rand.Int(), rand.Int()}

	p.On(testEvent, func(_args ...any) {
		assert.Equal(t, _args, args)
	})
	p.Emit(testEvent, args...)
	p.Emit(testEvent, args[0], args[1])
}

func TestEmitEventTwice(t *testing.T) {
	p := newPrompt()
	calledTimes := 0

	p.On(testEvent, func(args ...any) {
		calledTimes++
	})
	p.Emit(testEvent)
	p.Emit(testEvent)
	assert.Equal(t, 2, calledTimes)
}

func TestEmitEventOnce(t *testing.T) {
	p := newPrompt()
	calledTimes := 0

	p.Once(testEvent, func(args ...any) {
		calledTimes++
	})
	p.Emit(testEvent)
	p.Emit(testEvent)
	assert.Equal(t, 1, calledTimes)
}

func TestEmitUnsubscribedEvent(t *testing.T) {
	p := newPrompt()
	calledTimes := 0
	listener := func(args ...any) {
		calledTimes++
	}

	p.On(testEvent, listener)
	p.Off(testEvent, listener)
	p.Emit(testEvent)
	assert.Equal(t, 0, calledTimes)
}

func TestParseKey(t *testing.T) {
	p := newPrompt()

	assert.Equal(t, core.Key{Name: core.KeyEnter}, *p.ParseKey('\n'))
	assert.Equal(t, core.Key{Name: core.KeyEnter}, *p.ParseKey('\r'))
	assert.Equal(t, core.Key{Name: core.KeySpace}, *p.ParseKey(' '))
	assert.Equal(t, core.Key{Name: core.KeyBackspace}, *p.ParseKey('\b'))
	assert.Equal(t, core.Key{Name: core.KeyBackspace}, *p.ParseKey(127))
	assert.Equal(t, core.Key{Name: core.KeyTab}, *p.ParseKey('\t'))
	assert.Equal(t, core.Key{Name: core.KeyCancel}, *p.ParseKey(3))
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
	p.TrackKeyValue(&core.Key{Name: core.KeyHome}, p.Value)
	assert.Equal(t, 0, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 0
	p.TrackKeyValue(&core.Key{Name: core.KeyEnd}, p.Value)
	assert.Equal(t, 3, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 3
	p.TrackKeyValue(&core.Key{Name: core.KeyLeft}, p.Value)
	assert.Equal(t, 2, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 0
	p.TrackKeyValue(&core.Key{Name: core.KeyLeft}, p.Value)
	assert.Equal(t, 0, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 2
	p.TrackKeyValue(&core.Key{Name: core.KeyRight}, p.Value)
	assert.Equal(t, 3, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 3
	p.TrackKeyValue(&core.Key{Name: core.KeyRight}, p.Value)
	assert.Equal(t, 3, p.CursorIndex)
}

func TestTrackBackspace(t *testing.T) {
	p := newPrompt()

	p.Value = "abc"
	p.CursorIndex = 3
	p.Value = p.TrackKeyValue(&core.Key{Name: core.KeyBackspace}, p.Value)
	assert.Equal(t, "ab", p.Value)
	assert.Equal(t, 2, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 2
	p.Value = p.TrackKeyValue(&core.Key{Name: core.KeyBackspace}, p.Value)
	assert.Equal(t, "ac", p.Value)
	assert.Equal(t, 1, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 1
	p.Value = p.TrackKeyValue(&core.Key{Name: core.KeyBackspace}, p.Value)
	assert.Equal(t, "bc", p.Value)
	assert.Equal(t, 0, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 0
	p.Value = p.TrackKeyValue(&core.Key{Name: core.KeyBackspace}, p.Value)
	assert.Equal(t, "abc", p.Value)
	assert.Equal(t, 0, p.CursorIndex)
}

func TestTrackState(t *testing.T) {
	p := newPrompt()

	p.PressKey(&core.Key{Name: core.KeyCancel})
	assert.Equal(t, core.StateCancel, p.State)

	p.PressKey(&core.Key{Name: core.KeyEnter})
	assert.Equal(t, core.StateSubmit, p.State)
}

func TestLimitLines(t *testing.T) {
	p := newPrompt()
	lines := make([]string, 10)
	for i := range lines {
		lines[i] = fmt.Sprint(i)
	}

	p.CursorIndex = 0
	frame := p.LimitLines(lines, 0)
	startLines := lines[0:5]
	startLines[len(startLines)-1] = "..."
	expected := strings.Join(startLines, "\r\n")
	assert.Equal(t, expected, frame)

	p.CursorIndex = 5
	frame = p.LimitLines(lines, 0)
	midLines := lines[3:8]
	midLines[0] = "..."
	midLines[len(midLines)-1] = "..."
	expected = strings.Join(midLines, "\r\n")
	assert.Equal(t, expected, frame)

	p.CursorIndex = 9
	frame = p.LimitLines(lines, 0)
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
	p.PressKey(&core.Key{Name: core.KeyEnter})
	assert.Equal(t, core.StateError, p.State)
	assert.Equal(t, "invalid value: foo", p.Error)
}

func TestEmitFinalizeOnSubmit(t *testing.T) {
	p := newPrompt()
	calledTimes := 0
	p.On(core.EventFinalize, func(args ...any) {
		calledTimes++
	})

	p.PressKey(&core.Key{Name: core.KeyEnter})
	assert.Equal(t, 1, calledTimes)
}
