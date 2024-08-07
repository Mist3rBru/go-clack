package core

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/core/validator"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
	"github.com/Mist3rBru/go-clack/third_party/sisteransi"

	"golang.org/x/term"
)

type Listener func(args ...any)

type Prompt[TValue any] struct {
	listeners map[Event][]Listener

	rl     *bufio.Reader
	input  *os.File
	output *os.File

	State       State
	Error       string
	Value       TValue
	CursorIndex int

	Validate           func(value TValue) error
	ValidationDuration time.Duration
	IsValidating       bool

	Render func(p *Prompt[TValue]) string
	Frame  string
}

type PromptParams[TValue any] struct {
	Input        *os.File
	Output       *os.File
	InitialValue TValue
	CursorIndex  int
	Validate     func(value TValue) error
	Render       func(p *Prompt[TValue]) string
}

func NewPrompt[TValue any](params PromptParams[TValue]) *Prompt[TValue] {
	v := validator.NewValidator("Prompt")
	v.ValidateRender(params.Render)

	if params.Input == nil {
		params.Input = os.Stdin
	}
	if params.Output == nil {
		params.Output = os.Stdout
	}

	return &Prompt[TValue]{
		listeners: make(map[Event][]Listener),

		input:  params.Input,
		output: params.Output,
		rl:     bufio.NewReader(params.Input),

		State:       InitialState,
		Value:       params.InitialValue,
		CursorIndex: params.CursorIndex,

		Validate: params.Validate,
		Render:   params.Render,
	}
}

// On registers a listener for the specified event.
func (p *Prompt[TValue]) On(event Event, listener Listener) {
	p.listeners[event] = append(p.listeners[event], listener)
}

// Once registers a one-time listener for the specified event.
func (p *Prompt[TValue]) Once(event Event, listener Listener) {
	var onceListener Listener
	onceListener = func(args ...any) {
		listener(args)
		p.Off(event, onceListener)
	}
	p.On(event, onceListener)
}

// Off removes a listener for the specified event.
func (p *Prompt[TValue]) Off(event Event, listener Listener) {
	listeners := p.listeners[event]
	for i, l := range listeners {
		if fmt.Sprintf("%p", l) == fmt.Sprintf("%p", listener) {
			p.listeners[event] = append(listeners[:i], listeners[i+1:]...)
			break
		}
	}
}

// Emit triggers the specified event with the given arguments.
func (p *Prompt[TValue]) Emit(event Event, args ...any) {
	for _, listener := range p.listeners[event] {
		listener(args...)
	}
}

// ParseKey parses a rune into a Key.
func (p *Prompt[TValue]) ParseKey(r rune) *Key {
	// TODO: parse Backtab(shift+tab) and other variations of shift and ctrl
	switch r {
	case '\r', '\n':
		return &Key{Name: EnterKey}
	case ' ':
		return &Key{Name: SpaceKey}
	case '\b', 127:
		return &Key{Name: BackspaceKey}
	case '\t':
		return &Key{Name: TabKey}
	case 3:
		return &Key{Name: CancelKey}
	case 27:
		next, err := p.rl.Peek(2)
		if err == nil && len(next) == 2 && next[0] == '[' {
			switch next[1] {
			case 'A':
				p.rl.Discard(2)
				return &Key{Name: UpKey}
			case 'B':
				p.rl.Discard(2)
				return &Key{Name: DownKey}
			case 'C':
				p.rl.Discard(2)
				return &Key{Name: RightKey}
			case 'D':
				p.rl.Discard(2)
				return &Key{Name: LeftKey}
			case 'H':
				p.rl.Discard(2)
				return &Key{Name: HomeKey}
			case 'F':
				p.rl.Discard(2)
				return &Key{Name: EndKey}
			}
		}
		return &Key{}
	default:
		char := string(r)
		return &Key{Char: char, Name: KeyName(char)}
	}
}

func (p *Prompt[TValue]) validate() error {
	if p.Validate == nil {
		return nil
	}

	p.IsValidating = true
	validationStart := time.Now()

	time.AfterFunc(400*time.Millisecond, func() {
		for p.IsValidating {
			p.State = ValidateState
			p.ValidationDuration = time.Since(validationStart)
			p.Emit(Event(ValidateState), p.ValidationDuration)
			p.render()
			time.Sleep(125 * time.Millisecond)
		}
	})

	err := p.Validate(p.Value)
	p.IsValidating = false

	return err
}

