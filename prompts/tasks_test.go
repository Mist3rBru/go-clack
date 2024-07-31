package prompts_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Mist3rBru/go-clack/prompts"
	"github.com/Mist3rBru/go-clack/prompts/symbols"
	"github.com/stretchr/testify/assert"
)

func TestTasksStart(t *testing.T) {
	startTimes := 0
	task := func(message func(msg string)) (string, error) {
		startTimes++
		time.Sleep(time.Millisecond)
		return "", nil
	}
	timer := &MockTimer{autoResolve: true}
	writer := &MockWriter{}

	prompts.Tasks([]prompts.Task{
		{Title: "Foo", Task: task},
		{Title: "Bar", Task: task},
		{Title: "Baz", Task: task},
	}, prompts.SpinnerOptions{
		Timer:  timer,
		Output: writer,
	})

	expectedList := []string{
		"◒ Foo",
		"◒ Bar",
		"◒ Baz",
	}
	for _, expected := range expectedList {
		assert.Equal(t, expected, writer.HaveBeenCalledWith(expected))
	}
}

func TestTasksSubmit(t *testing.T) {
	startTimes := 0
	task := func(message func(msg string)) (string, error) {
		startTimes++
		time.Sleep(time.Millisecond)
		return "", nil
	}
	timer := &MockTimer{autoResolve: true}
	writer := &MockWriter{}

	prompts.Tasks([]prompts.Task{
		{Title: "Foo", Task: task},
		{Title: "Bar", Task: task},
		{Title: "Baz", Task: task},
	}, prompts.SpinnerOptions{
		Timer:  timer,
		Output: writer,
	})

	expectedList := []string{
		symbols.STEP_SUBMIT + " Foo\n",
		symbols.STEP_SUBMIT + " Bar\n",
		symbols.STEP_SUBMIT + " Baz\n",
	}
	for _, expected := range expectedList {
		assert.Equal(t, expected, writer.HaveBeenCalledWith(expected))
	}
}

func TestTasksUpdateMessage(t *testing.T) {
	task := func(message func(msg string)) (string, error) {
		message("Bar")
		time.Sleep(time.Millisecond)
		return "", nil
	}
	timer := &MockTimer{autoResolve: false}
	writer := &MockWriter{}

	prompts.Tasks([]prompts.Task{{Title: "Foo", Task: task}}, prompts.SpinnerOptions{Timer: timer, Output: writer})
	time.Sleep(time.Millisecond)
	timer.ResolveAll()
	time.Sleep(time.Millisecond)

	assert.NotEmpty(t, writer.HaveBeenCalledWith("◒ Bar"))
}

func TestTasksWithDisabledTask(t *testing.T) {
	counter := 0
	task := func(message func(msg string)) (string, error) {
		counter++
		return "", nil
	}
	timer := &MockTimer{autoResolve: true}
	writer := &MockWriter{}

	prompts.Tasks([]prompts.Task{
		{Title: "Foo", Task: task, Disabled: true},
	}, prompts.SpinnerOptions{
		Timer:  timer,
		Output: writer,
	})

	assert.Equal(t, 0, counter)
}

func TestTasksTaskWithError(t *testing.T) {
	task := func(message func(msg string)) (string, error) {
		return "", errors.New("task error")
	}
	timer := &MockTimer{autoResolve: false}
	writer := &MockWriter{}

	prompts.Tasks([]prompts.Task{{Title: "Foo", Task: task}}, prompts.SpinnerOptions{Timer: timer, Output: writer})
	time.Sleep(time.Millisecond)
	timer.ResolveAll()
	time.Sleep(time.Millisecond)

	assert.NotEmpty(t, writer.HaveBeenCalledWith(fmt.Sprintf("%s task error\n", symbols.STEP_CANCEL)))
}
