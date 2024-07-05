package core

import (
	"os"

	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/core/validator"
)

type GroupMultiSelectOption[TValue comparable] struct {
	MultiSelectOption[TValue]
	IsGroup bool
	Options []*GroupMultiSelectOption[TValue]
}

type GroupMultiSelectPrompt[TValue comparable] struct {
	Prompt[[]TValue]
	Options        []*GroupMultiSelectOption[TValue]
	DisabledGroups bool
	Required       bool
}

type GroupMultiSelectPromptParams[TValue comparable] struct {
	Input          *os.File
	Output         *os.File
	Options        map[string][]MultiSelectOption[TValue]
	InitialValue   []TValue
	DisabledGroups bool
	Required       bool
	Validate       func(value []TValue) error
	Render         func(p *GroupMultiSelectPrompt[TValue]) string
}

func NewGroupMultiSelectPrompt[TValue comparable](params GroupMultiSelectPromptParams[TValue]) *GroupMultiSelectPrompt[TValue] {
	v := validator.NewValidator("GroupMultiSelectPrompt")
	v.ValidateRender(params.Render)
	v.ValidateOptions(len(params.Options))

	options := []*GroupMultiSelectOption[TValue]{}
	for groupName, groupOptions := range params.Options {
		group := &GroupMultiSelectOption[TValue]{
			MultiSelectOption: MultiSelectOption[TValue]{
				Label: groupName,
			},
			IsGroup: true,
			Options: make([]*GroupMultiSelectOption[TValue], len(groupOptions)),
		}
		options = append(options, group)
		for i, groupOption := range groupOptions {
			option := &GroupMultiSelectOption[TValue]{
				MultiSelectOption: MultiSelectOption[TValue]{
					Label:      groupOption.Label,
					Value:      groupOption.Value,
					IsSelected: groupOption.IsSelected,
				},
			}
			if value, ok := any(option.Value).(string); ok && value == "" {
				option.Value = any(option.Label).(TValue)
			}
			group.Options[i] = option
			options = append(options, option)
		}
	}

	var initialValue []TValue
	if len(params.InitialValue) > 0 {
		initialValue = params.InitialValue
	} else {
		for _, option := range options {
			if option.IsSelected {
				initialValue = append(initialValue, option.Value)
			}
		}
	}

	var p GroupMultiSelectPrompt[TValue]
	p = GroupMultiSelectPrompt[TValue]{
		Prompt: *NewPrompt(PromptParams[[]TValue]{
			Input:        params.Input,
			Output:       params.Output,
			InitialValue: initialValue,
			Validate:     WrapValidateSlice(params.Validate, &p.Required, "Please select at least one option. Press `space` to select"),
			Render:       WrapRender[[]TValue](&p, params.Render),
		}),
		Options:        options,
		DisabledGroups: params.DisabledGroups,
		Required:       params.Required,
	}

	if p.DisabledGroups {
		p.CursorIndex = 1
	}

	p.On(KeyEvent, func(args ...any) {
		key := args[0].(*Key)
		switch key.Name {
		case UpKey, LeftKey:
			p.CursorIndex = utils.MinMaxIndex(p.CursorIndex-1, len(p.Options))
			if p.DisabledGroups && p.Options[p.CursorIndex].IsGroup {
				p.CursorIndex = utils.MinMaxIndex(p.CursorIndex-1, len(p.Options))
			}
		case DownKey, RightKey:
			p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+1, len(p.Options))
			if p.DisabledGroups && p.Options[p.CursorIndex].IsGroup {
				p.CursorIndex = utils.MinMaxIndex(p.CursorIndex+1, len(p.Options))
			}
		case HomeKey:
			p.CursorIndex = 0
		case EndKey:
			p.CursorIndex = len(p.Options) - 1
		case SpaceKey:
			p.toggleOption()
		}
	})
	return &p
}

func (p *GroupMultiSelectPrompt[TValue]) IsGroupSelected(group *GroupMultiSelectOption[TValue]) bool {
	if p.DisabledGroups {
		return false
	}
	for _, option := range group.Options {
		if !option.IsSelected {
			return false
		}

	}
	return true
}

func (p *GroupMultiSelectPrompt[TValue]) toggleOption() {
	option := p.Options[p.CursorIndex]
	if option.IsGroup {
		if p.IsGroupSelected(option) {
			for _, option := range option.Options {
				option.IsSelected = false
			}
			p.Value = []TValue{}
			for _, option := range option.Options {
				if option.IsSelected {
					p.Value = append(p.Value, option.Value)
				}
			}
		} else {
			for _, option := range option.Options {
				if !option.IsSelected {
					option.IsSelected = true
					p.Value = append(p.Value, option.Value)
				}
			}
		}
	} else {
		if option.IsSelected {
			option.IsSelected = false
			p.Value = []TValue{}
			for _, _option := range p.Options {
				if _option.IsSelected {
					p.Value = append(p.Value, _option.Value)
				}
			}
		} else {
			option.IsSelected = true
			p.Value = append(p.Value, option.Value)
		}
	}
}
