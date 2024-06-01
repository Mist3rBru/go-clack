package core

import (
	"os"
	"regexp"
	"strings"

	"github.com/Mist3rBru/go-clack/core/utils"
)

type PathPrompt struct {
	Prompt
	OnlyShowDir bool
	Value       string
	Placeholder string
	Hint        string
	HintOptions []string
	HintIndex   int
}

type PathPromptParams struct {
	Input       *os.File
	Output      *os.File
	Value       string
	Placeholder string
	OnlyShowDir bool
	Validate    func(value string) error
	Render      func(p *PathPrompt) string
}

func NewPathPrompt(params PathPromptParams) *PathPrompt {
	var p *PathPrompt
	p = &PathPrompt{
		Prompt: *NewPrompt(PromptParams{
			Input:  params.Input,
			Output: params.Output,
			Value:  params.Value,
			Track:  true,
			Validate: func(value any) error {
				return params.Validate(p.Value)
			},
			Render: func(_p *Prompt) string {
				return params.Render(p)
			},
		}),
		Value:       params.Value,
		Placeholder: params.Placeholder,
		OnlyShowDir: params.OnlyShowDir,
		HintIndex:   -1,
	}
	if cwd, err := os.Getwd(); err == nil && params.Value == "" {
		p.Prompt.Value = cwd
		p.Value = cwd
		p.CursorIndex = len(cwd)
	}
	p.changeHint()

	p.On("key", func(args ...any) {
		key := args[0].(*Key)
		p.Value = p.Prompt.Value.(string)
		if key.Name == "Right" && p.CursorIndex >= len(p.Value) {
			p.completeValue()
		} else if key.Name == "Tab" {
			p.tabComplete()
		} else {
			p.changeHint()
		}
	})

	return p
}

func (p *PathPrompt) mapHintOptions() []string {
	options := []string{}
	dirPathRegex := regexp.MustCompile(`^(.*)/.*\s*`)
	dirPath := dirPathRegex.ReplaceAllString(p.Value, "$1")

	if strings.HasPrefix(dirPath, "~") {
		if homeDir, err := os.UserHomeDir(); err == nil {
			dirPath = strings.Replace(dirPath, "~", homeDir, 1)
		}
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return options
	}
	for _, entry := range entries {
		if (p.OnlyShowDir && !entry.IsDir()) || !strings.HasPrefix(entry.Name(), p.valueEnd()) {
			continue
		}
		if entry.IsDir() {
			options = append(options, entry.Name()+"/")
		} else {
			options = append(options, entry.Name())
		}
	}
	return options
}

func (p *PathPrompt) valueEnd() string {
	valueEndRegex := regexp.MustCompile("^.*/(.*)$")
	valueEnd := valueEndRegex.ReplaceAllString(p.Value, "$1")
	return valueEnd
}

func (p *PathPrompt) changeHint() {
	hintOptions := p.mapHintOptions()
	p.HintOptions = []string{}
	if len(hintOptions) > 0 {
		p.Hint = strings.Replace(hintOptions[0], p.valueEnd(), "", 1)
	} else {
		p.Hint = ""
	}
}

func (p *PathPrompt) ValueWithHint() string {
	var (
		value string
		hint  string
	)
	if p.CursorIndex >= len(p.Value) {
		value = p.Value
		if p.Hint == "" {
			hint = color["inverse"](" ")
		} else {
			hint = color["inverse"](string(p.Hint[0])) + color["dim"](p.Hint[1:])
		}
	} else {
		s1 := p.Value[0:p.CursorIndex]
		s2 := p.Value[p.CursorIndex:]
		value = s1 + color["inverse"](string(s2[0])) + s2[1:]
		hint = color["dim"](p.Hint)
	}
	return value + hint
}

func (p *PathPrompt) completeValue() {
	var complete string
	if p.Value == "" {
		complete = p.Placeholder
	} else {
		complete = p.Hint
	}
	p.Value += complete
	p.Prompt.Value = p.Value
	p.CursorIndex = len(p.Value)
	p.Hint = ""
	p.HintOptions = []string{}
	p.changeHint()
}

func (p *PathPrompt) tabComplete() {
	hintOption := p.mapHintOptions()
	if len(hintOption) == 1 {
		p.completeValue()
	} else if len(p.HintOptions) == 0 {
		p.HintOptions = hintOption
		p.HintIndex = 0
	} else {
		p.HintIndex = utils.MinMaxIndex(p.HintIndex+1, len(p.HintOptions))
		p.Hint = strings.Replace(p.HintOptions[p.HintIndex], p.valueEnd(), "", 1)
	}
}

func (p *PathPrompt) Run() (string, error) {
	_, err := p.Prompt.Run()
	if err != nil {
		return "", err
	}
	return p.Value, nil
}
