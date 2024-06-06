package prompts

import "context"

type Task struct {
	Title    string
	Task     func(message func(msg string)) (string, error)
	Disabled bool
}

func Tasks(ctx context.Context, tasks []Task, options SpinnerOptions) error {
	for _, task := range tasks {
		if task.Disabled {
			continue
		}
		s, err := Spinner(ctx, options)
		if err != nil {
			return err
		}
		s.Start(task.Title)
		result, err := task.Task(s.Message)
		if err == nil {
			s.Stop(result, 0)
		} else {
			s.Stop(err.Error(), 1)
		}
	}

	return nil
}
