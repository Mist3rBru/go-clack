package core

import (
	"os"

	"github.com/Mist3rBru/go-clack/core/utils"
)

type GroupSelectOption struct {
	Label   string
	Value   any
	IsGroup bool
	Options []*GroupSelectOption
}

type GroupMultiSelectPrompt struct {
	Prompt
	Options []*GroupSelectOption
	Value   []any
}

type GroupMultiSelectPromptParams struct {
	Input   *os.File
	Output  *os.File
	Options map[string][]SelectOption
	Value   []any
	Render  func(p *GroupMultiSelectPrompt) string
}

func NewGroupMultiSelectPrompt(params GroupMultiSelectPromptParams) *GroupMultiSelectPrompt {
	options := []*GroupSelectOption{}
	for groupName, groupOptions := range params.Options {
		group := &GroupSelectOption{
			Label:   groupName,
			IsGroup: true,
		}
		options = append(options, group)
		for _, groupOption := range groupOptions {
			option := &GroupSelectOption{
				Label:   groupOption.Label,
				Value:   groupOption.Value,
				IsGroup: false,
			}
			group.Options = append(group.Options, option)
			options = append(options, option)
		}
	}

	var p *GroupMultiSelectPrompt
	p = &GroupMultiSelectPrompt{
		Prompt: *NewPrompt(PromptParams{
			Input:       params.Input,
			Output:      params.Output,
			Value:       params.Value,
			CursorIndex: 0,
			Track:       false,
			Render: func(_p *Prompt) string {
				return params.Render(p)
			},
		}),
		Options: options,
		Value:   params.Value,
	}

	p.On("key", func(args ...any) {
		key := args[0].(*Key)
		switch key.Name {
		case "Up", "Left":
			p.CursorIndex = utils.MinMaxIndex(p.CursorIndex-1, len(p.Options))
		case "Down", "Right":
			p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+1, len(p.Options))
		case "Home":
			p.CursorIndex = 0
		case "End":
			p.CursorIndex = len(p.Options) - 1
		case "Space":
			p.toggleOption()
		}
	})
	return p
}

func (p *GroupMultiSelectPrompt) IsGroupSelected(group *GroupSelectOption) bool {
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

func (p *GroupMultiSelectPrompt) IsSelected(option *GroupSelectOption) (int, bool) {
	for i, v := range p.Value {
		if v == option.Value {
			return i, true
		}
	}
	return -1, false
}

func (p *GroupMultiSelectPrompt) toggleOption() {
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

func (p *GroupMultiSelectPrompt) Run() ([]any, error) {
	_, err := p.Prompt.Run()
	if err != nil {
		return nil, err
	}
	return p.Value, nil
}
