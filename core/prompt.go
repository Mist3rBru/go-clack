package core

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/Mist3rBru/go-clack/core/utils"

	"golang.org/x/term"
)

type Listener func(args ...any)

type Prompt[TValue any] struct {
	listeners map[Event][]Listener

	rl     *bufio.Reader
	input  *os.File
	output *os.File

	State       State
	Value       TValue
	Error       string
	CursorIndex int

	Validate func(value TValue) error
	Render   func(p *Prompt[TValue]) string
}

type PromptParams[TValue any] struct {
	Input       *os.File
	Output      *os.File
	Value       TValue
	CursorIndex int
	Validate    func(value TValue) error
	Render      func(p *Prompt[TValue]) string
}

func NewPrompt[TValue any](params PromptParams[TValue]) *Prompt[TValue] {
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

		State:       "initial",
		Value:       params.Value,
		CursorIndex: params.CursorIndex,

		Validate: params.Validate,
		Render:   params.Render,
	}
}

func (p *Prompt[TValue]) On(event Event, listener Listener) {
	p.listeners[event] = append(p.listeners[event], listener)
}

func (p *Prompt[TValue]) Once(event Event, listener Listener) {
	var onceListener Listener
	onceListener = func(args ...any) {
		listener(args)
		p.Off(event, onceListener)
	}
	p.On(event, onceListener)
}

func (p *Prompt[TValue]) Off(event Event, listener Listener) {
	listeners := p.listeners[event]
	for i, l := range listeners {
		if fmt.Sprintf("%p", l) == fmt.Sprintf("%p", listener) {
			p.listeners[event] = append(listeners[:i], listeners[i+1:]...)
			break
		}
	}
}

func (p *Prompt[TValue]) Emit(event Event, args ...any) {
	listeners := append([]Listener{}, p.listeners[event]...)
	for _, listener := range listeners {
		listener(args...)
	}
}

func (p *Prompt[TValue]) ParseKey(r rune) *Key {
	// TODO: parse Backtab(shift+tab) and other variations of shift and ctrl
	switch r {
	case '\r', '\n':
		return &Key{Name: KeyEnter}
	case ' ':
		return &Key{Name: KeySpace}
	case '\b', 127:
		return &Key{Name: KeyBackspace}
	case '\t':
		return &Key{Name: KeyTab}
	case 3:
		return &Key{Name: KeyCancel}
	case 27:
		next, err := p.rl.Peek(2)
		if err == nil && len(next) == 2 && next[0] == '[' {
			switch next[1] {
			case 'A':
				p.rl.Discard(2)
				return &Key{Name: KeyUp}
			case 'B':
				p.rl.Discard(2)
				return &Key{Name: KeyDown}
			case 'C':
				p.rl.Discard(2)
				return &Key{Name: KeyRight}
			case 'D':
				p.rl.Discard(2)
				return &Key{Name: KeyLeft}
			case 'H':
				p.rl.Discard(2)
				return &Key{Name: KeyHome}
			case 'F':
				p.rl.Discard(2)
				return &Key{Name: KeyEnd}
			}
		}
		return &Key{}
	default:
		char := string(r)
		return &Key{Char: char, Name: KeyName(char)}
	}
}

func (p *Prompt[TValue]) TrackKeyValue(key *Key, value string) string {
	switch key.Name {
	case KeyBackspace:
		if p.CursorIndex > 0 {
			if p.CursorIndex == len(value) {
				p.CursorIndex--
				value = value[0:p.CursorIndex]
			} else {
				p.CursorIndex--
				value = value[0:p.CursorIndex] + value[p.CursorIndex+1:]
			}
		}
	case KeyHome:
		p.CursorIndex = 0
	case KeyEnd:
		p.CursorIndex = len(value)
	case KeyLeft:
		if p.CursorIndex == 0 {
			break
		}
		p.CursorIndex--
	case KeyRight:
		if p.CursorIndex < len(value) {
			p.CursorIndex++
		}
	default:
		if len(key.Char) == 1 {
			value = value[0:p.CursorIndex] + key.Char + value[p.CursorIndex:]
			p.CursorIndex++
		}
	}
	return value
}

func (p *Prompt[TValue]) PressKey(key *Key) {
	if p.State == StateError || p.State == StateInitial {
		p.State = StateActive
	}

	p.Emit(EventKey, key)

	if key.Name == KeyEnter {
		if p.Validate != nil {
			err := p.Validate(p.Value)
			if err != nil {
				p.State = StateError
				p.Error = err.Error()
			}
		}
		if p.State != StateError {
			p.State = StateSubmit
		}
	}
	if key.Name == KeyCancel {
		p.State = StateCancel
	}
	if p.State == StateSubmit || p.State == StateCancel {
		p.Emit(EventFinalize)
	}
}

