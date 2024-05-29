package prompt

import (
	"bufio"
	"fmt"
	"go-clack/internal/utils"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"golang.org/x/term"
)

type Listener func(args ...any)

type Prompt struct {
	mu        sync.Mutex
	listeners map[string][]Listener

	rl     *bufio.Reader
	input  *os.File
	output *os.File

	State string
	Value any
	Track bool

	Render func(p *Prompt) string
}

func NewPrompt(input *os.File, output *os.File, track bool) *Prompt {
	return &Prompt{
		mu:        sync.Mutex{},
		listeners: make(map[string][]Listener),

		input:  input,
		output: output,
		rl:     bufio.NewReader(input),

		State: "initial",
		Track: track,
	}
}

func DefaultPrompt(track bool) *Prompt {
	return &Prompt{
		mu:        sync.Mutex{},
		listeners: make(map[string][]Listener),

		input:  os.Stdin,
		output: os.Stdout,
		rl:     bufio.NewReader(os.Stdin),

		State: "initial",
		Track: track,
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

func (p *Prompt) ParseKey(r rune, next []byte) (string, string) {
	switch r {
	case '\r', '\n':
		return "Enter", "\n"
	case 27:
		if len(next) == 2 && next[0] == '[' {
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
			}
		}
		return "Escape", ""
	case 3:
		return "Cancel", ""
	default:
		char := string(r)
		return char, char
	}
}

func (p *Prompt) onKeypress(key string, char string) {
	if p.State == "error" {
		p.SetState("active")
	}

	strValue, ok := p.Value.(string)
	if p.Track && ok {
		p.SetValue(strValue + char)
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

	p.Emit(p.State, p.Value)
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

	if len(diff) == 1 {
		p.write(utils.MoveCursor(diffLineIndex, 0))
		p.write(utils.EraseCurrentLine())
		lines := strings.Split(frame, "\n")
		p.write(lines[diffLineIndex])
		*prevFrame = frame
		p.write(utils.MoveCursorDown(len(lines) - diffLineIndex - 1))
		return
	}

	p.write(utils.MoveCursor(diffLineIndex, 0))
	p.write(utils.EraseDown())
	lines := strings.Split(frame, "\n")
	newLines := lines[diffLineIndex:]
	p.write(strings.Join(newLines, "\n"))
	*prevFrame = frame
}

func (p *Prompt) Prompt() (any, error) {
	oldState, err := term.MakeRaw(int(p.input.Fd()))
	if err != nil {
		return nil, err
	}
	defer term.Restore(int(p.input.Fd()), oldState)

	wg := sync.WaitGroup{}
	done := make(chan struct{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	p.Once("submit", func(args ...any) {
		p.write("\n")
		p.write(utils.ShowCursor())
		close(done)
	})

	p.Once("cancel", func(args ...any) {
		p.write("\n")
		p.write(utils.ShowCursor())
		p.SetValue(nil)
		close(done)
	})

	wg.Add(1)
	go func() {
		defer wg.Done()
		prevFrame := ""

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
				nextKeyByte, _ := p.rl.Peek(2)
				key, char := p.ParseKey(r, nextKeyByte)
				p.onKeypress(key, char)
				p.render(&prevFrame)
			}
		}
	}()

	wg.Wait()
	close(sig)

	if p.Value == nil {
		return nil, fmt.Errorf("Prompt canceled")
	}

	return p.Value, nil
}
