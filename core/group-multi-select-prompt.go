package core

import (
	"os"

	"github.com/Mist3rBru/go-clack/core/utils"
)

type GroupSelectOption[TValue comparable] struct {
	Label   string
	Value   TValue
	IsGroup bool
	Options []*GroupSelectOption[TValue]
}

type GroupMultiSelectPrompt[TValue comparable] struct {
	Prompt[[]TValue]
	Options []*GroupSelectOption[TValue]
}

type GroupMultiSelectPromptParams[TValue comparable] struct {
	Input   *os.File
	Output  *os.File
	Options map[string][]SelectOption[TValue]
	Value   []TValue
	Render  func(p *GroupMultiSelectPrompt[TValue]) string
}

func NewGroupMultiSelectPrompt[TValue comparable](params GroupMultiSelectPromptParams[TValue]) *GroupMultiSelectPrompt[TValue] {
	options := []*GroupSelectOption[TValue]{}
	for groupName, groupOptions := range params.Options {
		group := &GroupSelectOption[TValue]{
			Label:   groupName,
			IsGroup: true,
		}
		options = append(options, group)
		for _, groupOption := range groupOptions {
			option := &GroupSelectOption[TValue]{
				Label:   groupOption.Label,
				Value:   groupOption.Value,
				IsGroup: false,
			}
			group.Options = append(group.Options, option)
			options = append(options, option)
		}
	}

	var p *GroupMultiSelectPrompt[TValue]
	p = &GroupMultiSelectPrompt[TValue]{
		Prompt: *NewPrompt(PromptParams[[]TValue]{
			Input:  params.Input,
			Output: params.Output,
			Value:  params.Value,
			Render: func(_p *Prompt[[]TValue]) string {
				return params.Render(p)
			},
		}),
		Options: options,
	}

	p.On("key", func(args ...any) {
		key := args[0].(*Key)
		switch key.Name {
		case KeyUp, KeyLeft:
			p.CursorIndex = utils.MinMaxIndex(p.CursorIndex-1, len(p.Options))
		case KeyDown, KeyRight:
			p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+1, len(p.Options))
		case KeyHome:
			p.CursorIndex = 0
		case KeyEnd:
			p.CursorIndex = len(p.Options) - 1
		case KeySpace:
			p.toggleOption()
		}
	})
	return p
}

func (p *GroupMultiSelectPrompt[TValue]) IsGroupSelected(group *GroupSelectOption[TValue]) bool {
	counter := 0
	for _, option := range group.Options {
		for _, v := range p.Value {
			if option.Value == v {
				counter++
				break
			}
		}
	}
	return counter == len(group.Options)
}

func (p *GroupMultiSelectPrompt[TValue]) IsSelected(option *GroupSelectOption[TValue]) (int, bool) {
	for i, v := range p.Value {
		if v == option.Value {
			return i, true
		}
	}
	return -1, false
}

func (p *GroupMultiSelectPrompt[TValue]) toggleOption() {
	option := p.Options[p.CursorIndex]
	if option.IsGroup {
		if p.IsGroupSelected(option) {
			for _, option := range option.Options {
				if i, isSelected := p.IsSelected(option); isSelected {
					p.Value = append(p.Value[0:i], p.Value[i+1:]...)
				}
			}
		} else {
			for _, option := range option.Options {
				if _, IsSelected := p.IsSelected(option); !IsSelected {
					p.Value = append(p.Value, option.Value)
				}
			}
		}
	} else {
		if i, IsSelected := p.IsSelected(option); IsSelected {
			p.Value = append(p.Value[0:i], p.Value[i+1:]...)
		} else {
			p.Value = append(p.Value, option.Value)
		}
	}
}