func (p *Prompt[TValue]) write(str string) {
	p.output.WriteString(str)
}

func (p *Prompt[TValue]) LimitLines(lines []string, usedLines int) string {
	_, maxRows, err := term.GetSize(int(p.output.Fd()))
	if err != nil {
		maxRows = 5
	}
	maxItems := min(maxRows-usedLines, len(lines))

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
			result = append(result, color["dim"]("..."))
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\r\n")
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
	return " "
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
		if a == " " {
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
	minWidth := max(*options.MinWidth, 1)
	maxWith := min(*options.MaxWidth, terminalWidth)

	firstLine := getLineOptions(options, FirstLine)
	newLine := getLineOptions(options, NewLine)
	lastLine := getLineOptions(options, LastLine)

	formattedLines := []string{}
	for i, line := range lines {
		var opts FormatLineOptions
		if i == 0 && len(lines) == 1 {
			opts = mergeOptions(firstLine, lastLine)
		} else if i == 0 {
			opts = firstLine
		} else if i == 1 && len(lines) == 2 {
			opts = mergeOptions(lastLine, newLine)
		} else if i+1 == len(lines) {
			opts = lastLine
		} else {
			opts = newLine
		}

		emptySlots := utils.StrLength(opts.Start+opts.End) + utils.StrLength(opts.Style("")) + 2
		formatAndAddLine := func(line string) {
			styledLine := opts.Style(line)
			fullLine := styledLine + strings.Repeat(" ", max(minWidth-utils.StrLength(styledLine)-emptySlots, 0))
			formattedLine := fmt.Sprintf("%s %s %s", opts.Start, fullLine, opts.End)
			formattedLines = append(formattedLines, formattedLine)
		}

		currentLine := ""
		for _, word := range strings.Split(line, " ") {
			if word == "" {
				currentLine += " "
			} else if currentLine == "" && utils.StrLength(word)+emptySlots <= maxWith {
				currentLine = word
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

func (p *Prompt[TValue]) render(prevFrame *string) {
	frame := p.Render(p)

	if lines := strings.Split(frame, "\r\n"); len(lines) == 1 {
		frame = strings.Join(strings.Split(frame, "\n"), "\r\n")
	}

	if p.State == StateInitial {
		p.write(utils.HideCursor())
		p.write(frame)
		p.State = StateActive
		*prevFrame = frame
		return
	}

	if frame == *prevFrame {
		return
	}

	diff := utils.DiffLines(frame, *prevFrame)
	diffLineIndex := diff[0]
	prevFrameLines := strings.Split((*prevFrame), "\n")

	// Move to first diff line
	p.write(utils.MoveCursor(-(len(prevFrameLines) - 1), -999))
	p.write(utils.MoveCursor(diffLineIndex, 0))

	if len(diff) == 1 {
		p.write(utils.EraseCurrentLine())
		lines := strings.Split(frame, "\n")
		p.write(lines[diffLineIndex])
		*prevFrame = frame
		p.write(utils.MoveCursorDown(len(lines) - diffLineIndex - 1))
		return
	}

	p.write(utils.EraseDown())
	lines := strings.Split(frame, "\n")
	newLines := lines[diffLineIndex:]
	p.write(strings.Join(newLines, "\n"))
	*prevFrame = frame
}

func (p *Prompt[TValue]) Run() (TValue, error) {
	oldState, err := term.MakeRaw(int(p.input.Fd()))
	if err != nil {
		return p.Value, err
	}
	defer term.Restore(int(p.input.Fd()), oldState)

	done := make(chan struct{})
	closeCb := func(args ...any) {
		p.write(utils.ShowCursor())
		p.write("\r\n")
		close(done)
	}
	p.Once(EventSubmit, closeCb)
	p.Once(EventCancel, closeCb)

	prevFrame := ""
	p.render(&prevFrame)

outer:
	for {
		select {
		case <-done:
			break outer
		default:
			r, size, err := p.rl.ReadRune()
			if err != nil {
				continue
			}
			if size == 0 {
				continue
			}
			key := p.ParseKey(r)
			p.PressKey(key)
			p.render(&prevFrame)
			p.Emit(Event(p.State), p.Value)
		}
	}

	if p.State == StateCancel {
		return p.Value, fmt.Errorf("prompt canceled")
	}

	return p.Value, nil
}
