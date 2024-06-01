package core

import "github.com/Mist3rBru/go-clack/core/utils"

var (
	color = utils.CreateColors()
)

type SelectOption[TValue comparable] struct {
	Label string
	Value TValue
}

type PromptEvent string

const (
	PromptEventKey      PromptEvent = "key"
	PromptEventFinalize PromptEvent = "finalize"
	PromptEventCancel   PromptEvent = "cancel"
	PromptEventSubmit   PromptEvent = "submit"
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
