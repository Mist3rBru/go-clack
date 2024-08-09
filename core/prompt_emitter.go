package core

import "fmt"

type Event int

const (
	// KeyEvent is emitted after each user's input
	KeyEvent Event = iota
	// ValidateEvent is emitted when the input is being validated
	ValidateEvent
	// ErrorEvent is emitted if an error occurs during the validation process
	ErrorEvent
	// FinalizeEvent is emitted on user's submit or cancel, and before rendering the related state
	FinalizeEvent
	// CancelEvent is emitted after the user cancels the prompt, and after rendering the cancel state
	CancelEvent
	// SubmitEvent is emitted after the user submits the input, and after rendering the submit state
	SubmitEvent
)

type Listener func(args ...any)

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
