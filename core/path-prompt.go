package core

import (
	"os"
	"regexp"
	"strings"

	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/core/validator"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type PathPrompt struct {
	Prompt[string]
	OnlyShowDir bool
	Required    bool
	Placeholder string
	Hint        string
	HintOptions []string
	HintIndex   int
}

type PathPromptParams struct {
	Input        *os.File
	Output       *os.File
	InitialValue string
	OnlyShowDir  bool
	Required     bool
	Validate     func(value string) error
	Render       func(p *PathPrompt) string
}

func NewPathPrompt(params PathPromptParams) *PathPrompt {
	v := validator.NewValidator("PathPrompt")
	v.ValidateRender(params.Render)

	var p PathPrompt
	p = PathPrompt{
		Prompt: *NewPrompt(PromptParams[string]{
			Input:        params.Input,
			Output:       params.Output,
			InitialValue: params.InitialValue,
			CursorIndex:  len(params.InitialValue),
			Validate:     WrapValidateString(params.Validate, &p.Required, "Path does not exist! Please enter a valid path."),
			Render:       WrapRender[string](&p, params.Render),
		}),
		OnlyShowDir: params.OnlyShowDir,
		HintIndex:   -1,
		Required:    params.Required,
	}

	if cwd, err := os.Getwd(); err == nil && params.InitialValue == "" {
		p.Prompt.Value = cwd
		p.Value = cwd
		p.CursorIndex = len(cwd)
	}
	p.changeHint()

	p.On(KeyEvent, func(args ...any) {
		p.handleKeyPress(args[0].(*Key))
	})

	return &p
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

func (p *PathPrompt) ValueWithCursor() string {
	var (
		value string
		hint  string
	)
	if p.CursorIndex >= len(p.Value) {
		value = p.Value
		if p.Hint == "" {
			hint = picocolors.Inverse(" ")
		} else {
			hint = picocolors.Inverse(string(p.Hint[0])) + picocolors.Dim(p.Hint[1:])
		}
	} else {
		s1 := p.Value[0:p.CursorIndex]
		s2 := p.Value[p.CursorIndex:]
		value = s1 + picocolors.Inverse(string(s2[0])) + s2[1:]
		hint = picocolors.Dim(p.Hint)
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

func (p *PathPrompt) handleKeyPress(key *Key) {
	p.Value, p.CursorIndex = p.TrackKeyValue(key, p.Value, p.CursorIndex)
	if key.Name == RightKey && p.CursorIndex >= len(p.Value) {
		p.completeValue()
	} else if key.Name == TabKey {
		p.tabComplete()
	} else {
		p.changeHint()
	}
}
