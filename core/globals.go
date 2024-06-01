package core

import "github.com/Mist3rBru/go-clack/core/utils"

var (
	color = utils.CreateColors()
)

type SelectOption[TValue comparable] struct {
	Label string
	Value TValue
}
