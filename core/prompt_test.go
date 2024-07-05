package core_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
	"github.com/stretchr/testify/assert"
)

func newPrompt() *core.Prompt[string] {
	return core.NewPrompt(core.PromptParams[string]{
		Render: func(p *core.Prompt[string]) string { return "" },
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

	assert.Equal(t, core.Key{Name: core.EnterKey}, *p.ParseKey('\n'))
	assert.Equal(t, core.Key{Name: core.EnterKey}, *p.ParseKey('\r'))
	assert.Equal(t, core.Key{Name: core.SpaceKey}, *p.ParseKey(' '))
	assert.Equal(t, core.Key{Name: core.BackspaceKey}, *p.ParseKey('\b'))
	assert.Equal(t, core.Key{Name: core.BackspaceKey}, *p.ParseKey(127))
	assert.Equal(t, core.Key{Name: core.TabKey}, *p.ParseKey('\t'))
	assert.Equal(t, core.Key{Name: core.CancelKey}, *p.ParseKey(3))
	assert.Equal(t, core.Key{Name: "a", Char: "a"}, *p.ParseKey('a'))
}

func TestTrackValue(t *testing.T) {
	p := newPrompt()

	assert.Equal(t, "", p.Value)
	assert.Equal(t, 0, p.CursorIndex)

	p.TrackKeyValue(&core.Key{Char: "a"}, &p.Value)
	assert.Equal(t, "a", p.Value)
	assert.Equal(t, 1, p.CursorIndex)

	p.TrackKeyValue(&core.Key{Char: "b"}, &p.Value)
	assert.Equal(t, "ab", p.Value)
	assert.Equal(t, 2, p.CursorIndex)

	p.CursorIndex = 1
	p.TrackKeyValue(&core.Key{Char: "c"}, &p.Value)
	assert.Equal(t, "acb", p.Value)
	assert.Equal(t, 2, p.CursorIndex)
}

func TestTrackCursor(t *testing.T) {
	p := newPrompt()

	p.Value = "abc"
	p.CursorIndex = 3
	p.TrackKeyValue(&core.Key{Name: core.HomeKey}, &p.Value)
	assert.Equal(t, 0, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 0
	p.TrackKeyValue(&core.Key{Name: core.EndKey}, &p.Value)
	assert.Equal(t, 3, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 3
	p.TrackKeyValue(&core.Key{Name: core.LeftKey}, &p.Value)
	assert.Equal(t, 2, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 0
	p.TrackKeyValue(&core.Key{Name: core.LeftKey}, &p.Value)
	assert.Equal(t, 0, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 2
	p.TrackKeyValue(&core.Key{Name: core.RightKey}, &p.Value)
	assert.Equal(t, 3, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 3
	p.TrackKeyValue(&core.Key{Name: core.RightKey}, &p.Value)
	assert.Equal(t, 3, p.CursorIndex)
}

func TestTrackBackspace(t *testing.T) {
	p := newPrompt()

	p.Value = "abc"
	p.CursorIndex = 3
	p.TrackKeyValue(&core.Key{Name: core.BackspaceKey}, &p.Value)
	assert.Equal(t, "ab", p.Value)
	assert.Equal(t, 2, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 2
	p.TrackKeyValue(&core.Key{Name: core.BackspaceKey}, &p.Value)
	assert.Equal(t, "ac", p.Value)
	assert.Equal(t, 1, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 1
	p.TrackKeyValue(&core.Key{Name: core.BackspaceKey}, &p.Value)
	assert.Equal(t, "bc", p.Value)
	assert.Equal(t, 0, p.CursorIndex)

	p.Value = "abc"
	p.CursorIndex = 0
	p.TrackKeyValue(&core.Key{Name: core.BackspaceKey}, &p.Value)
	assert.Equal(t, "abc", p.Value)
	assert.Equal(t, 0, p.CursorIndex)
}

func TestTrackState(t *testing.T) {
	p := newPrompt()

	p.PressKey(&core.Key{Name: core.CancelKey})
	assert.Equal(t, core.CancelState, p.State)

	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, core.SubmitState, p.State)
}

func TestDiffLines(t *testing.T) {
	p := newPrompt()

	assert.Equal(t, []int{1, 2}, p.DiffLines("a", "a\nb\nc"))
	assert.Equal(t, []int{1, 2}, p.DiffLines("a\nb\nc", "a"))
	assert.Equal(t, []int{1, 2}, p.DiffLines("a\nb\nc", "a\nc\nb"))
	assert.Equal(t, []int{}, p.DiffLines("a\nb\nc", "a\nb\nc"))
}

func TestLimitLines(t *testing.T) {
	p := newPrompt()
	lines := make([]string, 20)
	for i := range lines {
		lines[i] = fmt.Sprint(i)
	}

	p.CursorIndex = 0
	frame := p.LimitLines(lines, 0)
	startLines := lines[0:10]
	startLines[len(startLines)-1] = "..."
	expected := strings.Join(startLines, "\n")
	assert.Equal(t, expected, frame)

	p.CursorIndex = 10
	frame = p.LimitLines(lines, 0)
	midLines := lines[3:13]
	midLines[0] = "..."
	midLines[len(midLines)-1] = "..."
	expected = strings.Join(midLines, "\n")
	assert.Equal(t, expected, frame)

	p.CursorIndex = 20
	frame = p.LimitLines(lines, 0)
	lasLines := lines[10:20]
	lasLines[0] = "..."
	expected = strings.Join(lasLines, "\n")
	assert.Equal(t, expected, frame)

	p.CursorIndex = 0
	frame = p.LimitLines(lines, 3)
	startLines = lines[0:7]
	startLines[len(startLines)-1] = "..."
	expected = strings.Join(startLines, "\n")
	assert.Equal(t, expected, frame)
}

func TestFormatLines(t *testing.T) {
	p := newPrompt()

	lines := []string{"a", "b"}
	frame := p.FormatLines(lines, core.FormatLinesOptions{
		FirstLine: core.FormatLineOptions{
			Sides: "-",
		},
		NewLine: core.FormatLineOptions{
			Sides: "*",
		},
	})
	expected := strings.Join([]string{
		fmt.Sprintf("- %s -", "a"),
		fmt.Sprintf("* %s *", "b"),
	}, "\r\n")
	assert.Equal(t, expected, frame)
}

func TestFormatLinesWithBlankSpaces(t *testing.T) {
	p := newPrompt()

	lines := []string{"a", " b", "  c", "x", " y", "  z"}
	frame := p.FormatLines(lines, core.FormatLinesOptions{
		Default: core.FormatLineOptions{
			Start: "|",
		},
	})
	expected := strings.Join([]string{
		"| a",
		"|  b",
		"|   c",
		"| x",
		"|  y",
		"|   z",
	}, "\r\n")
	assert.Equal(t, expected, frame)
}

func TestFormatLinesWithOverflowedWords(t *testing.T) {
	p := newPrompt()

	lines := []string{strings.Repeat("a", 81), strings.Repeat("b", 76), strings.Repeat("c", 79)}
	frame := p.FormatLines(lines, core.FormatLinesOptions{
		FirstLine: core.FormatLineOptions{
			Sides: "-",
		},
		Default: core.FormatLineOptions{
			Sides: "|",
		},
		LastLine: core.FormatLineOptions{
			Sides: "*",
		},
	})
	expected := strings.Join([]string{
		fmt.Sprintf("- %s -", strings.Repeat("a", 76)),
		fmt.Sprintf("- %s -", strings.Repeat("a", 5)),
		fmt.Sprintf("| %s |", strings.Repeat("b", 76)),
		fmt.Sprintf("* %s *", strings.Repeat("c", 76)),
		fmt.Sprintf("* %s *", strings.Repeat("c", 3)),
	}, "\r\n")
	assert.Equal(t, expected, frame)

	lines = []string{strings.Repeat("a", 180), "b"}
	frame = p.FormatLines(lines, core.FormatLinesOptions{
		FirstLine: core.FormatLineOptions{
			Sides: "-",
		},
		Default: core.FormatLineOptions{
			Sides: "*",
		},
		LastLine: core.FormatLineOptions{
			Sides: "-",
		},
	})
	expected = strings.Join([]string{
		fmt.Sprintf("- %s -", strings.Repeat("a", 76)),
		fmt.Sprintf("- %s -", strings.Repeat("a", 76)),
		fmt.Sprintf("- %s -", strings.Repeat("a", 28)),
		fmt.Sprintf("- %s -", "b"),
	}, "\r\n")
	assert.Equal(t, expected, frame)
}

func TestFormatLinesWithBoxFormat(t *testing.T) {
	p := newPrompt()

	lines := []string{strings.Repeat("a", 67), "b"}
	width := 70
	frame := p.FormatLines(lines, core.FormatLinesOptions{
		Default: core.FormatLineOptions{
			Sides: "|",
		},
		MinWidth: &width,
		MaxWidth: &width,
	})
	expected := strings.Join([]string{
		fmt.Sprintf("| %s%s |", strings.Repeat("a", 66), strings.Repeat(" ", 0)),
		fmt.Sprintf("| %s%s |", strings.Repeat("a", 1), strings.Repeat(" ", 65)),
		fmt.Sprintf("| %s%s |", strings.Repeat("b", 1), strings.Repeat(" ", 65)),
	}, "\r\n")
	for i, line := range strings.Split(frame, "\r\n") {
		assert.Equal(t, width, len(line), fmt.Sprintf("index: %d\nline: %s", i, line))
	}
	assert.Equal(t, expected, frame)
}

func TestFormatLinesWithComplexParams(t *testing.T) {
	p := newPrompt()

	lines := []string{strings.Repeat("a", 20) + " " + strings.Repeat("b", 80), "c"}
	frame := p.FormatLines(lines, core.FormatLinesOptions{
		FirstLine: core.FormatLineOptions{
			Start: "-",
		},
		NewLine: core.FormatLineOptions{
			End: "-",
		},
		LastLine: core.FormatLineOptions{
			Start: "=",
			End:   "=",
		},
		Default: core.FormatLineOptions{
			Sides: "*",
		},
	})
	expected := strings.Join([]string{
		fmt.Sprintf("- %s %s *", strings.Repeat("a", 20), strings.Repeat("b", 55)),
		fmt.Sprintf("- %s *", strings.Repeat("b", 25)),
		fmt.Sprintf("= %s =", strings.Repeat("c", 1)),
	}, "\r\n")
	assert.Equal(t, expected, frame)
}

func TestFormatLinesWithStyleCallback(t *testing.T) {
	p := newPrompt()

	lines := []string{strings.Repeat("a", 20) + " " + strings.Repeat("b", 80), "c"}
	frame := p.FormatLines(lines, core.FormatLinesOptions{
		Default: core.FormatLineOptions{
			Style: func(line string) string {
				return fmt.Sprintf("(%s)", line)
			},
		},
	})
	expected := strings.Join([]string{
		fmt.Sprintf("(%s %s)", strings.Repeat("a", 20), strings.Repeat("b", 57)),
		fmt.Sprintf("(%s)", strings.Repeat("b", 23)),
		fmt.Sprintf("(%s)", strings.Repeat("c", 1)),
	}, "\r\n")
	assert.Equal(t, expected, frame)
}

func TestFormatLinesWithBlackLine(t *testing.T) {
	p := newPrompt()

	lines := []string{picocolors.Inverse(" ")}
	width := 80
	frame := p.FormatLines(lines, core.FormatLinesOptions{
		Default: core.FormatLineOptions{
			Start: "|",
		},
		LastLine: core.FormatLineOptions{
			End: "|",
		},
		MinWidth: &width,
	})
	expected := fmt.Sprintf("| %s |", strings.Repeat(" ", 76))
	assert.Equal(t, expected, frame)
}

func TestValidateValue(t *testing.T) {
	p := newPrompt()
	p.Validate = func(value string) error {
		return fmt.Errorf("invalid value: %v", value)
	}

	p.Value = "foo"
	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, core.ErrorState, p.State)
	assert.Equal(t, "invalid value: foo", p.Error)
}

func TestEmitFinalizeOnSubmit(t *testing.T) {
	p := newPrompt()
	calledTimes := 0
	p.On(core.FinalizeEvent, func(args ...any) {
		calledTimes++
	})

	p.PressKey(&core.Key{Name: core.EnterKey})
	assert.Equal(t, 1, calledTimes)
}
