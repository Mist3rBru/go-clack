package prompts

type Task struct {
	Title    string
	Task     func(message func(msg string)) (string, error)
	Disabled bool
}

func Tasks(tasks []Task, options SpinnerOptions) {
	for _, task := range tasks {
		if task.Disabled {
			continue
		}
		s := Spinner(options)
		s.Start(task.Title)
		result, err := task.Task(s.Message)
		if err != nil {
			s.Stop(err.Error(), 1)
			continue
		}
		s.Stop(result, 0)
	}
}
