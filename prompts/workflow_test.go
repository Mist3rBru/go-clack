package prompts_test

import (
	"errors"
	"testing"

	"github.com/Mist3rBru/go-clack/prompts"
	"github.com/stretchr/testify/assert"
)

type WorkflowResult struct {
	Name string
	Age  int
}

func TestWorkflowStep(t *testing.T) {
	var r WorkflowResult

	prompts.Workflow(&r).
		Step("Name", func() (any, error) {
			return "test name", nil
		}).
		Run()

	assert.Equal(t, "test name", r.Name)
}

func TestWorkflowStepError(t *testing.T) {
	var r WorkflowResult
	var calledTimes int

	err := prompts.Workflow(&r).
		Step("Name", func() (any, error) {
			calledTimes++
			return "", errors.New("test error")
		}).
		Step("Age", func() (any, error) {
			calledTimes++
			return 22, nil
		}).
		Run()

	assert.Equal(t, "test error", err.Error())
	assert.Equal(t, 1, calledTimes)
}

func TestWorkflowCaseInsensitive(t *testing.T) {
	var r WorkflowResult

	prompts.Workflow(&r).
		Step("name", func() (any, error) {
			return "test name", nil
		}).
		Step("age", func() (any, error) {
			return 22, nil
		}).
		Run()

	assert.Equal(t, "test name", r.Name)
	assert.Equal(t, 22, r.Age)
}
