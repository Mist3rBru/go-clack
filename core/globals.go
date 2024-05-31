package core

import "github.com/Mist3rBru/go-clack/core/utils"

var (
	color = utils.CreateColors()
)

type Listener func(args ...any)

type SelectOption struct {
	Label string
	Value any
}

type Key struct {
	Char  string
	Name  string
	Shift bool
	Ctrl  bool
}
