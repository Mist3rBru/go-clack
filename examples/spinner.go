package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Mist3rBru/go-clack/prompts"
)

func SpinnerExample() {
	prompts.Intro("spinner start...")

	s := prompts.Spinner(context.Background(), prompts.SpinnerOptions{})
	total := 10000
	progress := 0

	s.Start("")

	for progress < total {
		progress = min(total, progress+100)
		s.Message(fmt.Sprintf("Loading packages [%d/%d]", progress, total))

		time.Sleep(100 * time.Millisecond)
	}

	s.Stop("Done", 0)
	prompts.Outro("spinner stop...")

}