// PressKey handles key press events and updates the state of the prompt.
func (p *Prompt[TValue]) PressKey(key *Key) {
	if p.State == ErrorState || p.State == InitialState {
		p.State = ActiveState
	}

	p.Emit(KeyEvent, key)

	if key.Name == EnterKey {
		if err := p.validate(); err != nil {
			p.State = ErrorState
			p.Error = err.Error()
		} else {
			p.State = SubmitState
		}
	} else if key.Name == CancelKey {
		p.State = CancelState
	}

	if p.State == SubmitState || p.State == CancelState {
		p.Emit(FinalizeEvent)
	}

	p.render()
	p.Emit(Event(p.State), p.Value)
}

// TrackKeyValue updates the string value and cursor position based on key presses.
func (p *Prompt[TValue]) TrackKeyValue(key *Key, value string, cursorIndex int) (string, int) {
	switch key.Name {
	case BackspaceKey:
		if cursorIndex == 0 || len(value) == 0 {
			return value, cursorIndex
		}

		cursorIndex--
		if cursorIndex == len(value)-1 {
			return value[0:cursorIndex], cursorIndex
		}

		return value[0:cursorIndex] + value[cursorIndex+1:], cursorIndex
	case HomeKey:
		return value, 0
	case EndKey:
		return value, len(value)
	case LeftKey:
		return value, max(cursorIndex-1, 0)
	case RightKey:
		return value, min(cursorIndex+1, len(value))
	}

	if len(key.Char) == 1 {
		return value[0:cursorIndex] + key.Char + value[cursorIndex:], cursorIndex + 1
	}

	return value, cursorIndex
}

