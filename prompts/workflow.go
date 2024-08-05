package prompts

import (
	"fmt"
	"reflect"
	"strings"
)

type workflowStep struct {
	Name   string
	Prompt func() (any, error)
}

type workflowBuilder[TResult any] struct {
	steps  []*workflowStep
	result TResult
}

func (w *workflowBuilder[TResult]) Step(name string, prompt func() (any, error)) *workflowBuilder[TResult] {
	w.steps = append(w.steps, &workflowStep{Name: name, Prompt: prompt})
	return w
}

func (w *workflowBuilder[TResult]) Run() error {
	v := reflect.ValueOf(w.result).Elem()
	for i, step := range w.steps {
		stepResult, stepErr := step.Prompt()
		if stepErr != nil {
			return stepErr
		}

		stepResultVal := reflect.ValueOf(stepResult)
		if field := v.FieldByName(step.Name); field.IsValid() && field.CanSet() && stepResultVal.Type().AssignableTo(field.Type()) {
			field.Set(stepResultVal)
		} else if field := v.FieldByName(capitalize(step.Name)); field.IsValid() && field.CanSet() && stepResultVal.Type().AssignableTo(field.Type()) {
			field.Set(stepResultVal)
		} else {
			panic(fmt.Sprintf("workflow error: cannot set field `%s` from step `%d`", step.Name, i))
		}
	}
	return nil
}

func Workflow[TResult any](v TResult) *workflowBuilder[TResult] {
	return &workflowBuilder[TResult]{result: v}
}

func capitalize(str string) string {
	if str == "" {
		return str
	}
	return strings.ToUpper(string(str[0])) + strings.ToLower(str[1:])
}
