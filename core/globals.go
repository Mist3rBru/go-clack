package core

import "github.com/Mist3rBru/go-clack/core/utils"

var (
	color = utils.CreateColors()
)

type SelectOption[TValue comparable] struct {
	Label string
	Value TValue
}

type Event string

const (
	EventKey      Event = "key"
	EventFinalize Event = "finalize"
	EventCancel   Event = "cancel"
	EventSubmit   Event = "submit"
)

type State string

const (
	StateInitial State = "initial"
	StateActive  State = "active"
	StateError   State = "error"
	StateCancel  State = "cancel"
	StateSubmit  State = "submit"
)

type KeyName string

type Key struct {
	Name  KeyName
	Char  string
	Shift bool
	Ctrl  bool
}

const (
	KeyEnter     KeyName = "Enter"
	KeySpace     KeyName = "Space"
	KeyTab       KeyName = "Tab"
	KeyUp        KeyName = "Up"
	KeyDown      KeyName = "Down"
	KeyLeft      KeyName = "Left"
	KeyRight     KeyName = "Right"
	KeyCancel    KeyName = "Cancel"
	KeyHome      KeyName = "Home"
	KeyEnd       KeyName = "End"
	KeyBackspace KeyName = "Backspace"
)
