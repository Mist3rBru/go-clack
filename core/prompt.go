package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Mist3rBru/go-clack/core/utils"

	"golang.org/x/term"
)

type Prompt struct {
	mu        sync.Mutex
	listeners map[string][]Listener

	rl     *bufio.Reader
	input  *os.File
	output *os.File

	State       string
	Value       any
	Track       bool
	CursorIndex int

	Render func(p *Prompt) string
}

type PromptParams struct {
	Input       *os.File
	Output      *os.File
	Value       any
	CursorIndex int
	Track       bool
	Render      func(p *Prompt) string
}

func NewPrompt(params PromptParams) *Prompt {
	if strValue, ok := params.Value.(string); ok {
		params.CursorIndex = len(strValue)
	}
	return &Prompt{
		mu:        sync.Mutex{},
		listeners: make(map[string][]Listener),

		input:  params.Input,
		output: params.Output,
		rl:     bufio.NewReader(params.Input),

		State:       "initial",
		Value:       params.Value,
		CursorIndex: params.CursorIndex,
		Track:       params.Track,

		Render: params.Render,
	}
}

func (p *Prompt) On(event string, listener Listener) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.listeners[event] = append(p.listeners[event], listener)
}

func (p *Prompt) Once(event string, listener Listener) {
	var onceListener Listener
	onceListener = func(args ...any) {
		listener(args...)
		p.Off(event, onceListener)
	}
	p.On(event, onceListener)
}

func (p *Prompt) Off(event string, listener Listener) {
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

func (p *Prompt) Emit(event string, args ...any) {
	p.mu.Lock()
	listeners := append([]Listener(nil), p.listeners[event]...)
	p.mu.Unlock()
	for _, listener := range listeners {
		listener(args...)
	}
}

func (p *Prompt) SetState(state string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.State = state
}

func (p *Prompt) SetValue(value any) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Value = value
}

func (p *Prompt) ParseKey(r rune) *Key {
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
				return &Key{Name: "ArrowUp"}
			case 'B':
				p.rl.Discard(2)
				return &Key{Name: "ArrowDown"}
			case 'C':
				p.rl.Discard(2)
				return &Key{Name: "ArrowRight"}
			case 'D':
				p.rl.Discard(2)
				return &Key{Name: "ArrowLeft"}
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

func (p *Prompt) trackKeyValue(key *Key, value string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	switch key.Name {
	case "Backspace":
		if p.CursorIndex == 0 {
			break
		}
		if p.CursorIndex == len(value) {
			p.CursorIndex--
			p.Value = value[0:p.CursorIndex]
			break
		}
		p.CursorIndex--
		p.Value = value[0:p.CursorIndex] + value[p.CursorIndex+1:]
	case "Home":
		p.CursorIndex = 0
	case "End":
		p.CursorIndex = len(value)
	case "ArrowLeft":
		if p.CursorIndex == 0 {
			break
		}
		p.CursorIndex--
	case "ArrowRight":
		if p.CursorIndex == len(value) {
			break
		}
		p.CursorIndex++
	default:
		if key.Char != "" {
			p.Value = value[0:p.CursorIndex] + key.Char + value[p.CursorIndex:]
			p.CursorIndex++
		}
	}
}

func (p *Prompt) onKeypress(key *Key) {
	if p.State == "error" {
		p.SetState("active")
	}

	strValue, ok := p.Value.(string)
	if p.Track && ok {
		p.trackKeyValue(key, strValue)
	}

	p.Emit("key", key)
	if key.Name == "Enter" {
		if p.State != "error" {
			p.SetState("submit")
		}
	}
	if key.Name == "Cancel" {
		p.SetState("cancel")
	}
}

func (p *Prompt) write(str string) {
	p.output.Write([]byte(str))
}

type LimitLinesPamams struct {
	CursorIndex int
	Lines       []string
}

func (p *Prompt) LimitLines(params LimitLinesPamams) string {
	_, maxRows, _ := term.GetSize(int(p.output.Fd()))
	maxItems := min(maxRows, len(params.Lines))

	slidingWindowLocation := 0
	if params.CursorIndex >= maxItems-3 {
		slidingWindowLocation = max(min(params.CursorIndex-maxItems+3, len(params.Lines)-maxItems), 0)
	} else if params.CursorIndex < 2 {
		slidingWindowLocation = max(params.CursorIndex-2, 0)
	}

	result := []string{}
	shouldRenderTopEllipsis := maxItems < len(params.Lines) && slidingWindowLocation > 0
	shouldRenderBottomEllipsis := maxItems < len(params.Lines) && slidingWindowLocation+maxItems < len(params.Lines)

	for i, line := range params.Lines[slidingWindowLocation : slidingWindowLocation+maxItems] {
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

func (p *Prompt) render(prevFrame *string) {
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

func (p *Prompt) Run() (any, error) {
	oldState, err := term.MakeRaw(int(p.input.Fd()))
	if err != nil {
		return nil, err
	}
	defer term.Restore(int(p.input.Fd()), oldState)

	wg := sync.WaitGroup{}
	done := make(chan struct{})

	closeCb := func(args ...any) {
		p.write(utils.ShowCursor())
		p.write("\r\n")
		close(done)
	}
	p.Once("submit", closeCb)
	p.Once("cancel", closeCb)

	wg.Add(1)
	go func() {
		defer wg.Done()
		prevFrame := ""
		p.render(&prevFrame)

		for {
			select {
			case <-done:
				return
			default:
				r, size, err := p.rl.ReadRune()
				if err != nil {
					continue
				}
				if size == 0 {
					continue
				}
				key := p.ParseKey(r)
				p.onKeypress(key)
				p.render(&prevFrame)
				p.Emit(p.State, p.Value)
			}
		}
	}()

	wg.Wait()

	if p.Value == nil {
		return nil, fmt.Errorf("Prompt canceled")
	}

	return p.Value, nil
}