// LimitLines limits the number of lines to fit within the terminal size.
func (p *Prompt[TValue]) LimitLines(lines []string, usedLines int) string {
	_, maxRows, err := term.GetSize(int(p.output.Fd()))
	if err != nil {
		maxRows = 10
	}
	maxItems := max(min(maxRows-usedLines, len(lines)), 0)

	slidingWindowLocation := 0
	if p.CursorIndex >= maxItems-3 {
		slidingWindowLocation = max(min(p.CursorIndex-maxItems+3, len(lines)-maxItems), 0)
	} else if p.CursorIndex < 2 {
		slidingWindowLocation = max(p.CursorIndex-2, 0)
	}

	result := []string{}
	shouldRenderTopEllipsis := maxItems < len(lines) && slidingWindowLocation > 0
	shouldRenderBottomEllipsis := maxItems < len(lines) && slidingWindowLocation+maxItems < len(lines)

	for i, line := range lines[slidingWindowLocation : slidingWindowLocation+maxItems] {
		isTopLimit := i == 0 && shouldRenderTopEllipsis
		isBottomLimit := i == maxItems-1 && shouldRenderBottomEllipsis
		if isTopLimit || isBottomLimit {
			result = append(result, picocolors.Dim("..."))
			continue
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

type LineOption string

const (
	FirstLine LineOption = "firstLine"
	NewLine   LineOption = "newLine"
	LastLine  LineOption = "lastLine"
)

type FormatLineOptions struct {
	Start string
	End   string
	Sides string
	Style func(line string) string
}

type FormatLinesOptions struct {
	FirstLine FormatLineOptions
	NewLine   FormatLineOptions
	LastLine  FormatLineOptions
	Default   FormatLineOptions
	MinWidth  *int
	MaxWidth  *int
}

// getOptionOrDefault retrieves the option for the given line and option type, or returns a default value.
func getOptionOrDefault(line LineOption, opt string, options FormatLinesOptions) string {
	switch line {
	case FirstLine:
		switch opt {
		case "start":
			if options.FirstLine.Start != "" {
				return options.FirstLine.Start
			} else if options.FirstLine.Sides != "" {
				return options.FirstLine.Sides
			} else if options.Default.Start != "" {
				return options.Default.Start
			}
		case "end":
			if options.FirstLine.End != "" {
				return options.FirstLine.End
			} else if options.FirstLine.Sides != "" {
				return options.FirstLine.Sides
			} else if options.Default.End != "" {
				return options.Default.End
			}
		}
	case NewLine:
		switch opt {
		case "start":
			if options.NewLine.Start != "" {
				return options.NewLine.Start
			} else if options.NewLine.Sides != "" {
				return options.NewLine.Sides
			} else if options.Default.Start != "" {
				return options.Default.Start
			}
		case "end":
			if options.NewLine.End != "" {
				return options.NewLine.End
			} else if options.NewLine.Sides != "" {
				return options.NewLine.Sides
			} else if options.Default.End != "" {
				return options.Default.End
			}
		}
	case LastLine:
		switch opt {
		case "start":
			if options.LastLine.Start != "" {
				return options.LastLine.Start
			} else if options.LastLine.Sides != "" {
				return options.LastLine.Sides
			} else if options.Default.Start != "" {
				return options.Default.Start
			}
		case "end":
			if options.LastLine.End != "" {
				return options.LastLine.End
			} else if options.LastLine.Sides != "" {
				return options.LastLine.Sides
			} else if options.Default.End != "" {
				return options.Default.End
			}
		}
	}
	if options.Default.Sides != "" {
		return options.Default.Sides
	}
	return ""
}

func getStyleOrDefault(line LineOption, options FormatLinesOptions) func(line string) string {
	switch line {
	case FirstLine:
		if options.FirstLine.Style != nil {
			return options.FirstLine.Style
		}
	case NewLine:
		if options.NewLine.Style != nil {
			return options.NewLine.Style
		}
	case LastLine:
		if options.LastLine.Style != nil {
			return options.LastLine.Style
		}
	}
	if options.Default.Style != nil {
		return options.Default.Style
	}
	return func(line string) string { return line }
}

func getLineOptions(options FormatLinesOptions, line LineOption) FormatLineOptions {
	return FormatLineOptions{
		Start: getOptionOrDefault(line, "start", options),
		End:   getOptionOrDefault(line, "end", options),
		Style: getStyleOrDefault(line, options),
	}
}

func mergeOptions(primary, secondary FormatLineOptions) FormatLineOptions {
	ifEmpty := func(a, b string) string {
		if a == "" {
			return b
		}
		return a
	}

	ifNil := func(a, b func(line string) string) func(line string) string {
		if a == nil {
			return b
		}
		return a
	}

	return FormatLineOptions{
		Start: ifEmpty(primary.Start, secondary.Start),
		End:   ifEmpty(primary.End, secondary.End),
		Style: ifNil(primary.Style, secondary.Style),
	}
}

// FormatLines applies styles to multiple lines based on their type and the provided options.
func (p *Prompt[TValue]) FormatLines(lines []string, options FormatLinesOptions) string {
	terminalWidth, _, err := term.GetSize(int(p.output.Fd()))
	if err != nil {
		terminalWidth = 80
	}
	if options.MinWidth == nil {
		minWidth := 0
		options.MinWidth = &minWidth
	}
	if options.MaxWidth == nil {
		maxWidth := math.MaxInt
		options.MaxWidth = &maxWidth
	}
	minWidth := max(*options.MinWidth, 0)
	maxWith := min(*options.MaxWidth, terminalWidth)

	firstLine := getLineOptions(options, FirstLine)
	newLine := getLineOptions(options, NewLine)
	lastLine := getLineOptions(options, LastLine)

	formattedLines := []string{}
	for i, line := range lines {
		var opts FormatLineOptions
		if i == 0 && len(lines) == 1 {
			opts = mergeOptions(lastLine, firstLine)
		} else if i == 0 {
			opts = firstLine
		} else if i == 1 && len(lines) == 2 {
			opts = mergeOptions(lastLine, newLine)
		} else if i+1 == len(lines) {
			opts = lastLine
		} else {
			opts = newLine
		}

		// If empty slot > 0. Sum empty space
		startEmptySlots := utils.StrLength(opts.Start)
		if startEmptySlots != 0 {
			startEmptySlots++
		}
		endEmptySlots := utils.StrLength(opts.End)
		if endEmptySlots > 0 {
			endEmptySlots++
		}
		emptySlots := startEmptySlots + endEmptySlots + utils.StrLength(opts.Style(""))

		formatAndAddLine := func(line string) {
			var startSpace, endSpace string
			if startEmptySlots > 0 {
				startSpace = " "
			}
			styledLine := opts.Style(line)
			if minWidth > 0 && endEmptySlots > 0 {
				endSpace = strings.Repeat(" ", max(minWidth+1-startEmptySlots-utils.StrLength(styledLine)-endEmptySlots, 0))
			} else if minWidth > 0 {
				endSpace = strings.Repeat(" ", max(minWidth+1-startEmptySlots-utils.StrLength(styledLine), 0))
			} else if endEmptySlots > 0 {
				endSpace = " "
			}
			formattedLine := strings.Join([]string{opts.Start, startSpace, styledLine, endSpace, opts.End}, "")
			formattedLines = append(formattedLines, formattedLine)
		}

		if len(lines) == 1 && strings.Trim(lines[0], " ") == "" {
			formatAndAddLine("")
			break
		}

		var currentLine string
		for _, word := range strings.Split(line, " ") {
			if word == "" {
				currentLine += " "
			} else if strings.Trim(currentLine, " ") == "" && utils.StrLength(currentLine+word)+emptySlots <= maxWith {
				currentLine += word
			} else if utils.StrLength(currentLine+word)+emptySlots+1 <= maxWith {
				currentLine += " " + word
			} else if utils.StrLength(word)+emptySlots >= maxWith {
				var splitIndex int
				if utils.StrLength(currentLine) == 0 {
					splitIndex = maxWith - emptySlots
					formatAndAddLine(word[0:splitIndex])
				} else {
					splitIndex = maxWith - utils.StrLength(currentLine) - emptySlots - 1
					formatAndAddLine(currentLine + " " + word[0:splitIndex])
				}

				chunkLength := maxWith - emptySlots
				chunk := word[splitIndex:]
				for utils.StrLength(chunk) > chunkLength {
					formatAndAddLine(chunk[0:chunkLength])
					chunk = chunk[chunkLength:]
				}

				currentLine = chunk
			} else {
				formatAndAddLine(currentLine)
				currentLine = word
			}
		}
		formatAndAddLine(currentLine)
	}

	return strings.Join(formattedLines, "\r\n")
}

// DiffLines calculates the difference between an old and a new frame.
func (p *Prompt[TValue]) DiffLines(oldFrame, newFrame string) []int {
	var diff []int

	if oldFrame == newFrame {
		return diff
	}

	oldLines := strings.Split(oldFrame, "\n")
	newLines := strings.Split(newFrame, "\n")
	for i := range max(len(oldLines), len(newLines)) {
		if i >= len(oldLines) || i >= len(newLines) || oldLines[i] != newLines[i] {
			diff = append(diff, i)
		}
	}

	return diff
}

// render renders a new frame to the output.
func (p *Prompt[TValue]) render() {
	frame := p.Render(p)

	if lines := strings.Split(frame, "\r\n"); len(lines) == 1 {
		frame = strings.Join(strings.Split(frame, "\n"), "\r\n")
	}

	if p.State == InitialState {
		p.output.WriteString(sisteransi.HideCursor())
		p.output.WriteString(frame)
		p.Frame = frame
		return
	}

	if frame == p.Frame {
		return
	}

	diff := p.DiffLines(frame, p.Frame)
	diffLineIndex := diff[0]
	prevFrameLines := strings.Split((p.Frame), "\n")

	// Move to first diff line
	p.output.WriteString(sisteransi.MoveCursor(-(len(prevFrameLines) - 1), -999))
	p.output.WriteString(sisteransi.MoveCursor(diffLineIndex, 0))
	p.output.WriteString(sisteransi.EraseDown())
	lines := strings.Split(frame, "\n")
	newLines := lines[diffLineIndex:]
	p.output.WriteString(strings.Join(newLines, "\n"))
	p.Frame = frame
}

// Run runs the prompt and processes input.
func (p *Prompt[TValue]) Run() (TValue, error) {
	if flag.Lookup("test.v") == nil {
		oldState, err := term.MakeRaw(int(p.input.Fd()))
		if err != nil {
			return p.Value, err
		}
		defer term.Restore(int(p.input.Fd()), oldState)
	}

	done := make(chan struct{})
	closeCb := func(args ...any) {
		p.output.WriteString(sisteransi.ShowCursor())
		p.output.WriteString("\r\n")
		close(done)
	}
	p.Once(SubmitEvent, closeCb)
	p.Once(CancelEvent, closeCb)

	p.render()

outer:
	for {
		select {
		case <-done:
			break outer
		default:
			r, size, err := p.rl.ReadRune()
			if err != nil || size == 0 || p.IsValidating {
				continue
			}
			key := p.ParseKey(r)
			p.PressKey(key)
		}
	}

	if p.State == CancelState {
		return p.Value, ErrCancelPrompt
	}

	return p.Value, nil
}
