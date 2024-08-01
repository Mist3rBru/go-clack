package core_test

import (
	"errors"
	"testing"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/stretchr/testify/assert"
)

func TestWrapRender(t *testing.T) {
	prompt := core.NewPrompt(core.PromptParams[string]{
		InitialValue: "value",
	})
	render := func(p *core.Prompt[string]) string {
		return "frame: " + p.Value
	}
	renderWrapper := core.WrapRender[string](prompt, render)

	assert.Equal(t, "frame: value", renderWrapper(nil))
}

func TestWrapValidate(t *testing.T) {
	required := true
	errMsg := "test error"
	err := errors.New(errMsg)
	validate := core.WrapValidate[any](nil, &required, errMsg)

	assert.Equal(t, nil, validate(false))
	assert.Equal(t, nil, validate(true))

	assert.Equal(t, err, validate(0))
	assert.Equal(t, nil, validate(1))

	assert.Equal(t, err, validate(""))
	assert.Equal(t, nil, validate("a"))

	assert.Equal(t, err, validate([]any{}))
	assert.Equal(t, nil, validate([]any{"a"}))

	assert.Equal(t, err, validate(nil))
	assert.Equal(t, err, validate(struct{}{}))
	assert.Equal(t, nil, validate(struct{ a string }{a: "a"}))

	assert.Equal(t, err, validate((*any)(nil)))
	assert.Equal(t, nil, validate(&required))

	assert.Equal(t, err, validate([0]any{}))
	assert.Equal(t, err, validate([1]any{}))
	assert.Equal(t, nil, validate([1]any{1}))

	assert.Equal(t, err, validate(map[string]string{}))
	assert.Equal(t, nil, validate(map[string]string{"a": "b"}))

	assert.Equal(t, nil, validate(make(chan any)))

	type CustomType struct {
		Field string
	}
	assert.Equal(t, err, validate(CustomType{}))
	assert.Equal(t, nil, validate(CustomType{Field: "a"}))
}
