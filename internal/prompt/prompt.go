package prompt

import (
	"bufio"
	"fmt"
	"go-clack/internal/utils"
	"os"
	"strings"
	"sync"

	"golang.org/x/term"
)

var (
	color = utils.CreateColors()
)

type Listener func(args ...any)

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

func NewPrompt(input *os.File, output *os.File, track bool) *Prompt {
	return &Prompt{
		mu:        sync.Mutex{},
		listeners: make(map[string][]Listener),

		input:  input,
		output: output,
		rl:     bufio.NewReader(input),

		State:       "initial",
		Track:       track,
		CursorIndex: 0,
	}
}

func DefaultPrompt(track bool) *Prompt {
	return NewPrompt(os.Stdin, os.Stdout, track)
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

func (p *Prompt) ParseKey(r rune) (string, string) {
	switch r {
	case '\r', '\n':
		return "Enter", "\n"
	case '\b', 127:
		return "Backspace", ""
	case 27:
		next, err := p.rl.Peek(2)
		if err == nil && len(next) == 2 && next[0] == '[' {
			switch next[1] {
			case 'A':
				p.rl.Discard(2)
				return "ArrowUp", ""
			case 'B':
				p.rl.Discard(2)
				return "ArrowDown", ""
			case 'C':
				p.rl.Discard(2)
				return "ArrowRight", ""
			case 'D':
				p.rl.Discard(2)
				return "ArrowLeft", ""
			case 'H':
				p.rl.Discard(2)
				return "Home", ""
			case 'F':
				p.rl.Discard(2)
				return "End", ""
			}
		}
		return "", ""
	case 3:
		return "Cancel", ""
	default:
		char := string(r)
		return char, char
	}
}

func (p *Prompt) trackKeyValue(key, char, value string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	switch key {
	case "Enter", "Cancel":
		break
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
		p.Value = value[0:p.CursorIndex] + char + value[p.CursorIndex:]
		p.CursorIndex++
	}
}

func (p *Prompt) onKeypress(key string, char string) {
	if p.State == "error" {
		p.SetState("active")
	}

	strValue, ok := p.Value.(string)
	if p.Track && ok {
		p.trackKeyValue(key, char, strValue)
	}

	p.Emit("key", key, char)
	if key == "Enter" {
		if p.State != "error" {
			p.SetState("submit")
		}
	}
	if key == "Cancel" {
		p.SetState("cancel")
	}
}

func (p *Prompt) write(str string) {
	p.output.Write([]byte(str))
}

func (p *Prompt) render(prevFrame *string) {
	frame := p.Render(p)

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
		p.write(utils.MoveCursor(-1, 999))
		p.write(utils.ShowCursor())
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
				key, char := p.ParseKey(r)
				p.onKeypress(key, char)
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
