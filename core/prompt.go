package core

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/Mist3rBru/go-clack/core/utils"

	"golang.org/x/term"
)

type Listener func(args ...any)

type Key struct {
	Char  string
	Name  string
	Shift bool
	Ctrl  bool
}

type Prompt[TValue any] struct {
	mu        sync.Mutex
	listeners map[string][]Listener

	rl     *bufio.Reader
	input  *os.File
	output *os.File

	State       string
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
	return &Prompt[TValue]{
		mu:        sync.Mutex{},
		listeners: make(map[string][]Listener),

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

func (p *Prompt[TValue]) On(event string, listener Listener) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.listeners[event] = append(p.listeners[event], listener)
}

func (p *Prompt[TValue]) Once(event string, listener Listener) {
	var onceListener Listener
	onceListener = func(args ...any) {
		listener(args)
		p.Off(event, onceListener)
	}
	p.On(event, onceListener)
}

func (p *Prompt[TValue]) Off(event string, listener Listener) {
	p.mu.Lock()
	defer p.mu.Unlock()
	listeners := p.listeners[event]
	for i, l := range listeners {
		if fmt.Sprintf("%p", l) == fmt.Sprintf("%p", listener) {
			p.listeners[event] = append(listeners[:i], listeners[i+1:]...)
			break
		}
	}
}

func (p *Prompt[TValue]) Emit(event string, args ...any) {
	p.mu.Lock()
	listeners := append([]Listener{}, p.listeners[event]...)
	p.mu.Unlock()
	for _, listener := range listeners {
		listener(args...)
	}
}

func (p *Prompt[TValue]) SetState(state string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.State = state
}

func (p *Prompt[TValue]) SetValue(value TValue) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Value = value
}

func (p *Prompt[TValue]) SetError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Error = err.Error()
}

func (p *Prompt[TValue]) ParseKey(r rune) *Key {
	// TODO: parse Backtab(shift+tab) and other variations of shift and ctrl
	switch r {
	case '\r', '\n':
		return &Key{Name: "Enter"}
	case ' ':
		return &Key{Name: "Space"}
	case '\b', 127:
		return &Key{Name: "Backspace"}
	case '\t':
		return &Key{Name: "Tab"}
	case 3:
		return &Key{Name: "Cancel"}
	case 27:
		next, err := p.rl.Peek(2)
		if err == nil && len(next) == 2 && next[0] == '[' {
			switch next[1] {
			case 'A':
				p.rl.Discard(2)
				return &Key{Name: "Up"}
			case 'B':
				p.rl.Discard(2)
				return &Key{Name: "Down"}
			case 'C':
				p.rl.Discard(2)
				return &Key{Name: "Right"}
			case 'D':
				p.rl.Discard(2)
				return &Key{Name: "Left"}
			case 'H':
				p.rl.Discard(2)
				return &Key{Name: "Home"}
			case 'F':
				p.rl.Discard(2)
				return &Key{Name: "End"}
			}
		}
		return &Key{}
	default:
		char := string(r)
		return &Key{Char: char, Name: char}
	}
}

func (p *Prompt[TValue]) TrackKeyValue(key *Key, value string) string {
	switch key.Name {
	case "Backspace":
		if p.CursorIndex > 0 {
			if p.CursorIndex == len(value) {
				p.CursorIndex--
				value = value[0:p.CursorIndex]
			} else {
				p.CursorIndex--
				value = value[0:p.CursorIndex] + value[p.CursorIndex+1:]
			}
		}
	case "Home":
		p.CursorIndex = 0
	case "End":
		p.CursorIndex = len(value)
	case "Left":
		if p.CursorIndex == 0 {
			break
		}
		p.CursorIndex--
	case "Right":
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
	if p.State == "error" || p.State == "initial" {
		p.SetState("active")
	}

	p.Emit("key", key)

	if key.Name == "Enter" {
		if p.Validate != nil {
			err := p.Validate(p.Value)
			if err != nil {
				p.SetState("error")
				p.SetError(err)
			}
		}
		if p.State != "error" {
			p.SetState("submit")
		}
	}
	if key.Name == "Cancel" {
		p.SetState("cancel")
	}
	if p.State == "submit" || p.State == "cancel" {
		p.Emit("finalize")
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

func (p *Prompt[TValue]) render(prevFrame *string) {
	frame := p.Render(p)

	if lines := strings.Split(frame, "\r\n"); len(lines) == 1 {
		frame = strings.Join(strings.Split(frame, "\n"), "\r\n")
	}

	if p.State == "initial" {
		p.write(utils.HideCursor())
		p.write(frame)
		p.SetState("active")
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
	p.Once("submit", closeCb)
	p.Once("cancel", closeCb)

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
			p.Emit(p.State, p.Value)
		}
	}

	if ref := reflect.ValueOf(p.Value); ref.IsNil() {
		return p.Value, fmt.Errorf("Prompt canceled")
	}

	return p.Value, nil
}
